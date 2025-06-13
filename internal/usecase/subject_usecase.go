package usecase

import (
	"context"
	"fp-designpattern/internal/entity"
	"fp-designpattern/internal/model"
	"fp-designpattern/internal/model/converter"
	"fp-designpattern/internal/repository"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SubjectUsecase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	SubjectRepository *repository.SubjectRepository
}

func NewSubjectUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, subjectRepository *repository.SubjectRepository) *SubjectUsecase {
	return &SubjectUsecase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		SubjectRepository: subjectRepository,
	}
}

func (c *SubjectUsecase) Create(ctx context.Context, request *model.SubjectRequest) (*model.SubjectResponse, error) {
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
	subject := &entity.Subject{
		SubjectName: request.SubjectName,
	}

	if err := c.SubjectRepository.Create(tx, subject); err != nil {
		c.Log.Warnf("Failed to create subject: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SubjectToResponse(subject), nil
}
func (c *SubjectUsecase) Get(ctx context.Context, request *model.GetSubjectRequest) (*model.SubjectResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	subject := new(entity.Subject)
	if err := c.SubjectRepository.FindById(tx, subject, request.ID); err != nil {
		c.Log.Warnf("Failed find subject by id : %+v", err)
		return nil, fiber.ErrNotFound
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SubjectToResponse(subject), nil

}

func (c *SubjectUsecase) Update(ctx context.Context, request *model.UpdateSubjectRequest) (*model.SubjectResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	subject := new(entity.Subject)
	if err := c.SubjectRepository.FindById(tx, subject, request.ID); err != nil {
		c.Log.Warnf("Failed find subject by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if request.SubjectName != "" {
		subject.SubjectName = request.SubjectName
	}

	if err := c.SubjectRepository.Update(tx, subject); err != nil {
		c.Log.Warnf("Failed to update subject: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SubjectToResponse(subject), nil
}

func (c *SubjectUsecase) Search(ctx context.Context, request *model.SearchSubjectRequest) ([]model.SubjectResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warnf("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}
	subjects, total, err := c.SubjectRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search subject")
		return nil, 0, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.SubjectResponse, len(subjects))
	for i, subject := range subjects {
		responses[i] = *converter.SubjectToResponse(&subject)
	}
	return responses, total, nil
}
func (c *SubjectUsecase) Delete(ctx context.Context, request *model.DeleteSubjectRequest) (*model.SubjectResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find subject by id
	subject := new(entity.Subject)
	if err := c.SubjectRepository.FindById(tx, subject, request.ID); err != nil {
		c.Log.Warnf("Failed find subject by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Delete subject
	if err := c.SubjectRepository.Delete(tx, subject); err != nil {
		c.Log.Warnf("Failed delete subject : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SubjectToResponse(subject), nil
}
