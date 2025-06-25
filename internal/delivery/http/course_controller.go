package http

import (
	"fmt"
	"fp-designpattern/internal/model"
	"fp-designpattern/internal/usecase"
	"fp-designpattern/pkg/timezone"
	"math"
	"time"

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

func (c *CourseController) Get(ctx *fiber.Ctx) error {
	request := &model.GetCourseRequest{
		ID: ctx.Params("id"),
	}
	courseResponse, err := c.Usecase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to get subject: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.CourseResponse]{Data: courseResponse})
}

func (c *CourseController) List(ctx *fiber.Ctx) error {

	request := &model.SearchCourseRequest{
		CourseName: ctx.Query("course_name"),
		GradeLevel: ctx.QueryInt("grade_level"),
		SubjectID:  ctx.Query("subject_id"),
		Page:       ctx.QueryInt("page"),
		Size:       ctx.QueryInt("size"),
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

	return ctx.JSON(model.WebResponse[[]model.CourseListResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *CourseController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateCourseRequest)
	request.ID = ctx.Params("id")
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body: %v", err)
		return fiber.ErrBadRequest
	}
	courseResponse, err := c.Usecase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to update subject: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.CourseResponse]{Data: courseResponse})
}

func (c *CourseController) Delete(ctx *fiber.Ctx) error {
	request := &model.DeleteCourseRequest{
		ID: ctx.Params("id"),
	}
	courseResponse, err := c.Usecase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to delete subject: %v", err)
		return err
	}
	return ctx.JSON(model.WebResponse[*model.CourseResponse]{Data: courseResponse})
}

func (c *CourseController) UploadFile(ctx *fiber.Ctx) error {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		c.Log.Warnf("Failed to get file: %v", err)
		return fiber.ErrBadRequest
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.Log.Warnf("Failed to open file: %v", err)
		return fiber.ErrBadRequest
	}
	defer file.Close()

	now := time.Now().In(timezone.WIB)
	timestampPrefix := now.Format("20060102_150405")

	fileName := fmt.Sprintf("%s_%s", timestampPrefix, fileHeader.Filename)

	// Upload to folder "courses/"
	objectPath := fmt.Sprintf("courses/%s", fileName)

	url, err := c.Usecase.UploadFile(
		ctx.UserContext(),
		file,
		objectPath,
		fileHeader.Header.Get("Content-Type"),
	)
	if err != nil {
		c.Log.Warnf("Failed to upload file: %v", err)
		return fiber.ErrInternalServerError
	}

	// Return public URL
	return ctx.JSON(model.WebResponse[string]{Data: url})
}
