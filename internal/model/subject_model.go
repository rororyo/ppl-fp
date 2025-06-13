package model

import (
	"time"

	"github.com/google/uuid"
)

type SubjectResponse struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	SubjectName string     `json:"subject_name,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type SubjectRequest struct {
	SubjectName string `json:"subject_name,omitempty" validate:"required"`
}

type GetSubjectRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type SearchSubjectRequest struct {
	SubjectName string `json:"subject_name,omitempty"`
	Page        int    `json:"page,omitempty" validate:"min=1"`
	Size        int    `json:"size,omitempty" validate:"min=1,max=100"`
}

type UpdateSubjectRequest struct {
	ID          string `json:"-,omitempty" validate:"required"`
	SubjectName string `json:"subject_name,omitempty"`
}

type DeleteSubjectRequest struct {
	ID string `json:"-,omitempty" validate:"required"`
}
