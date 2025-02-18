package sentinel

import (
	"VirtualRegistryManagement/services"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type (
	Handler struct {
		serviceName string
		conn        *redis.Client
		channels    []string
	}
)

type Input struct {
	Channels []string
}

type Option func(handler *Handler)

func WithServiceName(name string) Option {
	return func(handler *Handler) {
		handler.serviceName = name
	}
}

func New(input *Input, opts ...Option) (handler *Handler) {
	handler = &Handler{
		channels: input.Channels,
	}
	for _, opt := range opts {
		opt(handler)
	}
	if handler.serviceName == "" {
		panic("Service must not be empty")
	} else {
		if services.ServiceMap[handler.serviceName] != nil {
			handler.conn = services.ServiceMap[handler.serviceName].Conn.(*redis.Client)
		} else {
			panic(fmt.Sprintf("Service %s not found", handler.serviceName))
		}
	}
	return handler
}
