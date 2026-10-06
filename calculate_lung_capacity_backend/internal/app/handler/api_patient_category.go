package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"calculate_lung_capacity/internal/app/auth"
	"calculate_lung_capacity/internal/app/ds"
	"calculate_lung_capacity/internal/app/repository"
	"calculate_lung_capacity/internal/app/schemes"
)

const dateFormat = "02.01.2006 15:04"

func (h *Handler) makeResponse(category ds.PatientCategory) schemes.PatientCategoryResponse {
	likesCount, err := h.Repository.GetLikesCount(category.PatientCategoryID)
	if err != nil {
		logrus.Error(err)
	}

	isMine := 0
	if category.PatientCategoryCreatorID == auth.GetCurrentUserID() {
		isMine = 1
	}

	isLiked := 0
	liked, err := h.Repository.IsLikedByUser(auth.GetCurrentUserID(), category.PatientCategoryID)
	if err != nil {
		logrus.Error(err)
	}
	if liked {
		isLiked = 1
	}

	publishedAt := ""
	if category.PatientCategoryPublishedAt.Valid {
		publishedAt = category.PatientCategoryPublishedAt.Time.Format(dateFormat)
	}

	return schemes.PatientCategoryResponse{
		PatientCategoryID:          category.PatientCategoryID,
		PatientCategoryTitle:       category.PatientCategoryTitle,
		PatientCategoryDescription: category.PatientCategoryDescription,
		PatientCategoryImageURL:    h.Minio.FileURL(category.PatientCategoryImageURL),
		PatientCategoryVideoURL:    h.Minio.FileURL(category.PatientCategoryVideoURL),
		PatientCategoryAge:         category.PatientCategoryAge,
		PatientCategoryHeight:      category.PatientCategoryHeight,
		PatientCategoryCreatedAt:   category.PatientCategoryCreatedAt.Format(dateFormat),
		PatientCategoryPublishedAt: publishedAt,
		PatientCategoryCreatorID:   category.PatientCategoryCreatorID,
		LikesCount:                 likesCount,
		IsLiked:                    isLiked,
		IsMine:                     isMine,
	}
}

func (h *Handler) ApiGetPatientCategories(ctx *gin.Context) {
	ageMin := repository.AgeMin
	ageMax := repository.AgeMax

	if query := ctx.Query("ageMin"); query != "" {
		parsed, err := strconv.Atoi(query)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		ageMin = parsed
	}

	if query := ctx.Query("ageMax"); query != "" {
		parsed, err := strconv.Atoi(query)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		ageMax = parsed
	}

	if ageMin > ageMax {
		ageMin, ageMax = ageMax, ageMin
	}

	categories, err := h.Repository.GetPublishedPatientCategories(ageMin, ageMax)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	responses := []schemes.PatientCategoryResponse{}
	for _, category := range categories {
		responses = append(responses, h.makeResponse(category))
	}

	draftID := uint(0)
	draft, err := h.Repository.GetDraftPatientCategory(auth.GetCurrentUserID())
	if err == nil {
		draftID = draft.PatientCategoryID
	}

	ctx.JSON(http.StatusOK, schemes.PatientCategoryListResponse{
		PatientCategories: responses,
		DraftID:           draftID,
	})
}

func (h *Handler) ApiGetPatientCategoryFeed(ctx *gin.Context) {
	idParam := ctx.Param("id")
	isNext := ctx.Query("next") == "true"

	var category ds.PatientCategory
	var err error

	if idParam == "" {
		category, err = h.Repository.GetFirstPublishedPatientCategory()
	} else {
		id, convErr := strconv.Atoi(idParam)
		if convErr != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if isNext {
			category, err = h.Repository.GetNextPatientCategory(id)
		} else {
			category, err = h.Repository.GetPatientCategoryByID(id)
		}
	}

	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	if category.PatientCategoryStatus != ds.StatusPublished {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, h.makeResponse(category))
}

func (h *Handler) ApiGetPatientCategoryDraft(ctx *gin.Context) {
	category, err := h.Repository.GetDraftPatientCategory(auth.GetCurrentUserID())
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, h.makeResponse(category))
}

func (h *Handler) ApiCreatePatientCategory(ctx *gin.Context) {
	var request schemes.CreatePatientCategoryRequest

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	userID := auth.GetCurrentUserID()

	_, err := h.Repository.GetDraftPatientCategory(userID)
	if err == nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	category, err := h.Repository.CreateDraftPatientCategory(userID, request.Title, "", "")
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	imageName := ""
	if header, fileErr := ctx.FormFile("image"); fileErr == nil {
		imageName, err = h.Minio.UploadFile(header, "image")
		if err != nil {
			logrus.Error(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}

	videoName := ""
	if header, fileErr := ctx.FormFile("video"); fileErr == nil {
		videoName, err = h.Minio.UploadFile(header, "video")
		if err != nil {
			logrus.Error(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}

	if imageName != "" || videoName != "" {
		err = h.Repository.UpdatePatientCategoryMedia(category.PatientCategoryID, userID, imageName, videoName)
		if err != nil {
			logrus.Error(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}

	created, err := h.Repository.GetPatientCategoryByID(int(category.PatientCategoryID))
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, h.makeResponse(created))
}

func (h *Handler) ApiPublishPatientCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	var request schemes.PublishPatientCategoryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = h.Repository.PublishPatientCategory(uint(id), auth.GetCurrentUserID(), request.Description, request.Age, request.Height)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	category, err := h.Repository.GetPatientCategoryByID(id)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, h.makeResponse(category))
}

func (h *Handler) ApiDeletePatientCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = h.Repository.DeletePatientCategory(id, auth.GetCurrentUserID())
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) ApiLikePatientCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	var request schemes.LikeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	category, err := h.Repository.GetPatientCategoryByID(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	userID := auth.GetCurrentUserID()

	if request.Value == 1 {
		err = h.Repository.SetLike(userID, category.PatientCategoryID)
	} else {
		err = h.Repository.RemoveLike(userID, category.PatientCategoryID)
	}

	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	updated, err := h.Repository.GetPatientCategoryByID(id)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, h.makeResponse(updated))
}
