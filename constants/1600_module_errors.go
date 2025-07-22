package constants

import tkErr "github.com/Zillaforge/toolkits/errors"

const (
	// 1600xxxx: Module

	ModuleInternalServerErrCode = 16000000
	ModuleInternalServerErrMsg  = "internal server error"
)

var (
	// 1600xxxx: Module

	// 16000000(internal server error)
	ModuleInternalServerErr = tkErr.Error(ModuleInternalServerErrCode, ModuleInternalServerErrMsg)
)
