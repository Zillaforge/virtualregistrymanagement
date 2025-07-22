package configs

import (
	cnt "VirtualRegistryManagement/constants"
	"math"

	tkCfg "github.com/Zillaforge/toolkits/configs"
	"github.com/Zillaforge/toolkits/mviper"
)

func init() {
	mviper.SetDefault("version", cnt.Version, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("kind", cnt.PascalCaseName, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("location_id", getLocation(), "-", tkCfg.TypeSystem, tkCfg.RegionLocal)
	mviper.SetDefault("host_id", getHostname(), "-", tkCfg.TypeSystem, tkCfg.RegionLocal)

	mviper.SetDefault("VirtualRegistryManagement.developer", true, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("VirtualRegistryManagement.instance", cnt.PascalCaseName, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)

	// http
	mviper.SetDefault("VirtualRegistryManagement.http.host", "0.0.0.0:8109", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.http.access_control.allow_origins", []string{"*"}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.http.access_control.allow_credentials", true, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.http.access_control.allow_headers", []string{"Origin", "Content-Length", "Content-Type", "Authorization"}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.http.access_control.allow_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.http.access_control.expose_headers", []string{"host-id", "version-id", "location-id"}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// grpc
	mviper.SetDefault("VirtualRegistryManagement.grpc.host", "0.0.0.0:5109", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("VirtualRegistryManagement.grpc.unix_socket.path", "/run/VirtualPlatformservice.sock", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.grpc.unix_socket.conn_count", 20, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.grpc.write_buffer_size", 32*1024, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.grpc.read_buffer_size", 32*1024, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.grpc.max_receive_message_size", 1024*1024*4, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.grpc.max_send_message_size", math.MaxInt32, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.grpc.hosts", []string{"127.0.0.1:5106"}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.grpc.conn_per_host", 3, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// tls
	mviper.SetDefault("VirtualRegistryManagement.tls.enable", false, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.tls.cert_path", "/var/lib/ASUS/vrm/tls/cert.pem", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.tls.key_path", "/var/lib/ASUS/vrm/tls/cert-key.pem", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// tracer
	mviper.SetDefault("VirtualRegistryManagement.tracer.enable", false, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.tracer.host", "http://127.0.0.1:14268/api/traces", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.tracer.timeout", 10, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// scopes
	mviper.SetDefault("VirtualRegistryManagement.scopes.memcache_ttl", 60, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("VirtualRegistryManagement.scopes.project_default_count", -1, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("VirtualRegistryManagement.scopes.project_default_size", -1, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("VirtualRegistryManagement.scopes.availability_district", "default", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("VirtualRegistryManagement.scopes.snapshot_size", "size", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)

	mviper.SetDefault("VirtualRegistryManagement.worker_pool.max_workers", 10, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("VirtualRegistryManagement.worker_pool.max_capacity", 1000, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)

	// filepath
	mviper.SetDefault("VirtualRegistryManagement.filepath", []interface{}{}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// log
	// system-log
	mviper.SetDefault("VirtualRegistryManagement.logger.system_log.path", "/var/log/ASUS/", "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.system_log.max_size", 100, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.system_log.max_backups", 5, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.system_log.max_age", 10, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.system_log.compress", false, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.system_log.mode", "error", "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.system_log.show_in_console", true, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	// access-log
	mviper.SetDefault("VirtualRegistryManagement.logger.access_log.path", "/var/log/ASUS/", "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.access_log.max_size", 100, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.access_log.max_backups", 5, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.access_log.max_age", 10, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.access_log.compress", false, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.access_log.default_source_ip", "127.0.0.1", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// metering service
	mviper.SetDefault("metering_service.account", "account", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.password", "password", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.host", "127.0.0.1:5672", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.manage_host", "127.0.0.1:15672", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.timeout", 10, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.rpc_timeout", 10, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.vhost", "/", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.operation_connection_num", 2, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.channel_num", 1, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.replica_num", 0, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.consumer_connection_num", 1, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// event-log
	mviper.SetDefault("VirtualRegistryManagement.logger.event_consume_log.path", "/var/log/ASUS/", "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.event_consume_log.max_size", 100, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.event_consume_log.max_backups", 5, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.event_consume_log.max_age", 10, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.logger.event_consume_log.compress", false, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)

	// storage
	mviper.SetDefault("storage.provider", "mariadb", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.auto_migrate", true, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.host", "mariadb-galera.pegasus-system:3306", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.account", "", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.password", "", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.name", "pt", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.timeout", 5, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.max_open_conns", 150, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.conn_max_lifetime", 10, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.max_idle_conns", 150, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// vps
	mviper.SetDefault("vps.host", "", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("vps.tls.enable", false, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("vps.tls.cert_path", "", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("vps.connection_per_host", 3, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)

	// authentication
	mviper.SetDefault("authentication.service", "", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	// openstack
	mviper.SetDefault("openstack", []interface{}{}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	// services
	mviper.SetDefault("services", []interface{}{}, "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
	// Event consume
	mviper.SetDefault("event_consume.service", "", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("event_consume.channels", []interface{}{}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// cryptography
	mviper.SetDefault("VirtualRegistryManagement.cryptography.support_type", []string{}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.cryptography.default_type", "", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagement.cryptography.default_encrypted", "", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// plugins
	mviper.SetDefault("event_publish.plugins", []interface{}{}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// littlebell
	mviper.SetDefault("littlebell.arn", "arn:aws:sns:default:14735dfa-5553-46cc-b4bd-405e711b223f:lbm-svc-event-publish-topic", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
	mviper.SetDefault("littlebell.host", "http://sns-service.pegasus-system:8092", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
	mviper.SetDefault("littlebell.region", "default", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
	mviper.SetDefault("littlebell.access_key", "", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
	mviper.SetDefault("littlebell.secret_key", "", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
	mviper.SetDefault("littlebell.credential.project_id", "", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
	mviper.SetDefault("littlebell.credential.user_id", "", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
}
