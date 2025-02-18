package admin

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

type GetExportInput struct {
	ID string `json:"-"`
	_  struct{}
}

type GetExportOutput struct {
	Export
	_ struct{}
}

func GetExport(c *gin.Context) {
	var (
		input      = &GetExportInput{ID: c.GetString(cnt.CtxExportID)}
		output     = &GetExportOutput{}
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

	getExportInput := &pb.GetInput{
		ID: input.ID,
	}
	getExportOutput, err := vrm.GetExport(getExportInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCExportNotFoundErr.Code():
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIExportNotFoundErr, input.ID)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Controller, "vrm.GetExport(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getExportInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Export.ExtractByProto(c, getExportOutput)
	utility.ResponseWithType(c, statusCode, output)
}
