package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
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
// @Success 201 {object} models.Laboratory
// @Failure 400 {object} map[string]string
// @Router /laboratories [post]
func CreateLaboratory(c *gin.Context) {
	var laboratoryDTO request.CreateLaboratoryDTO

	// Bind JSON to laboratoryDTO
	if err := c.ShouldBindJSON(&laboratoryDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha en formato string a time.Time
	testDate, err := time.Parse("2006-01-02", laboratoryDTO.TestDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Mapea a modelo Laboratory y usa la variable birthDate
	laboratory := models.Laboratory{
		IDConsultationVisit: laboratoryDTO.IDConsultationVisit,
		TestDate:            testDate, // Aquí estamos usando la variable testDate
		TestType:            laboratoryDTO.TestType,
		TestResults:         laboratoryDTO.TestResults,
		ExternalFilePath:    laboratoryDTO.ExternalFilePath,
	}

	// Guarda el nuevo laboratory en la base de datos
	if err := db.DB.Create(&laboratory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create laboratory"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"laboratory": laboratory})
}

// GetLaboratories godoc
// @Summary List all laboratories
// @Description Retrieve a list of all laboratories in the system
// @Tags laboratories
// @Produce json
// @Success 200 {array} models.Laboratory
// @Router /laboratories [get]
func GetLaboratories(c *gin.Context) {
	var laboratories []models.Laboratory
	db.DB.Preload("ConsultationVisit").Find(&laboratories)
	c.JSON(http.StatusOK, gin.H{"laboratories": laboratories})
}

// GetLaboratory godoc
// @Summary Get a laboratory by ID
// @Description Retrieve a single laboratory by its ID
// @Tags laboratories
// @Param id path string true "Laboratory ID"
// @Produce json
// @Success 200 {object} models.Laboratory
// @Failure 404 {object} map[string]string
// @Router /laboratories/{id} [get]
func GetLaboratory(c *gin.Context) {
	id := c.Param("id")
	var laboratory models.Laboratory
	if err := db.DB.Preload("ConsultationVisit").First(&laboratory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Laboratory not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"laboratory": laboratory})
}

// UpdateLaboratory godoc
// @Summary Update a laboratory
// @Description Update the information of an existing laboratory
// @Tags laboratories
// @Accept json
// @Produce json
// @Param id path string true "Laboratory ID"
// @Param laboratory body request.CreateLaboratoryDTO true "Updated laboratory data"
// @Success 200 {object} models.Laboratory
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /laboratories/{id} [put]
func UpdateLaboratory(c *gin.Context) {
	id := c.Param("id")
	var existingLaboratory models.Laboratory

	// Verificar que el usuario existe en la base de datos
	if err := db.DB.First(&existingLaboratory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Laboratory not found"})
		return
	}

	// Bind JSON al DTO para validar los datos
	var laboratoryDTO request.CreateLaboratoryDTO
	if err := c.ShouldBindJSON(&laboratoryDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha en formato string a time.Time
	testDate, err := time.Parse("2006-01-02", laboratoryDTO.TestDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Actualizar solo los campos permitidos
	existingLaboratory.IDConsultationVisit = laboratoryDTO.IDConsultationVisit
	existingLaboratory.TestDate = testDate // Usar la fecha convertida
	existingLaboratory.TestType = laboratoryDTO.TestType
	existingLaboratory.TestResults = laboratoryDTO.TestResults
	existingLaboratory.ExternalFilePath = laboratoryDTO.ExternalFilePath

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingLaboratory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update laboratory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"laboratory": existingLaboratory})
}

// DeleteLaboratory godoc
// @Summary Delete a laboratory
// @Description Remove a laboratory by its ID
// @Tags laboratories
// @Param id path string true "Laboratory ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /laboratories/{id} [delete]
func DeleteLaboratory(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Laboratory{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Laboratory not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
