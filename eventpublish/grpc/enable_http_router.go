package grpc

import (
	cnt "VirtualRegistryManagement/constants"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	pbac "pegasus-cloud.com/aes/toolkits/pbac/gin"
)

func (c *core) EnableHTTPRouter(rg *gin.RouterGroup) {
	getRouterOutput, err := c.handler.GetRouter()
	if err != nil {
		zap.L().With(
			zap.String(cnt.Plugin, "g.handler.GetRouter(...)"),
		).Error(err.Error())
		return
	}
	for _, route := range getRouterOutput.Response {
		getAccessLoggerInfo := cnt.GetAccessLoggerInfo(route.ActionName)
		if getAccessLoggerInfo == nil {
			insertNewAccessLoggerInfo := cnt.InsertNewAccessLoggerInfo(route.ActionName, int(route.ActionID))
			getAccessLoggerInfo = &insertNewAccessLoggerInfo
		}
		switch route.Method {
		case http.MethodPost:
			pbac.POST(rg, route.Path, c.http2grpc, getAccessLoggerInfo.Name, route.Administrator)
		case http.MethodGet:
			pbac.GET(rg, route.Path, c.http2grpc, getAccessLoggerInfo.Name, route.Administrator)
		case http.MethodPut:
			pbac.PUT(rg, route.Path, c.http2grpc, getAccessLoggerInfo.Name, route.Administrator)
		case http.MethodDelete:
			pbac.DELETE(rg, route.Path, c.http2grpc, getAccessLoggerInfo.Name, route.Administrator)
		}
	}
}
