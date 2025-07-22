package user

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/api"
	"VirtualRegistryManagement/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

type UpdateTagInput struct {
	ID   string  `json:"-"`
	Name *string `json:"name"`
	_    struct{}
}

type UpdateTagOutput struct {
	Tag
	_ struct{}
}

func UpdateTag(c *gin.Context) {
	var (
		input        = &UpdateTagInput{ID: c.GetString(cnt.CtxTagID)}
		output       = &UpdateTagOutput{}
		err          error
		requestID        = utility.MustGetContextRequestID(c)
		funcName         = tkUtils.NameOfFunction().Name()
		statusCode   int = http.StatusOK
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

	if err = c.ShouldBindWith(input, binding.JSON); err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "c.ShouldBindWith()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", input),
		).Error(err.Error())
		statusCode = http.StatusBadRequest
		err = api.Malformed(err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// Update 只有 Creator, TENANT_OWNER, TENANT_ADMIN 才可以
	// 1. 檢查是否為 TENANT_OWNER, TENANT_ADMIN
	// 2. 檢查是否為 Creator
	if role := c.GetString(cnt.CtxTenantRole); !supportRoles[role] &&
		c.GetString(cnt.CtxCreator) != c.GetString(cnt.CtxUserID) {
		statusCode = http.StatusUnauthorized
		err = tkErr.New(cnt.UserAPIUnauthorizedOpErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	updateInput := &pb.UpdateTagInput{
		ID:   input.ID,
		Name: input.Name,
	}
	updateOutput, err := vrm.UpdateTag(updateInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCTagNotFoundErr.Code():
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.UserAPITagNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "vrm.UpdateTag(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(c, updateOutput)
	utility.ResponseWithType(c, statusCode, output)
}
