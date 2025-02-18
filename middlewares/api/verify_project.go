package api

import (
	auth "VirtualRegistryManagement/authentication"
	authComm "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

type VerifyProjectInput struct {
	UserID    string
	ProjectID string
	Frozen    bool
	_         struct{}
}

// VerifyProject ...
func VerifyProject(c *gin.Context) {
	var (
		funcName   = tkUtils.NameOfFunction().Name()
		requestID  = utility.MustGetContextRequestID(c)
		statusCode = http.StatusOK
		err        error
		input      = &VerifyProjectInput{}
	)

	f := tracer.StartWithGinContext(c, funcName)

	defer f(tracer.Attributes{
		"input":      input,
		"err":        &err,
		"statusCode": &statusCode,
	})

	input.UserID = c.GetString(cnt.CtxUserID)
	input.ProjectID = c.Param(cnt.ParamProjectID)
	input.ProjectID = func() string {
		// 透過 ProjectHybrid取得的ProjectID
		if v := c.GetString(cnt.CtxProjectID); v != "" {
			return v
		}
		// 透過Path取得的ProjectID
		return c.Param(cnt.ParamProjectID)
	}()
	if !checkUpstreamFromKong(c) {
		authGetIn := &authComm.GetProjectInput{
			ID: input.ProjectID,
		}
		authGetOut, err := auth.Use().GetProject(c, authGetIn)
		if err != nil {
			if e, ok := tkErr.IsError(err); ok {
				switch e.Code() {
				case cnt.AuthProjectNotFoundErr.Code():
					statusCode = http.StatusForbidden
					err = tkErr.New(cnt.MidPermissionDeniedErr)
					utility.ResponseWithType(c, statusCode, err)
					return
				case cnt.AuthProjectIsFrozenErr.Code():
					statusCode = http.StatusForbidden
					err = tkErr.New(cnt.MidPermissionDeniedErr)
					utility.ResponseWithType(c, statusCode, err)
					return

				}
			}
			zap.L().With(
				zap.String(cnt.Middleware, "authentication.Use().GetProject(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", authGetIn),
			).Error(err.Error())
			statusCode = http.StatusInternalServerError
			err = tkErr.New(cnt.MidInternalServerErrorErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		input.ProjectID = authGetOut.ID
	}
	c.Set(cnt.CtxProjectID, input.ProjectID)

	c.Next()
}
