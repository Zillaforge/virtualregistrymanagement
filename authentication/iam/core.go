package iam

import (
	"VirtualRegistryManagement/services"

	"github.com/Zillaforge/pegasusiamclient/iam"
)

type Input struct {
	_ struct{}
}

type Provider struct {
	poolHandler *iam.PoolHandler
	_           struct{}
}

var _serviceName string = ""
var _connection *iam.PoolHandler

type Option func()

func WithServiceName(name string) Option {
	return func() {
		_serviceName = name
	}
}

func New(input *Input, opts ...Option) (provider *Provider) {
	for _, opt := range opts {
		opt()
	}
	_connection = services.ServiceMap[_serviceName].Conn.(*iam.PoolHandler)
	return &Provider{
		poolHandler: _connection,
	}
}
