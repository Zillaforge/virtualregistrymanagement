package tasks

import (
	cnt "VirtualRegistryManagement/constants"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

/*
SyncExports ...

errors:
- 18000000(internal server error)
*/
func SyncExports() (err error) {
	var (
		funcName  = tkUtils.NameOfFunction().Name()
		requestID = tracer.EmptyRequestID
	)

	ctx, f := tracer.StartWithContext(tracer.StartEntryContext(requestID), funcName)
	defer f(tracer.Attributes{
		"err": &err,
	})

	if err = vrm.SyncExports(ctx); err != nil {
		zap.L().With(
			zap.String(cnt.Task, "vrm.SyncExports(...)"),
			zap.String(cnt.RequestID, requestID),
		).Error(err.Error())
		return
	}
	return nil
}
