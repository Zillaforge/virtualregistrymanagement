package api

import (
	auth "VirtualRegistryManagement/authentication"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"net/http"

	authCom "VirtualRegistryManagement/authentication/common"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

// VerifyUserTokenInput ...
type VerifyUserTokenInput struct {
	UserID     string
	Token      string
	SAATUserID string
	_          struct{}
}

// VerifyUserToken ...
//
// contexts:
//   - CtxUserID
//
// errors:
//   - 12000001(permission denied)
//   - 12000002(the user has been frozen, please contact administrator)
func VerifyUserToken(c *gin.Context) {
	var (
		funcName   = tkUtils.NameOfFunction().Name()
		requestID  = utility.MustGetContextRequestID(c)
		statusCode = http.StatusOK
		err        error
		input      = &VerifyUserTokenInput{}
	)

	f := tracer.StartWithGinContext(c, funcName)

	defer f(tracer.Attributes{
		"input":      input,
		"err":        &err,
		"statusCode": &statusCode,
	})

	input.Token = c.GetHeader(cnt.HdrAuthorization)
	if c.GetHeader(cnt.HdrUserIDFromKong) == "" {
		verifyTokenInput := &authCom.VerifyTokenInput{
			Token: input.Token,
		}
		verifyTokenOutput, verifyTokenErr := auth.Use().VerifyToken(c, verifyTokenInput)
		if verifyTokenErr != nil {
			zap.L().With(
				zap.String(cnt.Middleware, "auth.Use().VerifyToken()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", verifyTokenInput),
			).Error(verifyTokenErr.Error())
			statusCode = http.StatusForbidden
			err = tkErr.New(cnt.MidPermissionDeniedErr).WithInner(verifyTokenErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		// 確認使用者是否被凍結，若是則回傳錯誤訊息
		if verifyTokenOutput.Frozen {
			err = tkErr.New(cnt.MidUserHasBeenFrozenErr)
			zap.L().With(
				zap.String(cnt.Middleware, "if verifyTokenOutput.Frozen"),
				zap.String(cnt.RequestID, requestID),
			).Error(err.Error())
			statusCode = http.StatusForbidden
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		input.UserID = verifyTokenOutput.UserID
		input.SAATUserID = verifyTokenOutput.SAATUserID
	} else {
		input.UserID = c.GetHeader(cnt.HdrUserIDFromKong)
		input.SAATUserID = c.GetHeader(cnt.HdrSAATUserIDFromKong)
	}
	c.Set(cnt.CtxUserID, input.UserID)
	c.Set(cnt.CtxSAATUserID, input.SAATUserID)
	c.Next()
}
