package userService

import (
	"newproject/internal/taskService"
	"time"
)

// Структура для User в GORM
type User struct {
	ID        uint               `json:"id" gorm:"primaryKey"`
	Email     string             `json:"email" gorm:"unique;not null"`
	Password  string             `json:"password" gorm:"not null"`
	Name      string             `json:"name" gorm:"not null"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Tasks     []taskService.Task `json:"tasks" gorm:"foreignKey:UserID"`
}
