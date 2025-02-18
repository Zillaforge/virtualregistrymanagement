package grpc

import (
	cnt "VirtualRegistryManagement/constants"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/toolkits/tracer"
)

func (c *core) CheckPluginVersion() bool {
	ctx, f := tracer.StartWithContext(tracer.StartEntryContext(tracer.EmptyRequestID), "Check Plugin Version")
	defer f(nil)
	checkPluginVersionOutput, err := c.handler.CheckPluginVersion(ctx)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Plugin, "c.handler.CheckPluginVersion(...)"),
		).Error(err.Error())
		return false
	}
	return checkPluginVersionOutput.IsMatch
}
