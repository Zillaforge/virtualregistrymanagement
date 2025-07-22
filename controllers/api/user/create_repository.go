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

type CreateRepositoryInput struct {
	Name            string `json:"name" binding:"required"`
	OperatingSystem string `json:"operatingSystem" binding:"required,oneof='' linux windows"`
	Description     string `json:"description"`

	Namespace string `json:"-"`
	ProjectID string `json:"-"`
	Creator   string `json:"-"`

	_ struct{}
}

type CreateRepositoryOutput struct {
	Repository
	_ struct{}
}

func CreateRepository(c *gin.Context) {
	var (
		input      = &CreateRepositoryInput{}
		output     = &CreateRepositoryOutput{}
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

	createRepositoryInput := &pb.RepositoryInfo{
		Name:            input.Name,
		OperatingSystem: input.OperatingSystem,
		Description:     input.Description,
		Namespace:       c.GetString(cnt.CtxNamespace),
		ProjectID:       c.GetString(cnt.CtxProjectID),
		Creator:         c.GetString(cnt.CtxUserID),
	}
	createRepositoryOutput, err := vrm.CreateRepository(createRepositoryInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCRepositoryExistErr.Code():
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.UserAPIRepositoryExistErr, input.Name)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Controller, "vrm.CreateRepository(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createRepositoryInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Repository.ExtractByProto(c, createRepositoryOutput)
	utility.ResponseWithType(c, statusCode, output)
}
