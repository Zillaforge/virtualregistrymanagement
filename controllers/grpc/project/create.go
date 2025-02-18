package project

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
)

func (m *Method) CreateProject(ctx context.Context, input *pb.ProjectInfo) (output *pb.ProjectInfo, err error) {
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

	createInput := &storCom.CreateProjectInput{
		Project: tables.Project{
			ID:             input.ID,
			LimitCount:     input.LimitCount,
			LimitSizeBytes: input.LimitSizeBytes,
		},
	}

	createOutput, err := storages.Use().CreateProject(ctx, createInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageProjectNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCProjectNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "stor.Use().CreateProject()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(ctx, &createOutput.Project)
	return
}
