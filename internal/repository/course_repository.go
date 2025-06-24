package repository

import (
	"fp-designpattern/internal/entity"
	"fp-designpattern/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourseRepository struct {
	Repository[entity.Course]
	Log *logrus.Logger
}

func NewCourseRepository(log *logrus.Logger) *CourseRepository {
	return &CourseRepository{
		Log: log,
	}
}

func (r *CourseRepository) FindById(db *gorm.DB, course *entity.Course, id string) error {
	return db.Preload("Subject").Where("id = ?", id).First(course).Error
}

func (r *CourseRepository) Search(db *gorm.DB, request *model.SearchCourseRequest) ([]entity.Course, int64, error) {
	// Query the actual data
	var courses []entity.Course
	if err := db.Preload("Subject").
		Scopes(r.FilterCourse(request)).
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	// Count total without pagination
	var total int64
	if err := db.Model(&entity.Course{}).
		Scopes(r.FilterCourse(request)).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

func (r *CourseRepository) FilterCourse(request *model.SearchCourseRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if courseName := request.CourseName; courseName != "" {
			tx = tx.Where("course_name LIKE ?", "%"+courseName+"%")
		}
		if subjectID := request.SubjectID; subjectID != "" {
			tx = tx.Where("subject_id = ?", subjectID)
		}
		if gradeLevel := request.GradeLevel; gradeLevel != 0 {
			tx = tx.Where("grade_level = ?", gradeLevel)
		}
		return tx
	}
}
