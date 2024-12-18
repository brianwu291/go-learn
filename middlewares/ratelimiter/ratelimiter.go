package rateLimiter

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/brianwu291/go-learn/cache"
	"github.com/gin-gonic/gin"
)

const (
	maxDuration = 15 * time.Minute

	rateLimitScript = `
    local key = KEYS[1]
    local limit = tonumber(ARGV[1])
    local duration = tonumber(ARGV[2])
    
    -- Get current value first
    local current = tonumber(redis.call('GET', key) or "0")
    
    -- Check if would exceed limit
    if current >= limit then
        local ttl = redis.call('TTL', key)
        return {current, ttl, 0}  -- return without incrementing
    end
    
    -- If we get here, we're under limit, safe to increment
    current = redis.call('INCR', key)
    
    -- Set expiry for new keys
    if current == 1 then
        redis.call('EXPIRE', key, duration)
    end
    
    local ttl = redis.call('TTL', key)
    return {current, ttl, 1}
  `
)

type (
	RateLimiter struct {
		cacheClient cache.Client
	}

	Config struct {
		Limit    int64
		Duration time.Duration
	}
)

func NewRateLimiter(cacheClient cache.Client) *RateLimiter {
	return &RateLimiter{
		cacheClient: cacheClient,
	}
}

func getSafePath(c *gin.Context) string {
	path := c.FullPath()
	if path == "" {
		// Get the raw path and remove query string
		path = c.Request.URL.EscapedPath()
		// Additional cleaning if needed
		path = strings.TrimSuffix(path, "/")
		if path == "" {
			path = "/"
		}
	}
	return path
}

func (rl *RateLimiter) formatKey(path, method, clientIdentifier string) string {
	return fmt.Sprintf("ratelimit:%s:%s:%s", path, method, clientIdentifier)
}

func (rl *RateLimiter) LimitRoute(config Config) gin.HandlerFunc {
	if config.Duration > maxDuration {
		config.Duration = maxDuration
	}

	return func(c *gin.Context) {
		path := getSafePath(c)
		key := rl.formatKey(path, c.Request.Method, c.ClientIP())

		// Execute atomic Lua script
		result, err := rl.cacheClient.Eval(
			c,
			rateLimitScript,
			[]string{key},
			[]interface{}{config.Limit, int(config.Duration.Seconds())},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limiting error",
				"msg":   err.Error(),
			})
			c.Abort()
			return
		}

		// Parse results from Lua script
		results := result.([]interface{})
		current := results[0].(int64)
		ttl := time.Duration(results[1].(int64)) * time.Second
		allowed := false
		if results[2].(int64) == 1 {
			allowed = true
		}

		// Calculate remaining requests
		remaining := config.Limit - current
		if remaining < 0 {
			remaining = 0
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		// Handle rate limit exceeded
		if !allowed {
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%.0f", ttl.Seconds()))
			c.Header("Retry-After", fmt.Sprintf("%.0f", ttl.Seconds()))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"retry_after": fmt.Sprintf("%.0f secs", ttl.Seconds()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
