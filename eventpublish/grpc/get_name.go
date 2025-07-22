package grpc

import (
	cnt "VirtualRegistryManagement/constants"

	"go.uber.org/zap"
	"github.com/Zillaforge/toolkits/tracer"
)

func (c *core) GetName() string {
	ctx, f := tracer.StartWithContext(tracer.StartEntryContext(tracer.EmptyRequestID), "Get Name")
	defer f(nil)
	getName, err := c.handler.GetName(ctx)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Plugin, "c.handler.GetName(...)"),
		).Error(err.Error())
		return ""
	}
	return getName.Name
}
