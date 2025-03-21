package taskService

import (
	"log"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task Task) (Task, error)
	GetAllTasks() ([]Task, error)
	UpdateTaskByID(id uint, task Task) (Task, error)
	DeleteTaskByID(id uint) error
	GetTasksByUserID(userID uint) ([]Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *taskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task Task) (Task, error) {
	result := r.db.Create(&task)
	if result.Error != nil {
		return Task{}, result.Error
	}
	return task, nil
}

func (r *taskRepository) GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) UpdateTaskByID(id uint, task Task) (Task, error) {
	var existing Task
	if err := r.db.First(&existing, id).Error; err != nil {
		log.Printf("Task with ID %d not found: %v", id, err)
		return Task{}, err
	}
	log.Printf("Updating task with ID %d: old task=%v, new task=%v", id, existing, task)

	existing.Task = task.Task
	existing.IsDone = task.IsDone
	existing.UserID = task.UserID

	err := r.db.Save(&existing).Error
	if err != nil {
		log.Printf("Error updating task with ID %d: %v", id, err)
	} else {
		log.Printf("Task with ID %d updated successfully: %v", id, existing)
	}

	return existing, err
}

func (r *taskRepository) DeleteTaskByID(id uint) error {
	return r.db.Delete(&Task{}, id).Error
}

func (r *taskRepository) GetTasksByUserID(userID uint) ([]Task, error) {
	var tasks []Task
	err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
