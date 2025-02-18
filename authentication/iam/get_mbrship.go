package iam

import (
	authCom "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	iamCnt "pegasus-cloud.com/aes/pegasusiamclient/constants"
	"pegasus-cloud.com/aes/pegasusiamclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

/*
GetMembership get membership information from iam server

errors:
- 17000000(internal server error)
- 17000006(not found membership)
*/
func (h *Provider) GetMembership(ctx context.Context, input *authCom.GetMembershipInput) (output *authCom.GetMembershipOutput, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().Name()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(tracer.Attributes{
		"input":  input,
		"output": output,
		"err":    &err,
	})
	getMembershipInput := &pb.MemUserProjectInput{
		UserID:    input.UserId,
		ProjectID: input.ProjectId,
	}
	data, err := h.poolHandler.GetMembership(getMembershipInput, ctx)
	if err != nil {
		if tkerr, ok := tkErr.IsError(err); ok {
			switch tkerr.Code() {
			case iamCnt.GRPCMembershipDoesNotExistErrCode:
				// 找不到成員關係
				return nil, tkErr.New(cnt.AuthMembershipNotFoundErr)
			}
		}
		zap.L().With(
			zap.String(cnt.Authentication, "h.poolHandler.GetMembership(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getMembershipInput),
		).Error(err.Error())
		return nil, tkErr.New(cnt.AuthInternalServerErr).WithInner(err)
	}

	output = &authCom.GetMembershipOutput{
		TenantRole: data.TenantRole,
		Frozen:     data.Frozen,
		Extra:      utility.Bytes2Extra(data.Extra),
	}
	return
}
