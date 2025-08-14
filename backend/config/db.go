package config

import (
	"context"
	"log"
	"os"

	"github.com/eyoba-bisru/url_shortener/backend/models"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func ConnectDB() {

	dsn := os.Getenv("DATABASE_URL")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Initialize Redis connection
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	DB.AutoMigrate(&models.URL{})
}
