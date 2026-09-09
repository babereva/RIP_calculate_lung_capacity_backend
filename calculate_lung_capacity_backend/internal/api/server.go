package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"calculate_lung_capacity/internal/app/handler"
	"calculate_lung_capacity/internal/app/repository"
)

func StartServer() {
	log.Println("Starting lung capacity server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Fatal(err)
	}
	h := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	// плитка категорий пациентов с фильтром по возрасту пациента
	r.GET("/patient-categories", h.GetPatientCategories)
	// лента: из панели вкладок — без идентификатора, по клику из плитки — с ним
	r.GET("/lung-capacity-feed", h.GetLungCapacityFeed)
	r.GET("/lung-capacity-feed/:id", h.GetLungCapacityFeed)
	// страница добавления: черновик категории
	r.GET("/patient-category-draft", h.GetPatientCategoryDraft)

	if err := r.Run(":8080"); err != nil {
		logrus.Fatal(err)
	}
	log.Println("Lung capacity server down")
}