package api

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

/*
Check model exists
endpoint: /project/:project-id/model/:model-id

errors:
- 12000000(internal server error)
- 12000006(model (%s) not found)
*/
func VerifyMemberAclHasInProject(c *gin.Context) {
	var (
		input      = &ResourceIDInput{ID: c.Param(cnt.ParamMemberAclID)}
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
	if match, _ := regexp.MatchString(uuidRegexpString, input.ID); !match {
		err = tkErr.New(cnt.MidRepositoryNotFoundErr, input.ID)
		zap.L().With(
			zap.String(cnt.Middleware, "regexp.MatchString(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.String(cnt.ParamMemberAclID, input.ID),
		).Error(err.Error())
		statusCode = http.StatusNotFound
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// call GetRepository
	namespace := c.GetString(cnt.CtxNamespace)
	listInput := &pb.ListRegistriesInput{
		Where: []string{
			"member-acl-id=" + input.ID,
		},
		Namespace: &namespace,
	}
	listOutput, err := vrm.ListRegistries(listInput, c)
	if err != nil {
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Middleware, "vrm.ListRegistries(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.MidInternalServerErrorErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}
	if listOutput.Count == 0 {
		err = tkErr.New(cnt.MidTagNotFoundErr)
		zap.L().With(
			zap.String(cnt.Middleware, "vrm.ListRegistries(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listInput),
		).Info(err.Error())
		statusCode = http.StatusNotFound
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	if listOutput.Data[0].ProjectID != c.GetString(cnt.CtxProjectID) {
		statusCode = http.StatusForbidden
		err = tkErr.New(cnt.MidTagIsReadOnlyErr, input.ID)
		utility.ResponseWithType(c, statusCode, err)
		return
	}
	c.Set(cnt.CtxCreator, listOutput.Data[0].Creator)
	c.Set(cnt.CtxMemberAclID, input.ID)

	c.Next()
}
