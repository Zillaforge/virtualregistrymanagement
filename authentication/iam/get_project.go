package iam

import (
	authCom "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"context"
	"fmt"

	"go.uber.org/zap"
	iamCnt "github.com/Zillaforge/pegasusiamclient/constants"
	iamPb "github.com/Zillaforge/pegasusiamclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (h *Provider) GetProject(ctx context.Context, input *authCom.GetProjectInput) (output *authCom.GetProjectOutput, err error) {
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
	iamGetIn := &iamPb.ProjectID{ID: input.ID}

	iamGetOut := &iamPb.ProjectInfo{}
	err = authCom.RetrieveFromCache(fmt.Sprintf("project.%s", input.ID), input.Cacheable, iamGetOut, func() error {
		iamGetOut, err = h.poolHandler.GetProject(iamGetIn, ctx)
		return err
	}, func() interface{} {
		return iamGetOut
	})

	if err != nil {
		if _err, ok := tkErr.IsError(err); ok {
			switch _err.Code() {
			case iamCnt.GRPCProjectDoesNotExistErrCode:
				return nil, tkErr.New(cnt.AuthProjectNotFoundErr)
			}
		}
		zap.L().With(
			zap.String(cnt.Authentication, "h.poolHandler.GetProject"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", iamGetIn),
		).Warn(err.Error())
		return nil, tkErr.New(cnt.AuthInternalServerErr)
	}

	if iamGetOut.Frozen {
		zap.L().With(
			zap.String(cnt.Authentication, "h.poolHandler.GetProject"),
			zap.Any("input", iamGetIn),
			zap.Bool("output", iamGetOut.Frozen),
		).Warn(tkErr.New(cnt.AuthProjectIsFrozenErr, iamGetIn.ID).Message())
	}
	output = &authCom.GetProjectOutput{
		ID:          iamGetOut.ID,
		DisplayName: iamGetOut.DisplayName,
		Frozen:      iamGetOut.Frozen,
		Extra:       utility.Bytes2Extra(iamGetOut.Extra),
	}
	return
}
