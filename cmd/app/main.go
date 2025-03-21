package main

import (
	"log"
	"net/http"
	"newproject/internal/database"
	"newproject/internal/handlers"
	"newproject/internal/taskService"
	"newproject/internal/userService"
	"newproject/internal/web/tasks"
	"newproject/internal/web/users"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	database.InitDB()
	if err := database.DB.AutoMigrate(&userService.User{}, &taskService.Task{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	taskRepo := taskService.NewTaskRepository(database.DB)
	userRepo := userService.NewUserRepository(database.DB)

	taskService := taskService.NewTaskService(taskRepo)
	userService := userService.NewUserService(userRepo, taskService)

	taskHandler := handlers.NewTaskHandler(taskService, userService)
	userHandler := handlers.NewUserHandler(userService)

	userStrictHandler := users.NewStrictHandler(userHandler, nil)
	users.RegisterHandlers(e, userStrictHandler)

	taskStrictHandler := tasks.NewStrictHandler(taskHandler, nil)
	tasks.RegisterHandlers(e, taskStrictHandler)

	e.GET("/users/:user_id/tasks", func(c echo.Context) error {
		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		}

		tasks, err := taskService.GetTasksByUserID(uint(userID))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, tasks)
	})

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start with err: %v", err)
	}
}
