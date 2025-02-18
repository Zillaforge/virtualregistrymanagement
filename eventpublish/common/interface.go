package common

import "github.com/gin-gonic/gin"

type EventIntf interface {
	GetName() string
	GetVersion() string
	SetConfig(conf []byte)
	CheckPluginVersion() bool
	InitPlugin() bool
	Reconcile(action string, meta map[string]string, req interface{}, resp interface{})
	CallGRPCRouter(operator string, hdr map[string]string, payload []byte) (map[string]string, []byte, error)
	EnableHTTPRouter(rg *gin.RouterGroup)
}
