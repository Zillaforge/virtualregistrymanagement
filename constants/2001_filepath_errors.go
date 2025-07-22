package constants

import tkErr "github.com/Zillaforge/toolkits/errors"

const (
	// 2001xxxx: filepath
	FilepathTypeAndSchemeIsRequiredErrCode = 20010000
	FilepathTypeAndSchemeIsRequiredErrMsg  = "type and scheme is required. it can not be empty"
	FilepathSchemeIsRepeatedErrCode        = 20010001
	FilepathSchemeIsRepeatedErrMsg         = "scheme is repeated"
	FilepathTypeIsNotSupportedErrCode      = 20010002
	FilepathTypeIsNotSupportedErrMsg       = "filepath type is not supported"
	FilepathFilepathIsExistErrCode         = 20010003
	FilepathFilepathIsExistErrMsg          = "there is already a file path with the same name as the specified one"
)

var (
	// 2001xxxx: filepath

	// 20010000(type and scheme is required. it can not be empty)
	FilepathTypeAndSchemeIsRequiredErr = tkErr.Error(FilepathTypeAndSchemeIsRequiredErrCode, FilepathTypeAndSchemeIsRequiredErrMsg)
	// 20010001(scheme is repeated)
	FilepathSchemeIsRepeatedErr = tkErr.Error(FilepathSchemeIsRepeatedErrCode, FilepathSchemeIsRepeatedErrMsg)
	// 20010002(filepath type is not supported)
	FilepathTypeIsNotSupportedErr = tkErr.Error(FilepathTypeIsNotSupportedErrCode, FilepathTypeIsNotSupportedErrMsg)
	// 20010003(there is already a file path with the same name as the specified one)
	FilepathFilepathIsExistErr = tkErr.Error(FilepathFilepathIsExistErrCode, FilepathFilepathIsExistErrMsg)
)
