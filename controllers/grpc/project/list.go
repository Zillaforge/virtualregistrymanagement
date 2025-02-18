package project

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

func (m *Method) ListProjects(ctx context.Context, input *pb.ListInput) (output *pb.ListProjectsOutput, err error) {
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

	whereInput := storCom.ListProjectWhere{}
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
		return &pb.ListProjectsOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	listProjectsInput := &storCom.ListProjectsInput{
		Pagination: storCom.Paginate(input.Limit, input.Offset),
		Where:      whereInput,
	}
	listProjectsOutput, err := storages.Use().ListProjects(ctx, listProjectsInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListProjects(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listProjectsInput),
		).Error(err.Error())
		return &pb.ListProjectsOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	output = &pb.ListProjectsOutput{
		Count: listProjectsOutput.Count,
	}
	for _, Project := range listProjectsOutput.Projects {
		output.Data = append(output.Data, m.storage2pb(ctx, &Project))
	}
	return
}
