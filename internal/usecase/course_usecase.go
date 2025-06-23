package usecase

import (
	"blessing-be/internal/entity"
	"blessing-be/internal/model"
	"blessing-be/internal/model/converter"
	"blessing-be/internal/repository"
	"context"
	"encoding/json"
	"mime/multipart"

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
	FileRepository    *repository.GCSFileRepository
}

func NewCourseUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, courseRepository *repository.CourseRepository, subjectRepository *repository.SubjectRepository, fileRepository *repository.GCSFileRepository) *CourseUsecase {
	return &CourseUsecase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		CourseRepository:  courseRepository,
		SubjectRepository: subjectRepository,
		FileRepository:    fileRepository,
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

func (c *CourseUsecase) Get(ctx context.Context, request *model.GetCourseRequest) (*model.CourseResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	course := new(entity.Course)
	if err := c.CourseRepository.FindById(tx, course, request.ID); err != nil {
		c.Log.Warnf("Failed find subject by id : %+v", err)
		return nil, fiber.ErrNotFound
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CourseToResponse(course), nil

}

func (c *CourseUsecase) Search(ctx context.Context, request *model.SearchCourseRequest) ([]model.CourseListResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warnf("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}
	courses, total, err := c.CourseRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search subject")
		return nil, 0, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.CourseListResponse, len(courses))
	for i, course := range courses {
		responses[i] = *converter.CourseToListResponse(&course)
	}
	return responses, total, nil
}

func (c *CourseUsecase) Update(ctx context.Context, request *model.UpdateCourseRequest) (*model.CourseResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	course := new(entity.Course)
	if err := c.CourseRepository.FindById(tx, course, request.ID); err != nil {
		c.Log.Warnf("Failed find subject by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if request.CourseName != "" {
		course.CourseName = request.CourseName
	}

	if request.Content != nil {
		contentJSON, err := json.Marshal(request.Content)
		if err != nil {
			c.Log.Warnf("Failed to marshal content: %+v", err)
			return nil, fiber.ErrInternalServerError
		}
		course.Content = contentJSON
	}

	if request.GradeLevel != 0 {
		course.GradeLevel = request.GradeLevel
	}

	if request.SubjectID != "" {
		subject := new(entity.Subject)
		if err := c.SubjectRepository.FindById(tx, subject, request.SubjectID); err != nil {
			c.Log.Warnf("Failed find subject by id : %+v", err)
			return nil, fiber.ErrNotFound
		}
		course.SubjectID = subject.ID
	}

	if err := c.CourseRepository.Update(tx, course); err != nil {
		c.Log.Warnf("Failed to update subject: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CourseToResponse(course), nil
}

func (c *CourseUsecase) Delete(ctx context.Context, request *model.DeleteCourseRequest) (*model.CourseResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find course by id
	course := new(entity.Course)
	if err := c.CourseRepository.FindById(tx, course, request.ID); err != nil {
		c.Log.Warnf("Failed find subject by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Delete course
	if err := c.CourseRepository.Delete(tx, course); err != nil {
		c.Log.Warnf("Failed delete subject : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CourseToResponse(course), nil
}

func (u *CourseUsecase) UploadFile(ctx context.Context, file multipart.File, fileName string, contentType string) (string, error) {
	return u.FileRepository.UploadFile(file, fileName, contentType)
}
