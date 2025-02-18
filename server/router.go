package server

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/eventpublish"
	mid "VirtualRegistryManagement/middlewares/api"
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	tkMid "pegasus-cloud.com/aes/toolkits/middleware"
	"pegasus-cloud.com/aes/toolkits/mviper"
)

func router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	if mviper.GetBool("VirtualRegistryManagement.developer") {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.New()
	router.UseRawPath = true
	router.UnescapePathValues = false
	router.Use(mid.GinLogger)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     mviper.GetStringSlice("VirtualRegistryManagement.http.access_control.allow_origins"),
		AllowMethods:     mviper.GetStringSlice("VirtualRegistryManagement.http.access_control.allow_methods"),
		AllowHeaders:     mviper.GetStringSlice("VirtualRegistryManagement.http.access_control.allow_headers"),
		ExposeHeaders:    mviper.GetStringSlice("VirtualRegistryManagement.http.access_control.expose_headers"),
		AllowCredentials: mviper.GetBool("VirtualRegistryManagement.http.access_control.allow_credentials"),
	}))

	router.Use(mid.APIOperationMiddleware, tkMid.RequestIDGenerator, mid.SetExtraHeaders, mid.AccessLoggerMiddleware)
	// Base Path: /vrm/api/v1
	rootRG := router.Group(path.Join(cnt.APIPrefix, cnt.APIVersion))
	enableUserVirtualRegistryManagementRouter(rootRG.Group(""))
	enableAdminVirtualRegistryManagementRouter(rootRG.Group("admin"))
	eventpublish.EnableHTTPRouters(rootRG)

	return router
}
