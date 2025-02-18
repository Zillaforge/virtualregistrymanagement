package iam

import (
	"VirtualRegistryManagement/authentication/common"
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

// VerifyToken ...
func (h *Provider) VerifyToken(ctx context.Context, input *common.VerifyTokenInput) (output *common.VerifyTokenOutput, err error) {
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

	verifyTokenInput := &pb.VerifyTokenInput{Token: input.Token}
	data, err := h.poolHandler.VerifyToken(verifyTokenInput, ctx)
	if err != nil {
		if tkerr, ok := tkErr.IsError(err); ok {
			switch tkerr.Code() {
			case iamCnt.GPRCUserDoesNotExistErrCode:
				// 若 h.poolHandler.VerifySystemAdminToken 回傳的錯誤為 iamCnt.GPRCUserDoesNotExistErrCode
				// 表示指定使用者不存在
				return nil, tkErr.New(cnt.AuthUserNotFoundErr)
			case iamCnt.GRPCIncorrectFormatOfAuthenticationErrCode:
				return nil, tkErr.New(cnt.AuthIncorrectFormatErr)
			}
		}
		zap.L().With(
			zap.String(cnt.Authentication, "h.poolHandler.VerifyToken(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", verifyTokenInput),
		).Error(err.Error())
		return nil, tkErr.New(cnt.AuthInternalServerErr)
	}
	output = &common.VerifyTokenOutput{
		UserID:     data.User.ID,
		Frozen:     data.User.Frozen,
		Account:    data.User.Account,
		SAATUserID: data.SAATUserID,
	}
	return
}
