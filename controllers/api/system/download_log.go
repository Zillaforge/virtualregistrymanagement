package system

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/api"
	util "VirtualRegistryManagement/utility"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

type DownloadLogInput struct {
	Filename   string `form:"filename" binding:"required"`
	From       string `form:"from" binding:"required"`
	Attachment bool   `form:"attachment"`
}

func DownloadLog(c *gin.Context) {
	var (
		err        error
		statusCode = http.StatusOK
		input      = &DownloadLogInput{}
		requestID  = util.MustGetContextRequestID(c)
		funcName   = tkUtils.NameOfFunction().Name()
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"err":        &err,
		"input":      &input,
		"statusCode": &statusCode,
	})

	if err = c.ShouldBindWith(input, binding.Query); err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "c.ShouldBindWith()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", input),
		).Error(err.Error())
		statusCode = http.StatusBadRequest
		err = api.Malformed(err)
		util.ResponseWithType(c, statusCode, err)
		return
	}

	verifiedFilename, verifiedFrom := false, false
	for _, f := range getLogFiles(requestID) {
		if f.Name == input.Filename {
			verifiedFilename = true
			break
		}
	}

	for _, pp := range walkPP {
		if pp == input.From {
			verifiedFrom = true
			break
		}
	}

	if !verifiedFilename || !verifiedFrom {
		zap.L().With(
			zap.String(cnt.Controller, "if verified"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", input),
		).Error("filename or from are not support")
		statusCode = http.StatusBadRequest
		err = tkErr.New(cnt.AdminAPILogFileNotFoundErr)
		util.ResponseWithType(c, statusCode, err)
		return
	}
	if input.Attachment {
		// c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", input.Filename))
	}

	c.File(filepath.Join(mviper.GetString(input.From), input.Filename))

}
