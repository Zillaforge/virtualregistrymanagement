package grpc

import (
	cnt "VirtualRegistryManagement/constants"

	"go.uber.org/zap"
	"github.com/Zillaforge/toolkits/tracer"
)

func (c *core) GetVersion() (v string) {
	ctx, f := tracer.StartWithContext(tracer.StartEntryContext(tracer.EmptyRequestID), "Get Version")
	defer f(nil)
	getVersionOutput, err := c.handler.GetVersion(ctx)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Plugin, "c.handler.GetVersion(...)"),
		).Error(err.Error())
		return ""
	}
	return getVersionOutput.Version
}
