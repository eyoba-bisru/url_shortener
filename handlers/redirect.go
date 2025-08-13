package handlers

import (
	"net/http"

	"github.com/eyoba-bisru/url_shortener/config"
	"github.com/eyoba-bisru/url_shortener/models"
	"github.com/gin-gonic/gin"
)

func Redirect(c *gin.Context) {
	shortCode := c.Param("code")

	// 1. Try Redis first
	originalURL, err := config.RedisClient.Get(config.Ctx, shortCode).Result()
	if err == nil {
		c.Redirect(http.StatusFound, originalURL)
		return
	}

	// 2. Fallback to DB
	var url models.URL
	if err := config.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	// Cache it in Redis for next time
	config.RedisClient.Set(config.Ctx, shortCode, url.OriginalURL, cacheTTL)

	// Optional: increment visit count
	config.DB.Model(&url).Update("visit_count", url.VisitCount+1)

	c.Redirect(http.StatusFound, url.OriginalURL)
}
