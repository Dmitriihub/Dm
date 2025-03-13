package userService

import (
	"fmt"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUser(user User) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUserByID(id uint, user User) (User, error)
	DeleteUserByID(id uint) error
	GetUserByID(id uint, user *User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user User) (User, error) {
	err := r.db.Create(&user).Error
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

	// Обновляем поля пользователя
	existingUser.Email = user.Email
	existingUser.Password = user.Password

	// Сохраняем обновленного пользователя
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
