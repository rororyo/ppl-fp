package config

import (
	"blessing-be/internal/delivery/http"
	"blessing-be/internal/delivery/http/middleware"
	"blessing-be/internal/delivery/http/route"
	"blessing-be/internal/repository"
	"blessing-be/internal/usecase"

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
	courseRepository := repository.NewCourseRepository(config.Log)
	userCourseRepository := repository.NewUserCourseRepository(config.Log)
	// ctx := context.Background()
	// gcsClient, err := storage.NewClient(ctx)
	// if err != nil {
	// 	config.Log.Fatalf("Failed to create GCS client: %v", err)
	// }
	// bucketName := config.Config.GetString("GCS_BUCKET_NAME")
	// if bucketName == "" {
	// 	config.Log.Fatal("GCS_BUCKET_NAME not set in config")
	// }

	// fileRepository := repository.NewGCSFileRepository(gcsClient, bucketName)
	fileRepository := repository.NewGCSFileRepository(nil, "")
	if fileRepository == nil {
		config.Log.Fatal("Failed to create GCS file repository")
	}
	quizRepository := repository.NewQuizRepository(config.Log)
	questionRepository := repository.NewQuestionRepository(config.Log)
	questionOptionRepository := repository.NewQuestionOptionRepository(config.Log)
	quizAnswerRepository := repository.NewQuizAnswerRepository(config.Log)
	userQuizSessionRepository := repository.NewUserQuizSessionRepository(config.Log)
	userAnswerRepository := repository.NewUserAnswerRepository(config.Log)

	//setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository)
	subjectUseCase := usecase.NewSubjectUsecase(config.DB, config.Log, config.Validate, subjectRepository)
	courseUseCase := usecase.NewCourseUsecase(config.DB, config.Log, config.Validate, courseRepository, subjectRepository, fileRepository)
	userCourseUseCase := usecase.NewUserCourseUsecase(config.DB, config.Log, config.Validate, courseRepository, userRepository, userCourseRepository)
	quizUseCase := usecase.NewQuizUsecase(config.DB, config.Log, config.Validate, quizRepository, courseRepository)
	questionUseCase := usecase.NewQuestionUsecase(config.DB, config.Log, config.Validate, questionRepository, quizRepository)
	questionOptionUseCase := usecase.NewQuestionOptionUsecase(config.DB, config.Log, config.Validate, questionOptionRepository, questionRepository)
	quizAnswerUseCase := usecase.NewQuizAnswerUsecase(config.DB, config.Log, config.Validate, quizAnswerRepository, questionOptionRepository)
	userQuizSessionUseCase := usecase.NewUserQuizSessionUsecase(config.DB, config.Log, config.Validate, userQuizSessionRepository, userRepository, quizRepository)
	userAnswerUseCase := usecase.NewUserAnswerUsecase(config.DB, config.Log, config.Validate, userAnswerRepository, userQuizSessionRepository, questionOptionRepository, questionRepository)

	//setup controllers
	userController := http.NewUserController(userUseCase, courseUseCase, userCourseUseCase, config.Log)
	subjectController := http.NewSubjectController(subjectUseCase, config.Log)
	courseController := http.NewCourseController(courseUseCase, config.Log)
	userCourseController := http.NewUserCourseController(userCourseUseCase, config.Log)
	quizController := http.NewQuizController(quizUseCase, config.Log)
	questionController := http.NewQuestionController(questionUseCase, config.Log)
	questionOptionController := http.NewQuestionOptionController(questionOptionUseCase, config.Log)
	quizAnswerController := http.NewQuizAnswerController(quizAnswerUseCase, config.Log)
	userQuizSessionController := http.NewUserQuizSessionController(userQuizSessionUseCase, config.Log)
	userAnswerController := http.NewUserAnswerController(userAnswerUseCase, config.Log)

	//setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)
	routeConfig := route.RouteConfig{
		App:                       config.App,
		UserController:            userController,
		SubjectController:         subjectController,
		CourseController:          courseController,
		UserCourseController:      userCourseController,
		QuizController:            quizController,
		QuestionController:        questionController,
		QuestionOptionController:  questionOptionController,
		QuizAnswerController:      quizAnswerController,
		UserQuizSessionController: userQuizSessionController,
		UserAnswerController:      userAnswerController,
		AuthMiddleware:            authMiddleware,
	}

	routeConfig.Setup()
}
