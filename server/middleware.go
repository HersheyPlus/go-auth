package server

import (
    "log"
    "net/http"
    "sync"
    "time"
	"github.com/HersheyPlus/go-auth/config"
    "github.com/gin-gonic/gin"
	"strings"
	"strconv"
)

// Simple in-memory rate limiter
type RateLimiterStore struct {
    sync.RWMutex
    requests map[string]int
    lastReset map[string]time.Time
}

var limiterStore = &RateLimiterStore{
    requests:  make(map[string]int),
    lastReset: make(map[string]time.Time),
}

func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method

        // Process request
        c.Next()

        // Log details
        latency := time.Since(start)
        statusCode := c.Writer.Status()
        clientIP := c.ClientIP()

        log.Printf("[%d] %v | %s | %s %s",
            statusCode,
            latency,
            clientIP,
            method,
            path,
        )
    }
}

func RateLimiter(cfg config.RateLimitConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        clientIP := c.ClientIP()

        limiterStore.Lock()
        defer limiterStore.Unlock()

        now := time.Now()
        lastReset, exists := limiterStore.lastReset[clientIP]

        // Reset counter if duration has passed
        if !exists || now.Sub(lastReset) > cfg.Duration {
            limiterStore.requests[clientIP] = 0
            limiterStore.lastReset[clientIP] = now
        }

        // Check rate limit
        if limiterStore.requests[clientIP] >= cfg.Requests {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
            })
            c.Abort()
            return
        }

        // Increment request counter
        limiterStore.requests[clientIP]++
        c.Next()
    }
}

func CORSMiddleware(cfg config.CORSConfig) gin.HandlerFunc {
    // Pre-process configuration for optimization
    allowedOriginsMap := make(map[string]bool)
    for _, origin := range cfg.AllowedOrigins {
        allowedOriginsMap[origin] = true
    }
    
    allowedMethodsStr := strings.Join(cfg.AllowedMethods, ", ")
    allowedHeadersStr := strings.Join(cfg.AllowedHeaders, ", ")
    exposedHeadersStr := strings.Join(cfg.ExposedHeaders, ", ")
    maxAgeStr := strconv.Itoa(cfg.MaxAge)
    
    // Wildcard check
    allowAllOrigins := len(cfg.AllowedOrigins) == 0 || (len(cfg.AllowedOrigins) == 1 && cfg.AllowedOrigins[0] == "*")

    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // Skip CORS for same-origin requests
        if origin == "" {
            c.Next()
            return
        }

        // Check if origin is allowed
        if !allowAllOrigins {
            if _, allowed := allowedOriginsMap[origin]; !allowed {
                c.Next()
                return
            }
        }

        // Set CORS headers
        headers := c.Writer.Header()
        
        if allowAllOrigins {
            headers.Set("Access-Control-Allow-Origin", "*")
        } else {
            headers.Set("Access-Control-Allow-Origin", origin)
            if cfg.AllowCredentials {
                headers.Set("Vary", "Origin")
            }
        }

        // Only set credentials header if not using wildcard origin
        if cfg.AllowCredentials && !allowAllOrigins {
            headers.Set("Access-Control-Allow-Credentials", "true")
        }

        // Set other CORS headers
        if allowedMethodsStr != "" {
            headers.Set("Access-Control-Allow-Methods", allowedMethodsStr)
        }
        
        if allowedHeadersStr != "" {
            headers.Set("Access-Control-Allow-Headers", allowedHeadersStr)
        }
        
        if exposedHeadersStr != "" {
            headers.Set("Access-Control-Expose-Headers", exposedHeadersStr)
        }

        // Handle preflight requests
        if c.Request.Method == "OPTIONS" {
            if cfg.MaxAge > 0 {
                headers.Set("Access-Control-Max-Age", maxAgeStr)
            }
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}