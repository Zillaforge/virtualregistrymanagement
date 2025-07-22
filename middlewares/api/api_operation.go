package api

import (
	"github.com/gin-gonic/gin"
	tkMid "github.com/Zillaforge/toolkits/middleware"
	pbac "github.com/Zillaforge/toolkits/pbac/gin"
)

func APIOperationMiddleware(c *gin.Context) {
	var actionName = "Unknown"
	hasRouter, entity := pbac.Actions.Checker(c)
	if hasRouter {
		actionName = entity.Action
	}
	c.Set(tkMid.CtxOperationName, actionName)
}
