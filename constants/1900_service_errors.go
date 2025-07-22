package constants

import tkErr "github.com/Zillaforge/toolkits/errors"

const (
	// 1900xxxx: Service

	ServiceInternalServerErrCode    = 19000000
	ServiceInternalServerErrMsg     = "internal server error"
	ServiceNameIsRequiredErrCode    = 19000001
	ServiceNameIsRequiredErrMsg     = "service name is required. it can not be empty"
	ServiceNameMustBeAStringErrCode = 19000002
	ServiceNameMustBeAStringErrMsg  = "service name must be a string"
	ServiceNameIsRepeatedErrCode    = 19000003
	ServiceNameIsRepeatedErrMsg     = "service name is repeated"
)

var (
	// 1900xxxx: Service

	// 19000000(internal server error)
	ServiceInternalServerErr = tkErr.Error(ServiceInternalServerErrCode, ServiceInternalServerErrMsg)
	// 19000001(service name is required. it can not be empty)
	ServiceNameIsRequiredErr = tkErr.Error(ServiceNameIsRequiredErrCode, ServiceNameIsRequiredErrMsg)
	// 19000002(service name must be a string)
	ServiceNameMustBeAStringErr = tkErr.Error(ServiceNameMustBeAStringErrCode, ServiceNameMustBeAStringErrMsg)
	// 19000003(service name is repeated)
	ServiceNameIsRepeatedErr = tkErr.Error(ServiceNameIsRepeatedErrCode, ServiceNameIsRepeatedErrMsg)
)
