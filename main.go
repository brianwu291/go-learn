package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	gin "github.com/gin-gonic/gin"
	dotEnv "github.com/joho/godotenv"

	cache "github.com/brianwu291/go-learn/cache"
	redis "github.com/brianwu291/go-learn/db/redis"

	ratelimiter "github.com/brianwu291/go-learn/middlewares/ratelimiter"

	financialhandler "github.com/brianwu291/go-learn/handlers/financial"
	financialservice "github.com/brianwu291/go-learn/services/financial"

	fakestorehandler "github.com/brianwu291/go-learn/handlers/fakestore"
	fakestorerepo "github.com/brianwu291/go-learn/repos/fakestore"
	fakestoreservice "github.com/brianwu291/go-learn/services/fakestore"
)

func main() {
	err := dotEnv.Load()
	if err != nil {
		fmt.Printf("error loading .env file: %+v", err.Error())
		return
	}

	// Using %+v - shows field names
	// Using %#v - shows type information and field names

	redisDBStr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		fmt.Printf("error on parsing REDIS_DB env value: %+v. REDIS_DB: %+v", err.Error(), redisDBStr)
		return
	}

	cacheConfig := &cache.Config{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	}
	cacheClient, err := redis.NewClient(cacheConfig)
	if err != nil {
		fmt.Printf("Failed to initialize Redis cache: %v", err)
		return
	}

	// Initialize rate limiter
	rateLimiter := ratelimiter.NewRateLimiter(cacheClient)

	r := gin.Default()
	financialService := financialservice.NewFinancialService()
	financialHandler := financialhandler.NewFinancialHandler(financialService)

	fakeStoreRepo := fakestorerepo.NewFakeStoreRepo()
	fakeStoreService := fakestoreservice.NewFakeStoreService(fakeStoreRepo)
	fakeStoreHandler := fakestorehandler.NewFakeStoreHandler(fakeStoreService)

	r.POST("/calculate",
		rateLimiter.LimitRoute(ratelimiter.Config{
			Limit:    100,
			Duration: time.Minute * 5,
		}),
		financialHandler.Calculate)

	r.GET("/fake-store/all/category",
		rateLimiter.LimitRoute(ratelimiter.Config{
			Limit:    50,
			Duration: time.Minute * 1,
		}),
		fakeStoreHandler.GetAllCategories)

	r.GET("/ping",
		rateLimiter.LimitRoute(ratelimiter.Config{
			Limit:    8,
			Duration: time.Second * 20,
		}),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

	r.Run()
}
