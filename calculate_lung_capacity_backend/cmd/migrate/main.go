package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"calculate_lung_capacity/internal/app/ds"
	"calculate_lung_capacity/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("не удалось подключиться к базе")
	}

	err = db.AutoMigrate(
		&ds.Physician{},
		&ds.PatientCategory{},
		&ds.CategoryLike{},
	)
	if err != nil {
		panic("не удалось создать таблицы")
	}
}
