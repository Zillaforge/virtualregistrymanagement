package repository

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
)

func (m *Method) CreateRepository(ctx context.Context, input *pb.RepositoryInfo) (output *pb.RepositoryDetail, err error) {
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

	createInput := &storCom.CreateRepositoryInput{
		Repository: tables.Repository{
			ID:              input.ID,
			Name:            input.Name,
			Namespace:       input.Namespace,
			OperatingSystem: input.OperatingSystem,
			Description:     input.Description,
			Creator:         input.Creator,
			ProjectID:       input.ProjectID,
		},
	}

	createOutput, err := storages.Use().CreateRepository(ctx, createInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageRepositoryNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCRepositoryNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().CreateRepository()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(ctx, &createOutput.Repository)
	return
}
