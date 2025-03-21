package userService

import (
	"errors"
	"fmt"
	"log"
	"newproject/internal/models"
	"newproject/internal/taskService"

	"gorm.io/gorm"
)

type UserService struct {
	repo        UserRepository
	taskService *taskService.TaskService
}

func NewUserService(repo UserRepository, taskService *taskService.TaskService) *UserService {
	return &UserService{
		repo:        repo,
		taskService: taskService,
	}
}

// CreateUser создает нового пользователя
func (s *UserService) CreateUser(user models.User) (models.User, error) {
	existingUser, err := s.repo.GetUserByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Error checking user existence: %v", err)
		return models.User{}, fmt.Errorf("error checking user existence: %w", err)
	}

	if existingUser != nil {
		return models.User{}, fmt.Errorf("user with email %s already exists", user.Email)
	}

	userForRepo := toUserRepo(user)
	createdUser, err := s.repo.CreateUser(userForRepo)
	if err != nil {
		log.Printf("Error creating user in repository: %v", err)
		return models.User{}, fmt.Errorf("error creating user in repository: %w", err)
	}
	return toUserModel(createdUser), nil
}

// GetAllUsers возвращает всех пользователей
func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %w", err)
	}

	modelUsers := make([]models.User, len(users))
	for i, u := range users {
		modelUsers[i] = toUserModel(u)
	}

	return modelUsers, nil
}

// DeleteUserByID удаляет пользователя по ID
func (s *UserService) DeleteUserByID(id uint) error {
	return s.repo.DeleteUserByID(id)
}

// UpdateUserByID обновляет пользователя по ID
func (s *UserService) UpdateUserByID(id uint, user models.User) (models.User, error) {
	userForRepo := toUserRepo(user)
	updatedUser, err := s.repo.UpdateUserByID(id, userForRepo)
	if err != nil {
		return models.User{}, fmt.Errorf("error updating user: %w", err)
	}
	return toUserModel(updatedUser), nil
}

// GetUserTasks возвращает задачи пользователя
func (s *UserService) GetUserTasks(userID uint) ([]taskService.Task, error) {
	tasks, err := s.taskService.GetTasksByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user tasks: %w", err)
	}
	return tasks, nil
}

// Преобразование из models.User в User
func toUserRepo(u models.User) User {
	return User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

// Преобразование из User в models.User
func toUserModel(u User) models.User {
	return models.User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}
