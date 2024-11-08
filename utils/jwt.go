package utils

import (
    "errors"
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
	"github.com/HersheyPlus/go-auth/config"
)

// Custom claims structure
type Claims struct {
    UserID    string   `json:"user_id"`
    Username  string `json:"username"`
    TokenType string `json:"token_type"` // "access" or "refresh"
    jwt.RegisteredClaims
}

type TokenDetails struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    AccessUuid   string    `json:"access_uuid"`
    RefreshUuid  string    `json:"refresh_uuid"`
    AtExpires    time.Time `json:"at_expires"`
    RtExpires    time.Time `json:"rt_expires"`
}

// GenerateTokenPair generates both access and refresh tokens
func GenerateTokenPair(userID string, username string, cfg *config.JWTConfig) (*TokenDetails, error) {
    td := &TokenDetails{
        AccessUuid:  GenerateUUID(),
        RefreshUuid: GenerateUUID(),
        AtExpires:   time.Now().Add(cfg.AccessTokenExpiry),
        RtExpires:   time.Now().Add(cfg.RefreshTokenExpiry),
    }

    // Generate Access Token
    accessToken, err := generateToken(
        userID,
        username,
        td.AccessUuid,
        "access",
        td.AtExpires,
        cfg.SecretKey,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to generate access token: %w", err)
    }
    td.AccessToken = accessToken

    // Generate Refresh Token
    refreshToken, err := generateToken(
        userID,
        username,
        td.RefreshUuid,
        "refresh",
        td.RtExpires,
        cfg.RefreshKey,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to generate refresh token: %w", err)
    }
    td.RefreshToken = refreshToken

    return td, nil
}

// generateToken creates a new token depending on token type
func generateToken(
    userID string,
    username string,
    uuid string,
    tokenType string,
    expiry time.Time,
    secret string,
) (string, error) {
    claims := Claims{
        UserID:    userID,
        Username:  username,
        TokenType: tokenType,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiry),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Subject:   userID,
            ID:        uuid,
            Issuer:    "your-app-name",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", fmt.Errorf("failed to sign token: %w", err)
    }

    return signedToken, nil
}

// ValidateToken validates the token and returns the claims
func ValidateToken(tokenString string, secret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, fmt.Errorf("token has expired")
        }
        return nil, fmt.Errorf("failed to parse token: %w", err)
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token claims")
}

// ExtractTokenMetadata extracts metadata from token
func ExtractTokenMetadata(tokenString string, secret string) (*Claims, error) {
    claims, err := ValidateToken(tokenString, secret)
    if err != nil {
        return nil, err
    }

    // Validate token type
    if claims.TokenType == "" {
        return nil, fmt.Errorf("invalid token type")
    }

    return claims, nil
}

// GenerateUUID generates a unique identifier for tokens
func GenerateUUID() string {
    return fmt.Sprintf("%d%s", time.Now().UnixNano(), RandomString(8))
}

// RandomString generates a random string of specified length
func RandomString(n int) string {
    const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
    }
    return string(b)
}