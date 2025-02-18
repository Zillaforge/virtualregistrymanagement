package constants

import tkErr "pegasus-cloud.com/aes/toolkits/errors"

const (
	// 1102xxxx: Server

	ServerInternalServerErrCode = 11020000
	ServerInternalServerErrMsg  = "internal server error"
)

var (
	// 1102xxxx: Server

	// 11020000(internal server error)
	ServerInternalServerErr = tkErr.Error(ServerInternalServerErrCode, ServerInternalServerErrMsg)
)
