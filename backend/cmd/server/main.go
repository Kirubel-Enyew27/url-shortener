package main

import (
	"fmt"
	"url-shortener/internal/config"
	"url-shortener/internal/handlers"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	store := storage.New()
	service := services.New(store)
	handler := handlers.New(service)

	r := gin.Default()

	r.POST("/api/shorten", handler.Shorten)
	r.GET("/api/urls", handler.GetAll)
	r.GET("/:code", handler.Redirect)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	r.Run(addr)
}
