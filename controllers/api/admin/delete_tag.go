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
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

type DeleteTagInput struct {
	ID string `json:"-"`
	_  struct{}
}

type DeleteTagOutput struct {
	_ struct{}
}

func DeleteTag(c *gin.Context) {
	var (
		input      = &DeleteTagInput{ID: c.GetString(cnt.CtxTagID)}
		output     = &DeleteTagOutput{}
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

	if c.GetBool(cnt.CtxTagProtect) {
		statusCode = http.StatusForbidden
		err = tkErr.New(cnt.AdminAPIResourceIsProtectedErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	deleteInput := &pb.DeleteInput{
		Where: []string{"ID=" + input.ID},
	}
	_, err = vrm.DeleteTag(deleteInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "vrm.DeleteTag(...)"),
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
