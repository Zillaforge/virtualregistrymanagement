package common

import (
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility/querydecoder"
)

// Create ...
type (
	CreateExportInput struct {
		Export tables.Export
		_      struct{}
	}

	CreateExportOutput struct {
		Export tables.Export
		_      struct{}
	}
)

// Get ...
type (
	GetExportInput struct {
		ID string
		_  struct{}
	}

	GetExportOutput struct {
		Export tables.Export
		_      struct{}
	}
)

// List ...
type (
	ListExportWhere struct {
		RepositoryID *string `where:"repository-id"`
		TagID        *string `where:"tag-id"`
		Type         *string `where:"type"`
		SnapshotID   *string `where:"snapshot-id"`
		VolumeID     *string `where:"volume-id"`
		ImageID      *string `where:"image-id"`
		Status       *string `where:"status"`
		Creator      *string `where:"creator"`
		ProjectID    *string `where:"project-id"`
		Namespace    *string `where:"namespace"`
		querydecoder.Query
		_ struct{}
	}

	ListExportsInput struct {
		Pagination *Pagination
		Where      ListExportWhere
		_          struct{}
	}

	ListExportsOutput struct {
		Exports []tables.Export
		Count   int64
		_       struct{}
	}
)

// Update ...
type (
	ExportUpdateInfo struct {
		SnapshotStatus *string
		VolumeStatus   *string
		ImageStatus    *string
		Status         *string
		_              struct{}
	}

	UpdateExportInput struct {
		ID         string
		UpdateData *ExportUpdateInfo
		_          struct{}
	}

	UpdateExportOutput struct {
		Export tables.Export
		_      struct{}
	}
)

// Delete ...
type (
	DeleteExportWhere struct {
		ID *string
		querydecoder.Query
		_ struct{}
	}

	DeleteExportInput struct {
		Where DeleteExportWhere
		_     struct{}
	}

	DeleteExportOutput struct {
		ID []string
		_  struct{}
	}
)
