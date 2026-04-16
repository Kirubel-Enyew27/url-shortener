package main

import (
	"fmt"
	"net/http"
	"strings"
	"url-shortener/internal/config"
	"url-shortener/internal/handlers"
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
	r.Use(corsMiddleware(cfg.AllowedHosts))

	routers.SetupRoutes(r, handler)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	r.Run(addr)
}

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	allowAll := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowAll {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if origin != "" && isAllowedOrigin(origin, allowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if strings.EqualFold(origin, allowed) {
			return true
		}
	}
	return false
}
