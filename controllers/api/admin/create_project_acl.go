package admin

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/api"
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
	TagID     string  `json:"tagId" binding:"required"`
	ProjectID *string `json:"projectId"`
	_         struct{}
}

type CreateProjectAclOutput struct {
	ProjectAcl ProjectAcl `json:"projectAcl"`
	_          struct{}
}

func CreateProjectAcl(c *gin.Context) {
	var (
		input      = &CreateProjectAclInput{}
		output     = &CreateProjectAclOutput{}
		err        error
		requestID      = utility.MustGetContextRequestID(c)
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusCreated
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	if tagID := c.GetString(cnt.CtxTagID); tagID != "" {
		input.TagID = c.GetString(cnt.CtxTagID)
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

	createProjectAclBatchInput := &pb.ProjectAclBatchInfo{}
	createProjectAclBatchInput.Data = append(createProjectAclBatchInput.Data, &pb.ProjectAclInfo{
		TagID:     input.TagID,
		ProjectID: input.ProjectID,
	})

	createProjectAclBatchOutput, err := vrm.CreateProjectAclBatch(createProjectAclBatchInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "vrm.CreateProjectAclBatch(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createProjectAclBatchInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	for _, projectAcl := range createProjectAclBatchOutput.Data {
		output.ProjectAcl.ExtractByProto(c, projectAcl)
	}
	utility.ResponseWithType(c, statusCode, output)
}
