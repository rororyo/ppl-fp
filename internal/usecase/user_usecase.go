package usecase

import (
	"context"
	"fp-designpattern/internal/entity"
	"fp-designpattern/internal/model"
	"fp-designpattern/internal/model/converter"
	"fp-designpattern/internal/repository"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (c *UserUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByToken(tx, user, request.Token); err != nil {
		c.Log.Warnf("Failed find user by token : %+v", err)
		return nil, fiber.ErrNotFound
	}

	auth := &model.Auth{
		ID:   user.ID.String(),
		Role: user.Role,
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return auth, nil
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("validation error: %v", err)
		return nil, fiber.ErrBadRequest
	}
	total, err := c.UserRepository.CountByEmail(tx, request.Email)
	if err != nil {
		c.Log.Warnf("Failed to count by email: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if total > 0 {
		c.Log.Warnf("email already exist")
		return nil, fiber.ErrConflict
	}
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to hash password: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	gradeInt, err := strconv.Atoi(request.GradeLevel)
	if err != nil {
		c.Log.Warnf("Failed to convert string to int: %v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "grade_level must be a number")
	}

	user := &entity.User{
		Email:      request.Email,
		Username:   request.Username,
		Password:   string(password),
		GradeLevel: gradeInt,
		Role:       "user",
	}
	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed to create user to database: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		c.Log.Warnf("Failed find user by email : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Failed compare password : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	user.Token = uuid.New().String()
	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed update user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToTokenResponse(user), nil
}

func (c *UserUseCase) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil

}

func (c *UserUseCase) Get(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return false, fiber.ErrNotFound
	}
	user.Token = ""

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed update user : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}

func (c *UserUseCase) Search(ctx context.Context, request *model.SearchUserRequest) ([]model.UserResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warnf("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}
	users, total, err := c.UserRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search user")
		return nil, 0, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *converter.UserToResponse(&user)
	}
	return responses, total, nil
}

func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find user by id
	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Update user
	if request.Username != "" {
		user.Username = request.Username
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warnf("Failed to hash password : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
		user.Password = string(hashedPassword)
	}
	if request.PhoneNumber != "" {
		user.PhoneNumber = request.PhoneNumber
	}
	if request.GradeLevel != "" {
		gradeInt, err := strconv.Atoi(request.GradeLevel)
		if err != nil {
			c.Log.Warnf("Failed to convert string to int: %v", err)
			return nil, fiber.NewError(fiber.StatusBadRequest, "grade_level must be a number")
		}
		user.GradeLevel = gradeInt
	}
	if request.BirthDate != nil {
		user.BirthDate = *request.BirthDate
	}

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed update user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Delete(ctx context.Context, request *model.DeleteUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find user by id
	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Delete user
	if err := c.UserRepository.Delete(tx, user); err != nil {
		c.Log.Warnf("Failed delete user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}
