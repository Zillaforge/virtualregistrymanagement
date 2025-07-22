package memberacl

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

func (m *Method) ListMemberAcls(ctx context.Context, input *pb.ListNamespaceInput) (output *pb.ListMemberAclsOutput, err error) {
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

	whereInput := storCom.ListMemberAclWhere{
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
		return &pb.ListMemberAclsOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	listMemberAclsInput := &storCom.ListMemberAclsInput{
		Pagination: storCom.Paginate(input.Limit, input.Offset),
		Where:      whereInput,
	}
	listMemberAclsOutput, err := storages.Use().ListMemberAcls(ctx, listMemberAclsInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListMemberAcls(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listMemberAclsInput),
		).Error(err.Error())
		return &pb.ListMemberAclsOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	output = &pb.ListMemberAclsOutput{
		Count: listMemberAclsOutput.Count,
		Data:  []*pb.MemberAclDetail{},
	}
	for _, memberAcl := range listMemberAclsOutput.MemberAcls {
		output.Data = append(output.Data, m.storage2pb(ctx, &memberAcl.MemberAcl))
	}
	return
}
