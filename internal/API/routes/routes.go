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
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Login Route
	r.POST("/login", authServices.Login)

	// Rutas protegidas con autenticación
	auth := r.Group("/")
	auth.Use(authMiddleware.AuthMiddleware())
	{
		auth.POST("/logout", authServices.Logout)
		auth.PUT("/users/:id/password", authMiddleware.ValidatePasswordAccess(), services2.UpdatePassword)

		// Rutas protegidas con roles específicos
		admin := auth.Group("/")
		admin.Use(authMiddleware.RoleMiddleware(db.AdminID))
		{
			admin.POST("/roles", services2.CreateRole)
			admin.GET("/roles", services2.GetRoles)
			admin.GET("/roles/:id", services2.GetRole)
			admin.PUT("/roles/:id", services2.UpdateRole)
			admin.DELETE("/roles/:id", services2.DeleteRole)

			admin.POST("/identifications", services2.CreateIdentification)
			admin.GET("/identifications", services2.GetIdentifications)
			admin.GET("/identifications/:id", services2.GetIdentification)
			admin.PUT("/identifications/:id", services2.UpdateIdentification)
			admin.DELETE("/identifications/:id", services2.DeleteIdentification)

			admin.POST("/users", services2.CreateUser)
			admin.GET("/users", services2.GetUsers)
			admin.GET("/users/:id", services2.GetUser)
			admin.PUT("/users/:id", services2.UpdateUser)
			admin.DELETE("/users/:id", services2.DeleteUser)

			admin.POST("/demographics_data", services2.CreateDemographicData)
			admin.GET("/demographics_data", services2.GetDemographicsData)
			admin.GET("/demographics_data/:id", services2.GetDemographicData)
			admin.PUT("/demographics_data/:id", services2.UpdateDemographicData)
			admin.DELETE("/demographics_data/:id", services2.DeleteDemographicData)

			admin.POST("/patients", services2.CreatePatient)
			admin.GET("/patients", services2.GetPatients)
			admin.GET("/patients/:id", services2.GetPatient)
			admin.PUT("/patients/:id", services2.UpdatePatient)
			admin.DELETE("/patients/:id", services2.DeletePatient)

			admin.DELETE("/user-and-patients/:id", services2.DeleteUserAndPatient)

			admin.POST("/doctor-and-user", services2.CreateDoctorAndUser)
			admin.GET("/doctors-and-users", services2.GetDoctorsAndUsers)
			admin.DELETE("/doctors-and-users/:id", services2.DeleteDoctorAndUser)

			admin.POST("/executive-and-user", services2.CreateExecutiveAndUser)
			admin.GET("/executives-and-users", services2.GetExecutivesAndUsers)
			admin.DELETE("/executives-and-users/:id", services2.DeleteExecutiveAndUser)

			admin.POST("/hospital_employees", services2.CreateHospitalEmployee)
			admin.GET("/hospital_employees", services2.GetHospitalEmployees)
			admin.GET("/hospital_employees/:id", services2.GetHospitalEmployee)
			admin.PUT("/hospital_employees/:id", services2.UpdateHospitalEmployee)
			admin.DELETE("/hospital_employees/:id", services2.DeleteHospitalEmployee)

			admin.PUT("/users/:id/passwords", services2.AdminUpdatePassword)

		}

		// Rutas para médicos (Rol 7)
		doctor := auth.Group("/")
		doctor.Use(authMiddleware.RoleMiddleware(db.MedicoID))
		{
			doctor.POST("/medical_records", services2.CreateMedicalRecord)
			doctor.GET("/medical_records", services2.GetMedicalRecords)
			//doctor.GET("/medical_records/:id", services2.
			//)
			doctor.PUT("/medical_records/:id", services2.UpdateMedicalRecord)
			doctor.DELETE("/medical_records/:id", services2.DeleteMedicalRecord)

			doctor.POST("/medical-records-and-related", services2.CreateMedicalRecordAndRelated)
			doctor.PUT("/medical-records-and-related/:id", services2.UpdateMedicalRecordAndRelated)
			doctor.GET("/medical-records-and-related", services2.GetMedicalRecordsAndRelated)

			doctor.POST("/treatments_prescriptions", services2.CreateTreatmentPrescription)
			doctor.GET("/treatments_prescriptions", services2.GetTreatmentsPrescriptions)
			doctor.PUT("/treatments_prescriptions/:id", services2.UpdateTreatmentPrescription)
			doctor.DELETE("/treatments_prescriptions/:id", services2.DeleteTreatmentPrescription)

			doctor.POST("/consultation_visits", services2.CreateConsultationVisit)
			doctor.GET("/consultation_visits", services2.GetConsultationVisits)
			doctor.PUT("/consultation_visits/:id", services2.UpdateConsultationVisit)
			doctor.DELETE("/consultation_visits/:id", services2.DeleteConsultationVisit)

			doctor.POST("/medical-records/:id/image", services2.UploadPatientImage)
			doctor.GET("/medical-records/:id/image", services2.GetPatientImage)
			doctor.PUT("/medical-records/:id/image", services2.UpdatePatientImage)
			doctor.DELETE("/medical-records/:id/image", services2.DeletePatientImage)

			doctor.POST("/laboratories/:id/file", services2.UploadLaboratoryFile)
			doctor.GET("/laboratories/:id/file", services2.GetLaboratoryFile)
			doctor.PUT("/laboratories/:id/file", services2.UpdateLaboratoryFile)
			doctor.DELETE("/laboratories/:id/file", services2.DeleteLaboratoryFile)

			doctor.POST("/diagnostic_images/:id/file", services2.UploadDiagnosticImage)
			doctor.GET("/diagnostic_images/:id/file", services2.GetDiagnosticImageFile)
			doctor.PUT("/diagnostic_images/:id/file", services2.UpdateDiagnosticImageFile)
			doctor.DELETE("/diagnostic_images/:id/file", services2.DeleteDiagnosticImageFile)

			doctor.POST("/diagnostic_images", services2.CreateDiagnosticImage)
			doctor.GET("/diagnostic_images", services2.GetDiagnosticImages)
			doctor.GET("/diagnostic_images/:id", services2.GetDiagnosticImage)
			doctor.PUT("/diagnostic_images/:id", services2.UpdateDiagnosticImage)
			doctor.DELETE("/diagnostic_images/:id", services2.DeleteDiagnosticImage)

			doctor.POST("/laboratories", services2.CreateLaboratory)
			doctor.GET("/laboratories", services2.GetLaboratories)
			doctor.GET("/laboratories/:id", services2.GetLaboratory)
			doctor.PUT("/laboratories/:id", services2.UpdateLaboratory)
			doctor.DELETE("/laboratories/:id", services2.DeleteLaboratory)

			doctor.POST("/clinical_notes", services2.CreateClinicalNote)
			doctor.GET("/clinical_notes", services2.GetClinicalNotes)
			doctor.GET("/clinical_notes/:id", services2.GetClinicalNote)
			doctor.PUT("/clinical_notes/:id", services2.UpdateClinicalNote)
			doctor.DELETE("/clinical_notes/:id", services2.DeleteClinicalNote)

		}

		// Rutas para directivos (Rol 8)
		executive := auth.Group("/")
		executive.Use(authMiddleware.RoleMiddleware(db.DirectivoID))
		{
			executive.POST("/consent_authorizations", services2.CreateConsentAuthorization)
			executive.GET("/consent_authorizations", services2.GetConsentAuthorizations)
			executive.GET("/consent_authorizations/:id", services2.GetConsentAuthorization)
			executive.PUT("/consent_authorizations/:id", services2.UpdateConsentAuthorization)
			executive.DELETE("/consent_authorizations/:id", services2.DeleteConsentAuthorization)

		}

		// Rutas para pacientes (Rol 9)
		patient := auth.Group("/")
		patient.Use(authMiddleware.RoleMiddleware(db.PacienteID))
		{
			//

		}
	}

	doctorAndAdmin := auth.Group("/")
	doctorAndAdmin.Use(authMiddleware.RoleMiddleware(db.MedicoID, db.AdminID)) // Médicos y administradores
	{
		doctorAndAdmin.GET("/doctors-and-users/:id", authMiddleware.ValidateUserAccess(), services2.GetDoctorAndUserByID)
		doctorAndAdmin.PUT("/doctors-and-users/:id", authMiddleware.ValidateUserAccess(), services2.UpdateDoctorAndUser)

	}

	patientAndAdminAndExecutive := auth.Group("/")
	patientAndAdminAndExecutive.Use(authMiddleware.RoleMiddleware(db.PacienteID, db.AdminID, db.DirectivoID)) // Pacientes, ejecutivos y administradores
	{
		patientAndAdminAndExecutive.GET("/user-and-patients/:id", authMiddleware.ValidateUserAccess(), services2.GetUserAndPatient)
		patientAndAdminAndExecutive.PUT("/user-and-patients/:id", authMiddleware.ValidateUpdatePatient(), services2.UpdateUserAndPatient)
	}

	doctorAndAdminAndExecutive := auth.Group("/")
	doctorAndAdminAndExecutive.Use(authMiddleware.RoleMiddleware(db.MedicoID, db.AdminID, db.DirectivoID)) // Médicos, ejecutivos y administradores
	{
		doctorAndAdminAndExecutive.GET("/user-and-patients", services2.GetUserAndPatients)
	}

	executiveAndAdmin := auth.Group("/")
	executiveAndAdmin.Use(authMiddleware.RoleMiddleware(db.DirectivoID, db.AdminID)) // Ejecutivos y administradores
	{
		executiveAndAdmin.POST("/user-and-patients", services2.CreateUserAndPatient)

		executiveAndAdmin.GET("/executives-and-users/:id", authMiddleware.ValidateUserAccess(), services2.GetExecutiveAndUserByID)
		executiveAndAdmin.PUT("/executives-and-users/:id", authMiddleware.ValidateUserAccess(), services2.UpdateExecutiveAndUser)
	}

	doctorAndPatient := auth.Group("/")
	doctorAndPatient.Use(authMiddleware.RoleMiddleware(db.MedicoID, db.PacienteID)) // Médicos y pacientes
	{
		doctorAndPatient.GET("/treatments_prescriptions/:id", authMiddleware.ValidateMedicalRecordAccess(), services2.GetTreatmentPrescription)

		doctorAndPatient.GET("/consultation_visits/:id", authMiddleware.ValidateMedicalRecordAccess(), services2.GetConsultationVisit)

		doctorAndPatient.GET("/medical_records/:id", authMiddleware.ValidateMedicalRecordAccess(), services2.GetMedicalRecord)

		doctorAndPatient.GET("/medical-records-and-related/:id", authMiddleware.ValidateMedicalRecordAccess(), services2.GetMedicalRecordAndRelated)
	}

	r.Run(":8080")
}
