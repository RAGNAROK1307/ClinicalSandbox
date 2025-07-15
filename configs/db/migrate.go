package db

import (
	"ClinicalSandBox/internal/API/models"
)

// AutoMigrate ejecuta el proceso de migración automática de esquemas de base de datos
// utilizando GORM. Esta función asegura que la tabla correspondiente al modelo `Role`
// exista en la base de datos y esté actualizada según la estructura definida en el modelo.
//
// En este caso, se migra únicamente el modelo `Role`, creando la tabla si no existe
// o actualizando su estructura si ha cambiado.
//
// Nota: Esta función debe ser llamada después de establecer una conexión exitosa
// con la base de datos mediante `ConnectDB()`.
func AutoMigrate() {
	DB.AutoMigrate(&models.Role{})
}
