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
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

type CreateTagInput struct {
	RepositoryID    string `json:"repositoryId" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Type            string `json:"type" binding:"required,oneof=common increase"`
	DiskFormat      string `json:"diskFormat" binding:"required,oneof=ami ari aki vhd vmdk raw qcow2 vdi iso"`
	ContainerFormat string `json:"containerFormat" binding:"required,oneof=ami ari aki bare ovf"`
	Extra           []byte `json:"extra"`

	_ struct{}
}

type CreateTagOutput struct {
	Tag
	_ struct{}
}

func CreateTag(c *gin.Context) {
	var (
		input      = &CreateTagInput{}
		output     = &CreateTagOutput{}
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
			RepositoryID: input.RepositoryID,
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
				err = tkErr.New(cnt.AdminAPITagExistErr, input.Name)
			case cCnt.GRPCExceedAllowedQuotaErr.Code():
				statusCode = http.StatusRequestEntityTooLarge
				err = tkErr.New(cnt.AdminAPIExceedAllowedQuotaErr)
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
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Tag.ExtractByProto(c, createTagOutput)
	utility.ResponseWithType(c, statusCode, output)
}
