package repository

import (
	"blessing-be/internal/entity"
	"blessing-be/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserCourseRepository struct {
	Repository[entity.UserCourse]
	Log *logrus.Logger
}

func NewUserCourseRepository(log *logrus.Logger) *UserCourseRepository {
	return &UserCourseRepository{
		Log: log,
	}
}

func (r *UserCourseRepository) FindById(db *gorm.DB, userCourse *entity.UserCourse, id string) error {
	return db.Preload("Course").Preload("User").Where("id = ?", id).First(userCourse).Error
}
func (r *UserCourseRepository) FindByCourseIdAndUserId(db *gorm.DB, userCourse *entity.UserCourse, request *model.GetUserCourseRequest) error {
	return db.
		Preload("Course").
		Preload("User").
		Preload("Course.Subject").
		Where("course_id = ? AND user_id = ?", request.CourseID, request.UserID).
		First(userCourse).Error
}

func (r *UserCourseRepository) Search(db *gorm.DB, request *model.SearchUserCourseRequest) ([]entity.UserCourse, int64, error) {
	var userCourses []entity.UserCourse
	if err := db.
		Preload("Course").
		Preload("Course.Subject").
		Preload("User").
		Scopes(r.FilterUserCourse(request)).
		Order("accessed_at DESC").
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&userCourses).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.
		Model(&entity.UserCourse{}).
		Scopes(r.FilterUserCourse(request)).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return userCourses, total, nil
}

func (r *UserCourseRepository) FilterUserCourse(request *model.SearchUserCourseRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if userID := request.UserID; userID != "" {
			_, err := uuid.Parse(userID)
			if err == nil {
				tx = tx.Where("users_courses.user_id = ?", userID)
			}
		}

		if courseID := request.CourseID; courseID != "" {
			_, err := uuid.Parse(courseID)
			if err == nil {
				tx = tx.Where("users_courses.course_id = ?", courseID)
			}
		}
		if subjectID := request.SubjectID; subjectID != "" {
			_, err := uuid.Parse(subjectID)
			if err == nil {
				tx = tx.Joins("JOIN courses ON courses.id = users_courses.course_id").
					Where("courses.subject_id = ?", request.SubjectID)
			}
		}

		return tx
	}
}
