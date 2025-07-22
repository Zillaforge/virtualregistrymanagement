package repository

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/grpc"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/utility"
	"VirtualRegistryManagement/utility/querydecoder"
	"context"

	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
)

func (m *Method) DeleteRepository(ctx context.Context, input *pb.DeleteInput) (output *pb.DeleteOutput, err error) {
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

	// binding the where parameter
	whereInput := storCom.DeleteRepositoryWhere{}
	if err = querydecoder.ShouldBindWhereSlice(&whereInput, input.Where); err != nil {
		if e, ok := tkErr.IsError(grpc.WhereErrorParser(err)); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErr.Code():
				return output, e
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "querydecoder.ShouldBindWhereSlice(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input.Where", input.Where),
		).Error(err.Error())
		return output, tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
	}

	// delete openstack resource
	listInput := &storCom.ListRegistriesInput{
		Pagination: storCom.Paginate(-1, 0),
		Where: storCom.ListRegistryWhere{
			RepositoryID: whereInput.ID,
			Creator:      whereInput.Creator,
			ProjectID:    whereInput.ProjectID,
		},
	}
	listOutput, _err := storages.Use().ListRegistries(ctx, listInput)
	if _err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListRegistries()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listInput),
		).Error(_err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(_err)
		return
	}
	grpc.DeleteOpstkResource(ctx, &grpc.OpstkResourceInput{
		Registries:   listOutput,
		RepositoryID: whereInput.ID,
		ProjectID:    whereInput.ProjectID,
		Creator:      whereInput.Creator,
	})

	deleteInput := &storCom.DeleteRepositoryInput{
		Where: whereInput,
	}

	deleteOutput, err := storages.Use().DeleteRepository(ctx, deleteInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageRepositoryNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCRepositoryNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().DeleteRepository()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", deleteInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = &pb.DeleteOutput{
		ID: deleteOutput.ID,
	}
	return
}
