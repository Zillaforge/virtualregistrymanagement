package constants

import tkErr "pegasus-cloud.com/aes/toolkits/errors"

const (
	// 1100xxxx: CMD

	CMDInternalServerErrCode = 11000000
	CMDInternalServerErrMsg  = "internal server error"
)

var (
	// 1100xxxx: CMD

	// 11000000(internal server error)
	CMDInternalServerErr = tkErr.Error(CMDInternalServerErrCode, CMDInternalServerErrMsg)
)
