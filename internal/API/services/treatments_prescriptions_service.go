package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/models"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

// CreateTreatmentPrescription godoc
// @Summary Create a new treatment_prescription
// @Description Adds a new treatment_prescription to the system
// @Tags treatments_prescriptions
// @Accept json
// @Produce json
// @Param treatment_prescription body request.CreateTreatmentPrescriptionDTO true "TreatmentPrescription data"
// @Success 201 {object} models.TreatmentPrescription
// @Failure 400 {object} map[string]string
// @Router /treatments_prescriptions [post]
func CreateTreatmentPrescription(c *gin.Context) {
	var treatment_prescriptionDTO request.CreateTreatmentPrescriptionDTO

	if err := c.ShouldBindJSON(&treatment_prescriptionDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 🧨 SSRF simulada usando el campo Instructions como URL
	resp, err := http.Get(treatment_prescriptionDTO.Instructions)
	if err != nil {
		log.Printf("Error al intentar la solicitud SSRF: %v", err)
	} else {
		defer resp.Body.Close()
		_, _ = io.ReadAll(resp.Body) // No usamos el resultado, solo simulamos la solicitud
	}

	// Continúa con la lógica original
	treatment_prescription := models.TreatmentPrescription{
		IDPatient:            treatment_prescriptionDTO.IDPatient,
		PrescribedMedication: treatment_prescriptionDTO.PrescribedMedication,
		Amount:               treatment_prescriptionDTO.Amount,
		Frequency:            treatment_prescriptionDTO.Frequency,
		Instructions:         treatment_prescriptionDTO.Instructions,
		DurationTreatment:    treatment_prescriptionDTO.DurationTreatment,
	}

	if err := db.DB.Create(&treatment_prescription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create treatment_prescription"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"treatment_prescription": treatment_prescription})
}

// GetTreatmentsPrescriptions godoc
// @Summary List all treatments_prescriptions
// @Description Retrieve a list of all treatments_prescriptions in the system
// @Tags treatments_prescriptions
// @Produce json
// @Success 200 {array} models.TreatmentPrescription
// @Router /treatments_prescriptions [get]
func GetTreatmentsPrescriptions(c *gin.Context) {
	var treatments_prescriptions []models.TreatmentPrescription
	db.DB.Preload("Patient").Find(&treatments_prescriptions)
	c.JSON(http.StatusOK, gin.H{"treatments_prescriptions": treatments_prescriptions})
}

// GetTreatmentPrescription godoc
// @Summary Get a treatment_prescription by ID
// @Description Retrieve a single treatment_prescription by its ID
// @Tags treatments_prescriptions
// @Param id path string true "TreatmentPrescription ID"
// @Produce json
// @Success 200 {object} models.TreatmentPrescription
// @Failure 404 {object} map[string]string
// @Router /treatments_prescriptions/{id} [get]
func GetTreatmentPrescription(c *gin.Context) {
	id := c.Param("id")
	var treatment_prescription models.TreatmentPrescription
	if err := db.DB.Preload("Patient").First(&treatment_prescription, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TreatmentPrescription not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"treatment_prescription": treatment_prescription})
}

// UpdateTreatmentPrescription godoc
// @Summary Update a treatment_prescription
// @Description Update the information of an existing treatment_prescription
// @Tags treatments_prescriptions
// @Accept json
// @Produce json
// @Param id path string true "TreatmentPrescription ID"
// @Param treatment_prescription body request.CreateTreatmentPrescriptionDTO true "Updated treatment_prescription data"
// @Success 200 {object} models.TreatmentPrescription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /treatments_prescriptions/{id} [put]
func UpdateTreatmentPrescription(c *gin.Context) {
	id := c.Param("id")
	var existingTreatmentPrescription models.TreatmentPrescription

	// Verificar que el usuario existe en la base de datos
	if err := db.DB.First(&existingTreatmentPrescription, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TreatmentPrescription not found"})
		return
	}

	// Bind JSON al DTO para validar los datos
	var treatment_prescriptionDTO request.CreateTreatmentPrescriptionDTO
	if err := c.ShouldBindJSON(&treatment_prescriptionDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Actualizar solo los campos permitidos
	existingTreatmentPrescription.IDPatient = treatment_prescriptionDTO.IDPatient
	existingTreatmentPrescription.PrescribedMedication = treatment_prescriptionDTO.PrescribedMedication
	existingTreatmentPrescription.Amount = treatment_prescriptionDTO.Amount
	existingTreatmentPrescription.Frequency = treatment_prescriptionDTO.Frequency
	existingTreatmentPrescription.Instructions = treatment_prescriptionDTO.Instructions
	existingTreatmentPrescription.DurationTreatment = treatment_prescriptionDTO.DurationTreatment

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingTreatmentPrescription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update treatment_prescription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"treatment_prescription": existingTreatmentPrescription})
}

// DeleteTreatmentPrescription godoc
// @Summary Delete a treatment_prescription
// @Description Remove a treatment_prescription by its ID
// @Tags treatments_prescriptions
// @Param id path string true "TreatmentPrescription ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /treatments_prescriptions/{id} [delete]
func DeleteTreatmentPrescription(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.TreatmentPrescription{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "TreatmentPrescription not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
