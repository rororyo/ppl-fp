package model

import (
	"time"

	"github.com/google/uuid"
)

type UserCourseResponse struct {
	ID         uuid.UUID      `json:"id"`
	User       UserResponse   `json:"user"`
	Course     CourseResponse `json:"course"`
	AccessedAt time.Time      `json:"accessed_at"`
}

type UserCourseListResponse struct {
	ID         uuid.UUID          `json:"id"`
	User       UserResponse       `json:"user"`
	Course     CourseListResponse `json:"course"`
	AccessedAt time.Time          `json:"accessed_at"`
}
type UserCourseRequest struct {
	UserID    string   `json:"user_id"`
	CourseIDs []string `json:"course_ids"`
}

type SearchUserCourseRequest struct {
	UserID     string    `json:"user_id"`
	CourseID   string    `json:"course_id"`
	SubjectID  string    `json:"subject_id"`
	AccessedAt time.Time `json:"accessed_at"`
	Page       int       `json:"page,omitempty" validate:"min=1"`
	Size       int       `json:"size,omitempty" validate:"min=1,max=100"`
}
type GetUserCourseRequest struct {
	CourseID string `json:"-" validate:"required,max=100"`
	UserID   string `json:"-" validate:"required,max=100"`
}

type DeleteUserCourseRequest struct {
	ID string `json:"-" validate:"required,max=100"`
}
