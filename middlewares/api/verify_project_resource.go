package api

import (
	cnt "VirtualRegistryManagement/constants"
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

// VerifyProjectResource ...
func VerifyProjectResource(c *gin.Context) {
	var (
		funcName   = tkUtils.NameOfFunction().Name()
		requestID  = utility.MustGetContextRequestID(c)
		statusCode = http.StatusOK
		err        error
		input      = &ResourceIDInput{ID: c.GetString(cnt.CtxProjectID)}
	)

	f := tracer.StartWithGinContext(c, funcName)

	defer f(tracer.Attributes{
		"input":      input,
		"err":        &err,
		"statusCode": &statusCode,
	})

	getProjectInput := &pb.GetInput{
		ID: input.ID,
	}
	getOutput, err := vrm.GetProject(getProjectInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCProjectNotFoundErr.Code():
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.MidProjectNotFoundErr, input.ID)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Middleware, "vrm.GetProject(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getProjectInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.MidInternalServerErrorErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}
	c.Set(cnt.CtxProjectSoftLimitSize, getOutput.LimitSizeBytes)
	c.Set(cnt.CtxProjectSoftLimitCount, getOutput.LimitCount)

	where := []string{"project-id=" + input.ID}
	listRepositoriesInput := &pb.ListNamespaceInput{
		Limit:  -1,
		Offset: 0,
		Where:  where,
	}
	listRepositoriesOutput, err := vrm.ListRepositories(listRepositoriesInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Middleware, "vrm.ListRepositories(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listRepositoriesInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.MidInternalServerErrorErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	var (
		currentCount int64 = 0
		currentSize  int64 = 0
	)
	for _, repository := range listRepositoriesOutput.Data {
		currentCount += int64(len(repository.Tags))
		for _, tag := range repository.Tags {
			currentSize += int64(tag.Size)
		}
	}

	c.Set(cnt.CtxProjectUsedSize, currentSize)
	c.Set(cnt.CtxProjectUsedCount, currentCount)

	if getOutput.LimitCount == -1 {
		c.Set(cnt.CtxProjectCountFlag, false)
	} else {
		c.Set(cnt.CtxProjectCountFlag, currentCount >= getOutput.LimitCount)
	}

	if getOutput.LimitSizeBytes == -1 {
		c.Set(cnt.CtxProjectSizeFlag, false)
	} else {
		c.Set(cnt.CtxProjectSizeFlag, currentSize >= getOutput.LimitSizeBytes)
	}
	c.Next()

}
