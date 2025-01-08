package main

import (
	"fmt"
	"net/http"
	"time"

	gin "github.com/gin-gonic/gin"
	dotEnv "github.com/joho/godotenv"

	cache "github.com/brianwu291/go-learn/cache"
	postgres "github.com/brianwu291/go-learn/db/postgres"
	redis "github.com/brianwu291/go-learn/db/redis"
	utils "github.com/brianwu291/go-learn/utils"

	ratelimiter "github.com/brianwu291/go-learn/middlewares/ratelimiter"

	financialhandler "github.com/brianwu291/go-learn/handlers/financial"
	financialservice "github.com/brianwu291/go-learn/services/financial"

	fakestorehandler "github.com/brianwu291/go-learn/handlers/fakestore"
	fakestorerepo "github.com/brianwu291/go-learn/repos/fakestore"
	fakestoreservice "github.com/brianwu291/go-learn/services/fakestore"

	constants "github.com/brianwu291/go-learn/constants"
	types "github.com/brianwu291/go-learn/types"
)

var (
	StrictAPIConfig = ratelimiter.Config{
		Limit:                   100,
		Duration:                5 * time.Minute,
		ClientIdentifierOptions: []ratelimiter.ClientIdentifierOption{ratelimiter.ClientIP, ratelimiter.UserAgent},
	}

	NormalAPIConfig = ratelimiter.Config{
		Limit:                   1000,
		Duration:                15 * time.Minute,
		ClientIdentifierOptions: []ratelimiter.ClientIdentifierOption{ratelimiter.ClientIP},
	}

	PublicAPIConfig = ratelimiter.Config{
		Limit:                   5000,
		Duration:                time.Hour,
		ClientIdentifierOptions: []ratelimiter.ClientIdentifierOption{ratelimiter.ClientIP},
	}
)

func main() {
	err := dotEnv.Load()
	if err != nil {
		fmt.Printf("error loading .env file: %+v", err.Error())
		return
	}

	// Using %+v - shows field names
	// Using %#v - shows type information and field names

	postgresDB, dbInitErr := postgres.NewDatabase()
	if dbInitErr != nil {
		fmt.Printf("failed to initialize postgres database: %+v\n", dbInitErr)
		return
	}
	defer postgresDB.Close()

	redisDB := utils.GetEnvAsInt("REDIS_DB", 0)	
	cacheConfig := &cache.Config{
		Host:     utils.GetEnv("REDIS_HOST", "localhost"),
		Port:     utils.GetEnv("REDIS_PORT", "6379"),
		Password: utils.GetEnv("REDIS_PASSWORD", ""),
		DB:       redisDB,
	}
	cacheClient, err := redis.NewClient(cacheConfig)
	if err != nil {
		fmt.Printf("failed to initialize Redis cache: %v\n", err)
		return
	}

	// Initialize rate limiter
	rateLimiter := ratelimiter.NewRateLimiter(cacheClient)

	r := gin.Default()

	r.Use(postgresDB.Middleware())
	r.GET("/health", func(c *gin.Context) {
		if healthyErr := postgresDB.Health(); healthyErr != nil {
			fmt.Printf("error on postgres healthy check: %+v. error: %+v\n", healthyErr.Error(), healthyErr)
			internalServerErr := fmt.Errorf(constants.InternalServerErrorMessage)
			c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponse{Message: internalServerErr.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	financialService := financialservice.NewFinancialService()
	financialHandler := financialhandler.NewFinancialHandler(financialService)

	fakeStoreRepo := fakestorerepo.NewFakeStoreRepo()
	fakeStoreService := fakestoreservice.NewFakeStoreService(cacheClient, fakeStoreRepo)
	fakeStoreHandler := fakestorehandler.NewFakeStoreHandler(fakeStoreService)

	r.POST("/calculate",
		rateLimiter.LimitRoute(StrictAPIConfig),
		financialHandler.Calculate)

	r.GET("/fake-store/all/categories",
		rateLimiter.LimitRoute(NormalAPIConfig),
		fakeStoreHandler.GetAllCategories)

	r.GET("/fake-store/all/categories/products",
		rateLimiter.LimitRoute(PublicAPIConfig),
		fakeStoreHandler.GetAllCategoriesProducts)

	r.Run()
}
