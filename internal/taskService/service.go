package taskService

import (
	"fmt"
)

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTask создает задачу
func (s *TaskService) CreateTask(task Task) (Task, error) {
	if task.UserID == 0 {
		return Task{}, fmt.Errorf("user_id is required")
	}
	return s.repo.CreateTask(task)
}

// GetAllTasks возвращает все задачи
func (s *TaskService) GetAllTasks() ([]Task, error) {
	return s.repo.GetAllTasks()
}

// UpdateTaskByID обновляет задачу по ID
func (s *TaskService) UpdateTaskByID(id uint, task Task) (Task, error) {
	return s.repo.UpdateTaskByID(id, task)
}

// DeleteTaskByID удаляет задачу по ID
func (s *TaskService) DeleteTaskByID(id uint) error {
	return s.repo.DeleteTaskByID(id)
}

// GetTasksByUserID возвращает задачи пользователя по user_id
func (s *TaskService) GetTasksByUserID(userID uint) ([]Task, error) {
	return s.repo.GetTasksByUserID(userID)
}
