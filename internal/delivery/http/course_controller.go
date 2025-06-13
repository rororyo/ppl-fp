package http

import (
	"fp-designpattern/internal/model"
	"fp-designpattern/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CourseController struct {
	Log     *logrus.Logger
	Usecase *usecase.CourseUsecase
}

func NewCourseController(usecase *usecase.CourseUsecase, logger *logrus.Logger) *CourseController {
	return &CourseController{
		Log:     logger,
		Usecase: usecase,
	}
}
func (c *CourseController) Create(ctx *fiber.Ctx) error {
	request := new(model.CourseRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}
	courseRepsonse, err := c.Usecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create subject: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.CourseResponse]{Data: courseRepsonse})
}
