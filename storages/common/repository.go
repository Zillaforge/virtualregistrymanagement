package common

import (
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility/querydecoder"
)

// Create ...
type (
	CreateRepositoryInput struct {
		Repository tables.Repository
		_          struct{}
	}

	CreateRepositoryOutput struct {
		Repository tables.Repository
		_          struct{}
	}
)

// Get ...
type (
	GetRepositoryInput struct {
		ID string
		_  struct{}
	}

	GetRepositoryOutput struct {
		Repository tables.Repository
		_          struct{}
	}
)

// List ...
type (
	ListRepositoryWhere struct {
		Namespace       *string  `where:"-"`
		OperatingSystem *string  `where:"os"`
		Creator         *string  `where:"creator"`
		ProjectID       *string  `where:"project-id"`
		ID              []string `where:"id"`
		querydecoder.Query
		_ struct{}
	}

	ListRepositoriesInput struct {
		Pagination *Pagination
		Where      ListRepositoryWhere
		_          struct{}
	}

	ListRepositoriesOutput struct {
		Repositories []tables.Repository
		Count        int64
		_            struct{}
	}
)

// Update ...
type (
	RepositoryUpdateInfo struct {
		Name        *string
		Description *string
		_           struct{}
	}

	UpdateRepositoryInput struct {
		ID         string
		UpdateData *RepositoryUpdateInfo
		_          struct{}
	}

	UpdateRepositoryOutput struct {
		Repository tables.Repository
		_          struct{}
	}
)

// Delete ...
type (
	DeleteRepositoryWhere struct {
		ID        *string
		Creator   *string
		ProjectID *string
		querydecoder.Query
		_ struct{}
	}

	DeleteRepositoryInput struct {
		Where DeleteRepositoryWhere
		_     struct{}
	}

	DeleteRepositoryOutput struct {
		ID []string
		_  struct{}
	}
)
