package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/eyoba-bisru/url_shortener/backend/config"
	"github.com/eyoba-bisru/url_shortener/backend/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	rand.New(rand.NewSource(time.Now().UnixNano()))

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5173"} // Specify allowed origins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour // Cache preflight requests for 12 hours

	r := gin.Default()

	r.Use(cors.New(corsConfig))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/shorten", handlers.ShortenURL)
	r.GET("/:code", handlers.Redirect)

	r.Run() // listen and serve on 0.0.0.0:8080
}
