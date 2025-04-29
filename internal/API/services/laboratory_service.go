package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/dto/response"
	"ClinicalSandBox/internal/API/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// CreateLaboratory godoc
// @Summary Create a new laboratory
// @Description Adds a new laboratory to the system
// @Tags laboratories
// @Accept json
// @Produce json
// @Param laboratory body request.CreateLaboratoryDTO true "Laboratory data"
// @Success 201 {object} response.LaboratoryResponseDTO
// @Failure 400 {object} map[string]string
// @Router /laboratories [post]
func CreateLaboratory(c *gin.Context) {
	var laboratoryDTO request.CreateLaboratoryDTO

	if err := c.ShouldBindJSON(&laboratoryDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	testDate, err := time.Parse("2006-01-02", laboratoryDTO.TestDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	laboratory := models.Laboratory{
		IDConsultationVisit: laboratoryDTO.IDConsultationVisit,
		TestDate:            testDate,
		TestType:            laboratoryDTO.TestType,
		TestResults:         laboratoryDTO.TestResults,
		//ExternalFilePath:    laboratoryDTO.ExternalFilePath,
	}

	if err := db.DB.Create(&laboratory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create laboratory"})
		return
	}

	// Convertir a DTO de respuesta
	responseDTO := response.LaboratoryResponseDTO{
		IDLaboratory:        laboratory.IDLaboratory,
		IDConsultationVisit: laboratory.IDConsultationVisit,
		TestDate:            laboratory.TestDate.Format("2006-01-02"),
		TestType:            laboratory.TestType,
		TestResults:         laboratory.TestResults,
		ExternalFilePath:    laboratory.ExternalFilePath,
	}

	c.JSON(http.StatusCreated, gin.H{"laboratory": responseDTO})
}

// GetLaboratories godoc
// @Summary List all laboratories
// @Description Retrieve a list of all laboratories in the system
// @Tags laboratories
// @Produce json
// @Success 200 {array} response.LaboratoryResponseDTO
// @Router /laboratories [get]
func GetLaboratories(c *gin.Context) {
	var laboratories []models.Laboratory
	db.DB.Preload("ConsultationVisit").Find(&laboratories)

	// Convertir a DTOs de respuesta
	var responseDTOs []response.LaboratoryResponseDTO
	for _, lab := range laboratories {
		responseDTOs = append(responseDTOs, response.LaboratoryResponseDTO{
			IDLaboratory:        lab.IDLaboratory,
			IDConsultationVisit: lab.IDConsultationVisit,
			TestDate:            lab.TestDate.Format("2006-01-02"),
			TestType:            lab.TestType,
			TestResults:         lab.TestResults,
			ExternalFilePath:    lab.ExternalFilePath,
		})
	}

	c.JSON(http.StatusOK, gin.H{"laboratories": responseDTOs})
}

// GetLaboratory godoc
// @Summary Get a laboratory by consultation ID
// @Description Retrieve a single laboratory by its consultation ID
// @Tags laboratories
// @Param id path string true "ConsultationVisit ID"
// @Produce json
// @Success 200 {object} response.LaboratoryResponseDTO
// @Failure 404 {object} map[string]string
// @Router /laboratories/{id} [get]
func GetLaboratory(c *gin.Context) {
	consultationID := c.Param("id")
	var laboratory models.Laboratory

	if err := db.DB.Preload("ConsultationVisit").
		Where("id_consulta = ?", consultationID).
		First(&laboratory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Laboratory not found for this consultation"})
		return
	}

	// Convertir a DTO de respuesta
	responseDTO := response.LaboratoryResponseDTO{
		IDLaboratory:        laboratory.IDLaboratory,
		IDConsultationVisit: laboratory.IDConsultationVisit,
		TestDate:            laboratory.TestDate.Format("2006-01-02"),
		TestType:            laboratory.TestType,
		TestResults:         laboratory.TestResults,
		ExternalFilePath:    laboratory.ExternalFilePath,
	}

	c.JSON(http.StatusOK, gin.H{"laboratory": responseDTO})
}

// UpdateLaboratory godoc
// @Summary Update a laboratory by consultation ID
// @Description Update the information of an existing laboratory by consultation ID
// @Tags laboratories
// @Accept json
// @Produce json
// @Param id path string true "ConsultationVisit ID"
// @Param laboratory body request.CreateLaboratoryDTO true "Updated laboratory data"
// @Success 200 {object} response.LaboratoryResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /laboratories/{id} [put]
func UpdateLaboratory(c *gin.Context) {
	consultationID := c.Param("id")
	var existingLaboratory models.Laboratory

	if err := db.DB.Where("id_consulta = ?", consultationID).
		First(&existingLaboratory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Laboratory not found for this consultation"})
		return
	}

	var laboratoryDTO request.CreateLaboratoryDTO
	if err := c.ShouldBindJSON(&laboratoryDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	testDate, err := time.Parse("2006-01-02", laboratoryDTO.TestDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Actualizar campos
	existingLaboratory.IDConsultationVisit = laboratoryDTO.IDConsultationVisit
	existingLaboratory.TestDate = testDate
	existingLaboratory.TestType = laboratoryDTO.TestType
	existingLaboratory.TestResults = laboratoryDTO.TestResults
	//existingLaboratory.ExternalFilePath = laboratoryDTO.ExternalFilePath

	if err := db.DB.Save(&existingLaboratory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update laboratory"})
		return
	}

	// Convertir a DTO de respuesta
	responseDTO := response.LaboratoryResponseDTO{
		IDLaboratory:        existingLaboratory.IDLaboratory,
		IDConsultationVisit: existingLaboratory.IDConsultationVisit,
		TestDate:            existingLaboratory.TestDate.Format("2006-01-02"),
		TestType:            existingLaboratory.TestType,
		TestResults:         existingLaboratory.TestResults,
		ExternalFilePath:    existingLaboratory.ExternalFilePath,
	}

	c.JSON(http.StatusOK, gin.H{"laboratory": responseDTO})
}

// DeleteLaboratory godoc
// @Summary Delete a laboratory by consultation ID
// @Description Remove a laboratory by its consultation ID
// @Tags laboratories
// @Param id path string true "ConsultationVisit ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /laboratories/{id} [delete]
func DeleteLaboratory(c *gin.Context) {
	consultationID := c.Param("id")

	// Buscar por id_consulta en lugar de ID directo
	var laboratory models.Laboratory
	if err := db.DB.Where("id_consulta = ?", consultationID).First(&laboratory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Laboratory not found for this consultation"})
		return
	}

	if err := db.DB.Delete(&laboratory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete laboratory"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
