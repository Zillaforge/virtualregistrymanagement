package admin

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/api"
	"VirtualRegistryManagement/modules/filepath"
	"VirtualRegistryManagement/utility"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

type UploadImageInput struct {
	Namespace string `json:"namespace" binding:"required"`
	Filepath  string `json:"filepath" binding:"required"`

	TagID           string `json:"tagId" binding:"required_without_all=RepositoryID Name OperatingSystem Description ProjectID Creator Version Type DiskFormat ContainerFormat,excluded_with=RepositoryID Name OperatingSystem Description ProjectID Creator Version Type DiskFormat ContainerFormat"`
	Version         string `json:"version" binding:"required_without=TagID,excluded_with=TagID"`
	DiskFormat      string `json:"diskFormat" binding:"required_without=TagID,excluded_with=TagID,oneof=ami ari aki vhd vmdk raw qcow2 vdi iso"`
	ContainerFormat string `json:"containerFormat" binding:"required_without=TagID,excluded_with=TagID,oneof=ami ari aki bare ovf"`
	Extra           []byte `json:"extra" binding:"excluded_with=TagID"`

	RepositoryID    string `json:"repositoryId" binding:"required_without_all=Name OperatingSystem Description ProjectID Creator TagID,excluded_with=Name OperatingSystem Description ProjectID Creator TagID"`
	Name            string `json:"name" binding:"required_without_all=RepositoryID TagID,excluded_with=RepositoryID TagID"`
	OperatingSystem string `json:"operatingSystem" binding:"required_without_all=RepositoryID TagID,excluded_with=RepositoryID TagID,oneof='' linux windows"`
	Description     string `json:"description"`
	ProjectID       string `json:"projectId" binding:"required_without_all=RepositoryID TagID,excluded_with=RepositoryID TagID"`
	Creator         string `json:"creator" binding:"required_without_all=RepositoryID TagID,excluded_with=RepositoryID TagID"`

	_ struct{}
}

type UploadImageOutput struct {
	Repository repositoryInfo `json:"repository"`
	Tag        tagInfo        `json:"tag"`
	_          struct{}
}

// UploadImage 會自動建立 repository & tag，如果有帶 repository-id 或 tag 會在下面建立
func UploadImage(c *gin.Context) {
	var (
		input      = &UploadImageInput{}
		output     = &UploadImageOutput{}
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

	if !namespaceIsLegal(input.Namespace) {
		statusCode = http.StatusBadRequest
		err = tkErr.New(cnt.AdminAPINamespaceNotFoundErr)
		utility.ResponseWithType(c, statusCode, err)
	}

	// get creator
	switch {
	case input.TagID != "":
		getTagInput := &pb.GetInput{
			ID: input.TagID,
		}
		getTagOutput, err := vrm.GetTag(getTagInput, c)
		if err != nil {
			// Expected errors
			if e, ok := tkErr.IsError(err); ok {
				switch e.Code() {
				case cCnt.GRPCTagNotFoundErr.Code():
					statusCode = http.StatusBadRequest
					err = tkErr.New(cnt.AdminAPITagNotFoundErr, input.TagID)
					utility.ResponseWithType(c, statusCode, err)
					return
				}
			}
			// Unexpected errors
			zap.L().With(
				zap.String(cnt.Controller, "vrm.GetTag(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", getTagInput),
			).Error(err.Error())
			statusCode = http.StatusInternalServerError
			err = tkErr.New(cnt.AdminAPIInternalServerErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		input.ProjectID = getTagOutput.Repository.ProjectID
		input.Creator = getTagOutput.Repository.Creator

		output.Tag.ExtractByProto(c, getTagOutput.Tag)
		output.Repository.ExtractByProto(c, getTagOutput.Repository)

	case input.RepositoryID != "":
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
		input.ProjectID = getRepositoryOutput.Repository.ProjectID
		input.Creator = getRepositoryOutput.Repository.Creator
	}

	path, err := filepath.Path(input.Filepath, []interface{}{input.ProjectID}, nil)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "filepath.Path()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("project-id", input.ProjectID),
			zap.String("filepath", input.Filepath),
		).Info(err.Error())
		statusCode = http.StatusInternalServerError
		err = api.Malformed(err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}
	// check file exist
	if _, err = os.Stat(path); err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "os.Stat()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("filepath", path),
		).Info(err.Error())
		if errors.Is(err, os.ErrNotExist) {
			statusCode = http.StatusBadRequest
			err = tkErr.New(cnt.AdminAPIFileNotFoundErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		statusCode = http.StatusBadRequest
		err = tkErr.New(cnt.AdminAPIFileNotFoundErr).WithInner(err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	switch {
	case input.RepositoryID != "" && input.TagID == "":
		createTagInput := &pb.CreateTagInput{
			Tag: &pb.TagInfo{
				RepositoryID: input.RepositoryID,
				Name:         input.Version,
				Type:         _TagTypeImage,
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

	case input.RepositoryID == "" && input.TagID == "":
		if !membershipIsLegal(input.ProjectID, input.Creator) {
			statusCode = http.StatusBadRequest
			err = tkErr.New(cnt.AdminAPIMembershipNotFoundErr)
			utility.ResponseWithType(c, statusCode, err)
		}
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
		createTagInput := &pb.CreateTagInput{
			Tag: &pb.TagInfo{
				RepositoryID: createRepositoryOutput.Repository.ID,
				Name:         input.Version,
				Type:         _TagTypeImage,
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
	}

	uploadTagInput := &pb.ImageDataInput{
		TagID:    output.Tag.ID,
		Filepath: input.Filepath,
	}
	err = vrm.UploadTag(uploadTagInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCNotSupportTypeErr.Code():
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPINotSupportTypeErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCFileNotFoundErr.Code():
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIFileNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return

			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "vrm.UploadTag(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", uploadTagInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	utility.ResponseWithType(c, statusCode, output)
}
