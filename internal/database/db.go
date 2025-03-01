package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Task struct {
	ID     int    `json:"id"`
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}

// переменная, через которую мы будем работать с БД
var DB *gorm.DB

func InitDB() {
	// в dsn вводим данные, которые мы указали при создании контейнера
	dsn := "host=localhost user=dmitrijelagin password=1234 dbname=tasks port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = DB.AutoMigrate(&Task{})
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
	fmt.Println("Database migrated successfully")
}

func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Ошибка: БД не инициализирована. Вызови InitDB() перед использованием.")
	}
	return DB
}
