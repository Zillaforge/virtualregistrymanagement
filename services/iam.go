package services

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"github.com/Zillaforge/pegasusiamclient/iam"
	iamUtil "github.com/Zillaforge/pegasusiamclient/utility"
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
		zap.L().Error(fmt.Sprintf("failed to initialize iam service [%s]", svcCfg.GetString("name")))
		return err
	}
	ServiceMap[svcCfg.GetString("name")] = &Service{
		Kind: _iamKind,
		Conn: _connection,
	}
	return nil
}
