package services

import (
	opstkExec "VirtualRegistryManagement/modules/openstack/execution"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitOpenstack(svcCfg *viper.Viper) (err error) {
	_connection, err := opstkExec.New(&opstkExec.Connection{
		IdentityEndpoint: svcCfg.GetString("identity_endpoint"),
		Username:         svcCfg.GetString("username"),
		Password:         svcCfg.GetString("password"),
		DomainName:       svcCfg.GetString("domain_name"),
		AllowReauth:      svcCfg.GetBool("allow_reauth"),
		AdminProject:     svcCfg.GetString("admin_project"),
		PidSource:        svcCfg.GetString("pid_source"),
	})
	if err != nil {
		zap.L().Error(fmt.Sprintf("failed to initialize openstack service [%s]", svcCfg.GetString("name")))
		return err
	}
	ServiceMap[svcCfg.GetString("name")] = &Service{
		Kind: _openstackKind,
		Conn: _connection,
	}
	return
}
