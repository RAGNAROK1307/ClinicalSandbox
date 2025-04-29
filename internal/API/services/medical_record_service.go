package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateMedicalRecord godoc
// @Summary Create a new medical_record
// @Description Adds a new medical_record to the system
// @Tags medical_records
// @Accept json
// @Produce json
// @Param medical_record body request.CreateMedicalRecordDTO true "MedicalRecord data"
// @Success 201 {object} models.MedicalRecord
// @Failure 400 {object} map[string]string
// @Router /medical_records [post]
func CreateMedicalRecord(c *gin.Context) {
	var medical_recordDTO request.CreateMedicalRecordDTO

	// Bind JSON to medical_recordDTO
	if err := c.ShouldBindJSON(&medical_recordDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Mapea a modelo MedicalRecord y usa la variable birthDate
	medical_record := models.MedicalRecord{
		IDPatient: medical_recordDTO.IDPatient,
		//IDConsultationVisit:     medical_recordDTO.IDConsultationVisit,
		//IDTreatmentPrescription: medical_recordDTO.IDTreatmentPrescription,
		PreviousDiagnoses:   medical_recordDTO.PreviousDiagnoses,
		ChronicDiseases:     medical_recordDTO.ChronicDiseases,
		Allergies:           medical_recordDTO.Allergies,
		CurrentMedications:  medical_recordDTO.CurrentMedications,
		SurgeryHistory:      medical_recordDTO.SurgeryHistory,
		FamilyHistory:       medical_recordDTO.FamilyHistory,
		CreationDateHistory: medical_recordDTO.CreationDateHistory,
		HistoryUpdateDate:   medical_recordDTO.HistoryUpdateDate,
	}

	// Guarda el nuevo medical_record en la base de datos
	if err := db.DB.Create(&medical_record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medical_record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"medical_record": medical_record})
}

// GetMedicalRecords godoc
// @Summary List all medical_records
// @Description Retrieve a list of all medical_records in the system
// @Tags medical_records
// @Produce json
// @Success 200 {array} models.MedicalRecord
// @Router /medical_records [get]
func GetMedicalRecords(c *gin.Context) {
	var medical_records []models.MedicalRecord
	db.DB.Preload("Patient").Find(&medical_records)
	c.JSON(http.StatusOK, gin.H{"medical_records": medical_records})
}

// GetMedicalRecord godoc
// @Summary Get a medical_record by ID
// @Description Retrieve a single medical_record by its ID
// @Tags medical_records
// @Param id path string true "MedicalRecord ID"
// @Produce json
// @Success 200 {object} models.MedicalRecord
// @Failure 404 {object} map[string]string
// @Router /medical_records/{id} [get]
func GetMedicalRecord(c *gin.Context) {
	id := c.Param("id")
	var medical_record models.MedicalRecord
	if err := db.DB.Preload("Patient").First(&medical_record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MedicalRecord not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"medical_record": medical_record})
}

// UpdateMedicalRecord godoc
// @Summary Update a medical_record
// @Description Update the information of an existing medical_record
// @Tags medical_records
// @Accept json
// @Produce json
// @Param id path string true "MedicalRecord ID"
// @Param medical_record body request.CreateMedicalRecordDTO true "Updated medical_record data"
// @Success 200 {object} models.MedicalRecord
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /medical_records/{id} [put]
func UpdateMedicalRecord(c *gin.Context) {
	id := c.Param("id")
	var existingMedicalRecord models.MedicalRecord

	// Verificar que el usuario existe en la base de datos
	if err := db.DB.First(&existingMedicalRecord, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MedicalRecord not found"})
		return
	}

	// Bind JSON al DTO para validar los datos
	var medical_recordDTO request.CreateMedicalRecordDTO
	if err := c.ShouldBindJSON(&medical_recordDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Actualizar solo los campos permitidos
	existingMedicalRecord.IDPatient = medical_recordDTO.IDPatient
	//existingMedicalRecord.IDConsultationVisit = medical_recordDTO.IDConsultationVisit
	//existingMedicalRecord.IDTreatmentPrescription = medical_recordDTO.IDTreatmentPrescription
	existingMedicalRecord.PreviousDiagnoses = medical_recordDTO.PreviousDiagnoses
	existingMedicalRecord.ChronicDiseases = medical_recordDTO.ChronicDiseases
	existingMedicalRecord.Allergies = medical_recordDTO.Allergies
	existingMedicalRecord.CurrentMedications = medical_recordDTO.CurrentMedications
	existingMedicalRecord.SurgeryHistory = medical_recordDTO.SurgeryHistory
	existingMedicalRecord.FamilyHistory = medical_recordDTO.FamilyHistory
	existingMedicalRecord.CreationDateHistory = medical_recordDTO.CreationDateHistory
	existingMedicalRecord.HistoryUpdateDate = medical_recordDTO.HistoryUpdateDate

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingMedicalRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update medical_record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"medical_record": existingMedicalRecord})
}

// DeleteMedicalRecord godoc
// @Summary Delete a medical_record
// @Description Remove a medical_record by its ID
// @Tags medical_records
// @Param id path string true "MedicalRecord ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /medical_records/{id} [delete]
func DeleteMedicalRecord(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.MedicalRecord{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MedicalRecord not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
