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
	"ClinicalSandBox/configs/db"
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
	r.MaxMultipartMemory = 5 << 20 // 5 MB

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Reemplaza con la URL de tu frontend si cambia http://127.0.0.1:5500
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/users", services2.CreateUser)
	r.GET("/users", services2.GetUsers)
	//r.GET("/users/:id", services2.GetUser)
	//r.PUT("/users/:id", services2.UpdateUser)
	//r.DELETE("/users/:id", services2.DeleteUser)
	r.POST("/executive-and-user", services2.CreateExecutiveAndUser)
	r.GET("/executives-and-users", services2.GetExecutivesAndUsers)

	r.POST("/medical-records-and-related", services2.CreateMedicalRecordAndRelated)
	//r.GET("/medical-records-and-related", services2.GetMedicalRecordsAndRelated)
	//r.GET("/medical-records-and-related/:id", services2.GetMedicalRecordAndRelated)
	r.GET("/user-and-patients", services2.GetUserAndPatients)
	r.GET("/user-and-patients/:id", services2.GetUserAndPatient)
	r.POST("/user-and-patients", services2.CreateUserAndPatient)
	r.PUT("/user-and-patients/:id", services2.UpdateUserAndPatient)
	r.GET("/doctors-and-users", services2.GetDoctorsAndUsers)
	r.POST("/doctor-and-user", services2.CreateDoctorAndUser)
	r.DELETE("/doctors-and-users/:id", services2.DeleteDoctorAndUser)
	r.DELETE("/executives-and-users/:id", services2.DeleteExecutiveAndUser)

	r.POST("/treatments_prescriptions", services2.CreateTreatmentPrescription)
	r.GET("/treatments_prescriptions", services2.GetTreatmentsPrescriptions)
	r.GET("/treatments_prescriptions/:id", services2.GetTreatmentPrescription)
	r.PUT("/treatments_prescriptions/:id", services2.UpdateTreatmentPrescription)
	r.DELETE("/treatments_prescriptions/:id", services2.DeleteTreatmentPrescription)

	r.POST("/consultation_visits", services2.CreateConsultationVisit)
	r.GET("/consultation_visits", services2.GetConsultationVisits)
	r.GET("/consultation_visits/:id", services2.GetConsultationVisit)
	r.PUT("/consultation_visits/:id", services2.UpdateConsultationVisit)
	r.DELETE("/consultation_visits/:id", services2.DeleteConsultationVisit)

	r.POST("/hospital_employees", services2.CreateHospitalEmployee)
	r.GET("/hospital_employees", services2.GetHospitalEmployees)
	r.GET("/medical-records-and-related", services2.GetMedicalRecordsAndRelated)
	r.POST("/medical-records/:id/image", services2.UploadPatientImage)
	r.GET("/medical-records/:id/image", services2.GetPatientImage)
	r.PUT("/medical-records/:id/image", services2.UpdatePatientImage)
	r.DELETE("/medical-records/:id/image", services2.DeletePatientImage)
	r.PUT("/medical_records/:id", services2.UpdateMedicalRecord)
	r.PUT("/medical-records-and-related/:id", services2.UpdateMedicalRecordAndRelated)

	r.POST("/diagnostic_images", services2.CreateDiagnosticImage)
	r.GET("/diagnostic_images", services2.GetDiagnosticImages)
	r.GET("/diagnostic_images/:id", services2.GetDiagnosticImage)
	r.PUT("/diagnostic_images/:id", services2.UpdateDiagnosticImage)
	r.DELETE("/diagnostic_images/:id", services2.DeleteDiagnosticImage)

	r.POST("/diagnostic_images/:id/file", services2.UploadDiagnosticImage)
	r.GET("/diagnostic_images/:id/file", services2.GetDiagnosticImageFile)
	r.PUT("/diagnostic_images/:id/file", services2.UpdateDiagnosticImageFile)
	r.DELETE("/diagnostic_images/:id/file", services2.DeleteDiagnosticImageFile)

	r.POST("/laboratories/:id/file", services2.UploadLaboratoryFile)
	r.GET("/laboratories/:id/file", services2.GetLaboratoryFile)
	r.PUT("/laboratories/:id/file", services2.UpdateLaboratoryFile)
	r.DELETE("/laboratories/:id/file", services2.DeleteLaboratoryFile)

	r.POST("/laboratories", services2.CreateLaboratory)
	r.GET("/laboratories", services2.GetLaboratories)
	r.GET("/laboratories/:id", services2.GetLaboratory)
	r.PUT("/laboratories/:id", services2.UpdateLaboratory)
	r.DELETE("/laboratories/:id", services2.DeleteLaboratory)

	r.POST("/clinical_notes", services2.CreateClinicalNote)
	r.GET("/clinical_notes", services2.GetClinicalNotes)
	r.GET("/clinical_notes/:id", services2.GetClinicalNote)
	r.PUT("/clinical_notes/:id", services2.UpdateClinicalNote)
	r.DELETE("/clinical_notes/:id", services2.DeleteClinicalNote)

	r.DELETE("/user-and-patients/:id", services2.DeleteUserAndPatient)

	//r.PUT("/users/:id/password", services2.UpdatePassword)

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

		auth.POST("/logout", authServices.Logout)

		auth.PUT("/users/:id/password", authMiddleware.ValidatePasswordAccess(), services2.UpdatePassword)

		// Rutas protegidas con roles específicos
		admin := auth.Group("/")
		admin.Use(authMiddleware.RoleMiddleware(db.AdminID)) // Supongamos que el rol 1 es el de administrador
		{
			admin.GET("/patients", services2.GetPatients)

			admin.PUT("/users/:id/passwords", services2.AdminUpdatePassword)

			//admin.POST("/users", services2.CreateUser)
			//admin.GET("/users", services2.GetUsers)
			admin.GET("/users/:id", services2.GetUser)
			admin.PUT("/users/:id", services2.UpdateUser)
			admin.DELETE("/users/:id", services2.DeleteUser)

			//admin.POST("/user-and-patients", services2.CreateUserAndPatient)
			//admin.GET("/user-and-patients", services2.GetUserAndPatients)
			//admin.GET("/user-and-patients/:id", services2.GetUserAndPatient)
			//admin.PUT("/user-and-patients/:id", services2.UpdateUserAndPatient)
			//admin.DELETE("/user-and-patients/:id", services2.DeleteUserAndPatient)

			//admin.POST("/doctor-and-user", services2.CreateDoctorAndUser)
			//admin.GET("/doctors-and-users", services2.GetDoctorsAndUsers)
			//admin.GET("/doctors-and-users/:id", services2.GetDoctorAndUserByID)
			//admin.PUT("/doctors-and-users/:id", services2.UpdateDoctorAndUser)
			//admin.DELETE("/doctors-and-users/:id", services2.DeleteDoctorAndUser)

			//admin.POST("/executive-and-user", services2.CreateExecutiveAndUser)
			//admin.GET("/executives-and-users", services2.GetExecutivesAndUsers)
			//admin.GET("/executives-and-users/:id", services2.GetExecutiveAndUserByID)
			//admin.PUT("/executives-and-users/:id", services2.UpdateExecutiveAndUser)
			//admin.DELETE("/executives-and-users/:id", services2.DeleteExecutiveAndUser)

		}

		// Rutas para médicos (Rol 7)
		doctor := auth.Group("/")
		doctor.Use(authMiddleware.RoleMiddleware(db.MedicoID))
		{
			//doctor.GET("/medical-records-and-related", services2.GetMedicalRecordsAndRelated)
			//doctor.GET("/doctors-and-users/:id", authMiddleware.ValidateUserAccess(), services2.GetDoctorAndUserByID)
		}

		// Rutas para directivos (Rol 8)
		executive := auth.Group("/")
		executive.Use(authMiddleware.RoleMiddleware(db.DirectivoID))
		{
			//executive.GET("/executives-and-users/:id", authMiddleware.ValidateUserAccess(), services2.GetExecutiveAndUserByID)

		}

		// Rutas para pacientes (Rol 9)
		patient := auth.Group("/")
		patient.Use(authMiddleware.RoleMiddleware(db.PacienteID))
		{
			//patient.GET("/user-and-patients/:id", authMiddleware.ValidateUserAccess(), services2.GetUserAndPatient)

		}
	}

	doctorAndAdmin := auth.Group("/")
	doctorAndAdmin.Use(authMiddleware.RoleMiddleware(db.MedicoID, db.AdminID)) // Médicos y administradores
	{
		doctorAndAdmin.GET("/doctors-and-users/:id", authMiddleware.ValidateUserAccess(), services2.GetDoctorAndUserByID)
		doctorAndAdmin.PUT("/doctors-and-users/:id", authMiddleware.ValidateUserAccess(), services2.UpdateDoctorAndUser)
	}

	patientAndAdmin := auth.Group("/")
	patientAndAdmin.Use(authMiddleware.RoleMiddleware(db.PacienteID, db.AdminID)) // Médicos y administradores
	{
		//patientAndAdmin.GET("/user-and-patients/:id", authMiddleware.ValidateUserAccess(), services2.GetUserAndPatient)
	}

	executiveAndAdmin := auth.Group("/")
	executiveAndAdmin.Use(authMiddleware.RoleMiddleware(db.DirectivoID, db.AdminID)) // Médicos y administradores
	{
		executiveAndAdmin.GET("/executives-and-users/:id", authMiddleware.ValidateUserAccess(), services2.GetExecutiveAndUserByID)
		executiveAndAdmin.PUT("/executives-and-users/:id", authMiddleware.ValidateUserAccess(), services2.UpdateExecutiveAndUser)
	}

	doctorAndPatient := auth.Group("/")
	doctorAndPatient.Use(authMiddleware.RoleMiddleware(db.MedicoID, db.PacienteID)) // Médicos y administradores
	{
		doctorAndPatient.GET("/medical-records-and-related/:id", authMiddleware.ValidateMedicalRecordAccess(), services2.GetMedicalRecordAndRelated)
	}

	r.Run(":8080")
}
