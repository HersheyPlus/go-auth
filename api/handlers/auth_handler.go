package handlers

import (
	"net/http"
	"github.com/HersheyPlus/go-auth/config"
	"github.com/HersheyPlus/go-auth/dto"
	"github.com/HersheyPlus/go-auth/models"
	"github.com/HersheyPlus/go-auth/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB  *gorm.DB
	Cfg *config.JWTConfig
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		DB:  db,
		Cfg: &cfg.JWT,
	}
}


func (h *AuthHandler) Register(c *gin.Context) {
	rb := dto.NewResponse(c)
	var req dto.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		rb.Error(http.StatusBadRequest, "Invalid request format")
		return
	}


	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		rb.Error(http.StatusInternalServerError, "Failed to hash password")
		return
	}

	newUser := models.User{
		Username: req.Username,
		FirstName: req.FirstName,
		LastName: req.LastName,
		Email:	req.Email,
		Phone: req.Phone,
		Password: hashedPassword,
	}

	if err := h.DB.Create(&newUser).Error; err != nil {
		rb.Error(http.StatusInternalServerError, "Failed to register user")
		return
	}
	
	dataReponse := dto.UserRegisterResponse{
		Username: newUser.Username,
		FirstName: newUser.FirstName,
		LastName: newUser.LastName,
		Phone: newUser.Phone,
		Email: newUser.Email,
	}

	rb.Success(http.StatusCreated, dataReponse, "User registered successfully")




}
