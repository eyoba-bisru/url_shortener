package handlers

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/eyoba-bisru/url_shortener/backend/config"
	"github.com/eyoba-bisru/url_shortener/backend/models"
	"github.com/gin-gonic/gin"
)

const base62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode(length int) string {
	code := make([]byte, length)
	for i := range code {
		code[i] = base62[rand.Intn(len(base62))]
	}
	return string(code)
}

const cacheTTL = 24 * time.Hour

func ShortenURL(c *gin.Context) {
	shortURLPrefix := os.Getenv("BASE_URL")

	var request struct {
		URL string `json:"url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	var existing models.URL
	if err := config.DB.Where("original_url = ?", request.URL).First(&existing).Error; err == nil {
		config.RedisClient.Set(config.Ctx, existing.ShortCode, existing.OriginalURL, cacheTTL)
		c.JSON(http.StatusOK, gin.H{"short_url": shortURLPrefix + existing.ShortCode})
		return
	}

	// Generate a unique code
	shortCode := generateShortCode(6)

	// Check if it is unique; if not retry
	for {
		var temp models.URL
		if err := config.DB.Where("short_code = ?", shortCode).First(&temp).Error; err != nil {
			break
		}
		shortCode = generateShortCode(6)
	}

	newURL := models.URL{
		OriginalURL: request.URL,
		ShortCode:   shortCode,
	}

	if err := config.DB.Create(&newURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL"})
		return
	}

	// Save to Redis
	config.RedisClient.Set(config.Ctx, shortCode, request.URL, cacheTTL)

	c.JSON(http.StatusOK, gin.H{
		"short_url": shortURLPrefix + shortCode,
	})
}
