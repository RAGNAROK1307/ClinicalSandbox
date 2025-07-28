package db

import (
	"ClinicalSandBox/internal/API/models"
	"fmt"
	"log"
	
	"gorm.io/gorm"
)

// Variables globales para almacenar los IDs de los roles
var (
	AdminID     uint
	MedicoID    uint
	DirectivoID uint
	PacienteID  uint
)

// SeedRoles verifica y corrige los roles en la base de datos si no existen o si tienen información incorrecta.
func SeedRoles() {
	roles := []models.Role{
		{RoleName: "Médico", Description: "Rol asignado a los médicos"},
		{RoleName: "Directivo", Description: "Rol asignado a los directivos"},
		{RoleName: "Paciente", Description: "Rol asignado a los pacientes"},
		{RoleName: "Administrador", Description: "Rol asignado al único Administrador"},
	}

	for _, role := range roles {
		var existingRole models.Role
		result := DB.Where("nombre_rol = ?", role.RoleName).First(&existingRole)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// Crear el rol SIN asignar manualmente IDRole
				if err := DB.Create(&role).Error; err != nil {
					log.Printf("❌ Error creando el rol %s: %v\n", role.RoleName, err)
					continue
				}
				fmt.Printf("✅ Rol %s creado exitosamente con ID %d.\n", role.RoleName, role.IDRole)

				// Asignar el ID generado automáticamente
				existingRole = role
			} else {
				log.Printf("❌ Error buscando el rol %s: %v\n", role.RoleName, result.Error)
				continue
			}
		} else {
			// Verificar si hay cambios en la descripción
			if existingRole.Description != role.Description {
				log.Printf("📝 Corrigiendo descripción del rol %s.\n", role.RoleName)
				if err := DB.Model(&existingRole).Update("descripcion", role.Description).Error; err != nil {
					log.Printf("❌ Error actualizando descripción del rol %s: %v\n", role.RoleName, err)
				} else {
					fmt.Printf("🔄 Descripción del rol %s actualizada correctamente.\n", role.RoleName)
				}
			}
		}

		// Guardar los ID en variables globales
		switch role.RoleName {
		case "Médico":
			MedicoID = existingRole.IDRole
		case "Directivo":
			DirectivoID = existingRole.IDRole
		case "Paciente":
			PacienteID = existingRole.IDRole
		case "Administrador":
			AdminID = existingRole.IDRole
		}
	}

	// Mostrar los valores de los IDs después de la ejecución
	fmt.Println("📌 IDs de roles asignados:")
	fmt.Printf("🩺 Médico: %d\n", MedicoID)
	fmt.Printf("📊 Directivo: %d\n", DirectivoID)
	fmt.Printf("🧑‍⚕️ Paciente: %d\n", PacienteID)
	fmt.Printf("🛠️ Administrador: %d\n", AdminID)
}
