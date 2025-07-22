package user

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

type CreateTagInput struct {
	Name            string `json:"name" binding:"required"`
	Type            string `json:"type" binding:"required,oneof=common increase"`
	DiskFormat      string `json:"diskFormat" binding:"required,oneof=ami ari aki vhd vmdk raw qcow2 vdi iso"`
	ContainerFormat string `json:"containerFormat" binding:"required,oneof=ami ari aki bare ovf"`
	Extra           []byte `json:"extra"`

	RepositoryID string `json:"-"`
	_            struct{}
}

type CreateTagOutput struct {
	Tag
	_ struct{}
}

func CreateTag(c *gin.Context) {
	var (
		input        = &CreateTagInput{}
		output       = &CreateTagOutput{}
		err          error
		requestID        = utility.MustGetContextRequestID(c)
		funcName         = tkUtils.NameOfFunction().Name()
		statusCode   int = http.StatusCreated
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

	// Create 只有 Creator, TENANT_OWNER, TENANT_ADMIN 才可以
	// 1. 檢查是否為 TENANT_OWNER, TENANT_ADMIN
	// 2. 檢查是否為 Creator
	if role := c.GetString(cnt.CtxTenantRole); !supportRoles[role] &&
		c.GetString(cnt.CtxCreator) != c.GetString(cnt.CtxUserID) {
		statusCode = http.StatusUnauthorized
		err = tkErr.New(cnt.UserAPIUnauthorizedOpErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

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

	createTagInput := &pb.CreateTagInput{
		Tag: &pb.TagInfo{
			Name:         input.Name,
			Type:         input.Type,
			Extra:        input.Extra,
			RepositoryID: c.GetString(cnt.CtxRepositoryID),
		},
		Image: &pb.OpenstackImageInfo{
			DiskFormat:      input.DiskFormat,
			ContainerFormat: input.ContainerFormat,
		},
	}
	createTagOutput, err := vrm.CreateTag(createTagInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCTagExistErr.Code():
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.UserAPITagExistErr, input.Name)
			case cCnt.GRPCExceedAllowedQuotaErr.Code():
				statusCode = http.StatusRequestEntityTooLarge
				err = tkErr.New(cnt.UserAPIExceedAllowedQuotaErr)
			}
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Controller, "vrm.CreateTag(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createTagInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Tag.ExtractByProto(c, createTagOutput)
	utility.ResponseWithType(c, statusCode, output)
}
