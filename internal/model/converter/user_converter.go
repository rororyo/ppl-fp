package converter

import (
	"fp-designpattern/internal/entity"
	"fp-designpattern/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:          &user.ID,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		GradeLevel:  user.GradeLevel,
		Role:        user.Role,
		AvatarUrl:   user.AvatarUrl,
		BirthDate:   &user.BirthDate,
		Token:       user.Token,
		CreatedAt:   &user.CreatedAt,
		UpdatedAt:   &user.UpdatedAt,
	}
}

func UserToTokenResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		Token: user.Token,
	}
}
