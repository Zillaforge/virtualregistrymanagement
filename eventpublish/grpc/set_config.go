package grpc

import (
	cnt "VirtualRegistryManagement/constants"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/eventpublishpluginclient/pb"
)

func (c *core) SetConfig(conf []byte) {
	setConfigInput := &pb.SetConfigRequest{Conf: conf}
	if err := c.handler.SetConfig(setConfigInput); err != nil {
		zap.L().With(
			zap.String(cnt.Plugin, "c.handler.SetConfig(...)"),
			zap.Any("input", setConfigInput),
		).Error(err.Error())
		return
	}
}
