package usecase

import (
	"context"
	"encoding/json"
	"fp-designpattern/internal/entity"
	"fp-designpattern/internal/model"
	"fp-designpattern/internal/model/converter"
	"fp-designpattern/internal/repository"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourseUsecase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	CourseRepository  *repository.CourseRepository
	SubjectRepository *repository.SubjectRepository
}

func NewCourseUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, courseRepository *repository.CourseRepository, subjectRepository *repository.SubjectRepository) *CourseUsecase {
	return &CourseUsecase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		CourseRepository:  courseRepository,
		SubjectRepository: subjectRepository,
	}
}

func (c *CourseUsecase) Create(ctx context.Context, request *model.CourseRequest) (*model.CourseResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		c.Log.Warnf("Failed to start transaction: %+v", tx.Error)
		return nil, fiber.ErrInternalServerError
	}
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	contentJSON, err := json.Marshal(request.Content)
	if err != nil {
		c.Log.Warnf("Failed to marshal content: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	subject := new(entity.Subject)
	if err := c.SubjectRepository.FindById(tx, subject, request.SubjectID); err != nil {
		c.Log.Warnf("Failed find subject by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	course := &entity.Course{
		CourseName: request.CourseName,
		Content:    contentJSON,
		GradeLevel: request.GradeLevel,
		SubjectID:  uuid.MustParse(request.SubjectID),
	}

	if err := c.CourseRepository.Create(tx, course); err != nil {
		c.Log.Warnf("Failed to create subject: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CourseToResponse(course), nil
}
