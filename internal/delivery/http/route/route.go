package route

import (
	"fp-designpattern/internal/delivery/http"
	"fp-designpattern/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                  *fiber.App
	UserController       *http.UserController
	SubjectController    *http.SubjectController
	CourseController     *http.CourseController
	UserCourseController *http.UserCourseController
	AuthMiddleware       fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	// users
	c.App.Post("/api/users/register", c.UserController.Register)
	c.App.Post("/api/users/login", c.UserController.Login)
	c.App.Get("/api/users/user/:id", c.UserController.Get)

	//subjects
	c.App.Get("api/subjects", c.SubjectController.List)
	c.App.Get("api/subjects/:id", c.SubjectController.Get)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)
	//authenticated users
	c.App.Get("/api/users/current", c.UserController.Current)
	c.App.Post("/api/users/logout", c.UserController.Logout)
	c.App.Put("api/users", c.UserController.Update)

	// accessable courses
	c.App.Get("/api/courses", c.UserCourseController.ListAccessable)
	c.App.Get("/api/courses/:id", c.UserCourseController.Get)

	// Admin-only
	adminOnly := c.App.Group("/api/admin", middleware.RequireRole("admin"))
	// users
	adminOnly.Get("/users", c.UserController.List)
	adminOnly.Put("/users/:id", c.UserController.AdminUpdate)
	adminOnly.Delete("/users/:id", c.UserController.Delete)
	// subjects
	adminOnly.Post("/subjects", c.SubjectController.Create)
	adminOnly.Put("/subjects/:id", c.SubjectController.Update)
	adminOnly.Delete("/subjects/:id", c.SubjectController.Delete)

	// courses
	adminOnly.Get("/courses", c.CourseController.List)
	adminOnly.Get("/courses/:id", c.CourseController.Get)
	adminOnly.Post("/courses/upload", c.CourseController.UploadFile)
	adminOnly.Post("/courses", c.CourseController.Create)
	adminOnly.Put("/courses/:id", c.CourseController.Update)
	adminOnly.Delete("/courses/:id", c.CourseController.Delete)

	// course permissions
	adminOnly.Get("/user-courses", c.UserCourseController.List)
	adminOnly.Post("/user-courses", c.UserCourseController.Create)
	adminOnly.Delete("/user-courses/:id", c.UserCourseController.Delete)

}
