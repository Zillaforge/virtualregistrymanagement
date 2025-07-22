package constants

import tkErr "github.com/Zillaforge/toolkits/errors"

const (
	// 9900xxxx: Plugin

	PluginInternalServerErrCode      = 99000000
	PluginInternalServerErrMsg       = "internal server error"
	PluginVersionDoesNotMatchErrCode = 99000001
	PluginVersionDoesNotMatchErrMsg  = "%s plugin version does not match"
)

var (
	// 9900xxxx: Plugin

	// 99000000(internal server error)
	PluginInternalServerErr = tkErr.Error(PluginInternalServerErrCode, PluginInternalServerErrMsg)
	// 99000001(plugin version does not match)
	PluginVersionDoesNotMatchErr = tkErr.Error(PluginVersionDoesNotMatchErrCode, PluginVersionDoesNotMatchErrMsg)
)
