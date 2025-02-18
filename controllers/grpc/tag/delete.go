package tag

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/grpc"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/utility"
	"VirtualRegistryManagement/utility/querydecoder"
	"context"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
)

func (m *Method) DeleteTag(ctx context.Context, input *pb.DeleteInput) (output *pb.DeleteOutput, err error) {
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
	whereInput := storCom.DeleteTagWhere{}
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
			TagID:        whereInput.ID,
			RepositoryID: whereInput.RepositoryID,
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
		RepositoryID: whereInput.RepositoryID,
		TagID:        whereInput.ID,
	})

	deleteInput := &storCom.DeleteTagInput{
		Where: whereInput,
	}

	deleteOutput, err := storages.Use().DeleteTag(ctx, deleteInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageTagNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCTagNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().DeleteTag()"),
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
