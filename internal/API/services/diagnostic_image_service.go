package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
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
// @Success 201 {object} models.DiagnosticImage
// @Failure 400 {object} map[string]string
// @Router /diagnostic_images [post]
func CreateDiagnosticImage(c *gin.Context) {
	var diagnostic_imageDTO request.CreateDiagnosticImageDTO

	// Bind JSON to diagnostic_imageDTO
	if err := c.ShouldBindJSON(&diagnostic_imageDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha en formato string a time.Time
	imageDate, err := time.Parse("2006-01-02", diagnostic_imageDTO.ImageDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Mapea a modelo DiagnosticImage y usa la variable birthDate
	diagnostic_image := models.DiagnosticImage{
		IDConsultationVisit: diagnostic_imageDTO.IDConsultationVisit,
		ImageDate:           imageDate, // Aquí estamos usando la variable imageDate
		ImageType:           diagnostic_imageDTO.ImageType,
		Description:         diagnostic_imageDTO.Description,
		ImageInterpretation: diagnostic_imageDTO.ImageInterpretation,
		ExternalFilePath:    diagnostic_imageDTO.ExternalFilePath,
	}

	// Guarda el nuevo diagnostic_image en la base de datos
	if err := db.DB.Create(&diagnostic_image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create diagnostic_image"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"diagnostic_image": diagnostic_image})
}

// GetDiagnosticImages godoc
// @Summary List all diagnostic_images
// @Description Retrieve a list of all diagnostic_images in the system
// @Tags diagnostic_images
// @Produce json
// @Success 200 {array} models.DiagnosticImage
// @Router /diagnostic_images [get]
func GetDiagnosticImages(c *gin.Context) {
	var diagnostic_images []models.DiagnosticImage
	db.DB.Preload("ConsultationVisit").Find(&diagnostic_images)
	c.JSON(http.StatusOK, gin.H{"diagnostic_images": diagnostic_images})
}

// GetDiagnosticImage godoc
// @Summary Get a diagnostic_image by ID
// @Description Retrieve a single diagnostic_image by its ID
// @Tags diagnostic_images
// @Param id path string true "DiagnosticImage ID"
// @Produce json
// @Success 200 {object} models.DiagnosticImage
// @Failure 404 {object} map[string]string
// @Router /diagnostic_images/{id} [get]
func GetDiagnosticImage(c *gin.Context) {
	id := c.Param("id")
	var diagnostic_image models.DiagnosticImage
	if err := db.DB.Preload("ConsultationVisit").First(&diagnostic_image, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DiagnosticImage not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"diagnostic_image": diagnostic_image})
}

// UpdateDiagnosticImage godoc
// @Summary Update a diagnostic_image
// @Description Update the information of an existing diagnostic_image
// @Tags diagnostic_images
// @Accept json
// @Produce json
// @Param id path string true "DiagnosticImage ID"
// @Param diagnostic_image body request.CreateDiagnosticImageDTO true "Updated diagnostic_image data"
// @Success 200 {object} models.DiagnosticImage
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /diagnostic_images/{id} [put]
func UpdateDiagnosticImage(c *gin.Context) {
	id := c.Param("id")
	var existingDiagnosticImage models.DiagnosticImage

	// Verificar que el usuario existe en la base de datos
	if err := db.DB.First(&existingDiagnosticImage, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DiagnosticImage not found"})
		return
	}

	// Bind JSON al DTO para validar los datos
	var diagnostic_imageDTO request.CreateDiagnosticImageDTO
	if err := c.ShouldBindJSON(&diagnostic_imageDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha en formato string a time.Time
	imageDate, err := time.Parse("2006-01-02", diagnostic_imageDTO.ImageDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Actualizar solo los campos permitidos
	existingDiagnosticImage.IDConsultationVisit = diagnostic_imageDTO.IDConsultationVisit
	existingDiagnosticImage.ImageDate = imageDate // Usar la fecha convertida
	existingDiagnosticImage.ImageType = diagnostic_imageDTO.ImageType
	existingDiagnosticImage.Description = diagnostic_imageDTO.Description
	existingDiagnosticImage.ImageInterpretation = diagnostic_imageDTO.ImageInterpretation
	existingDiagnosticImage.ExternalFilePath = diagnostic_imageDTO.ExternalFilePath

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingDiagnosticImage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update diagnostic_image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"diagnostic_image": existingDiagnosticImage})
}

// DeleteDiagnosticImage godoc
// @Summary Delete a diagnostic_image
// @Description Remove a diagnostic_image by its ID
// @Tags diagnostic_images
// @Param id path string true "DiagnosticImage ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /diagnostic_images/{id} [delete]
func DeleteDiagnosticImage(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.DiagnosticImage{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DiagnosticImage not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
