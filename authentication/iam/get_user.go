package iam

import (
	"VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"context"
	"fmt"

	"go.uber.org/zap"
	iamCnt "github.com/Zillaforge/pegasusiamclient/constants"
	protos "github.com/Zillaforge/pegasusiamclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (h *Provider) GetUser(ctx context.Context, input *common.GetUserInput) (output *common.GetUserOutput, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(tracer.Attributes{
		"input":  &input,
		"output": &output,
		"error":  &err,
	})
	iamGetIn := &protos.UserID{ID: input.ID}
	iamGetOut := &protos.UserInfo{}
	err = common.RetrieveFromCache(fmt.Sprintf("user.%s", input.ID), input.Cacheable, iamGetOut, func() error {
		iamGetOut, err = h.poolHandler.GetUser(iamGetIn, ctx)
		return err
	}, func() interface{} {
		return iamGetOut
	})
	if err != nil {

		if _err, ok := tkErr.IsError(err); ok {
			switch _err.Code() {
			case iamCnt.GPRCUserDoesNotExistErrCode:
				return nil, tkErr.New(cnt.AuthUserNotFoundErr)
			}
		}
		zap.L().With(
			zap.String(cnt.Authentication, "h.poolHandler.GetUser"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", iamGetIn),
		).Error(err.Error())
		return nil, tkErr.New(cnt.AuthInternalServerErr)
	}

	if iamGetOut.Frozen {
		zap.L().With(
			zap.String(cnt.Authentication, "if iamGetOut.Frozen"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", iamGetIn),
			zap.Any("output", iamGetOut),
		).Warn(tkErr.New(cnt.AuthUserIsFrozenErr, iamGetOut.ID).Message())
	}

	output = &common.GetUserOutput{
		ID:          iamGetOut.ID,
		DisplayName: iamGetOut.DisplayName,
		Account:     iamGetOut.Account,
		Email:       iamGetOut.Email,
		Frozen:      iamGetOut.Frozen,
		Extra:       utility.Bytes2Extra(iamGetOut.Extra),
	}

	return output, nil
}
