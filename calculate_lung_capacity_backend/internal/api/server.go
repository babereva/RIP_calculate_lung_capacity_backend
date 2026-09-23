package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"calculate_lung_capacity/internal/app/handler"
	"calculate_lung_capacity/internal/app/repository"
)

func StartServer() {
	log.Println("Starting patient category server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Fatal(err)
	}
	h := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	
	r.GET("/patient-categories", h.GetPatientCategories)
	
	r.GET("/patient-category-feed", h.GetPatientCategoryFeed)
	r.GET("/patient-category-feed/:id", h.GetPatientCategoryFeed)
	
	r.GET("/patient-category-draft", h.GetPatientCategoryDraft)

	if err := r.Run(":8080"); err != nil {
		logrus.Fatal(err)
	}
	log.Println("Patient category server down")
}