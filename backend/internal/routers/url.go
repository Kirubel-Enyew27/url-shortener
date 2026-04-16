package routers

import (
	"url-shortener/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handler *handlers.URLHandler) {
	router.POST("/api/shorten", handler.Shorten)
	router.GET("/api/urls", handler.GetAll)
	router.GET("/:code", handler.Redirect)
}
