package common

import (
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility/querydecoder"
)

// Create ...
type (
	CreateTagInput struct {
		Tag tables.Tag
		_   struct{}
	}

	CreateTagOutput struct {
		Tag tables.Tag
		_   struct{}
	}
)

// Get ...
type (
	GetTagInput struct {
		ID string
		_  struct{}
	}

	GetTagOutput struct {
		Tag tables.Tag
		_   struct{}
	}
)

// List ...
type (
	ListTagWhere struct {
		Namespace       *string  `where:"-" prefix:"Repository"`
		ProjectID       *string  `where:"project-id" prefix:"Repository"`
		RepositoryID    *string  `where:"repository-id" prefix:"tag"`
		Type            *string  `where:"type" prefix:"tag"`
		Status          *string  `where:"status" prefix:"tag"`
		ReferenceTarget *string  `where:"-" prefix:"tag"`
		ID              []string `where:"id" prefix:"tag"`
		querydecoder.Query
		_ struct{}
	}

	ListTagsInput struct {
		Pagination *Pagination
		Where      ListTagWhere
		_          struct{}
	}

	ListTagsOutput struct {
		Tags  []tables.Tag
		Count int64
		_     struct{}
	}
)

// Update ...
type (
	TagUpdateInfo struct {
		Name            *string
		ReferenceTarget *string
		Size            *uint64
		Status          *string
		Extra           *[]byte
		_               struct{}
	}

	UpdateTagInput struct {
		ID         string
		UpdateData *TagUpdateInfo
		_          struct{}
	}

	UpdateTagOutput struct {
		Tag tables.Tag
		_   struct{}
	}
)

// Delete ...
type (
	DeleteTagWhere struct {
		ID           *string
		RepositoryID *string
		querydecoder.Query
		_ struct{}
	}

	DeleteTagInput struct {
		Where DeleteTagWhere
		_     struct{}
	}

	DeleteTagOutput struct {
		ID []string
		_  struct{}
	}
)
