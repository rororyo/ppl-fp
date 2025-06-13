package repository

import (
	"fp-designpattern/internal/entity"

	"github.com/sirupsen/logrus"
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
