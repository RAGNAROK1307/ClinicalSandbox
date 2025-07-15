package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// DB es una variable global que representa la conexión a la base de datos.
// Esta será inicializada al momento de ejecutar la función ConnectDB().
var DB *gorm.DB

// ConnectDB establece la conexión con una base de datos PostgreSQL utilizando GORM.
// La cadena de conexión (DSN) incluye los parámetros necesarios como host, usuario, contraseña,
// nombre de la base de datos, puerto y configuración del modo SSL.
// En caso de error al conectarse, la aplicación se detendrá con un mensaje de error.
// Si la conexión es exitosa, se mostrará un mensaje en consola confirmando la conexión.
func ConnectDB() {
	dsn := "host=localhost user=postgres password=Santiago1307 dbname=bd_sandbox port=5432 sslmode=disable"
	var err error
	// Se intenta abrir la conexión usando el driver de PostgreSQL con GORM
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// Si ocurre un error, se detiene la ejecución y se imprime el error
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	// Confirmación en consola si la conexión es exitosa
	fmt.Println("Database connected!")
}
