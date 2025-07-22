package grpc

import (
	cnt "VirtualRegistryManagement/constants"

	"go.uber.org/zap"
	"github.com/Zillaforge/eventpublishpluginclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
)

func (c *core) CallGRPCRouter(operator string, hdr map[string]string, payload []byte) (map[string]string, []byte, error) {
	ctx, f := tracer.StartWithContext(tracer.StartEntryContext(tracer.EmptyRequestID), "Call gRPC Router")
	defer f(nil)
	callGRPCRouterInput := &pb.RPCRouterRequest{
		Operator: operator,
		Hdr:      hdr,
		Payload:  payload,
	}
	callGRPCRouterOutput, err := c.handler.CallGRPCRouter(callGRPCRouterInput, ctx)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Plugin, "c.handler.CallGRPCRouter(...)"),
			zap.Any("input", callGRPCRouterInput),
		).Error(err.Error())
		return nil, nil, tkErr.New(cnt.PluginInternalServerErr)
	}
	return callGRPCRouterOutput.Hdr, callGRPCRouterOutput.Payload, nil
}
