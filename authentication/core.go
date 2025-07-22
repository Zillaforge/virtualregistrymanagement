package authentication

import (
	authCom "VirtualRegistryManagement/authentication/common"
	iamAuth "VirtualRegistryManagement/authentication/iam"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/services"

	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/memkvdb"
	"github.com/Zillaforge/toolkits/mviper"
)

var _provider authCom.Provider

const _iam string = "iam"

func Init(serviceName string) (err error) {
	memkvdb.Init()

	switch services.ServiceMap[mviper.GetString("authentication.service")].Kind {
	case _iam:
		_provider = iamAuth.New(nil,
			iamAuth.WithServiceName(mviper.GetString("authentication.service")),
		)
	default:
		return tkErr.New(cnt.AuthThisAuthenticationTypeIsNotSupportedErr)
	}
	return nil
}

func Use() (provider authCom.Provider) {
	return _provider
}

// Replace replaces global provider by p
func Replace(provider authCom.Provider) {
	_provider = provider
}
