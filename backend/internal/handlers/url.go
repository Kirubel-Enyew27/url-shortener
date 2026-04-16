package handlers

import (
	"net/http"
	"strings"
	"url-shortener/internal/config"
	"url-shortener/internal/models"
	"url-shortener/internal/services"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	service *services.URLService
	config  *config.Config
}

func New(service *services.URLService, cfg *config.Config) *URLHandler {
	return &URLHandler{service: service, config: cfg}
}

func (h *URLHandler) Shorten(c *gin.Context) {
	var req models.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	url, err := h.service.Shorten(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ShortenResponse{
		ShortURL: h.shortURLBase(c) + "/" + url.ShortCode,
		Code:     url.ShortCode,
	})
}

func (h *URLHandler) GetAll(c *gin.Context) {
	urls := h.service.GetAll()
	c.JSON(http.StatusOK, models.URLListResponse{URLs: convertToSlice(urls)})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("code")

	longURL, err := h.service.Resolve(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.Redirect(http.StatusFound, longURL)
}

func convertToSlice(urls []*models.URL) []models.URL {
	result := make([]models.URL, len(urls))
	for i, url := range urls {
		result[i] = *url
	}
	return result
}

func (h *URLHandler) shortURLBase(c *gin.Context) string {
	if h.config != nil && h.config.BaseURL != "" {
		return h.config.BaseURL
	}

	scheme := c.GetHeader("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}

	host := c.GetHeader("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}

	return strings.TrimRight(scheme+"://"+host, "/")
}
