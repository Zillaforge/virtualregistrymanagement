package common

import (
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility/querydecoder"
)

// Create ...
type (
	CreateProjectAclBatchInput struct {
		ProjectAcls []tables.ProjectAcl
		_           struct{}
	}

	CreateProjectAclBatchOutput struct {
		ProjectAcls []tables.ProjectAcl
		Count       int64
		_           struct{}
	}
)

// Get ...
type (
	GetProjectAclInput struct {
		ID string
		_  struct{}
	}

	GetProjectAclOutput struct {
		ProjectAcl tables.ProjectAcl
		_          struct{}
	}
)

// List ...
type (
	ListProjectAclWhere struct {
		Namespace *string `where:"-"`
		TagID     *string `where:"tag-id"`
		ProjectID *string `where:"project-id"`
		querydecoder.Query
		_ struct{}
	}

	ListProjectAclsInput struct {
		Pagination *Pagination
		Where      ListProjectAclWhere
		_          struct{}
	}

	ListProjectAclsOutput struct {
		ProjectAcls []tables.Registry
		Count       int64
		_           struct{}
	}
)

// Delete ...
type (
	DeleteProjectAclWhere struct {
		ID        *string `where:"id"`
		TagID     *string `where:"tag-id"`
		ProjectID *string `where:"project-id"`
		querydecoder.Query
		_ struct{}
	}

	DeleteProjectAclInput struct {
		Where DeleteProjectAclWhere
		_     struct{}
	}

	DeleteProjectAclOutput struct {
		TagID []string
		_     struct{}
	}
)
