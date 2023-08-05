package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jordyf15/tweeter-api/middlewares"
	"github.com/jordyf15/tweeter-api/token/repository"
	"github.com/jordyf15/tweeter-api/token/usecase"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	router      *gin.Engine
	db          *gorm.DB
	redisClient *redis.Client
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	connectToDB()
	connectToRedis()
}

func main() {
	router = gin.Default()

	if allowedOriginsEnvValue := os.Getenv("ALLOWED_ORIGINS"); len(allowedOriginsEnvValue) > 0 {
		allowedOrigins := strings.Split(allowedOriginsEnvValue, ",")
		config := cors.DefaultConfig()
		config.AllowOrigins = allowedOrigins
		config.AllowHeaders = []string{"Origin", "Authorization"}

		router.Use(cors.New(config))
	}

	tokenRepo := repository.NewTokenRepository(db, redisClient)
	authMiddleware := middlewares.NewAuthMiddleware(usecase.NewTokenUsecase(tokenRepo))

	router.Use(authMiddleware.AuthenticateJWT)

	router.MaxMultipartMemory = 10 << 20
	initializeRoutes()
	if os.Getenv("ROUTER_PORT") != "" {
		router.Run(fmt.Sprintf(":%s", os.Getenv("ROUTER_PORT")))
	} else {
		router.Run()
	}
}

func connectToDB() {
	config := &gorm.Config{}
	if schemaStr := os.Getenv("DB_SCHEMA"); len(schemaStr) > 0 {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix: schemaStr + ".",
		}
	}

	dbURL := os.Getenv("DB_URL")
	conn, err := gorm.Open(postgres.Open(dbURL), config)
	sleepTime := 10 * time.Second
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Connection failed, attempting to reconnect in %v seconds\n", sleepTime)
		time.Sleep(sleepTime)
		fmt.Println("Attempting to reconnect now")

		conn, _ = gorm.Open(postgres.Open(dbURL), config)
		if sleepTime < time.Minute {
			sleepTime += 10 * time.Second
		}
	}

	db = conn
}

func connectToRedis() {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})
}
