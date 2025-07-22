package utility

import (
	"github.com/gin-gonic/gin"
	tkErr "github.com/Zillaforge/toolkits/errors"
)

type (
	// ResponseType ...
	ResponseType int
	// ErrResponse ...
	ErrResponse struct {
		ErrorCode int                    `json:"errorCode"`
		Message   string                 `json:"message,omitempty"`
		Meta      map[string]interface{} `json:"meta,omitempty"`
	}
)

const (
	// JSON ...
	JSON ResponseType = iota
	// XML ...
	XML
)

// RouteResponseType ...
var RouteResponseType ResponseType

// ResponseWithType ...
func ResponseWithType(c *gin.Context, statusCode int, body interface{}) {
	switch RouteResponseType {
	case XML:
		c.Abort()
		c.XML(statusCode, body)
	case JSON:
		fallthrough
	default:
		if err, ok := body.(error); ok {
			if e, ok := tkErr.IsError(err); ok {
				c.AbortWithStatusJSON(statusCode, &ErrResponse{
					ErrorCode: e.Code(),
					Message:   e.Message(),
				})
				return
			}
		}
		c.AbortWithStatusJSON(statusCode, body)
	}
}
