package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/dto/response"
	"ClinicalSandBox/internal/API/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// CreateMedicalRecordAndRelated godoc
// @Summary Create a new medical_record, treatment_prescription, and consultation_visit
// @Description Adds a new medical_record, treatment_prescription, and consultation_visit to the system in a single transaction
// @Tags medical_records
// @Accept json
// @Produce json
// @Param medical_record_and_related body request.CreateMedicalRecordAndRelatedDTO true "MedicalRecord, TreatmentPrescription, and ConsultationVisit data"
// @Success 201 {object} response.MedicalRecordAndRelatedResponseDTO
// @Failure 400 {object} map[string]string
// @Router /medical-records-and-related [post]
func CreateMedicalRecordAndRelated(c *gin.Context) {
	var medicalRecordAndRelatedDTO request.CreateMedicalRecordAndRelatedDTO

	// Bind JSON al DTO
	if err := c.ShouldBindJSON(&medicalRecordAndRelatedDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha y hora de la visita a time.Time
	dateTimeVisit, err := time.Parse("2006-01-02 15:04:05", medicalRecordAndRelatedDTO.DateTimeVisit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD HH:MM:SS"})
		return
	}

	// Iniciar una transacción
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Crear ConsultationVisit
	consultationVisit := models.ConsultationVisit{
		IDPatient:          medicalRecordAndRelatedDTO.IDPatient,
		IDHospitalEmployee: medicalRecordAndRelatedDTO.IDHospitalEmployee,
		DateTimeVisit:      dateTimeVisit,
		ReasonVisit:        medicalRecordAndRelatedDTO.ReasonVisit,
		MedicalNotes:       medicalRecordAndRelatedDTO.MedicalNotes,
		ExamResults:        medicalRecordAndRelatedDTO.ExamResults,
	}

	if err := tx.Create(&consultationVisit).Error; err != nil {
		log.Printf("Error creating consultation visit: %v", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create consultation visit"})
		return
	}

	// Crear TreatmentPrescription
	treatmentPrescription := models.TreatmentPrescription{
		IDPatient:            medicalRecordAndRelatedDTO.IDPatient,
		PrescribedMedication: medicalRecordAndRelatedDTO.PrescribedMedication,
		Amount:               medicalRecordAndRelatedDTO.Amount,
		Frequency:            medicalRecordAndRelatedDTO.Frequency,
		Instructions:         medicalRecordAndRelatedDTO.Instructions,
		DurationTreatment:    medicalRecordAndRelatedDTO.DurationTreatment,
	}

	if err := tx.Create(&treatmentPrescription).Error; err != nil {
		log.Printf("Error creating treatment prescription: %v", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create treatment prescription"})
		return
	}

	// Crear MedicalRecord
	medicalRecord := models.MedicalRecord{
		IDPatient:           medicalRecordAndRelatedDTO.IDPatient,
		PreviousDiagnoses:   medicalRecordAndRelatedDTO.PreviousDiagnoses,
		ChronicDiseases:     medicalRecordAndRelatedDTO.ChronicDiseases,
		Allergies:           medicalRecordAndRelatedDTO.Allergies,
		CurrentMedications:  medicalRecordAndRelatedDTO.CurrentMedications,
		SurgeryHistory:      medicalRecordAndRelatedDTO.SurgeryHistory,
		FamilyHistory:       medicalRecordAndRelatedDTO.FamilyHistory,
		CreationDateHistory: time.Now(),
		HistoryUpdateDate:   time.Now(),
	}

	if err := tx.Create(&medicalRecord).Error; err != nil {
		log.Printf("Error creating medical record: %v", err)
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medical record"})
		return
	}

	// Confirmar la transacción
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Mapear los modelos a los DTOs de respuesta
	medicalRecordResponse := response.MedicalRecordResponseDTO{
		IDPatient:           medicalRecord.IDPatient,
		PreviousDiagnoses:   medicalRecord.PreviousDiagnoses,
		ChronicDiseases:     medicalRecord.ChronicDiseases,
		Allergies:           medicalRecord.Allergies,
		CurrentMedications:  medicalRecord.CurrentMedications,
		SurgeryHistory:      medicalRecord.SurgeryHistory,
		FamilyHistory:       medicalRecord.FamilyHistory,
		CreationDateHistory: medicalRecord.CreationDateHistory,
		HistoryUpdateDate:   medicalRecord.HistoryUpdateDate,
	}

	treatmentPrescriptionsResponse := []response.TreatmentPrescriptionResponseDTO{
		{
			IDPatient:            treatmentPrescription.IDPatient,
			PrescribedMedication: treatmentPrescription.PrescribedMedication,
			Amount:               treatmentPrescription.Amount,
			Frequency:            treatmentPrescription.Frequency,
			Instructions:         treatmentPrescription.Instructions,
			DurationTreatment:    treatmentPrescription.DurationTreatment,
		},
	}

	consultationVisitsResponse := []response.ConsultationVisitResponseDTO{
		{
			IDPatient:                consultationVisit.IDPatient,
			IDHospitalEmployee:       consultationVisit.IDHospitalEmployee,
			HospitalEmployeeFullName: consultationVisit.HospitalEmployee.FullName,
			DateTimeVisit:            consultationVisit.DateTimeVisit,
			ReasonVisit:              consultationVisit.ReasonVisit,
			MedicalNotes:             consultationVisit.MedicalNotes,
			ExamResults:              consultationVisit.ExamResults,
		},
	}

	// Generar la respuesta
	responseDTO := response.MedicalRecordAndRelatedResponseDTO{
		MedicalRecord:          medicalRecordResponse,
		TreatmentPrescriptions: treatmentPrescriptionsResponse,
		ConsultationVisits:     consultationVisitsResponse,
	}

	c.JSON(http.StatusCreated, responseDTO)
}

// GetMedicalRecordsAndRelated godoc
// @Summary List all medical records and related data
// @Description Retrieve a list of all medical records with their related treatment prescriptions and consultation visits
// @Tags medical_records
// @Produce json
// @Success 200 {array} response.MedicalRecordAndRelatedResponseDTO
// @Router /medical-records-and-related [get]
func GetMedicalRecordsAndRelated(c *gin.Context) {
	var medicalRecords []models.MedicalRecord

	// Obtener todos los registros médicos con sus relaciones
	if err := db.DB.
		Preload("Patient").
		Preload("ConsultationVisits.HospitalEmployee").
		Preload("TreatmentPrescriptions").
		Find(&medicalRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch medical records"})
		return
	}

	var medicalRecordsResponse []response.MedicalRecordAndRelatedResponseDTO
	for _, medicalRecord := range medicalRecords {
		var consultationVisitsResponse []response.ConsultationVisitResponseDTO
		for _, consultationVisit := range medicalRecord.ConsultationVisits {
			consultationVisitsResponse = append(consultationVisitsResponse, response.ConsultationVisitResponseDTO{
				IDPatient:                consultationVisit.IDPatient,
				IDConsultationVisit:      consultationVisit.IDConsultationVisit, // Usar el ID generado
				IDHospitalEmployee:       consultationVisit.IDHospitalEmployee,
				HospitalEmployeeFullName: consultationVisit.HospitalEmployee.FullName,
				DateTimeVisit:            consultationVisit.DateTimeVisit,
				ReasonVisit:              consultationVisit.ReasonVisit,
				MedicalNotes:             consultationVisit.MedicalNotes,
				ExamResults:              consultationVisit.ExamResults,
			})
		}

		var treatmentPrescriptionsResponse []response.TreatmentPrescriptionResponseDTO
		for _, treatmentPrescription := range medicalRecord.TreatmentPrescriptions {
			treatmentPrescriptionsResponse = append(treatmentPrescriptionsResponse, response.TreatmentPrescriptionResponseDTO{
				IDPatient:               treatmentPrescription.IDPatient,
				IDTreatmentPrescription: treatmentPrescription.IDTreatmentPrescription, // Usar el ID generado
				PrescribedMedication:    treatmentPrescription.PrescribedMedication,
				Amount:                  treatmentPrescription.Amount,
				Frequency:               treatmentPrescription.Frequency,
				Instructions:            treatmentPrescription.Instructions,
				DurationTreatment:       treatmentPrescription.DurationTreatment,
			})
		}

		medicalRecordsResponse = append(medicalRecordsResponse, response.MedicalRecordAndRelatedResponseDTO{
			MedicalRecord: response.MedicalRecordResponseDTO{
				IDPatient:           medicalRecord.IDPatient,
				PreviousDiagnoses:   medicalRecord.PreviousDiagnoses,
				ChronicDiseases:     medicalRecord.ChronicDiseases,
				Allergies:           medicalRecord.Allergies,
				CurrentMedications:  medicalRecord.CurrentMedications,
				SurgeryHistory:      medicalRecord.SurgeryHistory,
				FamilyHistory:       medicalRecord.FamilyHistory,
				CreationDateHistory: medicalRecord.CreationDateHistory,
				HistoryUpdateDate:   medicalRecord.HistoryUpdateDate,
			},
			TreatmentPrescriptions: treatmentPrescriptionsResponse,
			ConsultationVisits:     consultationVisitsResponse,
			Patient: response.PatientResponseDTO{
				IDPatient:            medicalRecord.Patient.IDPatient,
				FullName:             medicalRecord.Patient.FullName,
				BirthDate:            medicalRecord.Patient.BirthDate,
				Gender:               medicalRecord.Patient.Gender,
				Address:              medicalRecord.Patient.Address,
				Phone:                medicalRecord.Patient.Phone,
				SocialSecurityNumber: medicalRecord.Patient.SocialSecurityNumber,
				//RoleName:             medicalRecord.Patient.Role.RoleName,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{"medical_records": medicalRecordsResponse})
}

// GetMedicalRecordAndRelated godoc
// @Summary Get a medical record and related data by patient ID
// @Description Retrieve a single medical record with its related treatment prescription and consultation visit by patient ID
// @Tags medical_records
// @Param id path string true "Patient ID"
// @Produce json
// @Success 200 {object} response.MedicalRecordAndRelatedResponseDTO
// @Failure 404 {object} map[string]string
// @Router /medical-records-and-related/{id} [get]
func GetMedicalRecordAndRelated(c *gin.Context) {
	idPatient := c.Param("id") // Obtener el ID del paciente desde la URL
	var medicalRecord models.MedicalRecord

	// Buscar el historial médico asociado al paciente
	if err := db.DB.
		Where("id_paciente = ?", idPatient).
		Preload("Patient").
		Preload("ConsultationVisits.HospitalEmployee").
		Preload("TreatmentPrescriptions").
		First(&medicalRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found for the patient"})
		return
	}

	var consultationVisitsResponse []response.ConsultationVisitResponseDTO
	for _, consultationVisit := range medicalRecord.ConsultationVisits {
		consultationVisitsResponse = append(consultationVisitsResponse, response.ConsultationVisitResponseDTO{
			IDPatient:                consultationVisit.IDPatient,
			IDConsultationVisit:      consultationVisit.IDConsultationVisit, // Usar el ID generado
			IDHospitalEmployee:       consultationVisit.IDHospitalEmployee,
			HospitalEmployeeFullName: consultationVisit.HospitalEmployee.FullName,
			DateTimeVisit:            consultationVisit.DateTimeVisit,
			ReasonVisit:              consultationVisit.ReasonVisit,
			MedicalNotes:             consultationVisit.MedicalNotes,
			ExamResults:              consultationVisit.ExamResults,
		})
	}

	var treatmentPrescriptionsResponse []response.TreatmentPrescriptionResponseDTO
	for _, treatmentPrescription := range medicalRecord.TreatmentPrescriptions {
		treatmentPrescriptionsResponse = append(treatmentPrescriptionsResponse, response.TreatmentPrescriptionResponseDTO{
			IDPatient:               treatmentPrescription.IDPatient,
			IDTreatmentPrescription: treatmentPrescription.IDTreatmentPrescription, // Usar el ID generado
			PrescribedMedication:    treatmentPrescription.PrescribedMedication,
			Amount:                  treatmentPrescription.Amount,
			Frequency:               treatmentPrescription.Frequency,
			Instructions:            treatmentPrescription.Instructions,
			DurationTreatment:       treatmentPrescription.DurationTreatment,
		})
	}

	// Generar la respuesta
	medicalRecordResponse := response.MedicalRecordAndRelatedResponseDTO{
		MedicalRecord: response.MedicalRecordResponseDTO{
			IDPatient:           medicalRecord.IDPatient,
			PatientImage:        medicalRecord.PatientImage,
			PreviousDiagnoses:   medicalRecord.PreviousDiagnoses,
			ChronicDiseases:     medicalRecord.ChronicDiseases,
			Allergies:           medicalRecord.Allergies,
			CurrentMedications:  medicalRecord.CurrentMedications,
			SurgeryHistory:      medicalRecord.SurgeryHistory,
			FamilyHistory:       medicalRecord.FamilyHistory,
			CreationDateHistory: medicalRecord.CreationDateHistory,
			HistoryUpdateDate:   medicalRecord.HistoryUpdateDate,
		},
		TreatmentPrescriptions: treatmentPrescriptionsResponse,
		ConsultationVisits:     consultationVisitsResponse,
		Patient: response.PatientResponseDTO{
			IDPatient:            medicalRecord.Patient.IDPatient,
			FullName:             medicalRecord.Patient.FullName,
			BirthDate:            medicalRecord.Patient.BirthDate,
			Gender:               medicalRecord.Patient.Gender,
			Address:              medicalRecord.Patient.Address,
			Phone:                medicalRecord.Patient.Phone,
			SocialSecurityNumber: medicalRecord.Patient.SocialSecurityNumber,
			//RoleName:             medicalRecord.Patient.Role.RoleName,
		},
	}

	c.JSON(http.StatusOK, gin.H{"medical_record": medicalRecordResponse})
}

// UpdateMedicalRecordAndRelated godoc
// @Summary Update a medical_record
// @Description Update the information of an existing medical_record
// @Tags medical_records
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Param medical_record body request.UpdateMedicalRecordDTO true "Updated medical_record data"
// @Success 200 {object} models.MedicalRecord
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /medical-records-and-related/{id} [put]
func UpdateMedicalRecordAndRelated(c *gin.Context) {
	// Obtener el ID del paciente de los parámetros
	idPatient := c.Param("id")

	// Buscar el historial médico por ID de paciente
	var existingMedicalRecord models.MedicalRecord
	if err := db.DB.Where("id_paciente = ?", idPatient).First(&existingMedicalRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found for this patient"})
		return
	}

	// Bind JSON al DTO (sin validación estricta para permitir actualizaciones parciales)
	var updateDTO request.UpdateMedicalRecordDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Actualizar solo los campos que vienen en la petición
	if updateDTO.PreviousDiagnoses != "" {
		existingMedicalRecord.PreviousDiagnoses = updateDTO.PreviousDiagnoses
	}
	if updateDTO.ChronicDiseases != "" {
		existingMedicalRecord.ChronicDiseases = updateDTO.ChronicDiseases
	}
	if updateDTO.Allergies != "" {
		existingMedicalRecord.Allergies = updateDTO.Allergies
	}
	if updateDTO.CurrentMedications != "" {
		existingMedicalRecord.CurrentMedications = updateDTO.CurrentMedications
	}
	if updateDTO.SurgeryHistory != "" {
		existingMedicalRecord.SurgeryHistory = updateDTO.SurgeryHistory
	}
	if updateDTO.FamilyHistory != "" {
		existingMedicalRecord.FamilyHistory = updateDTO.FamilyHistory
	}
	if !updateDTO.CreationDateHistory.IsZero() {
		existingMedicalRecord.CreationDateHistory = updateDTO.CreationDateHistory
	}

	// Siempre actualizar la fecha de modificación
	existingMedicalRecord.HistoryUpdateDate = time.Now()

	// Guardar los cambios
	if err := db.DB.Save(&existingMedicalRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update medical record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Medical record updated successfully",
		"data":    existingMedicalRecord,
	})
}
