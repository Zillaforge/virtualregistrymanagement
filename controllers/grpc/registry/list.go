package registry

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

func (m *Method) ListRegistries(ctx context.Context, input *pb.ListRegistriesInput) (output *pb.ListRegistriesOutput, err error) {
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

	whereInput := storCom.ListRegistryWhere{
		Namespace: input.Namespace,
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
		return &pb.ListRegistriesOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	flag := input.Flag
	if input.Flag == nil {
		flag = &pb.RegistryFlag{}
	}
	listRegistriesInput := &storCom.ListRegistriesInput{
		Pagination: storCom.Paginate(input.Limit, input.Offset),
		Where:      whereInput,
		Flag: storCom.Flag{
			UserID:    flag.UserID,
			ProjectID: flag.ProjectID,

			BelongUser:    flag.BelongUser,
			BelongProject: flag.BelongProject,
			ProjectLimit:  flag.ProjectLimit,
			ProjectPublic: flag.ProjectPublic,
			GlobalLimit:   flag.GlobalLimit,
			GlobalPublic:  flag.GlobalPublic,
		},
	}
	listRegistriesOutput, err := storages.Use().ListRegistries(ctx, listRegistriesInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListRegistries(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listRegistriesInput),
		).Error(err.Error())
		return &pb.ListRegistriesOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	output = &pb.ListRegistriesOutput{
		Count: listRegistriesOutput.Count,
	}
	for _, Registry := range listRegistriesOutput.Registries {
		output.Data = append(output.Data, m.storage2pb(&Registry))
	}
	return
}
