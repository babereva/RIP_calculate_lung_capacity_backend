package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"calculate_lung_capacity/internal/app/dsn"
	"calculate_lung_capacity/internal/app/handler"
	"calculate_lung_capacity/internal/app/minioclient"
	"calculate_lung_capacity/internal/app/repository"
)

func StartServer() {
	log.Println("Starting patient category server")

	if err := godotenv.Load(); err != nil {
		logrus.Warn("файл .env не найден, берём переменные окружения как есть")
	}

	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		logrus.Fatal(err)
	}

	minioStorage, err := minioclient.New()
	if err != nil {
		logrus.Fatal(err)
	}

	h := handler.NewHandler(repo, minioStorage)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/patient-categories", h.GetPatientCategories)
	r.GET("/patient-category-feed", h.GetPatientCategoryFeed)
	r.GET("/patient-category-feed/:id", h.GetPatientCategoryFeed)
	r.GET("/patient-category-draft", h.GetPatientCategoryDraft)

	r.POST("/patient-category-draft", h.CreatePatientCategoryDraft)
	r.POST("/patient-category-publish", h.PublishPatientCategory)
	r.POST("/patient-category-delete", h.DeletePatientCategory)

	api := r.Group("/api")

	api.GET("/patient-categories", h.ApiGetPatientCategories)
	api.GET("/patient-categories/feed", h.ApiGetPatientCategoryFeed)
	api.GET("/patient-categories/feed/:id", h.ApiGetPatientCategoryFeed)
	api.GET("/patient-categories/draft", h.ApiGetPatientCategoryDraft)
	api.POST("/patient-categories", h.ApiCreatePatientCategory)
	api.PUT("/patient-categories/:id/publish", h.ApiPublishPatientCategory)
	api.DELETE("/patient-categories/:id", h.ApiDeletePatientCategory)
	api.POST("/patient-categories/:id/like", h.ApiLikePatientCategory)

	api.POST("/users/register", h.ApiRegisterUser)
	api.POST("/users/login", h.ApiLoginUser)
	api.POST("/users/logout", h.ApiLogoutUser)

	if err := r.Run(":8080"); err != nil {
		logrus.Fatal(err)
	}

	log.Println("Patient category server down")
}
