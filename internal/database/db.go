package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User представляет модель пользователя
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}

// Task представляет модель задачи
type Task struct {
	ID     uint   `gorm:"primaryKey"`
	Task   string `gorm:"not null"`
	IsDone bool   `gorm:"default:false"`
	UserID uint   `gorm:"not null"` // Связь с пользователем
}

// DB — глобальная переменная для работы с базой данных
var DB *gorm.DB

// InitDB инициализирует подключение к базе данных и выполняет миграции
func InitDB() {
	// Настройки подключения к базе данных
	dsn := "host=localhost user=youruser password=yourpassword dbname=yourdb port=5432 sslmode=disable"
	var err error

	// Подключение к базе данных
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Автоматическое создание таблиц
	err = DB.AutoMigrate(&User{}, &Task{})
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	fmt.Println("Database migrated successfully")
}

// GetDB возвращает подключение к базе данных
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Ошибка: БД не инициализирована. Вызови InitDB() перед использованием.")
	}
	return DB
}
