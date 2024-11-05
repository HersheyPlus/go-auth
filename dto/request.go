package dto

import (
	"github.com/google/uuid"
)

type UserRegisterRequest struct {
    Username  string  `json:"username" binding:"required,min=3,max=100"`
    FirstName *string `json:"first_name,omitempty" binding:"omitempty,min=2,max=100"`
    LastName  *string `json:"last_name,omitempty" binding:"omitempty,min=2,max=100"`
    Phone     string  `json:"phone" binding:"required,min=10,max=20"`
    Email     string  `json:"email" binding:"required,email,max=100"`
    Password  string  `json:"password" binding:"required,min=8,max=72"`
}

type UserLoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
    Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=100"`
    FirstName *string `json:"first_name,omitempty" binding:"omitempty,min=2,max=100"`
    LastName  *string `json:"last_name,omitempty" binding:"omitempty,min=2,max=100"`
    Phone     *string `json:"phone,omitempty" binding:"omitempty,min=10,max=20"`
    Email     *string `json:"email,omitempty" binding:"omitempty,email,max=100"`
    Password  *string `json:"password,omitempty" binding:"omitempty,min=4,max=72"`
}

type UserIDRequest struct {
    UserID uuid.UUID `json:"user_id" binding:"required,uuid"`
}