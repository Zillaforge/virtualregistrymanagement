package api

import (
	auth "VirtualRegistryManagement/authentication"
	authCom "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

type VerifyAdminTokenInput struct {
	UserID        string
	Account       string
	IsSystemAdmin string
	_             struct{}
}

// VerifyAdminAuthentication is verify user belong project in iam
func VerifyAdminAuthentication(c *gin.Context) {

	var (
		funcName   = tkUtils.NameOfFunction().Name()
		requestID  = utility.MustGetContextRequestID(c)
		statusCode = http.StatusOK
		err        error
		input      = &VerifyAdminTokenInput{}
	)

	f := tracer.StartWithGinContext(c, funcName)

	defer f(tracer.Attributes{
		"input":      input,
		"err":        &err,
		"statusCode": &statusCode,
	})

	input.IsSystemAdmin = c.GetHeader(cnt.HdrSystemAdminFromKong)
	if input.IsSystemAdmin != "" {
		// 若 input.IsSystemAdmin = true，則將 Header 的值存進對應的變數中
		if input.IsSystemAdmin == "true" {
			input.UserID = c.GetHeader(cnt.HdrUserIDFromKong)
			input.Account = c.GetHeader(cnt.HdrUserAccountFromKong)
		} else {
			// 其餘皆回傳 Permission Denied 錯誤
			statusCode = http.StatusForbidden
			err = tkErr.New(cnt.MidPermissionDeniedErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
	} else {
		// 若 Request 沒有經過 Kong，則利用 Token 取得對應的 UserID 與 Account
		verifySystemAdminTokenInput := &authCom.VerifySystemAdminTokenInput{
			Token: getAuthorization(c),
		}
		verifySystemAdminTokenOutput, err := auth.Use().VerifySystemAdminToken(c, verifySystemAdminTokenInput)
		if err != nil {
			if err, ok := tkErr.IsError(err); ok {
				switch err.Code() {
				case cnt.AuthPermissionDeniedErr.Code():
					statusCode = http.StatusForbidden
					err = tkErr.New(cnt.MidPermissionDeniedErr)
					utility.ResponseWithType(c, statusCode, err)
					return
				case cnt.AuthIncorrectFormatErr.Code():
					statusCode = http.StatusBadRequest
					err = tkErr.New(cnt.MidIncorrectFormatErr)
					utility.ResponseWithType(c, statusCode, err)
					return
				}
			}
			zap.L().With(
				zap.String(cnt.Middleware, "auth.Use().VerifySystemAdminToken()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", verifySystemAdminTokenInput),
			).Error(err.Error())
			statusCode = http.StatusUnauthorized
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		input.UserID = verifySystemAdminTokenOutput.UserID
		input.Account = verifySystemAdminTokenOutput.Account
	}

	// set the UserID and UserAccount to context
	c.Set(cnt.CtxUserID, input.UserID)
	c.Set(cnt.CtxUserAccount, input.Account)
	c.Next()
}

func getAuthorization(c *gin.Context) string {
	if c.Request.Method == http.MethodGet && c.Query(cnt.QueryToken) != "" {
		return "Bearer " + c.Query(cnt.QueryToken)
	}
	return c.GetHeader(cnt.HdrAuthorization)
}
