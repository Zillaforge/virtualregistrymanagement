package admin

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/controllers/api"
	"VirtualRegistryManagement/utility"
	"encoding/json"
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

type CreateSnapshotInput struct {
	Version string `json:"version" binding:"required"`
	Extra   []byte `json:"extra"`

	RepositoryID    string `json:"repositoryId" binding:"required_without_all=Name OperatingSystem Description Creator,excluded_with=Name OperatingSystem Description Creator"`
	Name            string `json:"name" binding:"required_without_all=RepositoryID,excluded_with=RepositoryID"`
	OperatingSystem string `json:"operatingSystem" binding:"required_without_all=RepositoryID,excluded_with=RepositoryID"`
	Description     string `json:"description"`
	Creator         string `json:"creator" binding:"required_without_all=RepositoryID,excluded_with=RepositoryID"`

	Namespace string `json:"namespace" binding:"required"`

	VolumeID string `json:"-"`
	_        struct{}
}

type CreateSnapshotOutput struct {
	Repository repositoryInfo `json:"repository"`
	Tag        tagInfo        `json:"tag"`
	_          struct{}
}

// CreateSnapshot 會自動建立 repository & tag，如果有帶 repository-id 會在下面建立
func CreateSnapshot(c *gin.Context) {
	data := c.GetStringMapString(cnt.CtxVPSMetadata)
	extraData, _ := json.Marshal(&data)
	var (
		input = &CreateSnapshotInput{
			VolumeID: c.GetString(cnt.CtxVolumeID),
			Extra:    extraData,
		}
		output     = &CreateSnapshotOutput{}
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
		} else if getRepositoryOutput.Repository.OperatingSystem != c.GetString(cnt.CtxServerOS) {
			err = tkErr.New(cnt.AdminAPIOperatingSystemMismatchErr)
			zap.L().With(
				zap.String(cnt.Controller, "vrm.GetRepository(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", getRepositoryInput),
				zap.String("server os", c.GetString(cnt.CtxServerOS)),
			).Error(err.Error())
			statusCode = http.StatusBadRequest
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		input.Creator = getRepositoryOutput.Repository.Creator
	}

	if input.RepositoryID == "" {
		createRepositoryInput := &pb.RepositoryInfo{
			Name:            input.Name,
			OperatingSystem: input.OperatingSystem,
			Description:     input.Description,
			Namespace:       c.GetString(cnt.CtxNamespace),
			ProjectID:       c.GetString(cnt.CtxProjectID),
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
		output.Repository.ExtractByProto(c, createRepositoryOutput.Repository)
	}

	createTagInput := &pb.CreateTagInput{
		Tag: &pb.TagInfo{
			Name:         input.Name,
			Type:         _TagTypeVolumeSnapshot,
			Extra:        input.Extra,
			RepositoryID: output.Repository.ID,
		},
		Snapshot: &pb.OpenstackSnapshotInfo{
			VolumeID: input.VolumeID,
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

	utility.ResponseWithType(c, statusCode, output)

}
