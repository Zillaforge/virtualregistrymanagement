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

type ListTagsInput struct {
	pagination
	Where []string `json:"where" form:"where"`

	AdminRole     *bool `json:"adminRole" form:"adminRole"`
	Creator       *bool `json:"creator" form:"creator"`
	ProjectLimit  *bool `json:"projectLimit" form:"projectLimit"`
	ProjectPublic *bool `json:"projectPublic" form:"projectPublic"`
	GlobalLimit   *bool `json:"globalLimit" form:"globalLimit"`
	GlobalPublic  *bool `json:"globalPublic" form:"globalPublic"`

	Namespace string `json:"-"`
	ProjectID string `json:"-"`
	UserID    string `json:"-"`
	_         struct{}
}

type ListTagsOutput struct {
	Tags  []Tag `json:"tags"`
	Total int   `json:"total"`
	_     struct{}
}

func ListTags(c *gin.Context) {
	var (
		input = &ListTagsInput{
			Namespace: c.GetString(cnt.CtxNamespace),
			UserID:    c.GetString(cnt.CtxUserID),
			ProjectID: c.GetString(cnt.CtxProjectID),
		}
		output       = &ListTagsOutput{}
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
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	if err = c.ShouldBindQuery(input); err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "c.ShouldBindQuery(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("obj", input),
		).Error(err.Error())
		statusCode = http.StatusBadRequest
		err = api.Malformed(err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	if repositoryID := c.GetString(cnt.CtxRepositoryID); repositoryID != "" {
		input.Where = append(input.Where, "repository-id=="+repositoryID)
	}

	listRegistriesInput := &pb.ListRegistriesInput{
		Limit:  int32(-1),
		Offset: int32(0),
		Where:  input.Where,
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
	// tenant-owner and tenant-admin allow to get all in project
	if role := c.GetString(cnt.CtxTenantRole); supportRoles[role] {
		listRegistriesInput.Flag.BelongProject = true
	}
	if input.AdminRole != nil && !*input.AdminRole {
		listRegistriesInput.Flag.BelongProject = false
	}
	if input.Creator != nil && !*input.Creator {
		listRegistriesInput.Flag.BelongUser = false
	}
	if input.ProjectLimit != nil && !*input.ProjectLimit {
		listRegistriesInput.Flag.ProjectLimit = false
	}
	if input.ProjectPublic != nil && !*input.ProjectPublic {
		listRegistriesInput.Flag.ProjectPublic = false
	}
	if input.GlobalLimit != nil && !*input.GlobalLimit {
		listRegistriesInput.Flag.GlobalLimit = false
	}
	if input.GlobalPublic != nil && !*input.GlobalPublic {
		listRegistriesInput.Flag.GlobalPublic = false
	}

	listRegistriesOutput, err := vrm.ListRegistries(listRegistriesInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErr.Code():
				statusCode = http.StatusBadRequest
				if v, exist := e.Get("field"); exist {
					err = tkErr.New(cnt.UserAPIQueryNotSupportErr, "where", v)
				} else {
					err = tkErr.New(cnt.UserAPIIllegalWhereQueryFormatErr)
				}
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "vrm.ListRegistries(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listRegistriesInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}
	if listRegistriesOutput.Count == 0 {
		utility.ResponseWithType(c, statusCode, output)
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
		output.Tags = []Tag{}
		output.Total = 0
		utility.ResponseWithType(c, statusCode, output)
		return
	}

	namespace := c.GetString(cnt.CtxNamespace)
	listTagsInput := &pb.ListNamespaceInput{
		Limit:     int32(input.Limit),
		Offset:    int32(input.Offset),
		Where:     whereInput,
		Namespace: &namespace,
	}
	listTagsOutput, err := vrm.ListTags(listTagsInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErr.Code():
				statusCode = http.StatusBadRequest
				if v, exist := e.Get("field"); exist {
					err = tkErr.New(cnt.UserAPIQueryNotSupportErr, "where", v)
				} else {
					err = tkErr.New(cnt.UserAPIIllegalWhereQueryFormatErr)
				}
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "vrm.ListTags(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listTagsInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Tags = []Tag{}
	output.Total = int(listTagsOutput.Count)
	for _, data := range listTagsOutput.Data {
		m := Tag{}
		m.ExtractByProto(c, data)
		output.Tags = append(output.Tags, m)
	}
	utility.ResponseWithType(c, statusCode, output)
}
