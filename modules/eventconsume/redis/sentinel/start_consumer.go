package sentinel

import (
	eventConsumeCom "VirtualRegistryManagement/modules/eventconsume/common"
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (h *Handler) StartConsumer(ctx context.Context, input *eventConsumeCom.StartConsumerInput) {
	var (
		funcName = tkUtils.NameOfFunction().String()
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(nil)

	for _, channel := range h.channels {
		ch := h.conn.Subscribe(ctx, channel).Channel()
		router := input.Routers[channel]
		go func(c <-chan *redis.Message) {
			for {
				router(ctx, (<-c).Payload)
			}
		}(ch)
	}
}
