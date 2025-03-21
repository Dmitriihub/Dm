package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"newproject/internal/models"
	"newproject/internal/taskService"
	"newproject/internal/userService"
	openapi "newproject/internal/web/tasks"

	"gorm.io/gorm"
)

type TaskHandler struct {
	taskService *taskService.TaskService
	userService *userService.UserService
}

func NewTaskHandler(taskService *taskService.TaskService, userService *userService.UserService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		userService: userService,
	}
}

// GetUsers возвращает всех пользователей
func (h *TaskHandler) GetUsers(ctx context.Context) ([]openapi.User, error) {
	// Проверяем, что userService инициализирован
	if h.userService == nil {
		log.Println("userService is nil")
		return nil, fmt.Errorf("user service is not initialized")
	}

	// Получаем всех пользователей через userService
	users, err := h.userService.GetAllUsers()
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return nil, fmt.Errorf("error fetching users: %w", err)
	}

	// Преобразуем []models.User в []openapi.User
	var response []openapi.User
	for _, u := range users {
		response = append(response, openapi.User{
			Id:       int64Ptr(int64(u.ID)),
			Username: stringPtr(u.Name),
			Email:    stringPtr(u.Email),
		})
	}

	return response, nil
}

// PostUsers создает нового пользователя
func (h *TaskHandler) PostUsers(ctx context.Context, req openapi.NewUserRequest) (openapi.User, error) {
	if h.userService == nil {
		return openapi.User{}, fmt.Errorf("user service is not initialized")
	}

	if req.Username == nil || req.Email == nil || req.Password == nil {
		return openapi.User{}, fmt.Errorf("username, email, and password are required")
	}

	user := models.User{
		Name:     *req.Username,
		Email:    *req.Email,
		Password: *req.Password,
	}

	createdUser, err := h.userService.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return openapi.User{}, fmt.Errorf("error creating user: %w", err)
	}

	return openapi.User{
		Id:       int64Ptr(int64(createdUser.ID)),
		Username: stringPtr(createdUser.Name),
		Email:    stringPtr(createdUser.Email),
	}, nil
}

// GetTasks возвращает все задачи
func (h *TaskHandler) GetTasks(ctx context.Context) ([]openapi.Task, error) {
	if h.taskService == nil {
		return nil, fmt.Errorf("task service is not initialized")
	}

	taskList, err := h.taskService.GetAllTasks()
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		return nil, fmt.Errorf("error fetching tasks: %w", err)
	}

	var response []openapi.Task
	for _, t := range taskList {
		response = append(response, openapi.Task{
			Id:     int64Ptr(int64(t.ID)),
			IsDone: boolPtr(t.IsDone),
			Task:   stringPtr(t.Task),
			UserId: int64Ptr(int64(t.UserID)),
		})
	}

	return response, nil
}

// GetTasksByUserID возвращает задачи пользователя по user_id
func (h *TaskHandler) GetTasksByUserID(ctx context.Context, userID int64) ([]openapi.Task, error) {
	if h.taskService == nil {
		return nil, fmt.Errorf("task service is not initialized")
	}

	tasks, err := h.taskService.GetTasksByUserID(uint(userID))
	if err != nil {
		log.Printf("Error fetching tasks for user: %v", err)
		return nil, fmt.Errorf("error fetching tasks for user: %w", err)
	}

	var response []openapi.Task
	for _, t := range tasks {
		response = append(response, openapi.Task{
			Id:     int64Ptr(int64(t.ID)),
			IsDone: boolPtr(t.IsDone),
			Task:   stringPtr(t.Task),
			UserId: int64Ptr(int64(t.UserID)),
		})
	}

	return response, nil
}

// PostTasks создает новую задачу
func (h *TaskHandler) PostTasks(ctx context.Context, req openapi.NewTaskRequest) (openapi.Task, error) {
	if h.taskService == nil {
		return openapi.Task{}, fmt.Errorf("task service is not initialized")
	}

	if req.Task == nil || req.IsDone == nil || req.UserId == nil {
		return openapi.Task{}, fmt.Errorf("task, isDone, and userId are required")
	}

	task := taskService.Task{
		Task:   *req.Task,
		IsDone: *req.IsDone,
		UserID: uint(*req.UserId),
	}

	createdTask, err := h.taskService.CreateTask(task)
	if err != nil {
		log.Printf("Error creating task: %v", err)
		return openapi.Task{}, fmt.Errorf("error creating task: %w", err)
	}

	return openapi.Task{
		Id:     int64Ptr(int64(createdTask.ID)),
		Task:   stringPtr(createdTask.Task),
		IsDone: boolPtr(createdTask.IsDone),
		UserId: int64Ptr(int64(createdTask.UserID)),
	}, nil
}

// DeleteTasksId удаляет задачу по ID
func (h *TaskHandler) DeleteTasksId(ctx context.Context, id int64) error {
	if h.taskService == nil {
		return fmt.Errorf("task service is not initialized")
	}

	err := h.taskService.DeleteTaskByID(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("task not found")
	} else if err != nil {
		return fmt.Errorf("error deleting task: %w", err)
	}

	return nil
}

// PatchTasksId обновляет задачу по ID
func (h *TaskHandler) PatchTasksId(ctx context.Context, id int64, req openapi.PatchTasksIdJSONRequestBody) (openapi.Task, error) {
	// Проверяем, что taskService инициализирован
	if h.taskService == nil {
		log.Println("taskService is nil")
		return openapi.Task{}, fmt.Errorf("task service is not initialized")
	}

	log.Printf("Updating task with ID %d: task=%v, isDone=%v, userId=%v", id, req.Task, req.IsDone, req.UserId)

	// Обновляем задачу
	updatedTask, err := h.taskService.UpdateTaskByID(uint(id), taskService.Task{
		Task:   req.Task,
		IsDone: req.IsDone,
		UserID: uint(req.UserId),
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Task with ID %d not found", id)
		return openapi.Task{}, fmt.Errorf("task not found")
	} else if err != nil {
		log.Printf("Error updating task: %v", err)
		return openapi.Task{}, fmt.Errorf("error updating task: %w", err)
	}

	log.Printf("Task with ID %d updated successfully: %+v", id, updatedTask)

	return openapi.Task{
		Id:     int64Ptr(int64(updatedTask.ID)),
		Task:   stringPtr(updatedTask.Task),
		IsDone: boolPtr(updatedTask.IsDone),
		UserId: int64Ptr(int64(updatedTask.UserID)),
	}, nil
}

func int64Ptr(i int64) *int64    { return &i }
func boolPtr(b bool) *bool       { return &b }
func stringPtr(s string) *string { return &s }
