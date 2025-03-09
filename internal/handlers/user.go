package handlers

import (
	"errors"
	"go-manage-mysql/cmd/config"
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

	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		web.NewError(ctx, http.StatusBadRequest, config.ErrInvalidBody)
		return
	}

	if validate := validator.ValidateData(user, config.Create_ValidateFields); validate != nil {
		web.NewError(ctx, http.StatusBadRequest, validate.Error())
		return
	}

	create, createErr := h.Service.CreateUser(ctx, user)
	if createErr != nil {
		web.NewError(ctx, http.StatusInternalServerError, createErr.Error())
		return
	}

	ctx.JSON(http.StatusCreated, usersResponse(config.CreatedUserMessage, http.StatusCreated, create))
}

func (h *Handler) SearchUserHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	username := ctx.Query("username")
	if username == "" {
		web.NewError(ctx, http.StatusBadRequest, config.ErrInvalidQueryParam)
		return
	}

	search, searchErr := h.Service.SearchUser(ctx, username)
	if searchErr != nil {
		web.NewError(ctx, http.StatusInternalServerError, searchErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, usersResponse(config.SearchUserMessage, http.StatusOK, search))
}

func (h *Handler) UpdateUserHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var update models.User

	username := ctx.Query("username")
	if username == "" {
		web.NewError(ctx, http.StatusBadRequest, config.ErrInvalidQueryParam)
		return
	}

	if err := ctx.ShouldBindJSON(&update); err != nil {
		web.NewError(ctx, http.StatusBadRequest, config.ErrInvalidBody)
		return
	}

	if validate := validator.ValidateData(update, config.Update_ValidateFields); validate != nil {
		web.NewError(ctx, http.StatusBadRequest, validate.Error())
		return
	}

	if update := h.Service.UpdateUser(ctx, username, update); update != nil {
		web.NewError(ctx, http.StatusInternalServerError, update.Error())
		return
	}

	ctx.JSON(http.StatusOK, usersResponse(config.UpdateUserMessage, http.StatusOK, nil))
}

func (h *Handler) DeleteUserHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	username := ctx.Query("username")
	if username == "" {
		web.NewError(ctx, http.StatusBadRequest, config.ErrInvalidQueryParam)
		return
	}

	if delete := h.Service.DeleteUser(ctx, username); delete != nil {
		web.NewError(ctx, http.StatusInternalServerError, delete.Error())
		return
	}

	ctx.JSON(http.StatusOK, usersResponse(config.DeleteUserMessage, http.StatusOK, nil))
}

func (h *Handler) ChangePwdHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		web.NewError(ctx, http.StatusBadRequest, config.ErrInvalidBody)
		return
	}

	if validate := validator.ValidateData(user, config.ChangePwd_ValidateFields); validate != nil {
		web.NewError(ctx, http.StatusBadRequest, validate.Error())
		return
	}

	if changeErr := h.Service.ChangeUserPwd(ctx, user.Username, user.Password); changeErr != nil {
		web.NewError(ctx, http.StatusInternalServerError, changeErr.Error())
		return
	}

	ctx.JSON(http.StatusOK, usersResponse(config.ChangePwdMessage, http.StatusOK, nil))
}

func (h *Handler) LoginUserHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")

	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		web.NewError(ctx, http.StatusBadRequest, config.ErrInvalidBody)
		return
	}

	if user.Username == "" || user.Password == "" {
		web.NewError(ctx, http.StatusBadRequest, config.ErrAllFieldsAreRequired)
		return
	}

	err := h.Service.LoginUser(ctx, user.Username, user.Password)
	if err != nil {
		if errors.Is(err, config.ErrRecordNotFound) {
			web.NewError(ctx, http.StatusNotFound, config.ErrUserNotFound.Error())
			return
		} else {
			web.NewError(ctx, http.StatusUnauthorized, config.ErrUnauthorizedUser)
			return
		}
	}

	ctx.JSON(http.StatusOK, usersResponse("WELCOME "+user.Username, http.StatusOK, nil))
}

func usersResponse(msg string, status int, data interface{}) models.UserResponse {
	return models.UserResponse{
		Message: msg,
		Status:  status,
		Data:    data,
	}
}
