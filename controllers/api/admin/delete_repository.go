package admin

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

type DeleteRepositoryInput struct {
	ID string `json:"-"`
	_  struct{}
}

type DeleteRepositoryOutput struct {
	_ struct{}
}

func DeleteRepository(c *gin.Context) {
	var (
		input      = &DeleteRepositoryInput{ID: c.GetString(cnt.CtxRepositoryID)}
		output     = &DeleteRepositoryOutput{}
		err        error
		requestID      = utility.MustGetContextRequestID(c)
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusNoContent
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	if c.GetBool(cnt.CtxRepositoryProtect) {
		statusCode = http.StatusForbidden
		err = tkErr.New(cnt.AdminAPIResourceIsProtectedErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	listTagsInput := &pb.ListNamespaceInput{
		Where: []string{"repository-id=" + input.ID},
	}
	listTagsOutput, err := vrm.ListTags(listTagsInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "vrm.ListTags(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listTagsInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	for _, tag := range listTagsOutput.Data {
		if tag.Tag.Protect {
			statusCode = http.StatusForbidden
			err = tkErr.New(cnt.AdminAPIResourceIsProtectedErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
	}

	deleteInput := &pb.DeleteInput{
		Where: []string{"ID==" + input.ID},
	}
	_, err = vrm.DeleteRepository(deleteInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "vrm.DeleteRepository(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", deleteInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	utility.ResponseWithType(c, statusCode, output)
}
