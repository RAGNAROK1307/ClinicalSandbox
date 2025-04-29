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

// CreateDiagnosticImage godoc
// @Summary Create a new diagnostic_image
// @Description Adds a new diagnostic_image to the system
// @Tags diagnostic_images
// @Accept json
// @Produce json
// @Param diagnostic_image body request.CreateDiagnosticImageDTO true "DiagnosticImage data"
// @Success 201 {object} response.DiagnosticImageResponseDTO
// @Failure 400 {object} map[string]string
// @Router /diagnostic_images [post]
func CreateDiagnosticImage(c *gin.Context) {
	var diagnostic_imageDTO request.CreateDiagnosticImageDTO

	if err := c.ShouldBindJSON(&diagnostic_imageDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	imageDate, err := time.Parse("2006-01-02", diagnostic_imageDTO.ImageDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	diagnostic_image := models.DiagnosticImage{
		IDConsultationVisit: diagnostic_imageDTO.IDConsultationVisit,
		ImageDate:           imageDate,
		ImageType:           diagnostic_imageDTO.ImageType,
		Description:         diagnostic_imageDTO.Description,
		ImageInterpretation: diagnostic_imageDTO.ImageInterpretation,
		//ExternalFilePath:    diagnostic_imageDTO.ExternalFilePath,
	}

	if err := db.DB.Create(&diagnostic_image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create diagnostic_image"})
		return
	}

	// Convertir a DTO de respuesta
	responseDTO := response.DiagnosticImageResponseDTO{
		IDDiagnosticImage:   diagnostic_image.IDDiagnosticImage,
		IDConsultationVisit: diagnostic_image.IDConsultationVisit,
		ImageDate:           diagnostic_image.ImageDate.Format("2006-01-02"),
		ImageType:           diagnostic_image.ImageType,
		Description:         diagnostic_image.Description,
		ImageInterpretation: diagnostic_image.ImageInterpretation,
		ExternalFilePath:    diagnostic_image.ExternalFilePath,
	}

	c.JSON(http.StatusCreated, gin.H{"diagnostic_image": responseDTO})
}

// GetDiagnosticImages godoc
// @Summary List all diagnostic_images
// @Description Retrieve a list of all diagnostic_images in the system
// @Tags diagnostic_images
// @Produce json
// @Success 200 {array} response.DiagnosticImageResponseDTO
// @Router /diagnostic_images [get]
func GetDiagnosticImages(c *gin.Context) {
	var diagnosticImages []models.DiagnosticImage
	db.DB.Preload("ConsultationVisit").Find(&diagnosticImages)

	// Convertir a DTOs de respuesta
	var responseDTOs []response.DiagnosticImageResponseDTO
	for _, diagnosticImage := range diagnosticImages {
		responseDTOs = append(responseDTOs, response.DiagnosticImageResponseDTO{
			IDDiagnosticImage:   diagnosticImage.IDDiagnosticImage,
			IDConsultationVisit: diagnosticImage.IDConsultationVisit,
			ImageDate:           diagnosticImage.ImageDate.Format("2006-01-02"),
			ImageType:           diagnosticImage.ImageType,
			Description:         diagnosticImage.Description,
			ImageInterpretation: diagnosticImage.ImageInterpretation,
			//ExternalFilePath:    diagnosticImage.ExternalFilePath,
		})
	}

	c.JSON(http.StatusOK, gin.H{"diagnostic_images": responseDTOs})
}

// GetDiagnosticImage godoc
// @Summary Get a diagnostic_image by consultation ID
// @Description Retrieve a single diagnostic_image by its consultation ID
// @Tags diagnostic_images
// @Param id path string true "ConsultationVisit ID"
// @Produce json
// @Success 200 {object} response.DiagnosticImageResponseDTO
// @Failure 404 {object} map[string]string
// @Router /diagnostic_images/{id} [get]
func GetDiagnosticImage(c *gin.Context) {
	consultationID := c.Param("id")
	var diagnosticImage models.DiagnosticImage

	if err := db.DB.Preload("ConsultationVisit").
		Where("id_consulta = ?", consultationID).
		First(&diagnosticImage).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "DiagnosticImage not found for this consultation"})
		return
	}

	// Convertir a DTO de respuesta
	responseDTO := response.DiagnosticImageResponseDTO{
		IDDiagnosticImage:   diagnosticImage.IDDiagnosticImage,
		IDConsultationVisit: diagnosticImage.IDConsultationVisit,
		ImageDate:           diagnosticImage.ImageDate.Format("2006-01-02"),
		ImageType:           diagnosticImage.ImageType,
		Description:         diagnosticImage.Description,
		ImageInterpretation: diagnosticImage.ImageInterpretation,
		ExternalFilePath:    diagnosticImage.ExternalFilePath,
	}

	c.JSON(http.StatusOK, gin.H{"diagnostic_image": responseDTO})
}

// UpdateDiagnosticImage godoc
// @Summary Update a diagnostic_image by consultation ID
// @Description Update the information of an existing diagnostic_image by consultation ID
// @Tags diagnostic_images
// @Accept json
// @Produce json
// @Param id path string true "ConsultationVisit ID"
// @Param diagnostic_image body request.CreateDiagnosticImageDTO true "Updated diagnostic_image data"
// @Success 200 {object} response.DiagnosticImageResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /diagnostic_images/{id} [put]
func UpdateDiagnosticImage(c *gin.Context) {
	consultationID := c.Param("id")
	var existingDiagnosticImage models.DiagnosticImage

	if err := db.DB.Where("id_consulta = ?", consultationID).
		First(&existingDiagnosticImage).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "DiagnosticImage not found for this consultation"})
		return
	}

	var diagnostic_imageDTO request.CreateDiagnosticImageDTO
	if err := c.ShouldBindJSON(&diagnostic_imageDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	imageDate, err := time.Parse("2006-01-02", diagnostic_imageDTO.ImageDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Actualizar campos
	existingDiagnosticImage.IDConsultationVisit = diagnostic_imageDTO.IDConsultationVisit
	existingDiagnosticImage.ImageDate = imageDate
	existingDiagnosticImage.ImageType = diagnostic_imageDTO.ImageType
	existingDiagnosticImage.Description = diagnostic_imageDTO.Description
	existingDiagnosticImage.ImageInterpretation = diagnostic_imageDTO.ImageInterpretation
	//existingDiagnosticImage.ExternalFilePath = diagnostic_imageDTO.ExternalFilePath

	if err := db.DB.Save(&existingDiagnosticImage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update diagnostic_image"})
		return
	}

	// Convertir a DTO de respuesta
	responseDTO := response.DiagnosticImageResponseDTO{
		IDDiagnosticImage:   existingDiagnosticImage.IDDiagnosticImage,
		IDConsultationVisit: existingDiagnosticImage.IDConsultationVisit,
		ImageDate:           existingDiagnosticImage.ImageDate.Format("2006-01-02"),
		ImageType:           existingDiagnosticImage.ImageType,
		Description:         existingDiagnosticImage.Description,
		ImageInterpretation: existingDiagnosticImage.ImageInterpretation,
		ExternalFilePath:    existingDiagnosticImage.ExternalFilePath,
	}

	c.JSON(http.StatusOK, gin.H{"diagnostic_image": responseDTO})
}

// DeleteDiagnosticImage godoc
// @Summary Delete a diagnostic_image by consultation ID
// @Description Remove a diagnostic_image by its consultation ID
// @Tags diagnostic_images
// @Param id path string true "ConsultationVisit ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /diagnostic_images/{id} [delete]
func DeleteDiagnosticImage(c *gin.Context) {
	consultationID := c.Param("id")

	// Primero buscamos el registro para verificar que existe
	var diagnosticImage models.DiagnosticImage
	if err := db.DB.Where("id_consulta = ?", consultationID).
		First(&diagnosticImage).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diagnostic image not found for this consultation"})
		return
	}

	// Eliminar el registro usando el ID encontrado
	if err := db.DB.Delete(&diagnosticImage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete diagnostic image"})
		return
	}

	c.Status(http.StatusNoContent)
}
