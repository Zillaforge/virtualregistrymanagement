package eventrouter

import (
	cnt "VirtualRegistryManagement/constants"
	eventConsumeCom "VirtualRegistryManagement/modules/eventconsume/common"
	util "VirtualRegistryManagement/utility"
	"context"
	"encoding/json"

	iamCtl "VirtualRegistryManagement/controllers/iamconsumer"

	"go.uber.org/zap"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func IAMConsumerRouter(ctx context.Context, message interface{}) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = util.MustGetContextRequestID(ctx)
		input     = &eventConsumeCom.Data{}
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)

	var err error
	defer f(tracer.Attributes{
		"input": &input,
		"error": &err,
	})

	if err = json.Unmarshal([]byte(message.(string)), input); err != nil {
		zap.L().With(
			zap.String(cnt.Server, "json.Unmarshal(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("payloadData", input),
		).Error(err.Error())
		return
	}

	switch input.Action {
	case iamCtl.EventSyncProjects:
		iamCtl.SyncProjects(ctx, iamCtl.UnmarshalSyncProjects(input))
	case iamCtl.EventCreateProject, iamCtl.EventCreateProjectWithResp:
		iamCtl.CreateProject(ctx, iamCtl.UnmarshalCreateProject(input))
	case iamCtl.EventDeleteProject:
		iamCtl.DeleteProject(ctx, iamCtl.UnmarshalDeleteProject(input))
	case iamCtl.EventDeleteUser:
		iamCtl.DeleteUser(ctx, iamCtl.UnmarshalDeleteUser(input))
	case iamCtl.EventDeleteMembership:
		iamCtl.DeleteMembership(ctx, iamCtl.UnmarshalDeleteMembership(input))
	}
}
