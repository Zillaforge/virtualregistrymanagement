package common

import (
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility/querydecoder"
)

// Create ...
type (
	CreateMemberAclBatchInput struct {
		MemberAcls []tables.MemberAcl
		_          struct{}
	}

	CreateMemberAclBatchOutput struct {
		MemberAcls []tables.MemberAcl
		Count      int64
		_          struct{}
	}
)

// Get ...
type (
	GetMemberAclInput struct {
		ID string
		_  struct{}
	}

	GetMemberAclOutput struct {
		MemberAcl tables.MemberAcl
		_         struct{}
	}
)

// List ...
type (
	ListMemberAclWhere struct {
		Namespace *string `where:"-"`
		TagID     *string `where:"tag-id"`
		UserID    *string `where:"user-id"`
		querydecoder.Query
		_ struct{}
	}

	ListMemberAclsInput struct {
		Pagination *Pagination
		Where      ListMemberAclWhere
		_          struct{}
	}

	ListMemberAclsOutput struct {
		MemberAcls []tables.Registry
		Count      int64
		_          struct{}
	}
)

// Delete ...
type (
	DeleteMemberAclWhere struct {
		ID     *string
		TagID  *string
		UserID *string
		querydecoder.Query
		_ struct{}
	}

	DeleteMemberAclInput struct {
		Where DeleteMemberAclWhere
		_     struct{}
	}

	DeleteMemberAclOutput struct {
		ID []string
		_  struct{}
	}
)
