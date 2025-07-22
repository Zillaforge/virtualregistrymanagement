package project

import (
	"VirtualRegistryManagement/authentication"
	authCom "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/mviper"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	cCnt "github.com/Zillaforge/virtualregistrymanagementclient/constants"
	"github.com/Zillaforge/virtualregistrymanagementclient/pb"
)

const (
	_do_not_do_anything = "don't do anything"
	_create             = "create"
	_delete             = "delete"
)

func (m *Method) SyncProjects(ctx context.Context, input *emptypb.Empty) (output *emptypb.Empty, err error) {
	var (
		funcName   = tkUtils.NameOfFunction().String()
		requestID  = utility.MustGetContextRequestID(ctx)
		projectMap = make(map[string]string)
	)
	output = empty

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input":  &projectMap,
			"output": &output,
			"error":  &err,
		},
	)

	listInput := &pb.ListInput{
		Limit:  -1,
		Offset: 0,
	}
	listOutput, err := m.ListProjects(ctx, listInput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.EventConsume, "vrm.ListProjects(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listInput),
		).Error(err.Error())
		return
	}

	for _, project := range listOutput.Data {
		projectMap[project.ID] = _delete
	}

	// get all projects from IAM and put them into projectMap with key
	listProjectsFromIAMInput := &authCom.ListProjectsInput{
		Limit:  -1,
		Offset: 0,
	}
	listProjectsFromIAMOutput, err := authentication.Use().ListProjects(ctx, listProjectsFromIAMInput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Task, "authentication.Use().ListProjects(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listProjectsFromIAMInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	for _, project := range listProjectsFromIAMOutput.Projects {
		if _, ok := projectMap[project.ID]; ok {
			projectMap[project.ID] = _do_not_do_anything
		} else {
			projectMap[project.ID] = _create
		}
	}

	for pid, action := range projectMap {
		switch action {
		case _create:
			createInput := &pb.ProjectInfo{
				ID:             pid,
				LimitCount:     mviper.GetInt64("VirtualRegistryManagement.scopes.project_default_count"),
				LimitSizeBytes: mviper.GetInt64("VirtualRegistryManagement.scopes.project_default_size"),
			}
			if _, err = m.CreateProject(ctx, createInput); err != nil {
				zap.L().With(
					zap.String(cnt.EventConsume, "m.CreateProject(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.Any("input", createInput),
				).Error(err.Error())
				continue
			}
		case _delete:
			deleteInput := &pb.DeleteInput{
				Where: []string{"ID=" + pid},
			}
			if _, err = m.DeleteProject(ctx, deleteInput); err != nil {
				zap.L().With(
					zap.String(cnt.EventConsume, "m.DeleteProject(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.Any("input", deleteInput),
				).Error(err.Error())
				continue
			}
		case _do_not_do_anything:
			continue
		}
	}

	return
}
