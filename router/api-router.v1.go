package router

import (
	"go-file/common/config"
	"go-file/controller"
	"go-file/middleware"

	"github.com/gin-gonic/gin"
)

func setApiRouter(router *gin.Engine, conf *config.Config) {
	router.Use(middleware.GlobalAPIRateLimit())
	router.GET("/status", func(c *gin.Context) {
		controller.GetStatus(c, conf)
	})
	router.POST("/api/file", middleware.FileUploadPermissionCheck(), controller.UploadFile)
	router.POST("/api/image", middleware.ImageUploadPermissionCheck(), controller.UploadImage)
	router.GET("/api/notice", controller.GetNotice)
	basicAuth := router.Group("/api")
	basicAuth.Use(middleware.ApiAuth())
	{
		basicAuth.DELETE("/file", controller.DeleteFile)
		basicAuth.DELETE("/image", controller.DeleteImage)
		basicAuth.PUT("/user", middleware.NoTokenAuth(), controller.UpdateSelf)
		basicAuth.POST("/token", controller.GenerateNewUserToken)
	}
	adminAuth := router.Group("/api")
	adminAuth.Use(middleware.ApiAdminAuth())
	{
		adminAuth.POST("/user", controller.CreateUser)
		adminAuth.PUT("/manage_user", controller.ManageUser)
		adminAuth.GET("/option", controller.GetOptions)
		adminAuth.PUT("/option", controller.UpdateOption)
		statRouter := adminAuth.Group("/stat")
		{
			statRouter.GET("/ip", controller.GetIPs)
			statRouter.GET("/url", controller.GetURLs)
			statRouter.GET("/req", controller.GetReqs)
		}
	}
}
