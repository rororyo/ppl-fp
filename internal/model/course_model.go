package model

import (
	"time"

	"github.com/google/uuid"
)

type ContentBlock struct {
	Type string `json:"type"` // "text" or "image"
	Data string `json:"data"` // text content or image URL
}

type CourseResponse struct {
	ID         uuid.UUID       `json:"id"`
	CourseName string          `json:"course_name"`
	Content    []ContentBlock  `json:"content"`
	GradeLevel int             `json:"grade_level"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	Subject    SubjectResponse `json:"subject"`
}
type CourseListResponse struct {
	ID         uuid.UUID       `json:"id"`
	CourseName string          `json:"course_name"`
	GradeLevel int             `json:"grade_level"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	Subject    SubjectResponse `json:"subject"`
}

type CourseRequest struct {
	CourseName string         `json:"course_name"`
	Content    []ContentBlock `json:"content"`
	GradeLevel int            `json:"grade_level"`
	SubjectID  string         `json:"subject_id"`
}

type GetCourseRequest struct {
	ID string `json:"-"`
}
type DeleteCourseRequest struct {
	ID string `json:"-"`
}
type SearchCourseRequest struct {
	CourseName string `json:"course_name"`
	SubjectID  string `json:"subject_id"`
	GradeLevel int    `json:"grade_level"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
}

type UpdateCourseRequest struct {
	ID         string         `json:"-"`
	CourseName string         `json:"course_name"`
	Content    []ContentBlock `json:"content"`
	GradeLevel int            `json:"grade_level"`
	SubjectID  string         `json:"subject_id"`
}
