package api

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

/*
Check model exists
endpoint: /project/:project-id/model/:model-id

errors:
- 12000000(internal server error)
- 12000006(model (%s) not found)
*/
func VerifyRepositoryHasInProject(c *gin.Context) {
	var (
		input      = &ResourceIDInput{ID: c.Param(cnt.ParamRepositoryID)}
		err        error
		requestID      = utility.MustGetContextRequestID(c)
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusOK
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"error":      &err,
		"statusCode": &statusCode,
	})

	// validate input
	match, _ := regexp.MatchString(uuidRegexpString, input.ID)
	if !match {
		err = tkErr.New(cnt.MidRepositoryNotFoundErr, input.ID)
		zap.L().With(
			zap.String(cnt.Middleware, "regexp.MatchString(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.String(cnt.ParamRepositoryID, input.ID),
		).Error(err.Error())
		statusCode = http.StatusNotFound
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// call GetRepository
	getInput := &pb.GetInput{
		ID: input.ID,
	}
	getOutput, err := vrm.GetRepository(getInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCRepositoryNotFoundErr.Code():
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.MidRepositoryNotFoundErr, input.ID)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Middleware, "vrm.GetRepository(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.MidInternalServerErrorErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	if getOutput.Repository.ProjectID != c.GetString(cnt.CtxProjectID) {
		statusCode = http.StatusForbidden
		err = tkErr.New(cnt.MidRepositoryIsReadOnlyErr, input.ID)
		utility.ResponseWithType(c, statusCode, err)
		return
	}
	c.Set(cnt.CtxCreator, getOutput.Repository.Creator)
	c.Set(cnt.CtxRepositoryID, input.ID)
	c.Set(cnt.CtxRepositoryProtect, getOutput.Repository.Protect)

	c.Next()
}
