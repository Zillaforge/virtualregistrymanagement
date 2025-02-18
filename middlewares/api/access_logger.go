package api

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func AccessLoggerMiddleware(c *gin.Context) {
	// Start timer
	params := logger.NewAccessParams()
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	// Process request
	c.Next()

	if raw != "" {
		path = path + "?" + raw
	}

	params.Meta.Latency = time.Since(start).String()
	params.Meta.Path = path
	params.Meta.Method = c.Request.Method
	params.Meta.StatusCode = c.Writer.Status()

	params.
		SetUserID(c.GetString("ctxUserID")).
		SetSAATUserIDByContext(c).
		SetProjectIDByContext(c).
		SetRequestID(c.GetString(cnt.RequestID)).
		SetTriggerType(logger.API).
		SetSourceIP(c.ClientIP()).
		SetAccessLoggerInfoByPBAC(c).Writer()

}
