package configs

import (
	cnt "VirtualRegistryManagement/constants"
	"fmt"

	tkCfg "pegasus-cloud.com/aes/toolkits/configs"
	"pegasus-cloud.com/aes/toolkits/mviper"
)

func InitScheduler() {
	mviper.SetDefault("version", cnt.Version, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("kind", "Proxy", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)

	mviper.SetDefault("VirtualRegistryManagementScheduler.instance", fmt.Sprintf("%s-Scheduler", cnt.PascalCaseName), "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.core_grpc.hosts", []string{"0.0.0.0:5109"}, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.core_grpc.tls.enable", false, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.core_grpc.tls.cert_path", "/root/server-csr/RDS.pem", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.core_grpc.connection_per_host", 3, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)

	// tasks
	mviper.SetDefault("VirtualRegistryManagementScheduler.tasks.sync_projects.cron_expression", "0 */10 * * * * *", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.tasks.sync_tags.cron_expression", "0 */10 * * * * *", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.tasks.sync_exports.cron_expression", "0 * * * * * *", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.tasks.image_size.cron_expression", "0 */1 * * * * *", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.tasks.image_size.metering_service.exchange", "mts", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.tasks.image_size.metering_service.routing_key", "vrm-image-size", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.tasks.image_count.cron_expression", "0 */1 * * * * *", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.tasks.image_count.metering_service.exchange", "mts", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.tasks.image_count.metering_service.routing_key", "vrm-image-count", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.tasks.size_hard_limited_exceed.cron_expression", "0 */10 * * * * *", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("VirtualRegistryManagementScheduler.tasks.count_hard_limited_exceed.cron_expression", "0 */10 * * * * *", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// Event consume
	mviper.SetDefault("event_consume.service", "", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("event_consume.channels", []interface{}{}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
}
