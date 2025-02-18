package server

import (
	auth "VirtualRegistryManagement/authentication"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/eventpublish"
	epGRPC "VirtualRegistryManagement/eventpublish/grpc"
	"VirtualRegistryManagement/logger"
	midUnary "VirtualRegistryManagement/middlewares/unary"
	"VirtualRegistryManagement/modules/filepath"
	"VirtualRegistryManagement/modules/openstack"
	"VirtualRegistryManagement/services"
	"VirtualRegistryManagement/storages"
	"VirtualRegistryManagement/utility/workerpool"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/tracer"
	vpsUtil "pegasus-cloud.com/aes/virtualplatformserviceclient/utility"
	"pegasus-cloud.com/aes/virtualplatformserviceclient/vps"
	vrm "pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

var srvHTTP *http.Server
var srvGRPC *grpc.Server

func Run() {
	start()
	signAction()
	stop()
}

func start() {
	prepareUpstreamServices()
	startGRPCServer()
	startUpstreamServices()
	startHTTPServer()
}

func stop() {
	stopHTTPServer()
	stopGRPCServer()
	stopUpstreamServices()
}

func prepareUpstreamServices() {
	{ // logger
		// 初始化 Logger
		logger.Init(fmt.Sprintf("%s.log", mviper.GetString("VirtualRegistryManagement.instance")))
		// 初始化 Access Logger
		logger.InitAccessLogger(fmt.Sprintf("%s_access.log", mviper.GetString("VirtualRegistryManagement.instance")))
	}

	{ // tracer (jaeger)
		if mviper.GetBool("VirtualRegistryManagement.tracer.enable") {
			tracer.Init(&tracer.Config{
				ServiceName: cnt.Kind,
				Endpoint:    mviper.GetString("VirtualRegistryManagement.tracer.host"),
				Timeout:     mviper.GetInt("VirtualRegistryManagement.tracer.timeout"),
			})
		}
	}

	{ // storages (database)
		storages.New(mviper.GetString("storage.provider"))

		// migrate database to migrate map latest version
		if mviper.GetBool("storage.auto_migrate") {
			if err := storages.Exec().AutoMigration(); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func startUpstreamServices() {
	{ // services
		// 初始化 Services
		if err := services.InitServices(); err != nil {
			fmt.Println(tkErr.New(cnt.ServerInternalServerErr).WithInner(err))
			os.Exit(1)
		}
	}

	{ // authentication
		auth.Init(mviper.GetString("authentication.service"))
	}

	{
		// Initialize all of event publish plugins
		eventpublish.InitAllPlugins()
	}

	{ // worker pool
		workerpool.InitPool(
			mviper.GetInt("VirtualRegistryManagement.worker_pool.max_workers"),
			mviper.GetInt("VirtualRegistryManagement.worker_pool.max_capacity"),
		)
	}

	{ // openstack
		if err := openstack.Init(mviper.Get("openstack")); err != nil {
			fmt.Println(tkErr.New(cnt.ServerInternalServerErr).WithInner(err))
			os.Exit(1)
		}
	}

	{ // vps
		vps.ReplaceGlobals(&vps.ConnProvider{
			Mode: vps.TCPMode,
			TCPProvider: vps.TCPProvider{
				Host: mviper.GetString("vps.host"),
				TLS: vps.TLSConfig{
					Enable:   mviper.GetBool("vps.tls.enable"),
					CertPath: mviper.GetString("vps.tls.cert_path"),
				},
				ConnPerHost: mviper.GetInt("vps.connection_per_host"),
			},
			RouteResponseType: vpsUtil.JSON,
		})
	}

	{ // filepath
		filepath.Init(mviper.Get("VirtualRegistryManagement.filepath"))
	}
}

func startGRPCServer() {
	zap.L().Info(serverStartInfo("GRPC", mviper.GetString("VirtualRegistryManagement.grpc.host"), mviper.GetBool("VirtualRegistryManagement.tls.enable")))

	lis, err := net.Listen("tcp", mviper.GetString("VirtualRegistryManagement.grpc.host"))
	if err != nil {
		zap.L().Error(fmt.Sprintf("failed to listen of GRPC Port: %v, %s", err, mviper.GetString("VirtualRegistryManagement.grpc.host")))
	}
	grpcOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(tracer.RequestIDParser(), midUnary.ActionParser()),
		grpc.UnaryInterceptor(tracer.NewGRPCUnaryServerInterceptor()),
		grpc.WriteBufferSize(mviper.GetInt("VirtualRegistryManagement.grpc.write_buffer_size")),
		grpc.ReadBufferSize(mviper.GetInt("VirtualRegistryManagement.grpc.read_buffer_size")),
		grpc.MaxRecvMsgSize(mviper.GetInt("VirtualRegistryManagement.grpc.max_receive_message_size")),
		grpc.MaxSendMsgSize(mviper.GetInt("VirtualRegistryManagement.grpc.max_send_message_size")),
	}

	if mviper.GetBool("VirtualRegistryManagement.tls.enable") {
		c, err := credentials.NewServerTLSFromFile(mviper.GetString("VirtualRegistryManagement.tls.cert_path"), mviper.GetString("VirtualRegistryManagement.tls.key_path"))
		if err != nil {
			log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
		}
		grpcOptions = append(grpcOptions, grpc.Creds(c))
	}

	srvGRPC = grpc.NewServer(grpcOptions...)
	grpcRouters(srvGRPC)
	go func() {
		if err := srvGRPC.Serve(lis); err != nil {
			zap.L().Error(fmt.Sprintf("failed to GRPC Server: %v", err))
		}
	}()
	os.Remove(mviper.GetString("VirtualRegistryManagement.grpc.unix_socket.path"))
	lis2, err := net.Listen("unix", mviper.GetString("VirtualRegistryManagement.grpc.unix_socket.path"))
	if err != nil {
		zap.L().Error(fmt.Sprintf("failed to listen of GRPC Port: %v, %s", err, mviper.GetString("VirtualRegistryManagement.grpc.host")))
	}
	go func() {
		if err := srvGRPC.Serve(lis2); err != nil {
			zap.L().Error(fmt.Sprintf("failed to GRPC Server: %v", err))
		}
	}()
	vrm.Init(vrm.PoolProvider{
		Mode: vrm.UnixMode,
		UnixProvider: vrm.UnixProvider{
			SocketPath: mviper.GetString("VirtualRegistryManagement.grpc.unix_socket.path"),
			ConnCount:  mviper.GetInt("VirtualRegistryManagement.grpc.unix_socket.conn_count"),
		},
		WriteBufferSize:       mviper.GetInt("VirtualRegistryManagement.grpc.write_buffer_size"),
		ReadBufferSize:        mviper.GetInt("VirtualRegistryManagement.grpc.read_buffer_size"),
		MaxReceiveMessageSize: mviper.GetInt("VirtualRegistryManagement.grpc.max_receive_message_size"),
		MaxSendMessageSize:    mviper.GetInt("VirtualRegistryManagement.grpc.max_send_message_size"),
	})
}

func startHTTPServer() {
	zap.L().Info(serverStartInfo("HTTP",
		mviper.GetString("VirtualRegistryManagement.http.host"),
		mviper.GetBool("VirtualRegistryManagement.tls.enable")))

	if mviper.GetBool("VirtualRegistryManagement.tls.enable") {
		srvHTTP = &http.Server{
			Addr:    mviper.GetString("VirtualRegistryManagement.http.host"),
			Handler: router(),
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS13,
			},
		}
		go func() {
			if err := srvHTTP.ListenAndServeTLS(
				mviper.GetString("VirtualRegistryManagement.tls.cert_path"),
				mviper.GetString("VirtualRegistryManagement.tls.key_path")); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
	} else {
		srvHTTP = &http.Server{
			Addr:    mviper.GetString("VirtualRegistryManagement.http.host"),
			Handler: router(),
		}
		go func() {
			if err := srvHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
	}
}

func signAction() {
	quit := make(chan os.Signal)
	//lint:ignore SA1017 ...
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	//lint:ignore S1000 ...
	for {
		select {
		case sign := <-quit:
			switch sign {
			case syscall.SIGINT, syscall.SIGTERM:
				zap.L().Info("Shutdown Service.")
				return
			case syscall.SIGUSR2:
				stop()
				start()
			default:
				fmt.Println("Other Signal ", sign)
				return
			}
		}
	}
}

func stopHTTPServer() {
	srvHTTP.Shutdown(context.Background())
}
func stopGRPCServer() {
	srvGRPC.Stop()
}
func stopUpstreamServices() {
	epGRPC.ClosePlugins()
	tracer.Shutdown()
}

func serverStartInfo(serverName, host string, tlsEnable bool) string {
	tls := "disabled"
	if tlsEnable {
		tls = "enabled"
	}
	return fmt.Sprintf("%s server started at %s and TLS is %s", serverName, host, tls)
}
