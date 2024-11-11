package routes

import (
	_ "ClinicalSandBox/docs"
	"ClinicalSandBox/pkg/API/services"
	"ClinicalSandBox/pkg/auth/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"path/filepath"
)

// Routes configura las rutas de la API
func Routes() {
	r := gin.Default()

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/roles", services.CreateRole)
	r.GET("/roles", services.GetRoles)
	r.GET("/roles/:id", services.GetRole)
	r.PUT("/roles/:id", services.UpdateRole)
	r.DELETE("/roles/:id", services.DeleteRole)

	r.Run(":8080")
}

func SetupRoutes() {
	// Usa la ruta correcta para los archivos estáticos
	staticDir := filepath.Join("C:", "public", "web")

	// Verifica si la ruta es correcta
	absPath, err := filepath.Abs(staticDir)
	if err != nil {
		log.Fatalf("Error al obtener la ruta: %v", err)
	}
	log.Println("Ruta absoluta a archivos estáticos:", absPath)

	// Sirve los archivos estáticos desde public/web
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(absPath))))

	// Configura el manejador para la ruta de login
	http.HandleFunc("/login", routes.LoginHandler)

}
