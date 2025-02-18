package repository

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

func (m *Method) UpdateRepository(ctx context.Context, input *pb.UpdateRepositoryInput) (output *pb.RepositoryDetail, err error) {
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

	updateInput := &storCom.UpdateRepositoryInput{
		ID: input.ID,
		UpdateData: &storCom.RepositoryUpdateInfo{
			Name:        input.Name,
			Description: input.Description,
		},
	}

	updateOutput, err := storages.Use().UpdateRepository(ctx, updateInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageRepositoryNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCRepositoryNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().UpdateRepository()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(ctx, &updateOutput.Repository)
	return
}
