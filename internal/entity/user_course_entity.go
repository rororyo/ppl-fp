package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserCourse struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID `gorm:"column:user_id;not null;type:uuid;"`
	CourseID   uuid.UUID `gorm:"column:course_id;not null;type:uuid;"`
	AccessedAt time.Time `gorm:"column:accessed_at"`
	//Foreign Key
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Course Course `gorm:"foreignKey:CourseID;references:ID;constraint:OnDelete:CASCADE"`
}

func (UserCourse) TableName() string {
	return "users_courses"
}
