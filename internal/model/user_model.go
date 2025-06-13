package model

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	Username    string     `json:"username,omitempty"`
	Email       string     `json:"email,omitempty"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	GradeLevel  int        `json:"grade_level,omitempty"`
	Role        string     `json:"role,omitempty" validate:"required,oneof=admin user"`
	AvatarUrl   string     `json:"avatar_url,omitempty"`
	BirthDate   *time.Time `json:"birth_date,omitempty"`
	Token       string     `json:"token,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=255" json:"token"`
}

type RegisterUserRequest struct {
	Email      string `json:"email" validate:"required,max=100"`
	Password   string `json:"password" validate:"required,max=100"`
	Username   string `json:"username" validate:"required,max=100"`
	GradeLevel string `json:"grade_level" validate:"required,max=100"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LogoutUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type GetUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type SearchUserRequest struct {
	Username    string     `json:"username,omitempty"`
	Email       string     `json:"email,omitempty"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	GradeLevel  int        `json:"grade_level,omitempty"`
	BirthDate   *time.Time `json:"birth_date,omitempty"`
	Page        int        `json:"page,omitempty" validate:"min=1"`
	Size        int        `json:"size,omitempty" validate:"min=1,max=100"`
}

type UpdateUserRequest struct {
	ID          string     `json:"id" validate:"required,max=100"`
	Username    string     `json:"username,omitempty"`
	Email       string     `json:"email,omitempty"`
	Password    string     `json:"password,omitempty"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	GradeLevel  string     `json:"grade_level,omitempty"`
	BirthDate   *time.Time `json:"birth_date,omitempty"`
	AvatarUrl   string     `json:"avatar_url,omitempty"`
}

type DeleteUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
