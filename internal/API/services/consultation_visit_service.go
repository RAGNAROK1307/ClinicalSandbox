package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateConsultationVisit godoc
// @Summary Create a new consultation_visit
// @Description Adds a new consultation_visit to the system
// @Tags consultation_visits
// @Accept json
// @Produce json
// @Param consultation_visit body request.CreateConsultationVisitDTO true "ConsultationVisit data"
// @Success 201 {object} models.ConsultationVisit
// @Failure 400 {object} map[string]string
// @Router /consultation_visits [post]
func CreateConsultationVisit(c *gin.Context) {
	var consultation_visitDTO request.CreateConsultationVisitDTO

	// Bind JSON to consultation_visitDTO
	if err := c.ShouldBindJSON(&consultation_visitDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Mapea a modelo ConsultationVisit y usa la variable birthDate
	consultation_visit := models.ConsultationVisit{
		IDPatient:          consultation_visitDTO.IDPatient,
		IDHospitalEmployee: consultation_visitDTO.IDHospitalEmployee,
		DateTimeVisit:      consultation_visitDTO.DateTimeVisit,
		ReasonVisit:        consultation_visitDTO.ReasonVisit,
		MedicalNotes:       consultation_visitDTO.MedicalNotes,
		ExamResults:        consultation_visitDTO.ExamResults,
	}

	// Guarda el nuevo consultation_visit en la base de datos
	if err := db.DB.Create(&consultation_visit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create consultation_visit"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"consultation_visit": consultation_visit})
}

// GetConsultationVisits godoc
// @Summary List all consultation_visits
// @Description Retrieve a list of all consultation_visits in the system
// @Tags consultation_visits
// @Produce json
// @Success 200 {array} models.ConsultationVisit
// @Router /consultation_visits [get]
func GetConsultationVisits(c *gin.Context) {
	var consultation_visits []models.ConsultationVisit
	db.DB.Preload("Patient").Preload("HospitalEmployee").Find(&consultation_visits)
	c.JSON(http.StatusOK, gin.H{"consultation_visits": consultation_visits})
}

// GetConsultationVisit godoc
// @Summary Get a consultation_visit by ID
// @Description Retrieve a single consultation_visit by its ID
// @Tags consultation_visits
// @Param id path string true "ConsultationVisit ID"
// @Produce json
// @Success 200 {object} models.ConsultationVisit
// @Failure 404 {object} map[string]string
// @Router /consultation_visits/{id} [get]
func GetConsultationVisit(c *gin.Context) {
	id := c.Param("id")
	var consultation_visit models.ConsultationVisit
	if err := db.DB.Preload("Patient").Preload("HospitalEmployee").First(&consultation_visit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ConsultationVisit not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"consultation_visit": consultation_visit})
}

// UpdateConsultationVisit godoc
// @Summary Update a consultation_visit
// @Description Update the information of an existing consultation_visit
// @Tags consultation_visits
// @Accept json
// @Produce json
// @Param id path string true "ConsultationVisit ID"
// @Param consultation_visit body request.CreateConsultationVisitDTO true "Updated consultation_visit data"
// @Success 200 {object} models.ConsultationVisit
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /consultation_visits/{id} [put]
func UpdateConsultationVisit(c *gin.Context) {
	id := c.Param("id")
	var existingConsultationVisit models.ConsultationVisit

	// Verificar que el usuario existe en la base de datos
	if err := db.DB.First(&existingConsultationVisit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ConsultationVisit not found"})
		return
	}

	// Bind JSON al DTO para validar los datos
	var consultation_visitDTO request.CreateConsultationVisitDTO
	if err := c.ShouldBindJSON(&consultation_visitDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Actualizar solo los campos permitidos
	existingConsultationVisit.IDPatient = consultation_visitDTO.IDPatient
	existingConsultationVisit.IDHospitalEmployee = consultation_visitDTO.IDHospitalEmployee
	existingConsultationVisit.DateTimeVisit = consultation_visitDTO.DateTimeVisit
	existingConsultationVisit.ReasonVisit = consultation_visitDTO.ReasonVisit
	existingConsultationVisit.MedicalNotes = consultation_visitDTO.MedicalNotes
	existingConsultationVisit.ExamResults = consultation_visitDTO.ExamResults

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingConsultationVisit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update consultation_visit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"consultation_visit": existingConsultationVisit})
}

// DeleteConsultationVisit godoc
// @Summary Delete a consultation_visit
// @Description Remove a consultation_visit by its ID
// @Tags consultation_visits
// @Param id path string true "ConsultationVisit ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /consultation_visits/{id} [delete]
func DeleteConsultationVisit(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.ConsultationVisit{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ConsultationVisit not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
