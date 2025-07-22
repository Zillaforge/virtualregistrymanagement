package admin

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
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

type ListProjectsInput struct {
	pagination
	Namespace *string `json:"namespace" form:"namespace"`
	_         struct{}
}

type ListProjectsOutput struct {
	Projects []ProjectInfo `json:"projects"`
	Total    int           `json:"total"`
	_        struct{}
}

type ProjectInfo struct {
	ID             string `json:"id"`
	UsedSize       int64  `json:"usedSize"`
	UsedCount      int64  `json:"usedCount"`
	SoftLimitSize  int64  `json:"softLimitSize"`
	SoftLimitCount int64  `json:"softLimitCount"`
	_              struct{}
}

func ListProjects(c *gin.Context) {
	var (
		input      = &ListProjectsInput{}
		output     = &ListProjectsOutput{Projects: []ProjectInfo{}}
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

	listProjectsInput := &pb.ListInput{
		Limit:  int32(input.Limit),
		Offset: int32(input.Offset),
	}
	listProjectsOutput, err := vrm.ListProjects(listProjectsInput, c)
	output.Total = int(listProjectsOutput.Count)
	for _, project := range listProjectsOutput.Data {
		where := []string{"project-id=" + project.ID}
		listRepositoriesInput := &pb.ListNamespaceInput{
			Limit:     -1,
			Offset:    0,
			Where:     where,
			Namespace: input.Namespace,
		}
		listRepositoriesOutput, err := vrm.ListRepositories(listRepositoriesInput, c)
		if err != nil {
			zap.L().With(
				zap.String(cnt.Middleware, "vrm.ListRepositories(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", listRepositoriesInput),
			).Error(err.Error())
			statusCode = http.StatusInternalServerError
			err = tkErr.New(cnt.AdminAPIInternalServerErr)
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

		output.Projects = append(output.Projects,
			ProjectInfo{
				ID:             project.ID,
				UsedSize:       currentSize,
				UsedCount:      currentCount,
				SoftLimitSize:  project.LimitSizeBytes,
				SoftLimitCount: project.LimitCount,
			},
		)
	}

	utility.ResponseWithType(c, statusCode, output)
}
