package grpc

import (
	cnt "VirtualRegistryManagement/constants"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/eventpublishpluginclient/pb"
	"pegasus-cloud.com/aes/toolkits/tracer"
)

func (c *core) Reconcile(action string, meta map[string]string, req interface{}, resp interface{}) {
	reqValue, err := json.Marshal(req)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "json.Marshal(...)"),
			zap.String("action", action),
			zap.Any("metadata", meta),
			zap.Any("req", req),
			zap.Any("resp", resp),
		).Error(err.Error())
		return
	}

	respValue, err := json.Marshal(resp)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "json.Marshal(...)"),
			zap.String("action", action),
			zap.Any("metadata", meta),
			zap.Any("req", req),
			zap.Any("resp", resp),
		).Error(err.Error())
		return
	}

	input := &pb.ReconcileRequest{
		Action:   action,
		Metadata: meta,
		Request:  reqValue,
		Response: respValue,
	}
	requestID := tracer.EmptyRequestID
	if value, ok := meta[tracer.RequestID]; ok {
		requestID = value
	}
	ctx, f := tracer.StartWithContext(tracer.StartEntryContext(requestID), fmt.Sprintf("Send %s event to plugin via eventbus", input.Action))
	defer f(tracer.Attributes{
		tracer.RequestID: requestID,
	})
	if err := c.handler.Reconcile(input, ctx); err != nil {
		zap.L().With(zap.String("plugin.eventpublish.grpc", "g.handler.Reconcile(...)")).Error(err.Error())
		return
	}
}
