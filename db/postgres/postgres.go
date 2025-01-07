package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	utils "github.com/brianwu291/go-learn/utils"
)

const (
	defaultTimeout = 10 * time.Second
	dbKey          = "db"
)

type Database struct {
	Pool *pgxpool.Pool
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	Pool     struct {
		MaxConns int32
		MinConns int32
		MaxIdle  time.Duration
	}
}

func LoadConfigFromEnv() Config {
	maxPoolConns := utils.GetEnvAsInt("DB_MAX_POOL_CONS", 200)

	cfg := Config{
		Host:     utils.GetEnv("DB_HOST", "postgres"),
		Port:     utils.GetEnvAsInt("DB_PORT", 5432),
		User:     utils.GetEnv("DB_USER", ""),
		Password: utils.GetEnv("DB_PASSWORD", ""),
		DBName:   utils.GetEnv("DB_NAME", ""),
		SSLMode:  utils.GetEnv("DB_SSL_MODE", "disable"),
		Pool: struct {
			MaxConns int32
			MinConns int32
			MaxIdle  time.Duration
		}{
			MaxConns: int32(maxPoolConns),
			MinConns: int32(float64(maxPoolConns) * 0.05), // 5% of max
			MaxIdle:  15 * time.Minute,
		},
	}

	return cfg
}

func New(cfg Config) (*Database, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	// set pool configuration
	poolCfg.MaxConns = cfg.Pool.MaxConns
	poolCfg.MinConns = cfg.Pool.MinConns
	poolCfg.MaxConnIdleTime = cfg.Pool.MaxIdle

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func (db *Database) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	return db.Pool.Ping(ctx)
}

// adds the database to the context
func (db *Database) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(dbKey, db)
		c.Next()
	}
}

// retrieve the db from gin context
func FromContext(c *gin.Context) *Database {
	db, exists := c.Get(dbKey)
	if !exists {
		return nil
	}
	return db.(*Database)
}
