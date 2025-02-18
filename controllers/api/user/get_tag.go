package user

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

type GetTagInput struct {
	ID string `json:"-"`
	_  struct{}
}

type GetTagOutput struct {
	Tag
	_ struct{}
}

func GetTag(c *gin.Context) {
	var (
		input      = &GetTagInput{ID: c.GetString(cnt.CtxTagID)}
		output     = &GetTagOutput{}
		err        error
		requestID      = utility.MustGetContextRequestID(c)
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusOK
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	getTagInput := &pb.GetInput{
		ID: input.ID,
	}
	getTagOutput, err := vrm.GetTag(getTagInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCTagNotFoundErr.Code():
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.UserAPITagNotFoundErr, input.ID)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Controller, "vrm.GetTag(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getTagInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Tag.ExtractByProto(c, getTagOutput)
	utility.ResponseWithType(c, statusCode, output)
}
