package common

import (
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility/querydecoder"
)

// Create ...
type (
	CreateProjectInput struct {
		Project tables.Project
		_       struct{}
	}

	CreateProjectOutput struct {
		Project tables.Project
		_       struct{}
	}
)

// Get ...
type (
	GetProjectInput struct {
		ID string
		_  struct{}
	}

	GetProjectOutput struct {
		Project tables.Project
		_       struct{}
	}
)

// List ...
type (
	ListProjectWhere struct {
		querydecoder.Query
		_ struct{}
	}

	ListProjectsInput struct {
		Pagination *Pagination
		Where      ListProjectWhere
		_          struct{}
	}

	ListProjectsOutput struct {
		Projects []tables.Project
		Count    int64
		_        struct{}
	}
)

// Update ...
type (
	ProjectUpdateInfo struct {
		LimitCount     *int64
		LimitSizeBytes *int64
		_              struct{}
	}

	UpdateProjectInput struct {
		ID         string
		UpdateData *ProjectUpdateInfo
		_          struct{}
	}

	UpdateProjectOutput struct {
		Project tables.Project
		_       struct{}
	}
)

// Delete ...
type (
	DeleteProjectWhere struct {
		ID *string
		querydecoder.Query
		_ struct{}
	}

	DeleteProjectInput struct {
		Where DeleteProjectWhere
		_     struct{}
	}

	DeleteProjectOutput struct {
		ID []string
		_  struct{}
	}
)
