package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	utils "github.com/brianwu291/go-learn/utils"
)

const (
	dbKey = "db"
)

type Database struct {
	*gorm.DB
}

func NewDatabase() (*Database, error) {
	maxConns := utils.GetEnvAsInt("DB_MAX_POOL_CONS", 200)

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		utils.GetEnv("DB_USER", ""),
		utils.GetEnv("DB_PASSWORD", ""),
		utils.GetEnv("DB_HOST", "localhost"),
		utils.GetEnv("DB_PORT", "5432"),
		utils.GetEnv("DB_NAME", ""),
	)

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second * 30,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, connectErr := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if connectErr != nil {
		return nil, fmt.Errorf("failed to connect to database: %+v", connectErr)
	}

	sqlDB, getDBErr := db.DB()
	if getDBErr != nil {
		return nil, fmt.Errorf("failed to get database instance: %+v", getDBErr)
	}

	sqlDB.SetMaxOpenConns(maxConns)
	sqlDB.SetMaxIdleConns(int(float64(maxConns) * 0.05))
	sqlDB.SetConnMaxLifetime(15 * time.Minute)

	return &Database{db}, nil
}

// close the database connection
func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// health checks the database connection
func (db *Database) Health() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// add the database to the context
func (db *Database) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(dbKey, db)
		c.Next()
	}
}

// retrieve the database from the Gin context
func FromContext(c *gin.Context) (*Database, error) {
	db, exists := c.Get(dbKey)
	if !exists {
		return nil, fmt.Errorf("database not found in context")
	}

	if dbInstance, ok := db.(*Database); ok {
		return dbInstance, nil
	}
	return nil, fmt.Errorf("invalid database type in context")
}
