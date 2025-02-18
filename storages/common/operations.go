package common

import "context"

// Operations ...
type Operations interface {
	ProjectCRUDInterface
	RepositoryCRUDInterface
	TagCRUDInterface
	MemberAclCRUDInterface
	ProjectAclCRUDInterface
	RegistryCRUDInterface
	ExportCRUDInterface
}

type ProjectCRUDInterface interface {
	ListProjects(ctx context.Context, input *ListProjectsInput) (output *ListProjectsOutput, err error)
	CreateProject(ctx context.Context, input *CreateProjectInput) (output *CreateProjectOutput, err error)
	GetProject(ctx context.Context, input *GetProjectInput) (output *GetProjectOutput, err error)
	UpdateProject(ctx context.Context, input *UpdateProjectInput) (output *UpdateProjectOutput, err error)
	DeleteProject(ctx context.Context, input *DeleteProjectInput) (output *DeleteProjectOutput, err error)
}

type RepositoryCRUDInterface interface {
	ListRepositories(ctx context.Context, input *ListRepositoriesInput) (output *ListRepositoriesOutput, err error)
	CreateRepository(ctx context.Context, input *CreateRepositoryInput) (output *CreateRepositoryOutput, err error)
	GetRepository(ctx context.Context, input *GetRepositoryInput) (output *GetRepositoryOutput, err error)
	UpdateRepository(ctx context.Context, input *UpdateRepositoryInput) (output *UpdateRepositoryOutput, err error)
	DeleteRepository(ctx context.Context, input *DeleteRepositoryInput) (output *DeleteRepositoryOutput, err error)
}

type TagCRUDInterface interface {
	ListTags(ctx context.Context, input *ListTagsInput) (output *ListTagsOutput, err error)
	CreateTag(ctx context.Context, input *CreateTagInput) (output *CreateTagOutput, err error)
	GetTag(ctx context.Context, input *GetTagInput) (output *GetTagOutput, err error)
	UpdateTag(ctx context.Context, input *UpdateTagInput) (output *UpdateTagOutput, err error)
	DeleteTag(ctx context.Context, input *DeleteTagInput) (output *DeleteTagOutput, err error)
}

type MemberAclCRUDInterface interface {
	ListMemberAcls(ctx context.Context, input *ListMemberAclsInput) (output *ListMemberAclsOutput, err error)
	CreateMemberAclBatch(ctx context.Context, input *CreateMemberAclBatchInput) (output *CreateMemberAclBatchOutput, err error)
	GetMemberAcl(ctx context.Context, input *GetMemberAclInput) (output *GetMemberAclOutput, err error)
	DeleteMemberAcl(ctx context.Context, input *DeleteMemberAclInput) (output *DeleteMemberAclOutput, err error)
}

type ProjectAclCRUDInterface interface {
	ListProjectAcls(ctx context.Context, input *ListProjectAclsInput) (output *ListProjectAclsOutput, err error)
	CreateProjectAclBatch(ctx context.Context, input *CreateProjectAclBatchInput) (output *CreateProjectAclBatchOutput, err error)
	GetProjectAcl(ctx context.Context, input *GetProjectAclInput) (output *GetProjectAclOutput, err error)
	DeleteProjectAcl(ctx context.Context, input *DeleteProjectAclInput) (output *DeleteProjectAclOutput, err error)
}

type RegistryCRUDInterface interface {
	ListRegistries(ctx context.Context, input *ListRegistriesInput) (output *ListRegistriesOutput, err error)
}

type ExportCRUDInterface interface {
	ListExports(ctx context.Context, input *ListExportsInput) (output *ListExportsOutput, err error)
	CreateExport(ctx context.Context, input *CreateExportInput) (output *CreateExportOutput, err error)
	GetExport(ctx context.Context, input *GetExportInput) (output *GetExportOutput, err error)
	UpdateExport(ctx context.Context, input *UpdateExportInput) (output *UpdateExportOutput, err error)
	DeleteExport(ctx context.Context, input *DeleteExportInput) (output *DeleteExportOutput, err error)
}
