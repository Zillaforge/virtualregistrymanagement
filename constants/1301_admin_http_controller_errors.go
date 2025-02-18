package constants

import tkErr "pegasus-cloud.com/aes/toolkits/errors"

const (
	// 1301xxxx: Admin HTTP Controller

	AdminAPIInternalServerErrCode             = 13010000
	AdminAPIInternalServerErrMsg              = "internal server error"
	AdminAPILogFileNotFoundErrCode            = 13010001
	AdminAPILogFileNotFoundErrMsg             = "log file not found"
	AdminAPISenderMalformedInputErrCode       = 13010002
	AdminAPISenderMalformedInputErrMsg        = "input format is invalid (%s)"
	AdminAPIQueryNotSupportErrCode            = 13010003
	AdminAPIQueryNotSupportErrMsg             = "%s does not support %s"
	AdminAPIIllegalWhereQueryFormatErrCode    = 13010004
	AdminAPIIllegalWhereQueryFormatErrMsg     = "illegal query format with where"
	AdminAPINamespaceNotFoundErrCode          = 13010005
	AdminAPINamespaceNotFoundErrMsg           = "namespace not found"
	AdminAPIProjectNotFoundErrCode            = 13010006
	AdminAPIProjectNotFoundErrMsg             = "project not found"
	AdminAPIUserNotFoundErrCode               = 13010007
	AdminAPIUserNotFoundErrMsg                = "user not found"
	AdminAPIMembershipNotFoundErrCode         = 13010008
	AdminAPIMembershipNotFoundErrMsg          = "membership not found"
	AdminAPIRepositoryNotFoundErrCode         = 13010009
	AdminAPIRepositoryNotFoundErrMsg          = "repository (%s) not found"
	AdminAPITagNotFoundErrCode                = 13010010
	AdminAPITagNotFoundErrMsg                 = "tag (%s) not found"
	AdminAPITagNotBelongRepositoryErrCode     = 13010011
	AdminAPITagNotBelongRepositoryErrMsg      = "tag (%s) not belong repository"
	AdminAPIUserNotBelongProjectErrCode       = 13010012
	AdminAPIUserNotBelongProjectErrMsg        = "user (%s) not belong project"
	AdminAPIProjectLimitHasBeenReachedErrCode = 13010013
	AdminAPIProjectLimitHasBeenReachedErrMsg  = "project limit has been reached"
	AdminAPIImageNotFoundErrCode              = 13010014
	AdminAPIImageNotFoundErrMsg               = "image not found"
	AdminAPITagExistErrCode                   = 13010015
	AdminAPITagExistErrMsg                    = "tag (%s) exist"
	AdminAPIFileNotFoundErrCode               = 13010016
	AdminAPIFileNotFoundErrMsg                = "file not found"
	AdminAPIFileExistErrCode                  = 13010017
	AdminAPIFileExistErrMsg                   = "there is already a file path with the same name as the specified one"
	AdminAPIRepositoryExistErrCode            = 13010018
	AdminAPIRepositoryExistErrMsg             = "repository (%s) exist"
	AdminAPINotSupportTypeErrCode             = 13010019
	AdminAPINotSupportTypeErrMsg              = "not support type"
	AdminAPIResourceIsProtectedErrCode        = 13010020
	AdminAPIResourceIsProtectedErrMsg         = "resource is protected"
	AdminAPIExportNotFoundErrCode             = 13010021
	AdminAPIExportNotFoundErrMsg              = "export (%s) not found"
	AdminAPIExceedAllowedQuotaErrCode         = 13010022
	AdminAPIExceedAllowedQuotaErrMsg          = "request exceeds allowed quota"
	AdminAPIOperatingSystemMismatchErrCode    = 13010023
	AdminAPIOperatingSystemMismatchErrMsg     = "operating system mismatch"
)

var (
	// 1301xxxx: Admin HTTP Controller

	// 13010000(internal server error)
	AdminAPIInternalServerErr = tkErr.Error(AdminAPIInternalServerErrCode, AdminAPIInternalServerErrMsg)
	// 13010001(log file not found)
	AdminAPILogFileNotFoundErr = tkErr.Error(AdminAPILogFileNotFoundErrCode, AdminAPILogFileNotFoundErrMsg)
	// 13010002(input format is invalid)
	AdminAPISenderMalformedInputErr = tkErr.Error(AdminAPISenderMalformedInputErrCode, AdminAPISenderMalformedInputErrMsg)
	// 13010003(%s does not support %s)
	AdminAPIQueryNotSupportErr = tkErr.Error(AdminAPIQueryNotSupportErrCode, AdminAPIQueryNotSupportErrMsg)
	// 13010004(illegal query format with where)
	AdminAPIIllegalWhereQueryFormatErr = tkErr.Error(AdminAPIIllegalWhereQueryFormatErrCode, AdminAPIIllegalWhereQueryFormatErrMsg)
	// 13010005(namespace not found)
	AdminAPINamespaceNotFoundErr = tkErr.Error(AdminAPINamespaceNotFoundErrCode, AdminAPINamespaceNotFoundErrMsg)
	// 13010006(project not found)
	AdminAPIProjectNotFoundErr = tkErr.Error(AdminAPIProjectNotFoundErrCode, AdminAPIProjectNotFoundErrMsg)
	// 13010007(user not found)
	AdminAPIUserNotFoundErr = tkErr.Error(AdminAPIUserNotFoundErrCode, AdminAPIUserNotFoundErrMsg)
	// 13010008(membership not found)
	AdminAPIMembershipNotFoundErr = tkErr.Error(AdminAPIMembershipNotFoundErrCode, AdminAPIMembershipNotFoundErrMsg)
	// 13010009(repository (%s) not found)
	AdminAPIRepositoryNotFoundErr = tkErr.Error(AdminAPIRepositoryNotFoundErrCode, AdminAPIRepositoryNotFoundErrMsg)
	// 13010010(tag (%s) not found)
	AdminAPITagNotFoundErr = tkErr.Error(AdminAPITagNotFoundErrCode, AdminAPITagNotFoundErrMsg)
	// 13010011(tag (%s) not belong repository)
	AdminAPITagNotBelongRepositoryErr = tkErr.Error(AdminAPITagNotBelongRepositoryErrCode, AdminAPITagNotBelongRepositoryErrMsg)
	// 13010012(user (%s) not belong project)
	AdminAPIUserNotBelongProjectErr = tkErr.Error(AdminAPIUserNotBelongProjectErrCode, AdminAPIUserNotBelongProjectErrMsg)
	// 13010013(project limit has been reached)
	AdminAPIProjectLimitHasBeenReachedErr = tkErr.Error(AdminAPIProjectLimitHasBeenReachedErrCode, AdminAPIProjectLimitHasBeenReachedErrMsg)
	// 13010014(image not found)
	AdminAPIImageNotFoundErr = tkErr.Error(AdminAPIImageNotFoundErrCode, AdminAPIImageNotFoundErrMsg)
	// 13010015(tag (%s) exist)
	AdminAPITagExistErr = tkErr.Error(AdminAPITagExistErrCode, AdminAPITagExistErrMsg)
	// 13010016(file not found)
	AdminAPIFileNotFoundErr = tkErr.Error(AdminAPIFileNotFoundErrCode, AdminAPIFileNotFoundErrMsg)
	// 13010017(file exist)
	AdminAPIFileExistErr = tkErr.Error(AdminAPIFileExistErrCode, AdminAPIFileExistErrMsg)
	// 13010018(repository (%s) exist)
	AdminAPIRepositoryExistErr = tkErr.Error(AdminAPIRepositoryExistErrCode, AdminAPIRepositoryExistErrMsg)
	// 13010019(not support type)
	AdminAPINotSupportTypeErr = tkErr.Error(AdminAPINotSupportTypeErrCode, AdminAPINotSupportTypeErrMsg)
	// 13010020(resource is protected)
	AdminAPIResourceIsProtectedErr = tkErr.Error(AdminAPIResourceIsProtectedErrCode, AdminAPIResourceIsProtectedErrMsg)
	// 13010021(export (%s) not found)
	AdminAPIExportNotFoundErr = tkErr.Error(AdminAPIExportNotFoundErrCode, AdminAPIExportNotFoundErrMsg)
	// 13010022(request exceed allowed quota)
	AdminAPIExceedAllowedQuotaErr = tkErr.Error(AdminAPIExceedAllowedQuotaErrCode, AdminAPIExceedAllowedQuotaErrMsg)
	// 13020023(operating system mismatch)
	AdminAPIOperatingSystemMismatchErr = tkErr.Error(UserAPIOperatingSystemMismatchErrCode, UserAPIOperatingSystemMismatchErrMsg)
)
