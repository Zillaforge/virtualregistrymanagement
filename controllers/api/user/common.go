package user

import (
	"VirtualRegistryManagement/authentication"
	authCom "VirtualRegistryManagement/authentication/common"
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"context"
	"encoding/json"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	"pegasus-cloud.com/aes/virtualregistrymanagementclient/pb"
)

type (
	pagination struct {
		Limit  int `json:"limit" form:"limit,default=100" binding:"max=100"`
		Offset int `json:"offset" form:"offset,default=0" binding:"min=0"`
		_      struct{}
	}
)

type userInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Account     string `json:"account"`
}

func (data *userInfo) Fill(ctx context.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
		err       error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"data":  &data,
			"error": &err,
		},
	)

	authUserInput := &authCom.GetUserInput{ID: data.ID, Cacheable: true}
	authUserOutput, err := authentication.Use().GetUser(ctx, authUserInput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "authentication.Use().GetUser(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", authUserInput),
		).Warn(err.Error())
	} else {
		data.DisplayName = authUserOutput.DisplayName
		data.Account = authUserOutput.Account
	}
}

type projectInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

func (data *projectInfo) Fill(ctx context.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
		err       error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"data":  &data,
			"error": &err,
		},
	)

	authProjectInput := &authCom.GetProjectInput{ID: data.ID, Cacheable: true}
	authProjectOutput, err := authentication.Use().GetProject(ctx, authProjectInput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "authentication.Use().GetProject(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", authProjectInput),
		).Warn(err.Error())
	} else {
		data.DisplayName = authProjectOutput.DisplayName
	}
}

type (
	Repository struct {
		ID              string      `json:"id"`
		Name            string      `json:"name"`
		Namespace       string      `json:"namespace"`
		OperatingSystem string      `json:"operatingSystem"`
		Description     string      `json:"description"`
		Tags            []tagInfo   `json:"tags"` // for ui/ux create filter
		Count           int64       `json:"count"`
		Creator         userInfo    `json:"creator"`
		Project         projectInfo `json:"project"`
		CreatedAt       string      `json:"createdAt"`
		UpdatedAt       string      `json:"updatedAt"`
		_               struct{}
	}

	repositoryInfo struct {
		ID              string       `json:"id"`
		Name            string       `json:"name"`
		Namespace       string       `json:"namespace"`
		OperatingSystem string       `json:"operatingSystem"`
		Description     string       `json:"description"`
		Creator         *userInfo    `json:"creator,omitempty"`
		Project         *projectInfo `json:"project,omitempty"`
	}
)

func (data *repositoryInfo) ExtractByProto(ctx context.Context, input *pb.RepositoryInfo) repositoryInfo {
	var (
		funcName = tkUtils.NameOfFunction().String()
		err      error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"data":  &data,
			"error": &err,
		},
	)

	data.ID = input.ID
	data.Name = input.Name
	data.Namespace = input.Namespace
	data.OperatingSystem = input.OperatingSystem
	data.Description = input.Description
	data.Creator = &userInfo{
		ID: input.Creator,
	}
	data.Creator.Fill(ctx)

	data.Project = &projectInfo{
		ID: input.ProjectID,
	}
	data.Project.Fill(ctx)

	return *data
}

func (data *Repository) ExtractByProto(ctx context.Context, input *pb.RepositoryDetail) Repository {
	var (
		funcName = tkUtils.NameOfFunction().String()
		err      error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"data":  &data,
			"error": &err,
		},
	)

	data.ID = input.Repository.ID
	data.Name = input.Repository.Name
	data.Namespace = input.Repository.Namespace
	data.OperatingSystem = input.Repository.OperatingSystem
	data.Description = input.Repository.Description

	data.Tags = []tagInfo{}
	for _, tag := range input.Tags {
		t := tagInfo{}
		data.Tags = append(data.Tags, t.ExtractByProto(ctx, tag))
	}

	data.Count = int64(len(input.Tags))

	data.Creator = userInfo{
		ID: input.Repository.Creator,
	}
	data.Creator.Fill(ctx)
	data.Project = projectInfo{
		ID: input.Repository.ProjectID,
	}
	data.Project.Fill(ctx)

	data.CreatedAt = input.Repository.CreatedAt
	data.UpdatedAt = input.Repository.UpdatedAt

	return *data
}

const (
	_TagTypeImage          = "common"
	_TagTypeVolumeSnapshot = "increase"
)

type (
	Tag struct {
		ID              string                 `json:"id"`
		Name            string                 `json:"name"`
		Repository      repositoryInfo         `json:"repository"`
		ReferenceTarget string                 `json:"referenceTarget"`
		Type            string                 `json:"type"`
		Size            uint64                 `json:"size"`
		Status          string                 `json:"status"`
		Extra           map[string]interface{} `json:"extra"`
		CreatedAt       string                 `json:"createdAt"`
		UpdatedAt       string                 `json:"updatedAt"`
		_               struct{}
	}
	tagInfo struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
)

func (data *tagInfo) ExtractByProto(ctx context.Context, input *pb.TagInfo) tagInfo {
	var (
		funcName = tkUtils.NameOfFunction().String()
		err      error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"data":  &data,
			"error": &err,
		},
	)

	extra := map[string]interface{}{}
	if input.Extra != nil {
		json.Unmarshal(input.Extra, &extra)
	}

	data.ID = input.ID
	data.Name = input.Name
	data.Type = input.Type

	return *data
}

func (data *Tag) ExtractByProto(ctx context.Context, input *pb.TagDetail) Tag {
	var (
		funcName = tkUtils.NameOfFunction().String()
		err      error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"data":  &data,
			"error": &err,
		},
	)

	extra := map[string]interface{}{}
	if input.Tag.Extra != nil {
		json.Unmarshal(input.Tag.Extra, &extra)
	}

	data.ID = input.Tag.ID
	data.Name = input.Tag.Name

	data.Repository = repositoryInfo{
		ID:              input.Repository.ID,
		Name:            input.Repository.Name,
		Namespace:       input.Repository.Namespace,
		OperatingSystem: input.Repository.OperatingSystem,
		Description:     input.Repository.Description,
		Creator: &userInfo{
			ID: input.Repository.Creator,
		},
		Project: &projectInfo{
			ID: input.Repository.ProjectID,
		},
	}
	data.Repository.Creator.Fill(ctx)
	data.Repository.Project.Fill(ctx)

	data.ReferenceTarget = input.Tag.ReferenceTarget
	data.Type = input.Tag.Type
	data.Size = input.Tag.Size
	data.Status = input.Tag.Status
	data.Extra = extra
	data.CreatedAt = input.Tag.CreatedAt
	data.UpdatedAt = input.Tag.UpdatedAt

	return *data
}

type MemberAcl struct {
	ID   string   `json:"id"`
	Tag  tagInfo  `json:"tag"`
	User userInfo `json:"user"`
	_    struct{}
}

func (data *MemberAcl) ExtractByProto(ctx context.Context, input *pb.MemberAclDetail) MemberAcl {
	var (
		funcName = tkUtils.NameOfFunction().String()
		err      error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"data":  &data,
			"error": &err,
		},
	)

	data.ID = input.ID

	if input.Tag != nil {
		data.Tag = tagInfo{
			ID:   input.Tag.ID,
			Name: input.Tag.Name,
			Type: input.Tag.Type,
		}
	}
	data.User = userInfo{
		ID: input.UserID,
	}
	data.User.Fill(ctx)

	return *data
}

type ProjectAcl struct {
	ID      string      `json:"id"`
	Tag     tagInfo     `json:"tag"`
	Project projectInfo `json:"project"`
	_       struct{}
}

func (data *ProjectAcl) ExtractByProto(ctx context.Context, input *pb.ProjectAclDetail) ProjectAcl {
	var (
		funcName = tkUtils.NameOfFunction().String()
		err      error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"data":  &data,
			"error": &err,
		},
	)

	data.ID = input.ID
	if input.Tag != nil {
		data.Tag = tagInfo{
			ID:   input.Tag.ID,
			Name: input.Tag.Name,
			Type: input.Tag.Type,
		}
	}
	if input.ProjectID != nil {
		data.Project = projectInfo{
			ID: *input.ProjectID,
		}
		data.Project.Fill(ctx)
	}

	return *data
}

type Export struct {
	ID             string `json:"id"`
	RepositoryID   string `json:"repositoryId"`
	RepositoryName string `json:"repositoryName"`
	TagID          string `json:"tagId"`
	TagName        string `json:"tagName"`
	Type           string `json:"type"`
	Filepath       string `json:"filepath"`
	Status         string `json:"status"`
	Creator        string `json:"creator"`
	ProjectID      string `json:"projectId"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

func (data *Export) ExtractByProto(ctx context.Context, input *pb.ExportInfo) Export {
	var (
		funcName = tkUtils.NameOfFunction().String()
		err      error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"data":  &data,
			"error": &err,
		},
	)

	data.ID = input.ID
	data.RepositoryID = input.RepositoryID
	data.RepositoryName = input.RepositoryName
	data.TagID = input.TagID
	data.TagName = input.TagName
	data.Type = input.Type
	data.Filepath = input.Filepath
	data.Status = input.Status
	data.Creator = input.Creator
	data.ProjectID = input.ProjectID
	data.CreatedAt = input.CreatedAt
	data.UpdatedAt = input.UpdatedAt

	return *data
}
