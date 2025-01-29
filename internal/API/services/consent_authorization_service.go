package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// CreateConsentAuthorization godoc
// @Summary Create a new consent_authorization
// @Description Adds a new consent_authorization to the system
// @Tags consent_authorizations
// @Accept json
// @Produce json
// @Param consent_authorization body request.CreateConsentAuthorizationDTO true "ConsentAuthorization data"
// @Success 201 {object} models.ConsentAuthorization
// @Failure 400 {object} map[string]string
// @Router /consent_authorizations [post]
func CreateConsentAuthorization(c *gin.Context) {
	var consent_authorizationDTO request.CreateConsentAuthorizationDTO

	// Bind JSON to consent_authorizationDTO
	if err := c.ShouldBindJSON(&consent_authorizationDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha en formato string a time.Time
	consentDate, err := time.Parse("2006-01-02", consent_authorizationDTO.ConsentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Mapea a modelo ConsentAuthorization y usa la variable birthDate
	consent_authorization := models.ConsentAuthorization{
		IDPatient:        consent_authorizationDTO.IDPatient,
		ConsentType:      consent_authorizationDTO.ConsentType,
		ConsentDate:      consentDate, // Aquí estamos usando la variable consentDate
		Details:          consent_authorizationDTO.Details,
		ExternalFilePath: consent_authorizationDTO.ExternalFilePath,
	}

	// Guarda el nuevo consent_authorization en la base de datos
	if err := db.DB.Create(&consent_authorization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create consent_authorization"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"consent_authorization": consent_authorization})
}

// GetConsentAuthorizations godoc
// @Summary List all consent_authorizations
// @Description Retrieve a list of all consent_authorizations in the system
// @Tags consent_authorizations
// @Produce json
// @Success 200 {array} models.ConsentAuthorization
// @Router /consent_authorizations [get]
func GetConsentAuthorizations(c *gin.Context) {
	var consent_authorizations []models.ConsentAuthorization
	db.DB.Preload("Patient").Find(&consent_authorizations)
	c.JSON(http.StatusOK, gin.H{"consent_authorizations": consent_authorizations})
}

// GetConsentAuthorization godoc
// @Summary Get a consent_authorization by ID
// @Description Retrieve a single consent_authorization by its ID
// @Tags consent_authorizations
// @Param id path string true "ConsentAuthorization ID"
// @Produce json
// @Success 200 {object} models.ConsentAuthorization
// @Failure 404 {object} map[string]string
// @Router /consent_authorizations/{id} [get]
func GetConsentAuthorization(c *gin.Context) {
	id := c.Param("id")
	var consent_authorization models.ConsentAuthorization
	if err := db.DB.Preload("Patient").First(&consent_authorization, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ConsentAuthorization not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"consent_authorization": consent_authorization})
}

// UpdateConsentAuthorization godoc
// @Summary Update a consent_authorization
// @Description Update the information of an existing consent_authorization
// @Tags consent_authorizations
// @Accept json
// @Produce json
// @Param id path string true "ConsentAuthorization ID"
// @Param consent_authorization body request.CreateConsentAuthorizationDTO true "Updated consent_authorization data"
// @Success 200 {object} models.ConsentAuthorization
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /consent_authorizations/{id} [put]
func UpdateConsentAuthorization(c *gin.Context) {
	id := c.Param("id")
	var existingConsentAuthorization models.ConsentAuthorization

	// Verificar que el usuario existe en la base de datos
	if err := db.DB.First(&existingConsentAuthorization, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ConsentAuthorization not found"})
		return
	}

	// Bind JSON al DTO para validar los datos
	var consent_authorizationDTO request.CreateConsentAuthorizationDTO
	if err := c.ShouldBindJSON(&consent_authorizationDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha en formato string a time.Time
	consentDate, err := time.Parse("2006-01-02", consent_authorizationDTO.ConsentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Actualizar solo los campos permitidos
	existingConsentAuthorization.IDPatient = consent_authorizationDTO.IDPatient
	existingConsentAuthorization.ConsentType = consent_authorizationDTO.ConsentType
	existingConsentAuthorization.ConsentDate = consentDate // Usar la fecha convertida
	existingConsentAuthorization.Details = consent_authorizationDTO.Details
	existingConsentAuthorization.ExternalFilePath = consent_authorizationDTO.ExternalFilePath

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingConsentAuthorization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update consent_authorization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"consent_authorization": existingConsentAuthorization})
}

// DeleteConsentAuthorization godoc
// @Summary Delete a consent_authorization
// @Description Remove a consent_authorization by its ID
// @Tags consent_authorizations
// @Param id path string true "ConsentAuthorization ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /consent_authorizations/{id} [delete]
func DeleteConsentAuthorization(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.ConsentAuthorization{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ConsentAuthorization not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
