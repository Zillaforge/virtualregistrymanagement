package constants

import tkErr "pegasus-cloud.com/aes/toolkits/errors"

const (
	// 1200xxxx: Middleware

	MidInternalServerErrorErrCode     = 12000000
	MidInternalServerErrorErrMsg      = "internal server error"
	MidPermissionDeniedErrCode        = 12000001
	MidPermissionDeniedErrMsg         = "permission denied"
	MidUserHasBeenFrozenErrCode       = 12000002
	MidUserHasBeenFrozenErrMsg        = "the user has been frozen, please contact administrator"
	MidMembershipHasBeenFrozenErrCode = 12000003
	MidMembershipHasBeenFrozenErrMsg  = "the membership has been frozen, please contact tenant of admin"
	MidProjectNotFoundErrCode         = 12000004
	MidProjectNotFoundErrMsg          = "project (%s) not found"
	MidIncorrectFormatErrCode         = 12000005
	MidIncorrectFormatErrMsg          = "incorrect format of authentication"
	MidRepositoryNotFoundErrCode      = 12000006
	MidRepositoryNotFoundErrCodeMsg   = "repository (%s) not found"
	MidRepositoryIsReadOnlyErrCode    = 12000007
	MidRepositoryIsReadOnlyErrMsg     = "repository (%s) is read only for users of other projects"
	MidTagNotFoundErrCode             = 12000008
	MidTagNotFoundErrCodeMsg          = "tag (%s) not found"
	MidTagIsReadOnlyErrCode           = 12000009
	MidTagIsReadOnlyErrMsg            = "tag (%s) is read only for users of other projects"
	MidServerNotFoundErrCode          = 12000010
	MidServerNotFoundErrCodeMsg       = "server (%s) not found"
	MidProjectAclNotFoundErrCode      = 12000011
	MidProjectAclNotFoundErrMsg       = "project-acl (%s) not found"
	MidMemberAclNotFoundErrCode       = 12000012
	MidMemberAclNotFoundErrMsg        = "member-acl (%s) not found"
	MidExportNotFoundErrCode          = 12000013
	MidExportNotFoundErrCodeMsg       = "export (%s) not found"
	MidExportIsReadOnlyErrCode        = 12000014
	MidExportIsReadOnlyErrMsg         = "export (%s) is read only for users of other projects"
)

var (
	// 1200xxxx: Middleware

	// 12000000(internal server error)
	MidInternalServerErrorErr = tkErr.Error(MidInternalServerErrorErrCode, MidInternalServerErrorErrMsg)
	// 12000001(permission denied)
	MidPermissionDeniedErr = tkErr.Error(MidPermissionDeniedErrCode, MidPermissionDeniedErrMsg)
	// 12000002(the user has been frozen, please contact administrator)
	MidUserHasBeenFrozenErr = tkErr.Error(MidUserHasBeenFrozenErrCode, MidUserHasBeenFrozenErrMsg)
	// 12000003(the membership has been frozen, please contact tenant of admin)
	MidMembershipHasBeenFrozenErr = tkErr.Error(MidMembershipHasBeenFrozenErrCode, MidMembershipHasBeenFrozenErrMsg)
	// 12000004(project (%s) not found)
	MidProjectNotFoundErr = tkErr.Error(MidProjectNotFoundErrCode, MidProjectNotFoundErrMsg)
	// 12000005(incorrect format of authentication)
	MidIncorrectFormatErr = tkErr.Error(MidIncorrectFormatErrCode, MidIncorrectFormatErrMsg)
	// 12000006(repository (%s) not found)
	MidRepositoryNotFoundErr = tkErr.Error(MidRepositoryNotFoundErrCode, MidRepositoryNotFoundErrCodeMsg)
	// 12000007(repository (%s) is read only for users of other projects)
	MidRepositoryIsReadOnlyErr = tkErr.Error(MidRepositoryIsReadOnlyErrCode, MidRepositoryIsReadOnlyErrMsg)
	// 12000008(tag (%s) not found)
	MidTagNotFoundErr = tkErr.Error(MidTagNotFoundErrCode, MidTagNotFoundErrCodeMsg)
	// 12000009(tag (%s) is read only for users of other projects)
	MidTagIsReadOnlyErr = tkErr.Error(MidTagIsReadOnlyErrCode, MidTagIsReadOnlyErrMsg)
	// 12000010(server (%s) not found)
	MidServerNotFoundErr = tkErr.Error(MidServerNotFoundErrCode, MidServerNotFoundErrCodeMsg)
	// 12000011(project-acl (%s) not found)
	MidProjectAclNotFoundErr = tkErr.Error(MidProjectAclNotFoundErrCode, MidProjectAclNotFoundErrMsg)
	// 12000012(member-acl (%s) not found)
	MidMemberAclNotFoundErr = tkErr.Error(MidMemberAclNotFoundErrCode, MidMemberAclNotFoundErrMsg)
	// 12000013(export (%s) not found)
	MidExportNotFoundErr = tkErr.Error(MidExportNotFoundErrCode, MidExportNotFoundErrCodeMsg)
	// 12000014(export (%s) is read only for users of other projects)
	MidExportIsReadOnlyErr = tkErr.Error(MidExportIsReadOnlyErrCode, MidExportIsReadOnlyErrMsg)
)
