package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/dto/response"
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
// @Success 201 {object} response.ClinicalResponseNoteDTO
// @Failure 400 {object} map[string]string
// @Router /clinical_notes [post]
func CreateClinicalNote(c *gin.Context) {
	var clinicalNoteDTO request.CreateClinicalNoteDTO

	if err := c.ShouldBindJSON(&clinicalNoteDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	clinicalNote := models.ClinicalNote{
		IDConsultationVisit:         clinicalNoteDTO.IDConsultationVisit,
		ProgressNotes:               clinicalNoteDTO.ProgressNotes,
		ObservationsRecommendations: clinicalNoteDTO.ObservationsRecommendations,
	}

	if err := db.DB.Create(&clinicalNote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create clinical note"})
		return
	}

	// Convertir a DTO de respuesta
	responseDTO := response.ClinicalResponseNoteDTO{
		IDConsultationVisit:         clinicalNote.IDConsultationVisit,
		ProgressNotes:               clinicalNote.ProgressNotes,
		ObservationsRecommendations: clinicalNote.ObservationsRecommendations,
	}

	c.JSON(http.StatusCreated, gin.H{"clinical_note": responseDTO})
}

// GetClinicalNotes godoc
// @Summary List all clinical_notes
// @Description Retrieve a list of all clinical_notes in the system
// @Tags clinical_notes
// @Produce json
// @Success 200 {array} response.ClinicalResponseNoteDTO
// @Router /clinical_notes [get]
func GetClinicalNotes(c *gin.Context) {
	var clinicalNotes []models.ClinicalNote
	db.DB.Preload("ConsultationVisit").Find(&clinicalNotes)

	// Convertir a DTOs de respuesta
	var responseDTOs []response.ClinicalResponseNoteDTO
	for _, note := range clinicalNotes {
		responseDTOs = append(responseDTOs, response.ClinicalResponseNoteDTO{
			IDConsultationVisit:         note.IDConsultationVisit,
			ProgressNotes:               note.ProgressNotes,
			ObservationsRecommendations: note.ObservationsRecommendations,
		})
	}

	c.JSON(http.StatusOK, gin.H{"clinical_notes": responseDTOs})
}

// GetClinicalNote godoc
// @Summary Get a clinical_note by consultation ID
// @Description Retrieve a single clinical_note by its consultation ID
// @Tags clinical_notes
// @Param id path string true "ConsultationVisit ID"
// @Produce json
// @Success 200 {object} response.ClinicalResponseNoteDTO
// @Failure 404 {object} map[string]string
// @Router /clinical_notes/{id} [get]
func GetClinicalNote(c *gin.Context) {
	consultationID := c.Param("id")
	var clinicalNote models.ClinicalNote

	if err := db.DB.Preload("ConsultationVisit").
		Where("id_consulta = ?", consultationID).
		First(&clinicalNote).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Clinical note not found for this consultation"})
		return
	}

	// Convertir a DTO de respuesta
	responseDTO := response.ClinicalResponseNoteDTO{
		IDConsultationVisit:         clinicalNote.IDConsultationVisit,
		ProgressNotes:               clinicalNote.ProgressNotes,
		ObservationsRecommendations: clinicalNote.ObservationsRecommendations,
	}

	c.JSON(http.StatusOK, gin.H{"clinical_note": responseDTO})
}

// UpdateClinicalNote godoc
// @Summary Update a clinical_note by consultation ID
// @Description Update the information of an existing clinical_note by consultation ID
// @Tags clinical_notes
// @Accept json
// @Produce json
// @Param id path string true "ConsultationVisit ID"
// @Param clinical_note body request.CreateClinicalNoteDTO true "Updated clinical_note data"
// @Success 200 {object} response.ClinicalResponseNoteDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /clinical_notes/{id} [put]
func UpdateClinicalNote(c *gin.Context) {
	consultationID := c.Param("id")
	var existingClinicalNote models.ClinicalNote

	if err := db.DB.Where("id_consulta = ?", consultationID).
		First(&existingClinicalNote).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Clinical note not found for this consultation"})
		return
	}

	var clinicalNoteDTO request.CreateClinicalNoteDTO
	if err := c.ShouldBindJSON(&clinicalNoteDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Actualizar campos
	existingClinicalNote.IDConsultationVisit = clinicalNoteDTO.IDConsultationVisit
	existingClinicalNote.ProgressNotes = clinicalNoteDTO.ProgressNotes
	existingClinicalNote.ObservationsRecommendations = clinicalNoteDTO.ObservationsRecommendations

	if err := db.DB.Save(&existingClinicalNote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update clinical note"})
		return
	}

	// Convertir a DTO de respuesta
	responseDTO := response.ClinicalResponseNoteDTO{
		IDConsultationVisit:         existingClinicalNote.IDConsultationVisit,
		ProgressNotes:               existingClinicalNote.ProgressNotes,
		ObservationsRecommendations: existingClinicalNote.ObservationsRecommendations,
	}

	c.JSON(http.StatusOK, gin.H{"clinical_note": responseDTO})
}

// DeleteClinicalNote godoc
// @Summary Delete a clinical_note by consultation ID
// @Description Remove a clinical_note by its consultation ID
// @Tags clinical_notes
// @Param id path string true "ConsultationVisit ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /clinical_notes/{id} [delete]
func DeleteClinicalNote(c *gin.Context) {
	consultationID := c.Param("id")

	// Buscar por id_consulta en lugar de ID directo
	var clinicalNote models.ClinicalNote
	if err := db.DB.Where("id_consulta = ?", consultationID).First(&clinicalNote).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Clinical note not found for this consultation"})
		return
	}

	if err := db.DB.Delete(&clinicalNote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete clinical note"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
