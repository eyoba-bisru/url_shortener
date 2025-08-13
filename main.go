package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/eyoba-bisru/url_shortener/config"
	"github.com/eyoba-bisru/url_shortener/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// func redirect(c *gin.Context) {
// 	code := c.Param("code")

// 	longURL, exists := urlStore[code]

// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
// 		return
// 	}

// 	c.Redirect(http.StatusFound, longURL)
// }

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	rand.New(rand.NewSource(time.Now().UnixNano()))

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/shorten", handlers.ShortenURL)

	r.Run() // listen and serve on 0.0.0.0:8080
}
