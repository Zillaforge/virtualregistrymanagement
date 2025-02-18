package user

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
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	cCnt "pegasus-cloud.com/aes/virtualregistrymanagementclient/constants"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/vrm"
)

type UploadImageInput struct {
	Filepath string `json:"filepath" binding:"required"`

	TagID           string `json:"tagId" binding:"required_without_all=RepositoryID Name OperatingSystem Description Version Type,excluded_with=RepositoryID Name OperatingSystem Description Version Type"`
	Version         string `json:"version" binding:"required_without=TagID,excluded_with=TagID"`
	DiskFormat      string `json:"diskFormat" binding:"required_without=TagID,oneof='' ami ari aki vhd vmdk raw qcow2 vdi iso"`
	ContainerFormat string `json:"containerFormat" binding:"required_without=TagID,oneof='' ami ari aki bare ovf"`
	Extra           []byte `json:"extra" binding:"excluded_with=TagID"`

	RepositoryID    string `json:"repositoryId" binding:"required_without_all=Name OperatingSystem Description TagID,excluded_with=Name OperatingSystem Description TagID"`
	Name            string `json:"name" binding:"required_without_all=RepositoryID TagID,excluded_with=RepositoryID TagID"`
	OperatingSystem string `json:"operatingSystem" binding:"required_without_all=RepositoryID TagID,excluded_with=RepositoryID TagID,oneof='' linux windows"`
	Description     string `json:"description"`

	Namespace string `json:"-"`
	ProjectID string `json:"-"`
	Creator   string `json:"-"`
	_         struct{}
}

type UploadImageOutput struct {
	Repository repositoryInfo `json:"repository"`
	Tag        tagInfo        `json:"tag"`
	_          struct{}
}

// UploadImage 會自動建立 repository & tag，如果有帶 repository-id 或 tag 會在下面建立
func UploadImage(c *gin.Context) {
	var (
		input = &UploadImageInput{
			Namespace: c.GetString(cnt.CtxNamespace),
			ProjectID: c.GetString(cnt.CtxProjectID),
		}
		output       = &UploadImageOutput{}
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

	// check project limit
	if c.GetBool(cnt.CtxProjectCountFlag) || c.GetBool(cnt.CtxProjectSizeFlag) {
		err = tkErr.New(cnt.UserAPIProjectLimitHasBeenReachedErr)
		zap.L().With(
			zap.String(cnt.Controller, "project-count or project-size limit has been reached"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("obj", input),
		).Error(err.Error())
		statusCode = http.StatusForbidden
		utility.ResponseWithType(c, statusCode, err)
		return
	}

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
					err = tkErr.New(cnt.UserAPITagNotFoundErr, input.TagID)
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
			err = tkErr.New(cnt.UserAPIInternalServerErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
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
				case cCnt.GRPCRepositoryExistErr.Code():
					statusCode = http.StatusBadRequest
					err = tkErr.New(cnt.UserAPIRepositoryNotFoundErr, input.RepositoryID)
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
			err = tkErr.New(cnt.UserAPIInternalServerErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		input.Creator = getRepositoryOutput.Repository.Creator
	default:
		input.Creator = c.GetString(cnt.CtxUserID)
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
			err = tkErr.New(cnt.UserAPIFileNotFoundErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		statusCode = http.StatusBadRequest
		err = tkErr.New(cnt.UserAPIFileNotFoundErr).WithInner(err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// check permission, only creator and tenant-owner, tenant-admin allow upload to exist tag
	if role := c.GetString(cnt.CtxTenantRole); !supportRoles[role] &&
		input.Creator != c.GetString(cnt.CtxUserID) {
		zap.L().With(
			zap.String(cnt.Controller, "vrm.GetRepository(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.String("role", role),
			zap.String("token-user", c.GetString(cnt.CtxUserID)),
			zap.String("creator", input.Creator),
		).Info(cnt.UserAPIUnauthorizedOpErrMsg)
		statusCode = http.StatusUnauthorized
		err = tkErr.New(cnt.UserAPIUnauthorizedOpErr)
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
					err = tkErr.New(cnt.UserAPITagExistErr, input.Version)
				case cCnt.GRPCExceedAllowedQuotaErr.Code():
					statusCode = http.StatusRequestEntityTooLarge
					err = tkErr.New(cnt.UserAPIExceedAllowedQuotaErr)
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
			err = tkErr.New(cnt.UserAPIInternalServerErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		output.Tag.ExtractByProto(c, createTagOutput.Tag)
		output.Repository.ExtractByProto(c, createTagOutput.Repository)

	case input.RepositoryID == "" && input.TagID == "":
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
					err = tkErr.New(cnt.UserAPITagExistErr, input.Version)
				case cCnt.GRPCExceedAllowedQuotaErr.Code():
					statusCode = http.StatusRequestEntityTooLarge
					err = tkErr.New(cnt.UserAPIExceedAllowedQuotaErr)
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
			err = tkErr.New(cnt.UserAPIInternalServerErr)
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
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCNotSupportTypeErr.Code():
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.UserAPINotSupportTypeErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCFileNotFoundErr.Code():
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.UserAPIFileNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Controller, "vrm.UploadTag(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", uploadTagInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	utility.ResponseWithType(c, statusCode, output)
}
