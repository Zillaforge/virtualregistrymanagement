package project

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
)

func (m *Method) UpdateProject(ctx context.Context, input *pb.UpdateProjectInput) (output *pb.ProjectInfo, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input":  &input,
			"output": &output,
			"error":  &err,
		},
	)

	updateInput := &storCom.UpdateProjectInput{
		ID: input.ID,
		UpdateData: &storCom.ProjectUpdateInfo{
			LimitCount:     input.LimitCount,
			LimitSizeBytes: input.LimitSizeBytes,
		},
	}

	updateOutput, err := storages.Use().UpdateProject(ctx, updateInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageProjectNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCProjectNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().UpdateProject()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(ctx, &updateOutput.Project)
	return
}
