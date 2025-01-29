package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/models"
	"github.com/gin-gonic/gin" // Para el framework Gin
	"net/http"
	//"gorm.io/gorm"                 // Para trabajar con GORM, si no lo has hecho ya
	_ "net/http" // Para constantes HTTP como http.StatusNotFound, etc.
)

// CreateIdentification godoc
// @Summary Create a new identification
// @Description Adds a new identification to the system
// @Tags identifications
// @Accept json
// @Produce json
// @Param identification body request.CreateIdentificationDTO true "Identification data"
// @Failure 400 {object} map[string]string
// @Success 201 {object} models.Identification
// @Router /identifications [post]
func CreateIdentification(c *gin.Context) {
	var identificationDTO request.CreateIdentificationDTO

	// Bind JSON to identificationDTO
	if err := c.ShouldBindJSON(&identificationDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Mapea a modelo
	identification := models.Identification{
		DocumentType:   identificationDTO.DocumentType,
		DocumentNumber: identificationDTO.DocumentNumber,
	}

	// Guarda el nuevo rol en la base de datos
	if err := db.DB.Create(&identification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create identification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"identification": identification})
}

// GetIdentifications godoc
// @Summary List all identifications
// @Description Retrieve a list of all identifications in the system
// @Tags identifications
// @Produce json
// @Success 200 {array} models.Identification
// @Router /identifications [get]
func GetIdentifications(c *gin.Context) {
	var identifications []models.Identification
	db.DB.Find(&identifications)
	c.JSON(200, gin.H{"identifications": identifications})
}

// GetIdentification godoc
// @Summary Get a identification by ID
// @Description Retrieve a single identification by its ID
// @Tags identifications
// @Param id path string true "Identification ID"
// @Produce json
// @Success 200 {object} models.Identification
// @Failure 404 {object} map[string]string
// @Router /identifications/{id} [get]
func GetIdentification(c *gin.Context) {
	id := c.Param("id")
	var identification models.Identification
	if err := db.DB.First(&identification, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Identification not found"})
		return
	}
	c.JSON(200, gin.H{"identification": identification})
}

// UpdateIdentification godoc
// @Summary Update a identification
// @Description Update the information of an existing identification
// @Tags identifications
// @Accept json
// @Produce json
// @Param id path string true "Identification ID"
// @Param identification body request.CreateIdentificationDTO true "Updated identification data"
// @Success 200 {object} models.Identification
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /identifications/{id} [put]
func UpdateIdentification(c *gin.Context) {
	id := c.Param("id")
	var existingIdentification models.Identification

	// Verificar que el rol existe en la base de datos
	if err := db.DB.First(&existingIdentification, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Identification not found"})
		return
	}

	// Bind JSON al DTO para validar los datos
	var identificationDTO request.CreateIdentificationDTO
	if err := c.ShouldBindJSON(&identificationDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Actualizar solo los campos permitidos
	existingIdentification.DocumentType = identificationDTO.DocumentType
	existingIdentification.DocumentNumber = identificationDTO.DocumentNumber

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingIdentification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update identification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"identification": existingIdentification})
}

// DeleteIdentification godoc
// @Summary Delete a identification
// @Description Remove a identification by its ID
// @Tags identifications
// @Param id path string true "Identification ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /identifications/{id} [delete]
func DeleteIdentification(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Identification{}, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Identification not found"})
		return
	}
	c.JSON(204, nil)
}
