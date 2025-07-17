package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // Permite cargar variables de entorno desde un archivo .env
	"gorm.io/driver/postgres"  // Driver PostgreSQL para GORM
	"gorm.io/gorm"             // ORM para Go
)

// DB es la instancia global de la base de datos que se utilizará en todo el sistema.
var DB *gorm.DB

// ConnectDB establece la conexión con la base de datos PostgreSQL utilizando variables
// de entorno definidas en un archivo .env. Usa GORM como ORM y verifica que todas
// las variables requeridas estén disponibles.
//
// Si alguna variable falta o la conexión falla, el sistema se detiene con un error crítico.
func ConnectDB() {
	// Cargar variables del archivo .env ubicado en la raíz del proyecto.
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error cargando el archivo .env")
	}

	// Lista de variables de entorno necesarias para establecer la conexión.
	requiredVars := []string{
		"DB_HOST",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_PORT",
		"DB_SSLMODE",
	}

	// Validar que todas las variables requeridas estén definidas.
	for _, key := range requiredVars {
		if os.Getenv(key) == "" {
			log.Fatalf("Falta la variable de entorno: %s", key)
		}
	}

	// Construcción del DSN (Data Source Name) para PostgreSQL.
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	// Intentar abrir la conexión a la base de datos usando GORM y el driver de PostgreSQL.
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar con la base de datos: ", err)
	}

	// Confirmación en consola si la conexión es exitosa.
	fmt.Println("✅ Conexión a la base de datos establecida")
}
