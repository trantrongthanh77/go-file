package router

import (
	"go-file/common/config"
	"go-file/controller"
	"go-file/middleware"

	"github.com/gin-gonic/gin"
)

func SetRouter(router *gin.Engine, conf *config.Config) {
	router.Use(middleware.AllStat())
	setWebRouter(router)
	setApiRouter(router, conf)
	router.NoRoute(controller.Get404Page)
}
