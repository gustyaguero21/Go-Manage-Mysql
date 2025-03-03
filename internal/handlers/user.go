package handlers

import (
	"go-manage-mysql/internal/models"
	"go-manage-mysql/internal/services"
	"go-manage-mysql/internal/utils/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gustyaguero21/go-core/pkg/web"
)

type Handler struct {
	Service services.UserServices
}

func NewUserHandler(service services.UserServices) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) CreateUserHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.Set("Timeout", 5)

	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		web.NewError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	validateFields := []string{"name", "surname", "username", "phone", "email", "password"}
	if validate := validator.ValidateData(user, validateFields); validate != nil {
		web.NewError(ctx, http.StatusBadRequest, validate.Error())
		return
	}

	create, createErr := h.Service.CreateUser(ctx, user)
	if createErr != nil {
		web.NewError(ctx, http.StatusInternalServerError, createErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, usersResponse("user created successfully", http.StatusCreated, create))
}

func (h *Handler) SearchUserHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.Set("Timeout", 5)

	username := ctx.Query("username")
	if username == "" {
		web.NewError(ctx, http.StatusBadRequest, "invalid query param")
		return
	}

	search, searchErr := h.Service.SearchUser(ctx, username)
	if searchErr != nil {
		web.NewError(ctx, http.StatusInternalServerError, searchErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, usersResponse("user found", http.StatusOK, search))
}

func (h *Handler) UpdataUserHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.Set("Timeout", 5)

	var update models.User

	username := ctx.Query("username")
	if username == "" {
		web.NewError(ctx, http.StatusBadRequest, "invalid query param")
		return
	}

	if err := ctx.ShouldBindJSON(&update); err != nil {
		web.NewError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	validateFields := []string{"name", "surname", "phone", "email"}
	if validate := validator.ValidateData(update, validateFields); validate != nil {
		web.NewError(ctx, http.StatusBadRequest, validate.Error())
		return
	}

	if update := h.Service.UpdateUser(ctx, username, update); update != nil {
		web.NewError(ctx, http.StatusInternalServerError, update.Error())
		return
	}

	ctx.JSON(http.StatusOK, usersResponse("user updated successfully", http.StatusOK, nil))
}

func (h *Handler) DeleteUserHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	if username == "" {
		web.NewError(ctx, http.StatusBadRequest, "invalid query param")
		return
	}

	if delete := h.Service.DeleteUser(ctx, username); delete != nil {
		web.NewError(ctx, http.StatusInternalServerError, delete.Error())
		return
	}

	ctx.JSON(http.StatusOK, usersResponse("user deleted successfully", http.StatusOK, nil))
}

func usersResponse(msg string, status int, data interface{}) models.UserResponse {
	return models.UserResponse{
		Message: msg,
		Status:  status,
		Data:    data,
	}
}
