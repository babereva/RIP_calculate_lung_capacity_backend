package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"calculate_lung_capacity/internal/app/auth"
	"calculate_lung_capacity/internal/app/ds"
	"calculate_lung_capacity/internal/app/minioclient"
	"calculate_lung_capacity/internal/app/repository"
)

const (
	DefaultImageURL = "/static/img/default_image.png"
	DefaultVideoURL = "/static/img/default_video.mp4"
)

type Handler struct {
	Repository *repository.Repository
	Minio      *minioclient.Client
}

func NewHandler(r *repository.Repository, m *minioclient.Client) *Handler {
	return &Handler{Repository: r, Minio: m}
}

type PatientCategoryCard struct {
	Category   ds.PatientCategory
	ImageURL   string
	VideoURL   string
	LikesCount int
}

func (h *Handler) makeCard(category ds.PatientCategory) PatientCategoryCard {
	imageURL := h.Minio.FileURL(category.PatientCategoryImageURL)
	if imageURL == "" {
		imageURL = DefaultImageURL
	}

	videoURL := h.Minio.FileURL(category.PatientCategoryVideoURL)
	if videoURL == "" {
		videoURL = DefaultVideoURL
	}

	likesCount, err := h.Repository.GetLikesCount(category.PatientCategoryID)
	if err != nil {
		logrus.Error(err)
	}

	return PatientCategoryCard{
		Category:   category,
		ImageURL:   imageURL,
		VideoURL:   videoURL,
		LikesCount: likesCount,
	}
}

func (h *Handler) GetPatientCategories(ctx *gin.Context) {
	ageMin := repository.AgeMin
	ageMax := repository.AgeMax

	if query := ctx.Query("ageMin"); query != "" {
		parsed, err := strconv.Atoi(query)
		if err != nil {
			logrus.Error(err)
		} else {
			ageMin = parsed
		}
	}

	if query := ctx.Query("ageMax"); query != "" {
		parsed, err := strconv.Atoi(query)
		if err != nil {
			logrus.Error(err)
		} else {
			ageMax = parsed
		}
	}

	if ageMin > ageMax {
		ageMin, ageMax = ageMax, ageMin
	}

	categories, err := h.Repository.GetPublishedPatientCategories(ageMin, ageMax)
	if err != nil {
		logrus.Error(err)
	}

	cards := []PatientCategoryCard{}
	for _, category := range categories {
		cards = append(cards, h.makeCard(category))
	}

	ctx.HTML(http.StatusOK, "patient_categories.html", gin.H{
		"patientCategoryCards": cards,
		"ageMin":               ageMin,
		"ageMax":               ageMax,
		"ageScaleMin":          repository.AgeMin,
		"ageScaleMax":          repository.AgeMax,
		"activeTab":            "grid",
	})
}

func (h *Handler) GetPatientCategoryFeed(ctx *gin.Context) {
	idParam := ctx.Param("id")
	isNext := ctx.Query("next") == "true"

	var category ds.PatientCategory
	var err error

	if idParam == "" {
		category, err = h.Repository.GetFirstPublishedPatientCategory()
	} else {
		id, convErr := strconv.Atoi(idParam)
		if convErr != nil {
			logrus.Error(convErr)
			ctx.String(http.StatusBadRequest, "некорректный идентификатор категории пациентов")
			return
		}
		if isNext {
			category, err = h.Repository.GetNextPatientCategory(id)
		} else {
			category, err = h.Repository.GetPatientCategoryByID(id)
		}
	}

	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusNotFound, "категория пациентов не найдена")
		return
	}

	card := h.makeCard(category)

	ctx.HTML(http.StatusOK, "patient_category_feed.html", gin.H{
		"card":      card,
		"activeTab": "feed",
	})
}

func (h *Handler) GetPatientCategoryDraft(ctx *gin.Context) {
	category, err := h.Repository.GetDraftPatientCategory(auth.GetCurrentUserID())

	// черновика ещё нет, показываем пустую форму с кнопкой Далее
	if err != nil {
		ctx.HTML(http.StatusOK, "patient_category_draft.html", gin.H{
			"hasDraft":  false,
			"imageURL":  DefaultImageURL,
			"videoURL":  DefaultVideoURL,
			"activeTab": "draft",
		})
		return
	}

	card := h.makeCard(category)

	ctx.HTML(http.StatusOK, "patient_category_draft.html", gin.H{
		"hasDraft":  true,
		"card":      card,
		"imageURL":  card.ImageURL,
		"videoURL":  card.VideoURL,
		"activeTab": "draft",
	})
}

func (h *Handler) CreatePatientCategoryDraft(ctx *gin.Context) {
	title := ctx.PostForm("title")
	imageURL := ctx.PostForm("imageURL")
	videoURL := ctx.PostForm("videoURL")

	_, err := h.Repository.CreateDraftPatientCategory(auth.GetCurrentUserID(), title, imageURL, videoURL)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusInternalServerError, "не удалось создать черновик")
		return
	}

	ctx.Redirect(http.StatusFound, "/patient-category-draft")
}

func (h *Handler) PublishPatientCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusBadRequest, "некорректный идентификатор категории пациентов")
		return
	}

	description := ctx.PostForm("description")
	age, _ := strconv.Atoi(ctx.PostForm("age"))
	height, _ := strconv.Atoi(ctx.PostForm("height"))

	err = h.Repository.PublishPatientCategory(uint(id), auth.GetCurrentUserID(), description, age, height)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusInternalServerError, "не удалось опубликовать категорию пациентов")
		return
	}

	ctx.Redirect(http.StatusFound, "/patient-categories")
}

func (h *Handler) DeletePatientCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusBadRequest, "некорректный идентификатор категории пациентов")
		return
	}

	err = h.Repository.DeletePatientCategory(id, auth.GetCurrentUserID())
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusInternalServerError, "не удалось удалить категорию пациентов")
		return
	}

	ctx.Redirect(http.StatusFound, "/patient-categories")
}
