package utils

import (
	"errors"
	"fmt"
	"github.com/HersheyPlus/go-auth/config"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

var (
	ErrPasswordTooShort   = errors.New("password must be at least %d characters long")
	ErrPasswordTooLong    = errors.New("password must not exceed %d characters")
	ErrPasswordComplexity = errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	ErrEmptyPassword      = errors.New("password is required")
)


func HashPassword(password string, cfg *config.SecurityConfig) (string, error) {

	if cfg.BCryptCost < bcrypt.MinCost || cfg.BCryptCost > bcrypt.MaxCost {
		cfg.BCryptCost = bcrypt.DefaultCost
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cfg.BCryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

func ComparePasswords(hashedPassword string, plainPassword string) error {
	if hashedPassword == "" || plainPassword == "" {
		return ErrEmptyPassword
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return errors.New("invalid password")
		}
		return fmt.Errorf("error comparing passwords: %w", err)
	}

	return nil
}

func ValidatePassword(password string, cfg *config.SecurityConfig) error {
	if password == "" {
		return ErrEmptyPassword
	}
	
	if len(password) < cfg.MinPasswordLength {
		return fmt.Errorf(ErrPasswordTooShort.Error(), cfg.MinPasswordLength)
	}
	
	if len(password) > cfg.MaxPasswordLength {
		return fmt.Errorf(ErrPasswordTooLong.Error(), cfg.MaxPasswordLength)
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if cfg.PasswordRequirements.RequireUppercase && !hasUpper ||
		cfg.PasswordRequirements.RequireLowercase && !hasLower ||
		cfg.PasswordRequirements.RequireNumbers && !hasNumber ||
		cfg.PasswordRequirements.RequireSpecial && !hasSpecial {
		return ErrPasswordComplexity
	}

	return nil
}
