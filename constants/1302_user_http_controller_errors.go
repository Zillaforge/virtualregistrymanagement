package constants

import tkErr "github.com/Zillaforge/toolkits/errors"

const (
	// 1302xxxx: User HTTP Controller
	UserAPIInternalServerErrCode             = 13020000
	UserAPIInternalServerErrMsg              = "internal server error"
	UserAPIQueryNotSupportErrCode            = 13020001
	UserAPIQueryNotSupportErrMsg             = "%s does not support %s"
	UserAPIIllegalWhereQueryFormatErrCode    = 13020002
	UserAPIIllegalWhereQueryFormatErrMsg     = "illegal query format with where"
	UserAPIInvalidCertificateTypeErrCode     = 13020003
	UserAPIInvalidCertificateTypeErrMsg      = "invalid certificate type"
	UserAPIInvalidCertificateFormatErrCode   = 13020004
	UserAPIInvalidCertificateFormatErrMsg    = "invalid certificate format"
	UserAPIUnauthorizedOpErrCode             = 13020005
	UserAPIUnauthorizedOpErrMsg              = "unauthorized operation"
	UserAPIRepositoryNotFoundErrCode         = 13020006
	UserAPIRepositoryNotFoundErrMsg          = "repository (%s) not found"
	UserAPITagNotFoundErrCode                = 13020007
	UserAPITagNotFoundErrMsg                 = "tag (%s) not found"
	UserAPITagNotBelongRepositoryErrCode     = 13020008
	UserAPITagNotBelongRepositoryErrMsg      = "tag (%s) not belong repository"
	UserAPIUserNotBelongProjectErrCode       = 13020009
	UserAPIUserNotBelongProjectErrMsg        = "user (%s) not belong project"
	UserAPIProjectLimitHasBeenReachedErrCode = 13020010
	UserAPIProjectLimitHasBeenReachedErrMsg  = "project limit has been reached"
	UserAPIImageNotFoundErrCode              = 13020011
	UserAPIImageNotFoundErrMsg               = "image not found"
	UserAPITagExistErrCode                   = 13020012
	UserAPITagExistErrMsg                    = "tag (%s) exist"
	UserAPIFileNotFoundErrCode               = 13020013
	UserAPIFileNotFoundErrMsg                = "file not found"
	UserAPIFileExistErrCode                  = 13020014
	UserAPIFileExistErrMsg                   = "there is already a file path with the same name as the specified one"
	UserAPIRepositoryExistErrCode            = 13020015
	UserAPIRepositoryExistErrMsg             = "repository (%s) exist"
	UserAPINotSupportTypeErrCode             = 13020016
	UserAPINotSupportTypeErrMsg              = "not support type"
	UserAPIResourceIsProtectedErrCode        = 13020017
	UserAPIResourceIsProtectedErrMsg         = "resource is protected"
	UserAPIExportNotFoundErrCode             = 13020018
	UserAPIExportNotFoundErrMsg              = "export (%s) not found"
	UserAPIExceedAllowedQuotaErrCode         = 13020019
	UserAPIExceedAllowedQuotaErrMsg          = "request exceeds allowed quota"
	UserAPIOperatingSystemMismatchErrCode    = 13020020
	UserAPIOperatingSystemMismatchErrMsg     = "operating system mismatch"
)

var (
	// 1302xxxx: User HTTP Controller

	// 13020000(internal server error)
	UserAPIInternalServerErr = tkErr.Error(UserAPIInternalServerErrCode, UserAPIInternalServerErrMsg)
	// 13020001(%s does not support %s)
	UserAPIQueryNotSupportErr = tkErr.Error(UserAPIQueryNotSupportErrCode, UserAPIQueryNotSupportErrMsg)
	// 13020002(illegal query format with where)
	UserAPIIllegalWhereQueryFormatErr = tkErr.Error(UserAPIIllegalWhereQueryFormatErrCode, UserAPIIllegalWhereQueryFormatErrMsg)
	// 13000003(invalid certificate type)
	UserAPIInvalidCertificateTypeErr = tkErr.Error(UserAPIInvalidCertificateTypeErrCode, UserAPIInvalidCertificateTypeErrMsg)
	// 13000004(invalid certificate formate)
	UserAPIInvalidCertificateFormatErr = tkErr.Error(UserAPIInvalidCertificateFormatErrCode, UserAPIInvalidCertificateFormatErrMsg)
	// 13020005(unauthorized operation)
	UserAPIUnauthorizedOpErr = tkErr.Error(UserAPIUnauthorizedOpErrCode, UserAPIUnauthorizedOpErrMsg)
	// 13020006(repository (%s) not found)
	UserAPIRepositoryNotFoundErr = tkErr.Error(UserAPIRepositoryNotFoundErrCode, UserAPIRepositoryNotFoundErrMsg)
	// 13020007(tag (%s) not found)
	UserAPITagNotFoundErr = tkErr.Error(UserAPITagNotFoundErrCode, UserAPITagNotFoundErrMsg)
	// 13020008(tag (%s) not belong repository)
	UserAPITagNotBelongRepositoryErr = tkErr.Error(UserAPITagNotBelongRepositoryErrCode, UserAPITagNotBelongRepositoryErrMsg)
	// 13020009(user (%s) not belong project)
	UserAPIUserNotBelongProjectErr = tkErr.Error(UserAPIUserNotBelongProjectErrCode, UserAPIUserNotBelongProjectErrMsg)
	// 13020010(project limit has been reached)
	UserAPIProjectLimitHasBeenReachedErr = tkErr.Error(UserAPIProjectLimitHasBeenReachedErrCode, UserAPIProjectLimitHasBeenReachedErrMsg)
	// 13020011(image not found)
	UserAPIImageNotFoundErr = tkErr.Error(UserAPIImageNotFoundErrCode, UserAPIImageNotFoundErrMsg)
	// 13020012(tag (%s) exist)
	UserAPITagExistErr = tkErr.Error(UserAPITagExistErrCode, UserAPITagExistErrMsg)
	// 13020013(file not found)
	UserAPIFileNotFoundErr = tkErr.Error(UserAPIFileNotFoundErrCode, UserAPIFileNotFoundErrMsg)
	// 13020014(file exist)
	UserAPIFileExistErr = tkErr.Error(UserAPIFileExistErrCode, UserAPIFileExistErrMsg)
	// 13020015(repository (%s) exist)
	UserAPIRepositoryExistErr = tkErr.Error(UserAPIRepositoryExistErrCode, UserAPIRepositoryExistErrMsg)
	// 13020016(not support type)
	UserAPINotSupportTypeErr = tkErr.Error(UserAPINotSupportTypeErrCode, UserAPINotSupportTypeErrMsg)
	// 13020017(resource is protected)
	UserAPIResourceIsProtectedErr = tkErr.Error(UserAPIResourceIsProtectedErrCode, UserAPIResourceIsProtectedErrMsg)
	// 13020018(export (%s) not found)
	UserAPIExportNotFoundErr = tkErr.Error(UserAPIExportNotFoundErrCode, UserAPIExportNotFoundErrMsg)
	// 13020019(request exceed allowed quota)
	UserAPIExceedAllowedQuotaErr = tkErr.Error(UserAPIExceedAllowedQuotaErrCode, UserAPIExceedAllowedQuotaErrMsg)
	// 13020020(operating system mismatch)
	UserAPIOperatingSystemMismatchErr = tkErr.Error(UserAPIOperatingSystemMismatchErrCode, UserAPIOperatingSystemMismatchErrMsg)
)
