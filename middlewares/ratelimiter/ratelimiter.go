package rateLimiter

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/brianwu291/go-learn/cache"
	"github.com/gin-gonic/gin"
)

type (
	RateLimiter struct {
		cacheClient cache.Client
	}

	ClientIdentifierOption string
	Config                 struct {
		Limit                   int64
		Duration                time.Duration
		ClientIdentifierOptions []ClientIdentifierOption
	}
)

const (
	ClientIP        ClientIdentifierOption = "ClientIP"
	UserAgent       ClientIdentifierOption = "UserAgent"
	maxDuration                            = 15 * time.Minute
	defaultLimit                           = 100
	defaultDuration                        = 1 * time.Minute

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

var (
	ErrInvalidConfig = errors.New("invalid rate limit configuration")
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

func (rl *RateLimiter) formatKey(path string, method string, clientIdentifiers []string) string {
	return fmt.Sprintf("ratelimit:%s:%s:%s", path, method, strings.Join(clientIdentifiers, ":"))
}

func (rl *RateLimiter) getClientIdentifiers(c *gin.Context, identifierOptions []ClientIdentifierOption) []string {
	result := make([]string, 0, len(identifierOptions))
	for _, option := range identifierOptions {
		var identifier string
		switch option {
		case ClientIP:
			identifier = c.ClientIP()
		case UserAgent:
			identifier = c.Request.UserAgent()
		default:
			fmt.Printf("unsupported identifier option: %s", option)
			continue
		}

		if identifier != "" {
			result = append(result, identifier)
		}
	}
	if len(result) == 0 {
		// at least one client identifier with IP
		result = append(result, c.ClientIP())
	}
	return result
}

func (rl *RateLimiter) LimitRoute(config Config) gin.HandlerFunc {
	if err := config.validate(); err != nil {
		panic(err)
	}

	if config.Duration > maxDuration {
		config.Duration = maxDuration
	}

	return func(c *gin.Context) {
		path := getSafePath(c)
		clientIdentifiers := rl.getClientIdentifiers(c, config.ClientIdentifierOptions)
		key := rl.formatKey(path, c.Request.Method, clientIdentifiers)

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

func (config Config) validate() error {
	if config.Limit <= 0 {
		return fmt.Errorf("%w: limit must be positive", ErrInvalidConfig)
	}
	if config.Duration <= 0 {
		return fmt.Errorf("%w: duration must be positive", ErrInvalidConfig)
	}
	if len(config.ClientIdentifierOptions) == 0 {
		return fmt.Errorf("%w: at least one identifier option required", ErrInvalidConfig)
	}
	return nil
}
