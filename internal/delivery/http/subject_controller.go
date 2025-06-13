package http

import (
	"fp-designpattern/internal/model"
	"fp-designpattern/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type SubjectController struct {
	Log     *logrus.Logger
	Usecase *usecase.SubjectUsecase
}

func NewSubjectController(usecase *usecase.SubjectUsecase, logger *logrus.Logger) *SubjectController {
	return &SubjectController{
		Log:     logger,
		Usecase: usecase,
	}
}

func (c *SubjectController) Create(ctx *fiber.Ctx) error {
	request := new(model.SubjectRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}
	subjectResponse, err := c.Usecase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create subject: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.SubjectResponse]{Data: subjectResponse})
}

func (c *SubjectController) Get(ctx *fiber.Ctx) error {
	request := &model.GetSubjectRequest{
		ID: ctx.Params("id"),
	}
	subjectResponse, err := c.Usecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get subject: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.SubjectResponse]{Data: subjectResponse})
}
func (c *SubjectController) List(ctx *fiber.Ctx) error {

	request := &model.SearchSubjectRequest{
		SubjectName: ctx.Query("subject_name"),
		Page:        ctx.QueryInt("page"),
		Size:        ctx.QueryInt("size"),
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

	return ctx.JSON(model.WebResponse[[]model.SubjectResponse]{
		Data:   responses,
		Paging: paging,
	})
}
func (c *SubjectController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateSubjectRequest)
	request.ID = ctx.Params("id")
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}
	subjectResponse, err := c.Usecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update subject: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.SubjectResponse]{Data: subjectResponse})
}

func (c *SubjectController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteSubjectRequest{
		ID: ctx.Params("id"),
	}
	subjectResponse, err := c.Usecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete subject: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.SubjectResponse]{Data: subjectResponse})
}
