package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"newproject/internal/models"
	"newproject/internal/userService"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService *userService.UserService
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUserHandler(userService *userService.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUsers(ctx echo.Context) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error fetching users: %s", err))
	}

	return ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) PostUsers(ctx echo.Context) error {
	var request models.User
	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid input: %s", err))
	}

	user, err := h.userService.CreateUser(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error creating user: %s", err))
	}

	return ctx.JSON(http.StatusCreated, user)
}

func (h *UserHandler) DeleteUsersId(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	err = h.userService.DeleteUserByID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error deleting user: %s", err))
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *UserHandler) PatchUsersId(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	var request UpdateUserRequest
	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid input: %s", err))
	}

	updatedUser, err := h.userService.UpdateUserByID(uint(id), models.User{
		ID:    uint(id),
		Name:  request.Name,
		Email: request.Email,
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error updating user: %s", err))
	}

	return ctx.JSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) GetUsersIdTasks(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	tasks, err := h.userService.GetUserTasks(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error fetching tasks for user: %s", err))
	}

	return ctx.JSON(http.StatusOK, tasks)
}
