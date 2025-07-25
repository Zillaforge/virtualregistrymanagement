package admin

import (
	cnt "VirtualRegistryManagement/constants"
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

type GetRepositoryInput struct {
	ID string `json:"-"`
	_  struct{}
}

type GetRepositoryOutput struct {
	Repository
	_ struct{}
}

func GetRepository(c *gin.Context) {
	var (
		input      = &GetRepositoryInput{ID: c.GetString(cnt.CtxRepositoryID)}
		output     = &GetRepositoryOutput{}
		err        error
		requestID      = utility.MustGetContextRequestID(c)
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusOK
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	getRepositoryInput := &pb.GetInput{
		ID: input.ID,
	}
	getRepositoryOutput, err := vrm.GetRepository(getRepositoryInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCRepositoryNotFoundErr.Code():
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIRepositoryNotFoundErr, input.ID)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Controller, "vrm.GetRepository(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getRepositoryInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Repository.ExtractByProto(c, getRepositoryOutput)
	utility.ResponseWithType(c, statusCode, output)
}
