/*package routes

import (
	_ "ClinicalSandBox/docs"
	services2 "ClinicalSandBox/internal/API/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

// Routes configura las rutas de la API
func Routes() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Reemplaza con la URL de tu frontend si cambia http://127.0.0.1:5500
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/roles", services2.CreateRole)
	r.GET("/roles", services2.GetRoles)
	r.GET("/roles/:id", services2.GetRole)
	r.PUT("/roles/:id", services2.UpdateRole)
	r.DELETE("/roles/:id", services2.DeleteRole)

	r.POST("/identifications", services2.CreateIdentification)
	r.GET("/identifications", services2.GetIdentifications)
	r.GET("/identifications/:id", services2.GetIdentification)
	r.PUT("/identifications/:id", services2.UpdateIdentification)
	r.DELETE("/identifications/:id", services2.DeleteIdentification)

	r.POST("/users", services2.CreateUser)
	r.GET("/users", services2.GetUsers)
	r.GET("/users/:id", services2.GetUser)
	r.PUT("/users/:id", services2.UpdateUser)
	r.DELETE("/users/:id", services2.DeleteUser)

	r.POST("/demographics_data", services2.CreateDemographicData)
	r.GET("/demographics_data", services2.GetDemographicsData)
	r.GET("/demographics_data/:id", services2.GetDemographicData)
	r.PUT("/demographics_data/:id", services2.UpdateDemographicData)
	r.DELETE("/demographics_data/:id", services2.DeleteDemographicData)

	r.POST("/hospital_employees", services2.CreateHospitalEmployee)
	r.GET("/hospital_employees", services2.GetHospitalEmployees)
	r.GET("/hospital_employees/:id", services2.GetHospitalEmployee)
	r.PUT("/hospital_employees/:id", services2.UpdateHospitalEmployee)
	r.DELETE("/hospital_employees/:id", services2.DeleteHospitalEmployee)

	r.POST("/patients", services2.CreatePatient)
	r.GET("/patients", services2.GetPatients)
	r.GET("/patients/:id", services2.GetPatient)
	r.PUT("/patients/:id", services2.UpdatePatient)
	r.DELETE("/patients/:id", services2.DeletePatient)

	r.POST("/medical_records", services2.CreateMedicalRecord)
	r.GET("/medical_records", services2.GetMedicalRecords)
	r.GET("/medical_records/:id", services2.GetMedicalRecord)
	r.PUT("/medical_records/:id", services2.UpdateMedicalRecord)
	r.DELETE("/medical_records/:id", services2.DeleteMedicalRecord)

	r.POST("/consultation_visits", services2.CreateConsultationVisit)
	r.GET("/consultation_visits", services2.GetConsultationVisits)
	r.GET("/consultation_visits/:id", services2.GetConsultationVisit)
	r.PUT("/consultation_visits/:id", services2.UpdateConsultationVisit)
	r.DELETE("/consultation_visits/:id", services2.DeleteConsultationVisit)

	r.POST("/laboratories", services2.CreateLaboratory)
	r.GET("/laboratories", services2.GetLaboratories)
	r.GET("/laboratories/:id", services2.GetLaboratory)
	r.PUT("/laboratories/:id", services2.UpdateLaboratory)
	r.DELETE("/laboratories/:id", services2.DeleteLaboratory)

	r.POST("/diagnostic_images", services2.CreateDiagnosticImage)
	r.GET("/diagnostic_images", services2.GetDiagnosticImages)
	r.GET("/diagnostic_images/:id", services2.GetDiagnosticImage)
	r.PUT("/diagnostic_images/:id", services2.UpdateDiagnosticImage)
	r.DELETE("/diagnostic_images/:id", services2.DeleteDiagnosticImage)

	r.POST("/treatments_prescriptions", services2.CreateTreatmentPrescription)
	r.GET("/treatments_prescriptions", services2.GetTreatmentsPrescriptions)
	r.GET("/treatments_prescriptions/:id", services2.GetTreatmentPrescription)
	r.PUT("/treatments_prescriptions/:id", services2.UpdateTreatmentPrescription)
	r.DELETE("/treatments_prescriptions/:id", services2.DeleteTreatmentPrescription)

	r.POST("/clinical_notes", services2.CreateClinicalNote)
	r.GET("/clinical_notes", services2.GetClinicalNotes)
	r.GET("/clinical_notes/:id", services2.GetClinicalNote)
	r.PUT("/clinical_notes/:id", services2.UpdateClinicalNote)
	r.DELETE("/clinical_notes/:id", services2.DeleteClinicalNote)

	r.POST("/consent_authorizations", services2.CreateConsentAuthorization)
	r.GET("/consent_authorizations", services2.GetConsentAuthorizations)
	r.GET("/consent_authorizations/:id", services2.GetConsentAuthorization)
	r.PUT("/consent_authorizations/:id", services2.UpdateConsentAuthorization)
	r.DELETE("/consent_authorizations/:id", services2.DeleteConsentAuthorization)

	// En tu archivo de rutas
	r.POST("/user-and-patients", services2.CreateUserAndPatient)
	r.GET("/user-and-patients", services2.GetUserAndPatients)
	r.GET("/user-and-patients/:id", services2.GetUserAndPatient)
	r.PUT("/user-and-patients/:id", services2.UpdateUserAndPatient)
	r.DELETE("/user-and-patients/:id", services2.DeleteUserAndPatient)

	r.POST("/doctor-and-user", services2.CreateDoctorAndUser)
	r.GET("/doctors-and-users", services2.GetDoctorsAndUsers)
	r.GET("/doctors-and-users/:id", services2.GetDoctorAndUserByID)
	r.PUT("/doctors-and-users/:id", services2.UpdateDoctorAndUser)
	r.DELETE("/doctors-and-users/:id", services2.DeleteDoctorAndUser)

	r.POST("/executive-and-user", services2.CreateExecutiveAndUser)
	r.GET("/executives-and-users", services2.GetExecutivesAndUsers)
	r.GET("/executives-and-users/:id", services2.GetExecutiveAndUserByID)
	r.PUT("/executives-and-users/:id", services2.UpdateExecutiveAndUser)
	r.DELETE("/executives-and-users/:id", services2.DeleteExecutiveAndUser)

	r.Run(":8080")
}*/

package routes

import (
	_ "ClinicalSandBox/docs"
	services2 "ClinicalSandBox/internal/API/services"
	authMiddleware "ClinicalSandBox/internal/auth/middleware"
	authServices "ClinicalSandBox/internal/auth/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

func Routes() {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Reemplaza con la URL de tu frontend si cambia http://127.0.0.1:5500
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//r.POST("/users", services2.CreateUser)
	//r.GET("/users", services2.GetUsers)
	//r.GET("/users/:id", services2.GetUser)
	//r.PUT("/users/:id", services2.UpdateUser)
	//r.DELETE("/users/:id", services2.DeleteUser)

	// Login Route
	r.POST("/login", authServices.Login)

	// Rutas protegidas con autenticación
	auth := r.Group("/")
	auth.Use(authMiddleware.AuthMiddleware())

	{
		auth.POST("/roles", services2.CreateRole)
		auth.GET("/roles", services2.GetRoles)
		auth.GET("/roles/:id", services2.GetRole)
		auth.PUT("/roles/:id", services2.UpdateRole)
		auth.DELETE("/roles/:id", services2.DeleteRole)

		// Rutas protegidas con roles específicos
		admin := auth.Group("/")
		admin.Use(authMiddleware.RoleMiddleware(10)) // Supongamos que el rol 1 es el de administrador
		{
			admin.GET("/patients", services2.GetPatients)

			admin.POST("/users", services2.CreateUser)
			admin.GET("/users", services2.GetUsers)
			admin.GET("/users/:id", services2.GetUser)
			admin.PUT("/users/:id", services2.UpdateUser)
			admin.DELETE("/users/:id", services2.DeleteUser)

			admin.POST("/user-and-patients", services2.CreateUserAndPatient)
			admin.GET("/user-and-patients", services2.GetUserAndPatients)
			admin.GET("/user-and-patients/:id", services2.GetUserAndPatient)
			admin.PUT("/user-and-patients/:id", services2.UpdateUserAndPatient)
			admin.DELETE("/user-and-patients/:id", services2.DeleteUserAndPatient)

			admin.POST("/doctor-and-user", services2.CreateDoctorAndUser)
			admin.GET("/doctors-and-users", services2.GetDoctorsAndUsers)
			admin.GET("/doctors-and-users/:id", services2.GetDoctorAndUserByID)
			admin.PUT("/doctors-and-users/:id", services2.UpdateDoctorAndUser)
			admin.DELETE("/doctors-and-users/:id", services2.DeleteDoctorAndUser)

			admin.POST("/executive-and-user", services2.CreateExecutiveAndUser)
			admin.GET("/executives-and-users", services2.GetExecutivesAndUsers)
			admin.GET("/executives-and-users/:id", services2.GetExecutiveAndUserByID)
			admin.PUT("/executives-and-users/:id", services2.UpdateExecutiveAndUser)
			admin.DELETE("/executives-and-users/:id", services2.DeleteExecutiveAndUser)

		}

		// Otras rutas protegidas...
	}

	r.Run(":8080")
}
