package main

import (
	"fmt"
	"url-shortener/internal/config"
	"url-shortener/internal/handlers"
	"url-shortener/internal/middleware"
	"url-shortener/internal/routers"
	"url-shortener/internal/services"
	"url-shortener/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	store := storage.New()
	service := services.New(store)
	handler := handlers.New(service, cfg)

	r := gin.Default()
	r.Use(middleware.CorsMiddleware(cfg.AllowedHosts))

	routers.SetupRoutes(r, handler)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	r.Run(addr)
}

