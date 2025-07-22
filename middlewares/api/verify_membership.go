package api

import (
	auth "VirtualRegistryManagement/authentication"
	authComm "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

type VerifyMembershipTokenInput struct {
	UserID     string
	ProjectID  string
	TenantRole string
	Frozen     bool
	_          struct{}
}

// VerifyMembershipToken ...
//
// contexts:
//   - CtxTenantRole
//
// errors:
//   - 12000003(the membership has been frozen, please contact tenant of admin)
//   - 12000001(permission denied)
func VerifyMembership(c *gin.Context) {
	var (
		funcName   = tkUtils.NameOfFunction().Name()
		requestID  = utility.MustGetContextRequestID(c)
		statusCode = http.StatusOK
		err        error
		input      = &VerifyMembershipTokenInput{}
	)

	f := tracer.StartWithGinContext(c, funcName)

	defer f(tracer.Attributes{
		"input":      &input,
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
		authGetIn := &authComm.GetMembershipInput{
			UserId:    input.UserID,
			ProjectId: input.ProjectID,
		}
		authGetOut, err := auth.Use().GetMembership(c, authGetIn)
		if err != nil {
			if e, ok := tkErr.IsError(err); ok {
				switch e.Code() {
				case cnt.AuthMembershipNotFoundErr.Code():
					statusCode = http.StatusForbidden
					err = tkErr.New(cnt.MidPermissionDeniedErr)
					utility.ResponseWithType(c, statusCode, err)
					return
				}
			}
			zap.L().With(
				zap.String(cnt.Middleware, "authentication.Use().GetMembership(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", authGetIn),
			).Error(err.Error())
			statusCode = http.StatusInternalServerError
			err = tkErr.New(cnt.MidInternalServerErrorErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		input.Frozen = authGetOut.Frozen
		input.TenantRole = authGetOut.TenantRole
	} else {
		// Frozen在 upstream就已經會阻擋，故不再做判斷
		input.TenantRole = kongTenantRoleConvert(c.GetHeader(cnt.HdrUserRoleFromKong))
	}
	c.Set(cnt.CtxProjectID, input.ProjectID)

	// 確認成員是否被凍結，若是則回傳錯誤訊息
	if input.Frozen {
		err = tkErr.New(cnt.MidMembershipHasBeenFrozenErr)
		zap.L().With(
			zap.String(cnt.Middleware, "membership is frozen"),
			zap.String(cnt.RequestID, requestID),
		).Warn(err.Error())
		statusCode = http.StatusForbidden
		utility.ResponseWithType(c, statusCode, err)
		return
	}
	c.Set(cnt.CtxTenantRole, input.TenantRole)
	c.Next()
}

func checkUpstreamFromKong(c *gin.Context) bool {
	if c.GetHeader(cnt.HdrProjectIDFromKong) == "" ||
		c.GetHeader(cnt.HdrProjectActiveFromKong) == "" ||
		c.GetHeader(cnt.HdrUserRoleFromKong) == "" {
		return false
	}
	return true
}

func kongTenantRoleConvert(role string) string {
	switch strings.ToUpper(role) {
	case "TENANT-MEMBER", "TENANT_MEMBER":
		return cnt.TenantMember.String()
	case "TENANT-ADMIN", "TENANT_ADMIN", "SYSTEM-ADMIN", "SYSTEM_ADMIN":
		return cnt.TenantAdmin.String()
	case "TENANT-OWNER", "TENANT_OWNER":
		return cnt.TenantOwner.String()
	default:
		return ""
	}
}
