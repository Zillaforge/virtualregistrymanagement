package eventconsume

import (
	eventConsumeCom "VirtualRegistryManagement/modules/eventconsume/common"
	"VirtualRegistryManagement/modules/eventconsume/redis/sentinel"
	redisCtlConsumer "VirtualRegistryManagement/server/event_router"
	"VirtualRegistryManagement/services"
	"context"
	"fmt"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/tracer"
)

type Provider interface {
	StartConsumer(ctx context.Context, input *eventConsumeCom.StartConsumerInput)
}

var provider Provider

const (
	_empty         string = "empty"
	_redisSentinel string = "redis_sentinel"
	_fromIAM       string = "FromIAM"
)

func New(service string) {
	if service == _empty {
		zap.L().Info("disable consume events")
		return
	}
	zap.L().Info(fmt.Sprintf("enable consume events : %s", service))
	kind := service
	if services.ServiceMap[service] != nil {
		kind = services.ServiceMap[service].Kind
		zap.L().Info(fmt.Sprintf("EventConsume is %s(%s) service mode", kind, service))
	}
	switch kind {
	case _redisSentinel:
		ctx := tracer.StartEntryContext(tracer.EmptyRequestID)
		provider = sentinel.New(&sentinel.Input{
			Channels: mviper.GetStringSlice("event_consume.channels"),
		}, sentinel.WithServiceName(service))
		provider.StartConsumer(ctx, &eventConsumeCom.StartConsumerInput{
			Routers: map[string]func(ctx context.Context, message interface{}){
				_fromIAM: redisCtlConsumer.IAMConsumerRouter,
			},
		})
	default:
		panic(fmt.Errorf("event consumer does not support %s mode", kind))
	}
}

// Use returns eventconsume instance which is created by eventconsume.New() function
func Use() Provider {
	return provider
}

// Replace replaces global provider by p
func Replace(p Provider) {
	provider = p
}
