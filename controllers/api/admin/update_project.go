package admin

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/api"
	"VirtualRegistryManagement/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

type UpdateProjectInput struct {
	ID             string `json:"-"`
	LimitCount     *int64 `json:"softLimitCount"`
	LimitSizeBytes *int64 `json:"softLimitSize"`
	_              struct{}
}

type UpdateProjectOutput struct {
	ID             string `json:"id"`
	LimitCount     int64  `json:"softLimitCount"`
	LimitSizeBytes int64  `json:"softLimitSize"`
	_              struct{}
}

func UpdateProject(c *gin.Context) {
	var (
		input      = &UpdateProjectInput{ID: c.GetString(cnt.CtxProjectID)}
		output     = &UpdateProjectOutput{}
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

	updateInput := &pb.UpdateProjectInput{
		ID:             input.ID,
		LimitCount:     input.LimitCount,
		LimitSizeBytes: input.LimitSizeBytes,
	}
	updateOutput, err := vrm.UpdateProject(updateInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCProjectNotFoundErr.Code():
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.AdminAPIProjectNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "vrm.UpdateProject(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output = &UpdateProjectOutput{
		ID:             updateOutput.ID,
		LimitCount:     updateOutput.LimitCount,
		LimitSizeBytes: updateOutput.LimitSizeBytes,
	}
	utility.ResponseWithType(c, statusCode, output)
}
