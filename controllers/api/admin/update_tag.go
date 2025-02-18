package admin

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/api"
	"VirtualRegistryManagement/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

type UpdateTagInput struct {
	ID    string  `json:"-"`
	Name  *string `json:"name"`
	Extra *[]byte `json:"extra"`
	_     struct{}
}

type UpdateTagOutput struct {
	Tag
	_ struct{}
}

func UpdateTag(c *gin.Context) {
	var (
		input      = &UpdateTagInput{ID: c.GetString(cnt.CtxTagID)}
		output     = &UpdateTagOutput{}
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

	if err = c.ShouldBindWith(input, binding.JSON); err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "c.ShouldBindWith()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", input),
		).Error(err.Error())
		statusCode = http.StatusBadRequest
		err = api.Malformed(err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	updateInput := &pb.UpdateTagInput{
		ID:   input.ID,
		Name: input.Name,
	}
	if input.Extra != nil {
		updateInput.Extra = *input.Extra
	}
	updateOutput, err := vrm.UpdateTag(updateInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCTagNotFoundErr.Code():
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.AdminAPITagNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "vrm.UpdateTag(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(c, updateOutput)
	utility.ResponseWithType(c, statusCode, output)
}
