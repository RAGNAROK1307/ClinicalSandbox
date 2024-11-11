package db

import (
	"ClinicalSandBox/pkg/API/models"
)

func AutoMigrate() {
	DB.AutoMigrate(&models.Role{})
}
