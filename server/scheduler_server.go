package server

import (
	auth "VirtualRegistryManagement/authentication"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/logger"
	"VirtualRegistryManagement/modules/eventconsume"
	"VirtualRegistryManagement/modules/lbmevents"
	"VirtualRegistryManagement/services"
	"VirtualRegistryManagement/tasks"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/meteringtoolkits/metering"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/tracer"
	vrm "pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

func RunScheduler() {
	startScheduler()
	signActionForScheduler()
	stopScheduler()
}

func startScheduler() {
	prepareUpstreamServicesForScheduler()
	startUpstreamServicesForScheduler()
	prepareSchedulerServer()
	startSchedulerServer()
}

func stopScheduler() {
	stopSchedulerServer()
	stopUpstreamServicesForScheduler()
}

func signActionForScheduler() {
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
				stopScheduler()
				startScheduler()
			default:
				fmt.Println("Other Signal ", sign)
				return
			}
		}
	}
}

func prepareUpstreamServicesForScheduler() {
	{ // logger
		// 初始化 Logger
		logger.Init(fmt.Sprintf("%s.log", mviper.GetString("VirtualRegistryManagementScheduler.instance")))
		// 初始化 Event Consume Logger
		logger.InitEventConsumeLogger(fmt.Sprintf("%s_eventconsume.log", mviper.GetString("VirtualRegistryManagementScheduler.instance")))
	}

	{ // metering init
		metering.Init(&metering.AMQP{
			Account:                mviper.GetString("metering_service.account"),
			Password:               mviper.GetString("metering_service.password"),
			Host:                   mviper.GetString("metering_service.host"),
			ManageHost:             mviper.GetString("metering_service.manage_host"),
			Timeout:                mviper.GetInt("metering_service.timeout"),
			RPCTimeout:             mviper.GetInt("metering_service.rpc_timeout"),
			Vhost:                  mviper.GetString("metering_service.vhost"),
			OperationConnectionNum: mviper.GetInt("metering_service.operation_connection_num"),
			ChannelNum:             mviper.GetInt("metering_service.channel_num"),
			ReplicaNum:             mviper.GetInt("metering_service.replica_num"),
			ConsumerConnectionNum:  mviper.GetInt("metering_service.consumer_connection_num"),
		})
	}

	{ // tracer (jaeger)
		if mviper.GetBool("VirtualRegistryManagement.tracer.enable") {
			tracer.Init(&tracer.Config{
				ServiceName: mviper.GetString("VirtualRegistryManagementScheduler.instance"),
				Endpoint:    mviper.GetString("VirtualRegistryManagement.tracer.host"),
				Timeout:     mviper.GetInt("VirtualRegistryManagement.tracer.timeout"),
			})
		}
	}
}

func startUpstreamServicesForScheduler() {
}

func prepareSchedulerServer() {
	{ // 初始化 gRPC Connection
		vrm.Init(vrm.PoolProvider{
			Mode: vrm.TCPMode,
			TCPProvider: vrm.TCPProvider{
				Hosts: mviper.GetStringSlice("VirtualRegistryManagementScheduler.core_grpc.hosts"),
				TLS: vrm.TLSConfig{
					Enable:   mviper.GetBool("VirtualRegistryManagementScheduler.core_grpc.tls.enable"),
					CertPath: mviper.GetString("VirtualRegistryManagementScheduler.core_grpc.tls.cert_path"),
				},
				ConnPerHost: mviper.GetInt("VirtualRegistryManagementScheduler.core_grpc.connection_per_host"),
			},
		})
	}

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
	{ // littlebell
		// 啟動 LBM
		lbmevents.Init()
	}

	{ // event consume
		eventconsume.New(mviper.GetString("event_consume.service"))
	}

}

func startSchedulerServer() {
	// 啟動定期同步 Resources
	tasks.InitSchedulerTasks()
}

func stopSchedulerServer() {
}

func stopUpstreamServicesForScheduler() {
	tracer.Shutdown()
}
