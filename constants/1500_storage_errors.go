package constants

import tkErr "github.com/Zillaforge/toolkits/errors"

const (
	// 1500xxxx: Storage
	StorageInternalServerErrCode         = 15000000
	StorageInternalServerErrMsg          = "internal server error"
	StorageOneOfResourcesNotFoundErrCode = 15000001
	StorageOneOfResourcesNotFoundErrMsg  = "one of resource not found"
	StorageProjectExistErrCode           = 15000002
	StorageProjectExistErrMsg            = "project exist"
	StorageProjectNotFoundErrCode        = 15000003
	StorageProjectNotFoundErrMsg         = "project not found"
	StorageProjectInUseErrCode           = 15000004
	StorageProjectInUseErrMsg            = "project in use"
	StorageRepositoryExistErrCode        = 15000005
	StorageRepositoryExistErrMsg         = "repository exist"
	StorageRepositoryNotFoundErrCode     = 15000006
	StorageRepositoryNotFoundErrMsg      = "repository not found"
	StorageRepositoryInUseErrCode        = 15000007
	StorageRepositoryInUseErrMsg         = "repository in use"
	StorageTagExistErrCode               = 15000008
	StorageTagExistErrMsg                = "tag exist"
	StorageTagNotFoundErrCode            = 15000009
	StorageTagNotFoundErrMsg             = "tag not found"
	StorageTagInUseErrCode               = 15000010
	StorageTagInUseErrMsg                = "tag in use"
	StorageMemberAclExistErrCode         = 15000011
	StorageMemberAclExistErrMsg          = "member-acl exist"
	StorageMemberAclNotFoundErrCode      = 15000012
	StorageMemberAclNotFoundErrMsg       = "member-acl not found"
	StorageMemberAclInUseErrCode         = 15000013
	StorageMemberAclInUseErrMsg          = "member-acl in use"
	StorageProjectAclExistErrCode        = 15000014
	StorageProjectAclExistErrMsg         = "project-acl exist"
	StorageProjectAclNotFoundErrCode     = 15000015
	StorageProjectAclNotFoundErrMsg      = "project-acl not found"
	StorageProjectAclInUseErrCode        = 15000016
	StorageProjectAclInUseErrMsg         = "project-acl in use"
	StorageRegistryNotFoundErrCode       = 15000017
	StorageRegistryNotFoundErrMsg        = "registry not found"
	StorageExportExistErrCode            = 15000018
	StorageExportExistErrMsg             = "export exist"
	StorageExportNotFoundErrCode         = 15000019
	StorageExportNotFoundErrMsg          = "export not found"
	StorageExportInUseErrCode            = 15000020
	StorageExportInUseErrMsg             = "export in use"
)

var (
	// 15000000 (internal server error)
	StorageInternalServerErr = tkErr.Error(StorageInternalServerErrCode, StorageInternalServerErrMsg)
	// 15000001 (one of resource not found)
	StorageOneOfResourcesNotFoundErr = tkErr.Error(StorageOneOfResourcesNotFoundErrCode, StorageOneOfResourcesNotFoundErrMsg)
	// 15000002 (project exist)
	StorageProjectExistErr = tkErr.Error(StorageProjectExistErrCode, StorageProjectExistErrMsg)
	// 15000003 (project not found)
	StorageProjectNotFoundErr = tkErr.Error(StorageProjectNotFoundErrCode, StorageProjectNotFoundErrMsg)
	// 15000004 (project in use)
	StorageProjectInUseErr = tkErr.Error(StorageProjectInUseErrCode, StorageProjectInUseErrMsg)
	// 15000005 (repository exist)
	StorageRepositoryExistErr = tkErr.Error(StorageRepositoryExistErrCode, StorageRepositoryExistErrMsg)
	// 15000006 (repository not found)
	StorageRepositoryNotFoundErr = tkErr.Error(StorageRepositoryNotFoundErrCode, StorageRepositoryNotFoundErrMsg)
	// 15000007 (repository in use)
	StorageRepositoryInUseErr = tkErr.Error(StorageRepositoryInUseErrCode, StorageRepositoryInUseErrMsg)
	// 15000008 (tag exist)
	StorageTagExistErr = tkErr.Error(StorageTagExistErrCode, StorageTagExistErrMsg)
	// 15000009 (tag not found)
	StorageTagNotFoundErr = tkErr.Error(StorageTagNotFoundErrCode, StorageTagNotFoundErrMsg)
	// 15000010 (tag in use)
	StorageTagInUseErr = tkErr.Error(StorageTagInUseErrCode, StorageTagInUseErrMsg)
	// 15000011 (member-acl exist)
	StorageMemberAclExistErr = tkErr.Error(StorageMemberAclExistErrCode, StorageMemberAclExistErrMsg)
	// 15000012 (member-acl not found)
	StorageMemberAclNotFoundErr = tkErr.Error(StorageMemberAclNotFoundErrCode, StorageMemberAclNotFoundErrMsg)
	// 15000013 (member-acl in use)
	StorageMemberAclInUseErr = tkErr.Error(StorageMemberAclInUseErrCode, StorageMemberAclInUseErrMsg)
	// 15000014 (project-acl exist)
	StorageProjectAclExistErr = tkErr.Error(StorageProjectAclExistErrCode, StorageProjectAclExistErrMsg)
	// 15000015 (project-acl not found)
	StorageProjectAclNotFoundErr = tkErr.Error(StorageProjectAclNotFoundErrCode, StorageProjectAclNotFoundErrMsg)
	// 15000016 (project-acl in use)
	StorageProjectAclInUseErr = tkErr.Error(StorageProjectAclInUseErrCode, StorageProjectAclInUseErrMsg)
	// 15000017 (registry not found)
	StorageRegistryNotFoundErr = tkErr.Error(StorageRegistryNotFoundErrCode, StorageRegistryNotFoundErrMsg)
	// 15000018 (export exist)
	StorageExportExistErr = tkErr.Error(StorageExportExistErrCode, StorageExportExistErrMsg)
	// 15000019 (export not found)
	StorageExportNotFoundErr = tkErr.Error(StorageExportNotFoundErrCode, StorageExportNotFoundErrMsg)
	// 15000020 (export in use)
	StorageExportInUseErr = tkErr.Error(StorageExportInUseErrCode, StorageExportInUseErrMsg)
)
