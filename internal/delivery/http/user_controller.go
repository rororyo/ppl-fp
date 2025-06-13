package http

import (
	"fp-designpattern/internal/delivery/http/middleware"
	"fp-designpattern/internal/model"
	"fp-designpattern/internal/usecase"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log           *logrus.Logger
	UserUsecase   *usecase.UserUseCase
	CourseUsecase *usecase.CourseUsecase
}

func NewUserController(userUsecase *usecase.UserUseCase, courseUsecase *usecase.CourseUsecase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log:           logger,
		UserUsecase:   userUsecase,
		CourseUsecase: courseUsecase,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	err := ctx.BodyParser(request)

	if err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UserUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create user: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UserUsecase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to login user: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetUserRequest{
		ID: auth.ID,
	}

	response, err := c.UserUsecase.Current(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get current user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Get(ctx *fiber.Ctx) error {
	request := &model.GetUserRequest{
		ID: ctx.Params("id"),
	}
	response, err := c.UserUsecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to get current user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.LogoutUserRequest{
		ID: auth.ID,
	}

	response, err := c.UserUsecase.Logout(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to logout user")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: response})
}

func (c *UserController) List(ctx *fiber.Ctx) error {
	var birthDate *time.Time
	if birthDateStr := ctx.Query("birth_date"); birthDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", birthDateStr)
		if err != nil {
			c.Log.WithError(err).Warn("Invalid birth_date format")
			return fiber.NewError(fiber.StatusBadRequest, "Invalid birth_date format, expected YYYY-MM-DD")
		}
		birthDate = &parsedDate
	}

	request := &model.SearchUserRequest{
		Username:    ctx.Query("username"),
		Email:       ctx.Query("email"),
		PhoneNumber: ctx.Query("phone_number"),
		GradeLevel:  ctx.QueryInt("grade_level"),
		BirthDate:   birthDate,
		Page:        ctx.QueryInt("page"),
		Size:        ctx.QueryInt("size"),
	}

	responses, total, err := c.UserUsecase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search user")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.UserResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *UserController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	user := &model.LogoutUserRequest{
		ID: auth.ID,
	}
	request := new(model.UpdateUserRequest)
	request.ID = user.ID
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UserUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update user: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) AdminUpdate(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	request := new(model.UpdateUserRequest)
	request.ID = userID
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UserUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update user: %v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteUserRequest{
		ID: ctx.Params("id"),
	}
	response, err := c.UserUsecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to delete user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}
