package user

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

type DeleteExportInput struct {
	ID string `json:"-"`
	_  struct{}
}

type DeleteExportOutput struct {
	_ struct{}
}

func DeleteExport(c *gin.Context) {
	var (
		input        = &DeleteExportInput{ID: c.GetString(cnt.CtxExportID)}
		output       = &DeleteExportOutput{}
		err          error
		requestID        = utility.MustGetContextRequestID(c)
		funcName         = tkUtils.NameOfFunction().Name()
		statusCode   int = http.StatusNoContent
		supportRoles     = map[string]bool{
			cnt.TenantOwner.String(): true,
			cnt.TenantAdmin.String(): true,
		}
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	// Delete 只有 Creator, TENANT_OWNER, TENANT_ADMIN 才可以
	// 1. 檢查是否為 TENANT_OWNER, TENANT_ADMIN
	// 2. 檢查是否為 Creator
	if role := c.GetString(cnt.CtxTenantRole); !supportRoles[role] &&
		c.GetString(cnt.CtxCreator) != c.GetString(cnt.CtxUserID) {
		statusCode = http.StatusUnauthorized
		err = tkErr.New(cnt.UserAPIUnauthorizedOpErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	deleteInput := &pb.DeleteInput{
		Where: []string{"ID=" + input.ID},
	}
	_, err = vrm.DeleteExport(deleteInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "vrm.DeleteExport(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", deleteInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	utility.ResponseWithType(c, statusCode, output)
}
