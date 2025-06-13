package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username    string    `gorm:"column:username;"`
	Email       string    `gorm:"column:email;"`
	Password    string    `gorm:"column:password;"`
	PhoneNumber string    `gorm:"column:phone_number;"`
	GradeLevel  int       `gorm:"column:grade_level;"`
	Role        string    `gorm:"column:role;"`
	AvatarUrl   string    `gorm:"column:avatar_url;"`
	BirthDate   time.Time `gorm:"column:birth_date;type:date"`
	Token       string    `gorm:"column:token"`
	CreatedAt   time.Time `gorm:"column:created_at;"`
	UpdatedAt   time.Time `gorm:"column:updated_at;"`
}
