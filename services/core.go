package services

import (
	cnt "VirtualRegistryManagement/constants"
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/mviper"
)

type Service struct {
	Kind string
	Conn interface{}
}

const (
	_iamKind           string = "iam"
	_redisSentinelKind string = "redis_sentinel"
	_openstackKind     string = "openstack"
)

var ServiceMap = make(map[string]*Service)

func InitServices() (err error) {
	zap.L().Info("initial upstream services")
	for _, service := range cast.ToSlice(mviper.Get("services")) {
		if service == nil {
			continue
		}
		cfg := viper.New()
		cfg.MergeConfigMap(cast.ToStringMap(service))
		if !cfg.IsSet("name") || !cfg.IsSet("kind") {
			return tkErr.New(cnt.ServiceNameIsRequiredErr)
		}
		if _, exist := ServiceMap[cfg.GetString("name")]; exist {
			return tkErr.New(cnt.ServiceNameIsRepeatedErr)
		}

		switch cfg.GetString("kind") {
		case _iamKind:
			InitIAM(cfg)
		case _redisSentinelKind:
			UnmarshalRedisSentinel(cfg)
		case _openstackKind:
			InitOpenstack(cfg)
		}
	}
	var svcs = ""
	for k, v := range ServiceMap {
		svcs += fmt.Sprintf("%s[%s],", v.Kind, k)
	}

	zap.L().Info(fmt.Sprintf("available services: %s", strings.TrimSuffix(svcs, ",")))
	return nil
}
