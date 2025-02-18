package constants

import tkErr "pegasus-cloud.com/aes/toolkits/errors"

const (
	// 1101xxxx: Config

	ConfigReadConfigFailedErrCode  = 11010000
	ConfigReadConfigFailedErrMsg   = "config read failed"
	ConfigVersionMustBeSameErrCode = 11010001
	ConfigVersionMustBeSameErrMsg  = "config version must be same with service"
)

var (
	// 1101xxxx: Config

	// 11010000(config read failed)
	ConfigReadConfigFailedErr = tkErr.Error(ConfigReadConfigFailedErrCode, ConfigReadConfigFailedErrMsg)
	// 11010001(config version must be same with service)
	ConfigVersionMustBeSameErr = tkErr.Error(ConfigVersionMustBeSameErrCode, ConfigVersionMustBeSameErrMsg)
)
