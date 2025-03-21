package userService

import (
	"errors"
	"fmt"
	"log"
	"newproject/internal/taskService"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user User) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUserByID(id uint, user User) (User, error)
	DeleteUserByID(id uint) error
	GetUserByID(id uint, user *User) error
	GetTasksForUser(userID uint) ([]taskService.Task, error)
	GetUserByEmail(email string) (*User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		log.Printf("Error creating user in DB: %v", err)
	}
	return user, err
}

func (r *userRepository) GetAllUsers() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) UpdateUserByID(id uint, user User) (User, error) {
	var existingUser User
	err := r.db.First(&existingUser, id).Error
	if err != nil {
		return User{}, fmt.Errorf("user not found: %w", err)
	}

	existingUser.Email = user.Email
	existingUser.Password = user.Password

	err = r.db.Save(&existingUser).Error
	if err != nil {
		return User{}, fmt.Errorf("error updating user: %w", err)
	}

	return existingUser, nil
}

func (r *userRepository) DeleteUserByID(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

func (r *userRepository) GetUserByID(id uint, user *User) error {
	return r.db.First(user, id).Error
}

func (r *userRepository) GetTasksForUser(userID uint) ([]taskService.Task, error) {
	var tasks []taskService.Task
	err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

func (r *userRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	result := r.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
