package db

import (
	"ClinicalSandBox/internal/API/models"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Función local duplicada para evitar import cíclico
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// SeedAdminUser asegura que exista un único usuario con el rol Administrador y nombre 'AdministradorSimclec'
func SeedAdminUser() {
	const adminUsername = "AdministradorSimclec"
	const defaultPassword = "Simclec123**"

	var adminUsers []models.User

	// Buscar usuarios con rol de Administrador
	if err := DB.Where("id_rol = ?", AdminID).Find(&adminUsers).Error; err != nil {
		log.Printf("❌ Error buscando usuarios Administrador: %v\n", err)
		return
	}

	// Filtrar cuántos hay con ese nombre exacto
	var correctUser *models.User
	var usersToDelete []models.User

	for _, u := range adminUsers {
		if u.UserName == adminUsername {
			if correctUser == nil {
				correctUser = &u // Primer correcto encontrado
			} else {
				usersToDelete = append(usersToDelete, u) // Exceso
			}
		} else {
			usersToDelete = append(usersToDelete, u)
		}
	}

	// Eliminar usuarios duplicados o incorrectos
	for _, u := range usersToDelete {
		if err := DB.Delete(&u).Error; err != nil {
			log.Printf("⚠️ Error eliminando usuario duplicado %s (ID %d): %v\n", u.UserName, u.IDUser, err)
		} else {
			fmt.Printf("🗑️ Usuario duplicado %s (ID %d) eliminado.\n", u.UserName, u.IDUser)
		}
	}

	// Crear el usuario si no existe
	if correctUser == nil {
		hashedPassword, err := hashPassword(defaultPassword)
		if err != nil {
			log.Printf("❌ Error hasheando la contraseña por defecto: %v\n", err)
			return
		}

		newUser := models.User{
			IDRole:   AdminID,
			UserName: adminUsername,
			Password: hashedPassword,
		}

		if err := DB.Create(&newUser).Error; err != nil {
			log.Printf("❌ Error creando el usuario administrador: %v\n", err)
			return
		}

		fmt.Printf("✅ Usuario administrador '%s' creado exitosamente con ID %d.\n", newUser.UserName, newUser.IDUser)
	} else {
		fmt.Printf("🔎 Usuario administrador '%s' ya existe con ID %d.\n", correctUser.UserName, correctUser.IDUser)
	}
}
