package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"calculate_lung_capacity/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repository: r}
}


type PatientCategoryCard struct {
	Category   repository.PatientCategory
	LikesCount int
}


func (h *Handler) GetPatientCategories(ctx *gin.Context) {
	patientAgeQuery := ctx.Query("patientAge")

	patientAge := 0
	if patientAgeQuery != "" {
		parsedAge, err := strconv.Atoi(patientAgeQuery)
		if err != nil {
			logrus.Error(err)
		} else {
			patientAge = parsedAge
		}
	}

	categories, err := h.Repository.GetPublishedPatientCategories(patientAge)
	if err != nil {
		logrus.Error(err)
	}

	cards := []PatientCategoryCard{}
	for _, category := range categories {
		cards = append(cards, PatientCategoryCard{
			Category:   category,
			LikesCount: len(category.LikedByPhysicianIDs),
		})
	}

	ctx.HTML(http.StatusOK, "patient_categories.html", gin.H{
		"patientCategoryCards": cards,
		"patientAge":           patientAgeQuery,
		"minioBaseUrl":         repository.MinioBaseURL,
		"activeTab":            "grid",
	})
}


func (h *Handler) GetLungCapacityFeed(ctx *gin.Context) {
	idParam := ctx.Param("id")
	isNext := ctx.Query("next") == "true"

	var patientCategory repository.PatientCategory
	var err error

	if idParam == "" {
		patientCategory, err = h.Repository.GetNextPatientCategory(0)
	} else {
		id, convErr := strconv.Atoi(idParam)
		if convErr != nil {
			logrus.Error(convErr)
			ctx.String(http.StatusBadRequest, "некорректный идентификатор категории пациентов")
			return
		}
		if isNext {
			patientCategory, err = h.Repository.GetNextPatientCategory(id)
		} else {
			patientCategory, err = h.Repository.GetPatientCategoryByID(id)
		}
	}

	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusNotFound, "категория пациентов не найдена")
		return
	}

	ctx.HTML(http.StatusOK, "lung_capacity_feed.html", gin.H{
		"patientCategory": patientCategory,
		"likesCount":      len(patientCategory.LikedByPhysicianIDs),
		"minioBaseUrl":    repository.MinioBaseURL,
		"activeTab":       "feed",
	})
}


func (h *Handler) GetPatientCategoryDraft(ctx *gin.Context) {
	draftCategory, err := h.Repository.GetDraftPatientCategory()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "patient_category_draft.html", gin.H{
		"patientCategory": draftCategory,
		"minioBaseUrl":    repository.MinioBaseURL,
		"activeTab":       "draft",
	})
}