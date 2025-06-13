package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Course struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CourseName string         `gorm:"column:course_name;not null"`
	Content    datatypes.JSON `gorm:"column:content;type:jsonb;not null"`
	GradeLevel int            `gorm:"column:grade_level;not null"`
	CreatedAt  time.Time      `gorm:"column:created_at;default:now()"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;default:now()"`
	SubjectID  uuid.UUID      `gorm:"column:subject_id;not null;type:uuid"`
	//Foreign Key
	Subject Subject `gorm:"foreignKey:SubjectID;references:ID;constraint:OnDelete:CASCADE"`
}
