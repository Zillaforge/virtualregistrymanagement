package constants

import tkErr "pegasus-cloud.com/aes/toolkits/errors"

const (
	// 2000xxxx: Openstack
	OpenstackInternalServerErrCode                = 20000000
	OpenstackInternalServerErrMsg                 = "internal server error"
	OpenstackNamespaceAndServiceIsRequiredErrCode = 20000001
	OpenstackNamespaceAndServiceIsRequiredErrMsg  = "namespace and service is required. it can not be empty"
	OpenstackNamespaceIsRepeatedErrCode           = 20000002
	OpenstackNamespaceIsRepeatedErrMsg            = "namespace is repeated"
	OpenstackTypeIsNotSupportedErrCode            = 20000003
	OpenstackTypeIsNotSupportedErrMsg             = "service type is not supported in openstack"
	OpenstackResourceIsNotSupportedErrCode        = 20000004
	OpenstackResourceIsNotSupportedErrMsg         = "resource is not supported"
	OpenstackConnectionIsNotCreatedErrCode        = 20000005
	OpenstackConnectionIsNotCreatedErrMsg         = "connection is not created"
	OpenstackUploadFileNotFoundErrCode            = 20000006
	OpenstackUploadFileNotFoundErrMsg             = "upload file not found"
	OpenstackCreateFileFailedErrCode              = 20000007
	OpenstackCreateFileFailedErrMsg               = "create file failed"
	OpenstackImageToFileFailedErrCode             = 20000008
	OpenstackImageToFileFailedErrMsg              = "image to file failed"
	OpenstackFileIsExistErrCode                   = 20000009
	OpenstackFileIsExistErrMsg                    = "file is exist"
	OpenstackExceedAllowedQuotaErrCode            = 20000010
	OpenstackExceedAllowedQuotaErrMsg             = "request exceeds allowed quota"
)

var (
	// 2000xxxx: Openstack

	// 20000000(internal server error)
	OpenstackInternalServerErr = tkErr.Error(OpenstackInternalServerErrCode, OpenstackInternalServerErrMsg)
	// 20000001(namespace and service is required. it can not be empty)
	OpenstackNamespaceAndServiceIsRequiredErr = tkErr.Error(OpenstackNamespaceAndServiceIsRequiredErrCode, OpenstackNamespaceAndServiceIsRequiredErrMsg)
	// 20000002(namespace is repeated)
	OpenstackNamespaceIsRepeatedErr = tkErr.Error(OpenstackNamespaceIsRepeatedErrCode, OpenstackNamespaceIsRepeatedErrMsg)
	// 20000003(service type is not supported in openstack)
	OpenstackTypeIsNotSupportedErr = tkErr.Error(OpenstackTypeIsNotSupportedErrCode, OpenstackTypeIsNotSupportedErrMsg)
	// 20000004(resource is not supported)
	OpenstackResourceIsNotSupportedErr = tkErr.Error(OpenstackResourceIsNotSupportedErrCode, OpenstackResourceIsNotSupportedErrMsg)
	// 20000005(connection is not created)
	OpenstackConnectionIsNotCreatedErr = tkErr.Error(OpenstackConnectionIsNotCreatedErrCode, OpenstackConnectionIsNotCreatedErrMsg)
	// 20000006(upload file not found)
	OpenstackUploadFileNotFoundErr = tkErr.Error(OpenstackUploadFileNotFoundErrCode, OpenstackUploadFileNotFoundErrMsg)
	// 20000007(create file failed)
	OpenstackCreateFileFailedErr = tkErr.Error(OpenstackCreateFileFailedErrCode, OpenstackCreateFileFailedErrMsg)
	// 20000008(image to file failed)
	OpenstackImageToFileFailedErr = tkErr.Error(OpenstackImageToFileFailedErrCode, OpenstackImageToFileFailedErrMsg)
	// 20000009(file is exist)
	OpenstackFileIsExistErr = tkErr.Error(OpenstackFileIsExistErrCode, OpenstackFileIsExistErrMsg)
	// 20000010(request exceeds allowed quota)
	OpenstackExceedAllowedQuotaErr = tkErr.Error(OpenstackExceedAllowedQuotaErrCode, OpenstackExceedAllowedQuotaErrMsg)
)
