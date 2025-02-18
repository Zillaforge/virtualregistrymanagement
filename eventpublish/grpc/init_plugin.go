package grpc

import (
	cnt "VirtualRegistryManagement/constants"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/toolkits/tracer"
)

func (c *core) InitPlugin() bool {
	ctx, f := tracer.StartWithContext(tracer.StartEntryContext(tracer.EmptyRequestID), "Init Plugin")
	defer f(nil)
	initPluginOutput, err := c.handler.InitPlugin(ctx)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Plugin, "c.handler.InitPlugin(...)"),
		).Error(err.Error())
		return false
	}
	return initPluginOutput.IsEnable
}
