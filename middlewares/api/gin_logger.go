package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Zillaforge/toolkits/tracer"
)

// Gin Logger Format:
//
//	[GIN] time | status | latency | clientIP | requestID | method | path
//	Error #01: error1 message
//	Error #02: error2 message
//
// Example:
//
//	[GIN] 2021/01/01 - 10:23:56 | 200 |    36.53674ms |      172.40.0.1 | 370c49be-363c-4fc3-8bf4-3143a47050eb | GET | /pt/api/v1/version
//	[GIN] 2021/01/01 - 10:25:39 | 500 |   37.942766ms |      172.40.0.1 | 830d4714-3251-4523-95ef-a988ebce3273 | GET | /pt/api/v1/model/ee9c99c9-d526-42ad-9e9b-1418be1d6655
//	Error #01: E13000000 internal server error [rpc error: code = Code(15000006) desc = model not found]
func GinLogger(c *gin.Context) {
	const (
		green   = "\033[97;42m"
		white   = "\033[90;47m"
		yellow  = "\033[90;43m"
		red     = "\033[97;41m"
		blue    = "\033[97;44m"
		magenta = "\033[97;45m"
		cyan    = "\033[97;46m"
		reset   = "\033[0m"
	)

	gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, requestIDColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			requestIDColor = green
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %v %s| %s %-7s %s %#v\n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			requestIDColor, param.Keys[tracer.RequestID], resetColor,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	},
	)(c)
}
