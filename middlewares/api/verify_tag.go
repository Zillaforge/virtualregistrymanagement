package api

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

/*
errors:
- 12000000(internal server error)
- 12000006(model (%s) not found)
*/
func VerifyTag(c *gin.Context) {
	var (
		input      = &ResourceIDInput{ID: c.Param(cnt.ParamTagID)}
		err        error
		requestID      = utility.MustGetContextRequestID(c)
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusOK
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"error":      &err,
		"statusCode": &statusCode,
	})

	// validate input
	if match, _ := regexp.MatchString(uuidRegexpString, input.ID); !match {
		err = tkErr.New(cnt.MidTagNotFoundErr, input.ID)
		zap.L().With(
			zap.String(cnt.Middleware, "regexp.MatchString(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.String(cnt.ParamTagID, input.ID),
		).Error(err.Error())
		statusCode = http.StatusNotFound
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	getInput := &pb.GetInput{
		ID: input.ID,
	}
	getOutput, err := vrm.GetTag(getInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCTagNotFoundErr.Code():
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.MidTagNotFoundErr, input.ID)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Middleware, "vrm.GetTag(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.MidInternalServerErrorErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	c.Set(cnt.CtxTagID, input.ID)
	c.Set(cnt.CtxTagProtect, getOutput.Tag.Protect)

	c.Next()
}
