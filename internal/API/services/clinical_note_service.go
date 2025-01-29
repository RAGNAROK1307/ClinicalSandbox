package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateClinicalNote godoc
// @Summary Create a new clinical_note
// @Description Adds a new clinical_note to the system
// @Tags clinical_notes
// @Accept json
// @Produce json
// @Param clinical_note body request.CreateClinicalNoteDTO true "ClinicalNote data"
// @Success 201 {object} models.ClinicalNote
// @Failure 400 {object} map[string]string
// @Router /clinical_notes [post]
func CreateClinicalNote(c *gin.Context) {
	var clinical_noteDTO request.CreateClinicalNoteDTO

	// Bind JSON to clinical_noteDTO
	if err := c.ShouldBindJSON(&clinical_noteDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Mapea a modelo ClinicalNote y usa la variable birthDate
	clinical_note := models.ClinicalNote{
		IDConsultationVisit:         clinical_noteDTO.IDConsultationVisit,
		ProgressNotes:               clinical_noteDTO.ProgressNotes,
		ObservationsRecommendations: clinical_noteDTO.ObservationsRecommendations,
	}

	// Guarda el nuevo clinical_note en la base de datos
	if err := db.DB.Create(&clinical_note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create clinical_note"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"clinical_note": clinical_note})
}

// GetClinicalNotes godoc
// @Summary List all clinical_notes
// @Description Retrieve a list of all clinical_notes in the system
// @Tags clinical_notes
// @Produce json
// @Success 200 {array} models.ClinicalNote
// @Router /clinical_notes [get]
func GetClinicalNotes(c *gin.Context) {
	var clinical_notes []models.ClinicalNote
	db.DB.Preload("ConsultationVisit").Find(&clinical_notes)
	c.JSON(http.StatusOK, gin.H{"clinical_notes": clinical_notes})
}

// GetClinicalNote godoc
// @Summary Get a clinical_note by ID
// @Description Retrieve a single clinical_note by its ID
// @Tags clinical_notes
// @Param id path string true "ClinicalNote ID"
// @Produce json
// @Success 200 {object} models.ClinicalNote
// @Failure 404 {object} map[string]string
// @Router /clinical_notes/{id} [get]
func GetClinicalNote(c *gin.Context) {
	id := c.Param("id")
	var clinical_note models.ClinicalNote
	if err := db.DB.Preload("ConsultationVisit").First(&clinical_note, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ClinicalNote not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"clinical_note": clinical_note})
}

// UpdateClinicalNote godoc
// @Summary Update a clinical_note
// @Description Update the information of an existing clinical_note
// @Tags clinical_notes
// @Accept json
// @Produce json
// @Param id path string true "ClinicalNote ID"
// @Param clinical_note body request.CreateClinicalNoteDTO true "Updated clinical_note data"
// @Success 200 {object} models.ClinicalNote
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /clinical_notes/{id} [put]
func UpdateClinicalNote(c *gin.Context) {
	id := c.Param("id")
	var existingClinicalNote models.ClinicalNote

	// Verificar que el usuario existe en la base de datos
	if err := db.DB.First(&existingClinicalNote, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ClinicalNote not found"})
		return
	}

	// Bind JSON al DTO para validar los datos
	var clinical_noteDTO request.CreateClinicalNoteDTO
	if err := c.ShouldBindJSON(&clinical_noteDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Actualizar solo los campos permitidos
	existingClinicalNote.IDConsultationVisit = clinical_noteDTO.IDConsultationVisit
	existingClinicalNote.ProgressNotes = clinical_noteDTO.ProgressNotes
	existingClinicalNote.ObservationsRecommendations = clinical_noteDTO.ObservationsRecommendations

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingClinicalNote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update clinical_note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clinical_note": existingClinicalNote})
}

// DeleteClinicalNote godoc
// @Summary Delete a clinical_note
// @Description Remove a clinical_note by its ID
// @Tags clinical_notes
// @Param id path string true "ClinicalNote ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /clinical_notes/{id} [delete]
func DeleteClinicalNote(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.ClinicalNote{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ClinicalNote not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
