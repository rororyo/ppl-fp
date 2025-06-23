package usecase

import (
	"blessing-be/internal/entity"
	"blessing-be/internal/model"
	"blessing-be/internal/model/converter"
	"blessing-be/internal/repository"
	"blessing-be/pkg/timezone"
	"context"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserCourseUsecase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	CourseRepository     *repository.CourseRepository
	UserRepository       *repository.UserRepository
	UserCourseRepository *repository.UserCourseRepository
}

func NewUserCourseUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, courseRepository *repository.CourseRepository, userRepository *repository.UserRepository, userCourseRepository *repository.UserCourseRepository) *UserCourseUsecase {
	return &UserCourseUsecase{
		DB:                   db,
		Log:                  log,
		Validate:             validate,
		CourseRepository:     courseRepository,
		UserRepository:       userRepository,
		UserCourseRepository: userCourseRepository,
	}
}

func (c *UserCourseUsecase) Create(ctx context.Context, request *model.UserCourseRequest) ([]*model.UserCourseResponse, error) {
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

	// Check if CourseIDs is empty
	if len(request.CourseIDs) == 0 {
		c.Log.Warn("No courses provided in request")
		return []*model.UserCourseResponse{}, nil // Return empty slice but without error
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.UserID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	userCourses := make([]*entity.UserCourse, len(request.CourseIDs))
	for i, courseID := range request.CourseIDs {
		course := new(entity.Course)
		if err := c.CourseRepository.FindById(tx, course, courseID); err != nil {
			c.Log.Warnf("Course not found: %s, %+v", courseID, err)
			return nil, fiber.ErrNotFound
		}

		userCourses[i] = &entity.UserCourse{
			UserID:   uuid.MustParse(request.UserID),
			CourseID: uuid.MustParse(courseID),
		}
	}

	if err := c.UserCourseRepository.CreateBatch(tx, userCourses); err != nil {
		c.Log.Warnf("Failed to bulk create user courses: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	responses := make([]*model.UserCourseResponse, len(userCourses))
	for i, userCourse := range userCourses {
		responses[i] = converter.UserCourseToResponse(userCourse)
	}
	return responses, nil
}
func (c *UserCourseUsecase) Get(ctx context.Context, request *model.GetUserCourseRequest) (*model.UserCourseResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find user course by course id and user id
	userCourse := new(entity.UserCourse)
	if err := c.UserCourseRepository.FindByCourseIdAndUserId(tx, userCourse, request); err != nil {
		c.Log.Warnf("Failed to find user course: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Update AccessedAt to current WIB time
	userCourse.AccessedAt = time.Now().In(timezone.WIB)
	if err := tx.Save(userCourse).Error; err != nil {
		c.Log.Warnf("Failed to update AccessedAt: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserCourseToResponse(userCourse), nil
}

func (c *UserCourseUsecase) Search(ctx context.Context, request *model.SearchUserCourseRequest) ([]model.UserCourseListResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warnf("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}
	userCourses, total, err := c.UserCourseRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search user course")
		return nil, 0, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.UserCourseListResponse, len(userCourses))
	for i, userCourse := range userCourses {
		responses[i] = *converter.UserCourseListToResponse(&userCourse)
	}
	return responses, total, nil
}

func (c *UserCourseUsecase) Delete(ctx context.Context, request *model.DeleteUserCourseRequest) (*model.UserCourseResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find user course by id
	userCourse := new(entity.UserCourse)
	if err := c.UserCourseRepository.FindById(tx, userCourse, request.ID); err != nil {
		c.Log.Warnf("Failed find user course by id : %+v", err)
		return nil, fiber.ErrNotFound
	}
	// Delete course
	if err := c.UserCourseRepository.Delete(tx, userCourse); err != nil {
		c.Log.Warnf("Failed delete user course : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserCourseToResponse(userCourse), nil
}
