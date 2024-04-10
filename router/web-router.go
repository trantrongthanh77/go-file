package router

import (
	"go-file/common"
	"go-file/controller"
	"go-file/externalinterface/storage"
	"go-file/middleware"

	"github.com/gin-gonic/gin"
)

func setWebRouter(router *gin.Engine, cloudinary storage.Cloudinary) {
	router.Use(middleware.GlobalWebRateLimit())
	// Always available
	// All page must have username in context
	router.GET("/", middleware.ExtractUserInfo(), controller.GetIndexPage)
	router.GET("/public/static/:file", controller.GetStaticFile)
	router.GET("/public/lib/:file", controller.GetLibFile)
	router.GET("/public/icon/:file", controller.GetIconFile)
	router.GET("/login", middleware.ExtractUserInfo(), controller.GetLoginPage)
	router.POST("/login", middleware.CriticalRateLimit(), controller.Login)
	router.GET("/logout", controller.Logout)
	router.GET("/help", middleware.ExtractUserInfo(), controller.GetHelpPage)

	explorerController := controller.NewPageExplorer(cloudinary)
	// Download files
	fileDownloadAuth := router.Group("/")
	fileDownloadAuth.Use(middleware.DownloadRateLimit(), middleware.FileDownloadPermissionCheck())
	{
		fileDownloadAuth.GET("/upload/*filepath", controller.DownloadFile)
		fileDownloadAuth.GET("/explorer", middleware.ExtractUserInfo(), explorerController.GetExplorerPageOrFile)
	}

	imageDownloadAuth := router.Group("/")
	imageDownloadAuth.Use(middleware.DownloadRateLimit(), middleware.ImageDownloadPermissionCheck())
	{
		imageDownloadAuth.Static("/image", common.ImageUploadPath)
	}

	router.GET("/image", middleware.ExtractUserInfo(), controller.GetImagePage)

	router.GET("/video", middleware.ExtractUserInfo(), controller.GetVideoPage)

	basicAuth := router.Group("/")
	basicAuth.Use(middleware.WebAuth()) // WebAuth already has username in context
	{
		basicAuth.GET("/manage", controller.GetManagePage)
	}
}
