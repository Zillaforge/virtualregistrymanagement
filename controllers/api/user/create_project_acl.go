package user

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

type CreateProjectAclInput struct {
	TagID     string `json:"-"`
	ProjectID string `json:"-"`
	_         struct{}
}

type CreateProjectAclOutput struct {
	ProjectAcl ProjectAcl `json:"projectAcl"`
	_          struct{}
}

func CreateProjectAcl(c *gin.Context) {
	var (
		input = &CreateProjectAclInput{
			TagID:     c.GetString(cnt.CtxTagID),
			ProjectID: c.GetString(cnt.CtxProjectID),
		}
		output       = &CreateProjectAclOutput{}
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

	createProjectAclBatchInput := &pb.ProjectAclBatchInfo{}
	createProjectAclBatchInput.Data = append(createProjectAclBatchInput.Data, &pb.ProjectAclInfo{
		TagID:     input.TagID,
		ProjectID: &input.ProjectID,
	})

	createProjectAclBatchOutput, err := vrm.CreateProjectAclBatch(createProjectAclBatchInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "vrm.CreateProjectAclBatch(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createProjectAclBatchInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	for _, projectAcl := range createProjectAclBatchOutput.Data {
		output.ProjectAcl.ExtractByProto(c, projectAcl)
	}
	utility.ResponseWithType(c, statusCode, output)
}
