package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"calculate_lung_capacity/internal/app/schemes"
)

func (h *Handler) ApiRegisterUser(ctx *gin.Context) {
	var request schemes.RegisterUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	_, err := h.Repository.GetUserByUsername(request.Username)
	if err == nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	user, err := h.Repository.CreateUser(request.Username, request.Password)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, schemes.UserResponse{
		PatientCategoryUserID:   user.PatientCategoryUserID,
		PatientCategoryUsername: user.PatientCategoryUsername,
	})
}

func (h *Handler) ApiLoginUser(ctx *gin.Context) {
	var request schemes.LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) ApiLogoutUser(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
