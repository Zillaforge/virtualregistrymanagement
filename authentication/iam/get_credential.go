package iam

import (
	authCom "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/pegasusiamclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

/*
GetCredential ...

errors:
- 17000000(internal server error)
*/
func (h *Provider) GetCredential(ctx context.Context, input *authCom.GetCredentialInput) (output *authCom.GetCredentialOutput, err error) {
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

	getCredentialInput := &pb.CredUserProjectInput{
		UserID:    input.UserId,
		ProjectID: input.ProjectId,
	}
	getCredentialOutput, getCredentialErr := h.poolHandler.GetCredential(getCredentialInput, ctx)
	if getCredentialErr != nil {
		zap.L().With(
			zap.String(cnt.Authentication, "h.poolHandler.GetCredential(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getCredentialInput),
		).Error(getCredentialErr.Error())
		return nil, tkErr.New(cnt.AuthInternalServerErr).WithInner(getCredentialErr)
	}

	return &authCom.GetCredentialOutput{
		AccessKey: getCredentialOutput.Access,
		SecretKey: getCredentialOutput.Secret,
	}, nil
}
