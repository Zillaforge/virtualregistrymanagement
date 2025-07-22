package admin

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/api"
	"VirtualRegistryManagement/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

type DownloadImageInput struct {
	Filepath string `json:"filepath" binding:"required"`

	TagID string `json:"-"`
	_     struct{}
}

type DownloadImageOutput struct {
	_ struct{}
}

func DownloadImage(c *gin.Context) {
	var (
		input      = &DownloadImageInput{TagID: c.GetString(cnt.CtxTagID)}
		output     = &DownloadImageOutput{}
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

	if err = c.ShouldBindJSON(input); err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "c.ShouldBindJSON(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("obj", input),
		).Error(err.Error())
		statusCode = http.StatusBadRequest
		err = api.Malformed(err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	downloadTagInput := &pb.ImageDataInput{
		TagID:    input.TagID,
		Filepath: input.Filepath,
	}
	err = vrm.DownloadTag(downloadTagInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCNotSupportTypeErr.Code():
				statusCode = http.StatusBadRequest
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCFileExistErr.Code():
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIFileExistErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "vrm.DownloadTag(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", downloadTagInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	utility.ResponseWithType(c, statusCode, output)
}
