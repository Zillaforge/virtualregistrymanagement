package memberacl

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/storages"
	storCom "VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
)

func (m *Method) GetMemberAcl(ctx context.Context, input *pb.GetInput) (output *pb.MemberAclDetail, err error) {
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

	getInput := &storCom.GetMemberAclInput{
		ID: input.ID,
	}

	getOutput, err := storages.Use().GetMemberAcl(ctx, getInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageMemberAclNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCMemberAclNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetMemberAcl()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(ctx, &getOutput.MemberAcl)
	return
}
