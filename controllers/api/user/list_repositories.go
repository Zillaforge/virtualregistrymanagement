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

type ListRepositoriesInput struct {
	pagination
	Where []string `json:"where" form:"where"`

	Namespace string `json:"-"`
	ProjectID string `json:"-"`
	UserID    string `json:"-"`
	_         struct{}
}

type ListRepositoriesOutput struct {
	Repositories []Repository `json:"repositories"`
	Total        int          `json:"total"`
	_            struct{}
}

func ListRepositories(c *gin.Context) {
	var (
		input = &ListRepositoriesInput{
			Namespace: c.GetString(cnt.CtxNamespace),
			UserID:    c.GetString(cnt.CtxUserID),
			ProjectID: c.GetString(cnt.CtxProjectID),
		}
		output = &ListRepositoriesOutput{
			Repositories: []Repository{},
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

	listRegistriesInput := &pb.ListRegistriesInput{
		Limit:  int32(-1),
		Offset: int32(0),
		Where:  append(input.Where, "project-id="+input.ProjectID),
		Flag: &pb.RegistryFlag{
			UserID:    &input.UserID,
			ProjectID: &input.ProjectID,

			BelongUser:    true,
			ProjectLimit:  true,
			ProjectPublic: true,
		},
		Namespace: &input.Namespace,
	}
	// tenant-owner and tenant-admin allow to get all in project
	if role := c.GetString(cnt.CtxTenantRole); supportRoles[role] {
		listRegistriesInput.Flag.BelongProject = true
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

	existID := map[string]map[string]bool{}
	whereInput := []string{}
	for _, registry := range listRegistriesOutput.Data {
		if _, ok := existID[registry.RepositoryID]; !ok {
			existID[registry.RepositoryID] = map[string]bool{}
			whereInput = append(whereInput, "id="+registry.RepositoryID)
		}
		existID[registry.RepositoryID][registry.TagID] = true
	}

	if len(whereInput) == 0 {
		output.Repositories = []Repository{}
		output.Total = 0
		utility.ResponseWithType(c, statusCode, output)
		return
	}

	listRepositoriesInput := &pb.ListNamespaceInput{
		Limit:     int32(input.Limit),
		Offset:    int32(input.Offset),
		Where:     whereInput,
		Namespace: &input.Namespace,
	}
	listRepositoriesOutput, err := vrm.ListRepositories(listRepositoriesInput, c)
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
			zap.String(cnt.Controller, "vrm.ListRepositories(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listRepositoriesInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Total = int(listRepositoriesOutput.Count)
	for _, d := range listRepositoriesOutput.Data {
		data := &pb.RepositoryDetail{
			Repository: d.Repository,
		}
		for _, tag := range d.Tags {
			if _, ok := existID[d.Repository.ID][tag.ID]; ok {
				data.Tags = append(data.Tags, tag)
			}
		}

		m := Repository{}
		m.ExtractByProto(c, data)
		output.Repositories = append(output.Repositories, m)
	}
	utility.ResponseWithType(c, statusCode, output)
}
