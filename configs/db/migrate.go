package db

import (
	"ClinicalSandBox/internal/API/models"
)

func AutoMigrate() {
	DB.AutoMigrate(&models.Role{})
}
