package tasks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// Error defines model for Error.
type Error struct {
	Message *string `json:"message,omitempty"`
}

// NewTaskRequest defines model for NewTaskRequest.
type NewTaskRequest struct {
	IsDone *bool   `json:"is_done,omitempty"`
	Task   *string `json:"task,omitempty"`
	UserId *int64  `json:"user_id,omitempty"`
}

// NewUserRequest defines model for NewUserRequest.
type NewUserRequest struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	Username *string `json:"username,omitempty"`
}

type StrictMiddlewareFunc func(f echo.HandlerFunc) echo.HandlerFunc

type StrictHandler interface {
	GetTasks(ctx context.Context) ([]Task, error)
	GetTasksByUserID(ctx context.Context, userID int64) ([]Task, error)
	PostTasks(ctx context.Context, req NewTaskRequest) (Task, error)
	DeleteTasksId(ctx context.Context, id int64) error
	PatchTasksId(ctx context.Context, id int64, req PatchTasksIdJSONRequestBody) (Task, error)
	GetUsers(ctx context.Context) ([]User, error)                    // Добавлено
	PostUsers(ctx context.Context, req NewUserRequest) (User, error) // Добавлено
}

func NewStrictHandler(handler StrictHandler, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{handler: handler, middlewares: middlewares}
}

type strictHandler struct {
	handler     StrictHandler
	middlewares []StrictMiddlewareFunc
}

func (sh *strictHandler) GetTasks(ctx echo.Context) error {
	tasks, err := sh.handler.GetTasks(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, tasks)
}

func (sh *strictHandler) PostTasks(ctx echo.Context) error {
	var req NewTaskRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	task, err := sh.handler.PostTasks(ctx.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, task)
}

func (sh *strictHandler) DeleteTasksId(ctx echo.Context, id int64) error {
	err := sh.handler.DeleteTasksId(ctx.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (sh *strictHandler) PatchTasksId(ctx echo.Context, id int64) error {
	var req PatchTasksIdJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	task, err := sh.handler.PatchTasksId(ctx.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, task)
}

func (sh *strictHandler) GetUsers(ctx echo.Context) error {
	users, err := sh.handler.GetUsers(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, users)
}

func (sh *strictHandler) PostUsers(ctx echo.Context) error {
	var req NewUserRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	user, err := sh.handler.PostUsers(ctx.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, user)
}

// Task defines model for Task.
type Task struct {
	Id     *int64  `json:"id,omitempty"`
	IsDone *bool   `json:"is_done,omitempty"`
	Task   *string `json:"task,omitempty"`
	UserId *int64  `json:"user_id,omitempty"`
}
type PatchTasksIdJSONRequestBody struct {
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
	UserId int64  `json:"user_id"`
}

// User defines model for User.
type User struct {
	Email    *string `json:"email,omitempty"`
	Id       *int64  `json:"id,omitempty"`
	Username *string `json:"username,omitempty"`
}

// PostTasksJSONRequestBody defines body for PostTasks for application/json ContentType.
type PostTasksJSONRequestBody = NewTaskRequest

// PatchTasksIdJSONRequestBody defines body for PatchTasksId for application/json ContentType.
//type PatchTasksIdJSONRequestBody = Task

// PostUsersJSONRequestBody defines body for PostUsers for application/json ContentType.
type PostUsersJSONRequestBody = NewUserRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all tasks
	// (GET /tasks)
	GetTasks(ctx echo.Context) error
	// Create a new task
	// (POST /tasks)
	PostTasks(ctx echo.Context) error
	// Delete a task by ID
	// (DELETE /tasks/{id})
	DeleteTasksId(ctx echo.Context, id int64) error
	// Update a task by ID
	// (PATCH /tasks/{id})
	PatchTasksId(ctx echo.Context, id int64) error
	// Get all users
	// (GET /users)
	GetUsers(ctx echo.Context) error
	// Create a new user
	// (POST /users)
	PostUsers(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetTasks converts echo context to params.
func (w *ServerInterfaceWrapper) GetTasks(ctx echo.Context) error {
	err := w.Handler.GetTasks(ctx)
	return err
}

// PostTasks converts echo context to params.
func (w *ServerInterfaceWrapper) PostTasks(ctx echo.Context) error {
	err := w.Handler.PostTasks(ctx)
	return err
}

// DeleteTasksId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteTasksId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteTasksId(ctx, id)
	return err
}

// PatchTasksId converts echo context to params.
func (w *ServerInterfaceWrapper) PatchTasksId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithOptions("simple", "id", ctx.Param("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PatchTasksId(ctx, id)
	return err
}

// GetUsers converts echo context to params.
func (w *ServerInterfaceWrapper) GetUsers(ctx echo.Context) error {
	err := w.Handler.GetUsers(ctx)
	return err
}

// PostUsers converts echo context to params.
func (w *ServerInterfaceWrapper) PostUsers(ctx echo.Context) error {
	err := w.Handler.PostUsers(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/tasks", wrapper.GetTasks)
	router.POST(baseURL+"/tasks", wrapper.PostTasks)
	router.DELETE(baseURL+"/tasks/:id", wrapper.DeleteTasksId)
	router.PATCH(baseURL+"/tasks/:id", wrapper.PatchTasksId)
	router.GET(baseURL+"/users", wrapper.GetUsers)
	router.POST(baseURL+"/users", wrapper.PostUsers)

}
