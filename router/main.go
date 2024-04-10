package router

import (
	"go-file/common/config"
	"go-file/controller"
	"go-file/externalinterface/storage"
	"go-file/middleware"

	"github.com/gin-gonic/gin"
)

func SetRouter(router *gin.Engine, conf *config.Config, cloudinary storage.Cloudinary) {
	router.Use(middleware.AllStat())
	setWebRouter(router, cloudinary)
	setApiRouter(router, conf)
	router.NoRoute(controller.Get404Page)
}
