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

type VerifyTagHasInACLInput struct {
	ID        string
	Namespace string
	ProjectID string
	UserID    string
}

/*
Check model exists
endpoint: /project/:project-id/model/:model-id

errors:
- 12000000(internal server error)
- 12000006(model (%s) not found)
*/
func VerifyTagHasInACL(c *gin.Context) {
	var (
		input = &VerifyTagHasInACLInput{
			ID:        c.Param(cnt.ParamTagID),
			Namespace: c.GetString(cnt.CtxNamespace),
			ProjectID: c.GetString(cnt.CtxProjectID),
			UserID:    c.GetString(cnt.CtxUserID),
		}
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
		"error":      &err,
		"statusCode": &statusCode,
	})

	// validate input
	if match, _ := regexp.MatchString(uuidRegexpString, input.ID); !match {
		err = tkErr.New(cnt.MidRepositoryNotFoundErr, input.ID)
		zap.L().With(
			zap.String(cnt.Middleware, "regexp.MatchString(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.String(cnt.ParamTagID, input.ID),
		).Error(err.Error())
		statusCode = http.StatusNotFound
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	listRegistriesInput := &pb.ListRegistriesInput{
		Limit:  int32(-1),
		Offset: int32(0),
		Flag: &pb.RegistryFlag{
			UserID:    &input.UserID,
			ProjectID: &input.ProjectID,

			BelongUser:    true,
			ProjectLimit:  true,
			ProjectPublic: true,
			GlobalLimit:   true,
			GlobalPublic:  true,
		},
		Namespace: &input.Namespace,
	}
	if role := c.GetString(cnt.CtxTenantRole); supportRoles[role] {
		listRegistriesInput.Flag.BelongProject = true
	}
	listRegistriesOutput, err := vrm.ListRegistries(listRegistriesInput, c)
	if err != nil {
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Middleware, "vrm.ListRegistries(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listRegistriesInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.MidInternalServerErrorErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}
	if listRegistriesOutput.Count == 0 {
		zap.L().With(
			zap.String(cnt.Middleware, "listRegistriesOutput.Count == 0"),
			zap.String(cnt.RequestID, requestID),
			zap.String(cnt.ParamTagID, input.ID),
		).Error(err.Error())
		statusCode = http.StatusNotFound
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	existID := map[string]bool{}
	whereInput := []string{}
	for _, registry := range listRegistriesOutput.Data {
		if _, ok := existID[registry.TagID]; !ok && registry.TagID != "" {
			existID[registry.TagID] = true
			whereInput = append(whereInput, "id="+registry.TagID)
		}
	}

	if len(whereInput) == 0 {
		zap.L().With(
			zap.String(cnt.Middleware, "len(whereInput) == 0"),
			zap.String(cnt.RequestID, requestID),
			zap.String(cnt.ParamTagID, input.ID),
		).Error(err.Error())
		statusCode = http.StatusNotFound
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	namespace := c.GetString(cnt.CtxNamespace)
	listTagsInput := &pb.ListNamespaceInput{
		Limit:     int32(-1),
		Offset:    int32(0),
		Where:     whereInput,
		Namespace: &namespace,
	}
	listTagsOutput, err := vrm.ListTags(listTagsInput, c)
	if err != nil {
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Middleware, "vrm.ListTags(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listTagsInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.MidInternalServerErrorErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	var Tag *pb.TagDetail
	for _, tag := range listTagsOutput.Data {
		if tag.Tag.ID == input.ID {
			Tag = tag
		}
	}

	if Tag == nil {
		statusCode = http.StatusForbidden
		err = tkErr.New(cnt.MidTagIsReadOnlyErr, input.ID)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	c.Set(cnt.CtxCreator, Tag.Repository.Creator)
	c.Set(cnt.CtxTagProtect, Tag.Tag.Protect)
	c.Set(cnt.CtxTagID, input.ID)

	c.Next()

}
