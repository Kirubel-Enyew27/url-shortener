package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware(allowedOrigns []string) gin.HandlerFunc {
	allowAll := len(allowedOrigns) == 1 && allowedOrigns[0] == "*"

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowAll {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Vary", "Origin")
		} else if origin != "" && isAllowedOrigin(origin, allowedOrigns) {
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
