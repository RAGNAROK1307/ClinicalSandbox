package main

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/pkg/API/routes"
)

func main() {

	db.ConnectDB()
	db.AutoMigrate()
	routes.Routes()
	//routes.SetupRoutes()
	//log.Println("Server starting on port 8080...")
	//err := http.ListenAndServe(":8080", nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

/*func main() {

	// Iniciar la base de datos
	db.ConnectDB()
	db.AutoMigrate()

	// Iniciar las rutas de la API (Swagger, etc.)
	go routes.Routes()

	// Configurar y arrancar las rutas estáticas y de login
	go routes.SetupRoutes()

	// Arrancar el servidor para escuchar en el puerto 8080
	log.Println("Server starting on port 8080...")
	err := http.ListenAndServe(":8070", nil)
	if err != nil {
		log.Fatal(err)
	}
}*/
