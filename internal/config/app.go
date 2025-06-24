package config

import (
	"fp-designpattern/internal/delivery/http"
	"fp-designpattern/internal/delivery/http/middleware"
	"fp-designpattern/internal/delivery/http/route"
	"fp-designpattern/internal/repository"
	"fp-designpattern/internal/usecase"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	//setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	subjectRepository := repository.NewSubjectRepository(config.Log)
	fileRepository := repository.NewLocalFileRepository(
		"./public/images/courses",
		"/images/courses",
	)
	courseRepository := repository.NewCourseRepository(config.Log)
	userCourseRepository := repository.NewUserCourseRepository(config.Log)
	//setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository)
	subjectUseCase := usecase.NewSubjectUsecase(config.DB, config.Log, config.Validate, subjectRepository)
	courseUseCase := usecase.NewCourseUsecase(config.DB, config.Log, config.Validate, courseRepository, subjectRepository, fileRepository)
	userCourseUseCase := usecase.NewUserCourseUsecase(config.DB, config.Log, config.Validate, courseRepository, userRepository, userCourseRepository)
	//setup controllers
	userController := http.NewUserController(userUseCase, courseUseCase, config.Log)
	subjectController := http.NewSubjectController(subjectUseCase, config.Log)
	courseController := http.NewCourseController(courseUseCase, config.Log)
	userCourseController := http.NewUserCourseController(userCourseUseCase, config.Log)
	//setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)
	routeConfig := route.RouteConfig{
		App:                  config.App,
		UserController:       userController,
		SubjectController:    subjectController,
		CourseController:     courseController,
		UserCourseController: userCourseController,
		AuthMiddleware:       authMiddleware,
	}

	routeConfig.Setup()
}
