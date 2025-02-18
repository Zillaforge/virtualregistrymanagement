package services

import (
	opstkExec "VirtualRegistryManagement/modules/openstack/execution"

	"github.com/spf13/viper"
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
		return err
	}
	ServiceMap[svcCfg.GetString("name")] = &Service{
		Kind: _openstackKind,
		Conn: _connection,
	}
	return
}
