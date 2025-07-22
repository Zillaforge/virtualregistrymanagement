package tasks

import (
	"VirtualRegistryManagement/authentication"
	authCom "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
	"github.com/Zillaforge/virtualregistrymanagementclient/vrm"
)

const UNLIMITED = -1

type project struct {
	ID               string
	Name             string
	CurrentSizeBytes int64
	LimitSizeBytes   int64
	CurrentCount     int64
	LimitCount       int64

	_ struct{}
}

func getProjectQuota(ctx context.Context, id string) (output *project, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"id":     &id,
			"output": &output,
			"error":  &err,
		},
	)

	output = &project{
		ID: id,
	}

	getProjectInput := &pb.GetInput{
		ID: id,
	}
	getProjectOutput, err := vrm.GetProject(getProjectInput, ctx)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Task, "vrm.GetProject(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getProjectOutput),
		).Warn(err.Error())
		return
	}
	output.LimitSizeBytes = getProjectOutput.LimitSizeBytes
	output.LimitCount = getProjectOutput.LimitCount

	where := []string{"project-id=" + id}
	listRepositoriesInput := &pb.ListNamespaceInput{
		Limit:  -1,
		Offset: 0,
		Where:  where,
	}
	listRepositoriesOutput, err := vrm.ListRepositories(listRepositoriesInput, ctx)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Task, "vrm.ListRepositories(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listRepositoriesInput),
		).Warn(err.Error())
		err = tkErr.New(cnt.TaskInternalServerErr, err)
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

	output.CurrentSizeBytes = currentSize
	output.CurrentCount = currentCount

	authProjectInput := &authCom.GetProjectInput{ID: id, Cacheable: true}
	authProjectOutput, err := authentication.Use().GetProject(ctx, authProjectInput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Task, "authentication.Use().GetProject(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", authProjectInput),
		).Warn(err.Error())
		err = tkErr.New(cnt.TaskInternalServerErr, err)
	} else {
		output.Name = authProjectOutput.DisplayName
	}
	return
}
