package common

import (
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility/querydecoder"
)

// List ...
type (
	ListRegistryWhere struct {
		Namespace       *string `where:"-"`
		RepositoryID    *string `where:"repository-id"`
		TagID           *string `where:"tag-id"`
		Creator         *string
		ProjectID       *string `where:"project-id"`
		OperatingSystem *string
		Type            *string
		Status          *string
		MemberAclID     *string `where:"member-acl-id"`
		AllowUserID     *string `where:"allow-user-id"`
		ProjectAclID    *string
		AllowProjectID  *string
		querydecoder.Query
		_ struct{}
	}

	ListRegistriesInput struct {
		Pagination *Pagination
		Where      ListRegistryWhere
		Flag       Flag
		_          struct{}
	}

	ListRegistriesOutput struct {
		Registries []tables.Registry
		Count      int64
		_          struct{}
	}

	Flag struct { // Or Operator
		UserID    *string // query using
		ProjectID *string // query using

		BelongUser    bool // UserID required
		BelongProject bool // ProjectID required
		ProjectLimit  bool // hare to Specific User, *UserID required*
		ProjectPublic bool // Public in Project, *ProjectID required*
		GlobalLimit   bool // Share to Specific Project, *ProjectID required*
		GlobalPublic  bool // Public in Namespace
	}
)
