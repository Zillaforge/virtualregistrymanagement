package services

import (
	"github.com/spf13/viper"
	"pegasus-cloud.com/aes/pegasusiamclient/iam"
	iamUtil "pegasus-cloud.com/aes/pegasusiamclient/utility"
)

func InitIAM(svcCfg *viper.Viper) (err error) {
	svcCfg.SetDefault("connection_per_host", 3)
	var _connection *iam.PoolHandler
	_connection, err = iam.New(iam.PoolProvider{
		Mode: iam.TCPMode,
		TCPProvider: iam.TCPProvider{
			Hosts: svcCfg.GetStringSlice("hosts"),
			TLS: iam.TLSConfig{
				Enable:   svcCfg.GetBool("tls.enable"),
				CertPath: svcCfg.GetString("tls.cert_path"),
			},
			ConnPerHost: svcCfg.GetInt("connection_per_host"),
		},
		RouteResponseType: iamUtil.JSON,
	})
	if err != nil {
		return err
	}
	ServiceMap[svcCfg.GetString("name")] = &Service{
		Kind: _iamKind,
		Conn: _connection,
	}
	return nil
}
