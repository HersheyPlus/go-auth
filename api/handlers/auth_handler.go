package handlers

import (
	"net/http"
	"github.com/HersheyPlus/go-auth/config"
	"github.com/HersheyPlus/go-auth/dto"
	"github.com/HersheyPlus/go-auth/api/validators"
	"github.com/HersheyPlus/go-auth/models"
	"github.com/HersheyPlus/go-auth/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"time"
	"strings"
	"github.com/google/uuid"
)

type AuthHandler struct {
	DB  *gorm.DB
	Cfg *config.Config
	logger *log.Logger
	tokenStore *utils.TokenStore
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		DB:  db,
		Cfg: cfg,
		logger: log.New(log.Writer(), "AuthHandler: ", log.LstdFlags),
		tokenStore: utils.NewTokenStore(db),
	}
}


func (h *AuthHandler) Register(c *gin.Context) {
	rb := dto.NewResponse(c)
	var req dto.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		rb.Error(http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := validators.ValidateRegisterFields(&req, &h.Cfg.Security); err != nil {
		rb.Error(http.StatusBadRequest, err.Error())
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			h.logger.Printf("Recovered from panic: %v", r)
		}
	}()

	var count int64
	if err := tx.Model(&models.User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
		tx.Rollback()
		h.logger.Printf("Failed to check user existence: %v", err)
		rb.Error(http.StatusInternalServerError, "Failed to check user existence")
		return
	}

	if count > 0 {
		tx.Rollback()
		rb.Error(http.StatusConflict, "User with email already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password, &h.Cfg.Security)
	if err != nil {
		tx.Rollback()
		h.logger.Printf("Failed to hash password: %v", err)
		rb.Error(http.StatusBadRequest, err.Error())
		return
	}

	newUser := models.User{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  hashedPassword,
	}

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		h.logger.Printf("Failed to create user: %v", err)
		rb.Error(http.StatusInternalServerError, "Failed to register user")
		return
	}

	// Generate tokens for automatic login
    tokens, err := utils.GenerateTokenPair(newUser.UserID.String(), newUser.Username, &h.Cfg.JWT)
    if err != nil {
        tx.Rollback()
        h.logger.Printf("Failed to generate tokens: %v", err)
        rb.Error(http.StatusInternalServerError, "User registered but failed to generate login tokens")
        return
    }

    // Store refresh token
    if err := h.tokenStore.StoreToken(c.Request.Context(), newUser.UserID, tokens); err != nil {
        tx.Rollback()
        h.logger.Printf("Failed to store refresh token: %v", err)
        rb.Error(http.StatusInternalServerError, "Failed to complete registration process")
        return
    }
	
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		h.logger.Printf("Failed to commit transaction: %v", err)
		rb.Error(http.StatusInternalServerError, "Failed to complete registration")
		return
	}

	dataResponse := struct {
        User         dto.UserRegisterResponse `json:"user"`
        AccessToken  string                   `json:"access_token"`
        RefreshToken string                   `json:"refresh_token"`
        TokenType    string                   `json:"token_type"`
        ExpiresIn    int64                    `json:"expires_in"`
    }{
        User: dto.UserRegisterResponse{
            Username:  newUser.Username,
            FirstName: newUser.FirstName,
            LastName:  newUser.LastName,
            Phone:     newUser.Phone,
            Email:     newUser.Email,
			LastLogin: time.Now(),

        },
        AccessToken:  tokens.AccessToken,
        RefreshToken: tokens.RefreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    int64(h.Cfg.JWT.AccessTokenExpiry.Seconds()),
    }

	rb.Success(http.StatusCreated, dataResponse, "User registered successfully")
}


func (h *AuthHandler) Login(c *gin.Context) {
    rb := dto.NewResponse(c)
    var req dto.UserLoginRequest

    // Validate request binding
    if err := c.ShouldBindJSON(&req); err != nil {
        rb.ValidationError(http.StatusBadRequest, "Invalid request format", err.Error())
        return
    }

    // Start transaction
    tx := h.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            h.logger.Printf("Recovered from panic in Login: %v", r)
        }
    }()

    // Find user by email
    var user models.User
    if err := tx.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
        tx.Rollback()
        if err == gorm.ErrRecordNotFound {
            rb.Error(http.StatusUnauthorized, "Invalid email or password")
            return
        }
        h.logger.Printf("Database error during login: %v", err)
        rb.Error(http.StatusInternalServerError, "Failed to process login")
        return
    }

    // Verify password
    if err := utils.ComparePasswords(user.Password, req.Password); err != nil {
        tx.Rollback()
        h.logger.Printf("Failed password attempt for user %s", user.Email)
        rb.Error(http.StatusUnauthorized, "Invalid email or password")
        return
    }

    // Generate JWT token pair
    tokens, err := utils.GenerateTokenPair(user.UserID.String(), user.Username, &h.Cfg.JWT)
    if err != nil {
        tx.Rollback()
        h.logger.Printf("Failed to generate tokens: %v", err)
        rb.Error(http.StatusInternalServerError, "Failed to generate authentication tokens")
        return
    }

    // Store refresh token
    refreshToken := models.RefreshToken{
        UserID:    user.UserID,
        TokenUUID: tokens.RefreshUuid,
        ExpiresAt: tokens.RtExpires,
    }

    // Delete existing refresh tokens for user (optional, based on your requirements)
    if err := tx.Where("user_id = ?", user.UserID).Delete(&models.RefreshToken{}).Error; err != nil {
        tx.Rollback()
        h.logger.Printf("Failed to delete old refresh tokens: %v", err)
        rb.Error(http.StatusInternalServerError, "Failed to process login")
        return
    }

    // Save new refresh token
    if err := tx.Create(&refreshToken).Error; err != nil {
        tx.Rollback()
        h.logger.Printf("Failed to store refresh token: %v", err)
        rb.Error(http.StatusInternalServerError, "Failed to complete login process")
        return
    }

    // Update last login time
    if err := tx.Model(&user).Update("last_login", time.Now()).Error; err != nil {
        tx.Rollback()
        h.logger.Printf("Failed to update last login: %v", err)
        rb.Error(http.StatusInternalServerError, "Failed to update login information")
        return
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        h.logger.Printf("Failed to commit login transaction: %v", err)
        rb.Error(http.StatusInternalServerError, "Failed to complete login")
        return
    }

    // Prepare response
    response := dto.UserLoginResponse{
        UserID:       user.UserID,
        Username:     user.Username,
        Email:        user.Email,
        AccessToken:  tokens.AccessToken,
        RefreshToken: tokens.RefreshToken,
		LastLogin:   user.LastLogin,
        ExpiresIn:    int64(h.Cfg.JWT.AccessTokenExpiry.Seconds()),
    }

    rb.Success(http.StatusOK, response, "Login successful")
}

func (h *AuthHandler) Logout(c *gin.Context) {
    rb := dto.NewResponse(c)

    // Get token from Authorization header
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        rb.Error(http.StatusUnauthorized, "Authorization header is required")
        return
    }

    // Extract token from Bearer scheme
    tokenParts := strings.Split(authHeader, " ")
    if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
        rb.Error(http.StatusUnauthorized, "Invalid authorization header format")
        return
    }

    // Validate access token
    claims, err := utils.ValidateToken(tokenParts[1], h.Cfg.JWT.SecretKey)
    if err != nil {
        rb.Error(http.StatusUnauthorized, "Invalid or expired token")
        return
    }

    // Parse user ID
    userID, err := uuid.Parse(claims.Subject)
    if err != nil {
        h.logger.Printf("Failed to parse user ID from token: %v", err)
        rb.Error(http.StatusInternalServerError, "Invalid token data")
        return
    }

    // Start transaction
    tx := h.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            h.logger.Printf("Recovered from panic in Logout: %v", r)
        }
    }()

    // Delete refresh tokens for the user
    if err := tx.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error; err != nil {
        tx.Rollback()
        h.logger.Printf("Failed to delete refresh tokens: %v", err)
        rb.Error(http.StatusInternalServerError, "Failed to process logout")
        return
    }

    // Optional: Update last activity timestamp
    if err := tx.Model(&models.User{}).Where("user_id = ?", userID).
        Updates(map[string]interface{}{
            "updated_at": time.Now(),
            "last_login": time.Now(),
        }).Error; err != nil {
        tx.Rollback()
        h.logger.Printf("Failed to update user activity: %v", err)
        rb.Error(http.StatusInternalServerError, "Failed to update user information")
        return
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        h.logger.Printf("Failed to commit logout transaction: %v", err)
        rb.Error(http.StatusInternalServerError, "Failed to complete logout")
        return
    }

    rb.Success(http.StatusOK, nil, "Logged out successfully")
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
    rb := dto.NewResponse(c)
    
    // Get user ID from context (set by auth middleware)
    userID, exists := c.Get("userID")
    if !exists {
        rb.Error(http.StatusUnauthorized, "User not authenticated")
        return
    }

    var user models.User
    if err := h.DB.First(&user, "user_id = ?", userID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            rb.Error(http.StatusNotFound, "User not found")
            return
        }
        h.logger.Printf("Failed to fetch user profile: %v", err)
        rb.Error(http.StatusInternalServerError, "Failed to fetch user profile")
        return
    }


    response := dto.UserProfileResponse{
        UserID:    user.UserID,
        Username:  user.Username,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Phone:     user.Phone,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }

    rb.Success(http.StatusOK, response, "Profile retrieved successfully")
}