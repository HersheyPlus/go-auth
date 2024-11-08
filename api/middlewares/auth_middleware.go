package middlewares

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/HersheyPlus/go-auth/dto"
    "github.com/HersheyPlus/go-auth/utils"
    "github.com/HersheyPlus/go-auth/config"
)

// AuthMiddleware verifies JWT tokens in the Authorization header
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        rb := dto.NewResponse(c)

        // Get token from Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            rb.Error(http.StatusUnauthorized, "Authorization header is required")
            c.Abort()
            return
        }

        // Check Bearer scheme
        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
            rb.Error(http.StatusUnauthorized, "Invalid authorization header format")
            c.Abort()
            return
        }

        // Validate token
        claims, err := utils.ValidateToken(tokenParts[1], cfg.JWT.SecretKey)
        if err != nil {
            rb.Error(http.StatusUnauthorized, "Invalid or expired token")
            c.Abort()
            return
        }

        // Store user info in context
        c.Set("userID", claims.Subject)
        c.Set("username", claims.Username)
        c.Next()
    }
}


