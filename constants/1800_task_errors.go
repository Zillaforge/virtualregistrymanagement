package constants

import tkErr "pegasus-cloud.com/aes/toolkits/errors"

const (
	// 1800xxxx: Task

	TaskInternalServerErrCode = 18000000
	TaskInternalServerErrMsg  = "internal server error"
	TaskProjectExistErrCode   = 18000001
	TaskProjectExistErrMsg    = "project (%s) exist"
)

var (
	// 1800xxxx: Task

	// 18000000(internal server error)
	TaskInternalServerErr = tkErr.Error(TaskInternalServerErrCode, TaskInternalServerErrMsg)
	// 18000001(project (%s) exist)
	TaskProjectExistErr = tkErr.Error(TaskProjectExistErrCode, TaskProjectExistErrMsg)
)
