package main

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/routes"
)

// main es el punto de entrada principal de la aplicación ClinicalSandBox.
// Esta función ejecuta las siguientes tareas en orden:
// 1. Conexión a la base de datos mediante db.ConnectDB().
// 2. Migración automática de los modelos definidos con db.AutoMigrate().
// 3. Inserción de datos iniciales como roles predeterminados con db.SeedRoles().
// 4. Inicialización de las rutas HTTP del sistema con routes.Routes().
//
// Es fundamental que la base de datos esté disponible y correctamente configurada
// antes de ejecutar esta función para garantizar el correcto funcionamiento del sistema.
func main() {

	// Establece la conexión con la base de datos PostgreSQL
	db.ConnectDB()

	// Realiza la migración automática de esquemas de base de datos
	db.AutoMigrate()

	// Carga los datos iniciales requeridos (por ejemplo, roles del sistema)
	db.SeedRoles()

	// Inicializa y ejecuta las rutas de la API
	routes.Routes()
}
