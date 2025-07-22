package iam

import (
	authCom "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	iamCnt "github.com/Zillaforge/pegasusiamclient/constants"
	"github.com/Zillaforge/pegasusiamclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

// errors:
//   - 17000001(internal server error)
//   - 17000010(permission denied)
//   - 17000005(incorrect format of authentication)
func (h *Provider) VerifySystemAdminToken(ctx context.Context, input *authCom.VerifySystemAdminTokenInput) (output *authCom.VerifySystemAdminTokenOutput, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(tracer.Attributes{
		"input":  &input,
		"output": &output,
		"err":    &err,
	})
	verifySystemAdminTokenInput := &pb.VerifySystemAdminTokenInput{
		Token: input.Token,
	}
	verifySystemAdminTokenOutput, err := h.poolHandler.VerifySystemAdminToken(verifySystemAdminTokenInput, ctx)
	if err != nil {
		if tkerr, ok := tkErr.IsError(err); ok {
			switch tkerr.Code() {
			case iamCnt.GPRCYouAreNotSystemAdminErrCode,
				iamCnt.GRPCIllegalTokenErrCode,
				iamCnt.GRPCYouAreNotAnAdministratorErrCode:
				return nil, tkErr.New(cnt.AuthPermissionDeniedErr)
			case iamCnt.GRPCIncorrectFormatOfAuthenticationErrCode,
				iamCnt.GRPCIncorrectFormatOfAuthenticationNotFoundProjectErrCode:
				return nil, tkErr.New(cnt.AuthIncorrectFormatErr)
			}
		}
		zap.L().With(
			zap.String(cnt.Authentication, "h.poolHandler.VerifySystemAdminToken()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", verifySystemAdminTokenInput),
		).Error(err.Error())
		return nil, tkErr.New(cnt.AuthInternalServerErr)
	}
	return &authCom.VerifySystemAdminTokenOutput{
		UserID:  verifySystemAdminTokenOutput.User.ID,
		Account: verifySystemAdminTokenOutput.User.Account,
	}, nil
}
