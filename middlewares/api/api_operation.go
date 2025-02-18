package api

import (
	"github.com/gin-gonic/gin"
	tkMid "pegasus-cloud.com/aes/toolkits/middleware"
	pbac "pegasus-cloud.com/aes/toolkits/pbac/gin"
)

func APIOperationMiddleware(c *gin.Context) {
	var actionName = "Unknown"
	hasRouter, entity := pbac.Actions.Checker(c)
	if hasRouter {
		actionName = entity.Action
	}
	c.Set(tkMid.CtxOperationName, actionName)
}
