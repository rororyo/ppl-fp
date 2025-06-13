package entity

import (
	"time"

	"github.com/google/uuid"
)

type Subject struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SubjectName string    `gorm:"column:subject_name;"`
	CreatedAt   time.Time `gorm:"column:created_at;"`
	UpdatedAt   time.Time `gorm:"column:updated_at;"`
}
