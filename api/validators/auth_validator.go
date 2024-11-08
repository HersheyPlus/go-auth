package validators

import (
	"fmt"
	"github.com/HersheyPlus/go-auth/config"
	"github.com/HersheyPlus/go-auth/dto"
	"github.com/HersheyPlus/go-auth/utils"
)

func ValidateRegisterFields(req *dto.UserRegisterRequest, cfg *config.SecurityConfig) error {
	// Validate basic required fields
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Phone == "" {
		return fmt.Errorf("phone is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}

	// Validate password requirements
	if err := utils.ValidatePassword(req.Password, cfg); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}