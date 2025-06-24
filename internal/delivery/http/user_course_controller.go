package http

import (
	"fp-designpattern/internal/delivery/http/middleware"
	"fp-designpattern/internal/model"
	"fp-designpattern/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserCourseController struct {
	Log     *logrus.Logger
	Usecase *usecase.UserCourseUsecase
}

func NewUserCourseController(usecase *usecase.UserCourseUsecase, logger *logrus.Logger) *UserCourseController {
	return &UserCourseController{
		Log:     logger,
		Usecase: usecase,
	}
}
func (c *UserCourseController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.GetUserCourseRequest{
		CourseID: ctx.Params("id"),
	}
	request.UserID = auth.ID
	courseResponse, err := c.Usecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get subject: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.UserCourseResponse]{Data: courseResponse})
}
func (c *UserCourseController) Create(ctx *fiber.Ctx) error {
	request := new(model.UserCourseRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}
	userCourseRepsonse, err := c.Usecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create subject: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[[]*model.UserCourseResponse]{Data: userCourseRepsonse})
}

func (c *UserCourseController) List(ctx *fiber.Ctx) error {

	request := &model.SearchUserCourseRequest{
		UserID:   ctx.Query("user_id"),
		CourseID: ctx.Query("course_id"),
		Page:     ctx.QueryInt("page"),
		Size:     ctx.QueryInt("size"),
	}

	responses, total, err := c.Usecase.Search(ctx.UserContext(), request)
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

	return ctx.JSON(model.WebResponse[[]model.UserCourseListResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *UserCourseController) ListAccessable(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	authRequest := &model.GetUserRequest{
		ID: auth.ID,
	}
	request := &model.SearchUserCourseRequest{
		CourseID:  ctx.Query("course_id"),
		SubjectID: ctx.Query("subject_id"),
		Page:      ctx.QueryInt("page"),
		Size:      ctx.QueryInt("size"),
	}
	request.UserID = authRequest.ID

	responses, total, err := c.Usecase.Search(ctx.UserContext(), request)
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

	return ctx.JSON(model.WebResponse[[]model.UserCourseListResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *UserCourseController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteUserCourseRequest{
		ID: ctx.Params("id"),
	}
	userCourseResponse, err := c.Usecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete user course: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.UserCourseResponse]{Data: userCourseResponse})
}
