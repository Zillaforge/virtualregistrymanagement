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
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
)

func (m *Method) ListRepositories(ctx context.Context, input *pb.ListNamespaceInput) (output *pb.ListRepositoriesOutput, err error) {
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
	whereInput := storCom.ListRepositoryWhere{}
	if input.Namespace != nil {
		whereInput.Namespace = input.Namespace
	}
	if err := querydecoder.ShouldBindWhereSlice(&whereInput, input.Where); err != nil {
		if e, ok := tkErr.IsError(grpc.WhereErrorParser(err)); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErr.Code():
				return output, e
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "querydecoder.ShouldBindWhereSlice(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input.where", input.Where),
		).Error(err.Error())
		return &pb.ListRepositoriesOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	listRepositoriesInput := &storCom.ListRepositoriesInput{
		Pagination: storCom.Paginate(input.Limit, input.Offset),
		Where:      whereInput,
	}
	listRepositoriesOutput, err := storages.Use().ListRepositories(ctx, listRepositoriesInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListRepositories(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listRepositoriesInput),
		).Error(err.Error())
		return &pb.ListRepositoriesOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	output = &pb.ListRepositoriesOutput{
		Count: listRepositoriesOutput.Count,
	}
	for _, Repository := range listRepositoriesOutput.Repositories {
		output.Data = append(output.Data, m.storage2pb(ctx, &Repository))
	}
	return
}
