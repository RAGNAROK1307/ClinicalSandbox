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

// CreateConsentAuthorization godoc
// @Summary Create a new consent_authorization
// @Description Adds a new consent_authorization to the system
// @Tags consent_authorizations
// @Accept json
// @Produce json
// @Param consent_authorization body request.CreateConsentAuthorizationDTO true "ConsentAuthorization data"
// @Success 201 {object} response.ConsentAuthorizationResponseDTO
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

	// Mapea a modelo ConsentAuthorization
	consent_authorization := models.ConsentAuthorization{
		IDPatient:   consent_authorizationDTO.IDPatient,
		ConsentType: consent_authorizationDTO.ConsentType,
		ConsentDate: consentDate,
		Details:     consent_authorizationDTO.Details,
	}

	// Guarda el nuevo consent_authorization en la base de datos
	if err := db.DB.Create(&consent_authorization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create consent_authorization"})
		return
	}

	// Mapear a DTO de respuesta
	responseDTO := response.ConsentAuthorizationResponseDTO{
		IDConsentAuthorization: consent_authorization.IDConsentAuthorization,
		ConsentType:            consent_authorization.ConsentType,
		ConsentDate:            consent_authorization.ConsentDate.Format("2006-01-02"),
		Details:                consent_authorization.Details,
	}

	c.JSON(http.StatusCreated, gin.H{"consent_authorization": responseDTO})
}

// GetConsentAuthorizations godoc
// @Summary List all consent_authorizations
// @Description Retrieve a list of all consent_authorizations in the system
// @Tags consent_authorizations
// @Produce json
// @Success 200 {array} response.ConsentAuthorizationResponseDTO
// @Router /consent_authorizations [get]
func GetConsentAuthorizations(c *gin.Context) {
	var consent_authorizations []models.ConsentAuthorization
	db.DB.Preload("Patient").Find(&consent_authorizations)

	// Mapear a DTOs de respuesta
	var responseDTOs []response.ConsentAuthorizationResponseDTO
	for _, ca := range consent_authorizations {
		responseDTOs = append(responseDTOs, response.ConsentAuthorizationResponseDTO{
			IDConsentAuthorization: ca.IDConsentAuthorization,
			IDPatient:              ca.IDPatient,
			ConsentType:            ca.ConsentType,
			ConsentDate:            ca.ConsentDate.Format("2006-01-02"),
			Details:                ca.Details,
			FullName:               ca.Patient.FullName,
		})
	}

	c.JSON(http.StatusOK, gin.H{"consent_authorizations": responseDTOs})
}

// GetConsentAuthorization godoc
// @Summary Get consent_authorizations by patient ID
// @Description Retrieve consent_authorizations for a specific patient
// @Tags consent_authorizations
// @Param id path string true "Patient ID"
// @Produce json
// @Success 200 {array} response.ConsentAuthorizationResponseDTO
// @Failure 404 {object} map[string]string
// @Router /consent_authorizations/{id} [get]
func GetConsentAuthorization(c *gin.Context) {
	idPatient := c.Param("id")
	if idPatient == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID is required"})
		return
	}

	var consent_authorizations []models.ConsentAuthorization

	// Usar el nombre del campo exacto como está en el modelo (probablemente IDPatient)
	if err := db.DB.Preload("Patient").Where("id_paciente = ?", idPatient).Find(&consent_authorizations).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No consent authorizations found for this patient"})
		return
	}

	// Mapear a DTOs de respuesta
	var responseDTOs []response.ConsentAuthorizationResponseDTO
	for _, ca := range consent_authorizations {
		responseDTOs = append(responseDTOs, response.ConsentAuthorizationResponseDTO{
			IDConsentAuthorization: ca.IDConsentAuthorization,
			ConsentType:            ca.ConsentType,
			ConsentDate:            ca.ConsentDate.Format("2006-01-02"),
			Details:                ca.Details,
		})
	}

	c.JSON(http.StatusOK, gin.H{"consent_authorizations": responseDTOs})
}

// UpdateConsentAuthorization godoc
// @Summary Update a consent_authorization
// @Description Update the information of an existing consent_authorization
// @Tags consent_authorizations
// @Accept json
// @Produce json
// @Param id path string true "ConsentAuthorization ID"
// @Param consent_authorization body request.UpdateConsentAuthorizationDTO true "Updated consent_authorization data"
// @Success 200 {object} response.ConsentAuthorizationResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /consent_authorizations/{id} [put]
func UpdateConsentAuthorization(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ConsentAuthorization ID is required"})
		return
	}

	var existingConsentAuthorization models.ConsentAuthorization
	if err := db.DB.First(&existingConsentAuthorization, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ConsentAuthorization not found"})
		return
	}

	var consent_authorizationDTO request.UpdateConsentAuthorizationDTO
	if err := c.ShouldBindJSON(&consent_authorizationDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	consentDate, err := time.Parse("2006-01-02", consent_authorizationDTO.ConsentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Actualizar solo los campos permitidos
	existingConsentAuthorization.ConsentType = consent_authorizationDTO.ConsentType
	existingConsentAuthorization.ConsentDate = consentDate
	existingConsentAuthorization.Details = consent_authorizationDTO.Details

	if err := db.DB.Save(&existingConsentAuthorization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update consent_authorization"})
		return
	}

	responseDTO := response.ConsentAuthorizationResponseDTO{
		IDConsentAuthorization: existingConsentAuthorization.IDConsentAuthorization,
		ConsentType:            existingConsentAuthorization.ConsentType,
		ConsentDate:            existingConsentAuthorization.ConsentDate.Format("2006-01-02"),
		Details:                existingConsentAuthorization.Details,
	}

	c.JSON(http.StatusOK, gin.H{"consent_authorization": responseDTO})
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
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ConsentAuthorization ID is required"})
		return
	}

	if err := db.DB.Delete(&models.ConsentAuthorization{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ConsentAuthorization not found"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
