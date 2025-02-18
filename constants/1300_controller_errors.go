package constants

import tkErr "pegasus-cloud.com/aes/toolkits/errors"

const (
	// 1300xxxx: Controller

	ControllerInternalServerErrCode = 13010000
	ControllerInternalServerErrMsg  = "internal server error"
)

var (
	// 1300xxxx: Controller

	// 13000000(internal server error)
	ControllerInternalServerErr = tkErr.Error(ControllerInternalServerErrCode, ControllerInternalServerErrMsg)
)
