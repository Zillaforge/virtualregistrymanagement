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

type ImportImageInput struct {
	ImageID string `json:"imageId" binding:"required"`
	Version string `json:"version"`

	RepositoryID    string `json:"repositoryId" binding:"required_without_all=Name OperatingSystem Description,excluded_with=Name OperatingSystem Description"`
	Name            string `json:"name" binding:"required_without_all=RepositoryID,excluded_with=RepositoryID"`
	OperatingSystem string `json:"operatingSystem" binding:"required_without_all=RepositoryID,excluded_with=RepositoryID"`
	Description     string `json:"description"`
	ProjectID       string `json:"projectId" binding:"required_without_all=RepositoryID,excluded_with=RepositoryID"`
	Creator         string `json:"creator" binding:"required_without_all=RepositoryID,excluded_with=RepositoryID"`

	Namespace string `json:"namespace" binding:"required"`
	_         struct{}
}

type ImportImageOutput struct {
	Repository repositoryInfo `json:"repository"`
	Tag        tagInfo        `json:"tag"`
	_          struct{}
}

// ImportImage 會自動建立 repository & tag，如果有帶 repository-id 會在下面建立
func ImportImage(c *gin.Context) {
	var (
		input      = &ImportImageInput{}
		output     = &ImportImageOutput{}
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

	// get creator
	if input.RepositoryID != "" {
		getRepositoryInput := &pb.GetInput{
			ID: input.RepositoryID,
		}
		getRepositoryOutput, err := vrm.GetRepository(getRepositoryInput, c)
		if err != nil {
			// Expected errors
			if e, ok := tkErr.IsError(err); ok {
				switch e.Code() {
				case cCnt.GRPCRepositoryNotFoundErr.Code():
					statusCode = http.StatusBadRequest
					err = tkErr.New(cnt.AdminAPIRepositoryNotFoundErr, input.RepositoryID)
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
		input.Creator = getRepositoryOutput.Repository.Creator
	} else {
		createRepositoryInput := &pb.RepositoryInfo{
			Name:            input.Name,
			OperatingSystem: input.OperatingSystem,
			Description:     input.Description,
			Namespace:       input.Namespace,
			ProjectID:       input.ProjectID,
			Creator:         input.Creator,
		}
		createRepositoryOutput, err := vrm.CreateRepository(createRepositoryInput, c)
		if err != nil {
			// Expected errors
			if e, ok := tkErr.IsError(err); ok {
				switch e.Code() {
				case cCnt.GRPCRepositoryExistErr.Code():
					statusCode = http.StatusBadRequest
					err = tkErr.New(cnt.AdminAPIRepositoryExistErr, input.Name)
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
			err = tkErr.New(cnt.AdminAPIInternalServerErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		input.RepositoryID = createRepositoryOutput.Repository.ID
	}

	createTagInput := &pb.CreateTagInput{
		Tag: &pb.TagInfo{
			Name:            input.Version,
			RepositoryID:    input.RepositoryID,
			ReferenceTarget: input.ImageID,
			Type:            _TagTypeImage,
		},
	}
	createTagOutput, err := vrm.CreateTag(createTagInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCTagExistErr.Code():
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPITagExistErr, input.Version)
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
	output.Tag.ExtractByProto(c, createTagOutput.Tag)
	output.Repository.ExtractByProto(c, createTagOutput.Repository)

	// set acl to global public
	createProjectAclInput := &pb.ProjectAclBatchInfo{
		Data: []*pb.ProjectAclInfo{
			{TagID: createTagOutput.Tag.ID},
		},
	}
	_, err = vrm.CreateProjectAclBatch(createProjectAclInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "vrm.CreateProjectAclBatch(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createProjectAclInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	utility.ResponseWithType(c, statusCode, output)
}
