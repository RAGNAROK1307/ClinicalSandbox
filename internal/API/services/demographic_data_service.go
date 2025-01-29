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

// CreateDemographicData godoc
// @Summary Create a new demographic_data
// @Description Adds a new demographic_data to the system
// @Tags demographics_data
// @Accept json
// @Produce json
// @Param demographic_data body request.CreateDemographicDataDTO true "DemographicData data"
// @Failure 400 {object} map[string]string
// @Success 201 {object} models.DemographicData
// @Router /demographics_data [post]
func CreateDemographicData(c *gin.Context) {
	var demographic_dataDTO request.CreateDemographicDataDTO

	// Bind JSON to demographic_dataDTO
	if err := c.ShouldBindJSON(&demographic_dataDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Mapea a modelo
	demographic_data := models.DemographicData{
		RacialEthnicInformation: demographic_dataDTO.RacialEthnicInformation,
		MaritalStatus:           demographic_dataDTO.MaritalStatus,
		Nationality:             demographic_dataDTO.Nationality,
		EmploymentInformation:   demographic_dataDTO.EmploymentInformation,
	}

	// Guarda el nuevo rol en la base de datos
	if err := db.DB.Create(&demographic_data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create demographic_data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"demographic_data": demographic_data})
}

// GetDemographicsData godoc
// @Summary List all demographics_data
// @Description Retrieve a list of all demographics_data in the system
// @Tags demographics_data
// @Produce json
// @Success 200 {array} models.DemographicData
// @Router /demographics_data [get]
func GetDemographicsData(c *gin.Context) {
	var demographics_data []models.DemographicData
	db.DB.Find(&demographics_data)
	c.JSON(200, gin.H{"demographics_data": demographics_data})
}

// GetDemographicData godoc
// @Summary Get a demographic_data by ID
// @Description Retrieve a single demographic_data by its ID
// @Tags demographics_data
// @Param id path string true "DemographicData ID"
// @Produce json
// @Success 200 {object} models.DemographicData
// @Failure 404 {object} map[string]string
// @Router /demographics_data/{id} [get]
func GetDemographicData(c *gin.Context) {
	id := c.Param("id")
	var demographic_data models.DemographicData
	if err := db.DB.First(&demographic_data, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "DemographicData not found"})
		return
	}
	c.JSON(200, gin.H{"demographic_data": demographic_data})
}

// UpdateDemographicData godoc
// @Summary Update a demographic_data
// @Description Update the information of an existing demographic_data
// @Tags demographics_data
// @Accept json
// @Produce json
// @Param id path string true "DemographicData ID"
// @Param demographic_data body request.CreateDemographicDataDTO true "Updated demographic_data data"
// @Success 200 {object} models.DemographicData
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /demographics_data/{id} [put]
func UpdateDemographicData(c *gin.Context) {
	id := c.Param("id")
	var existingDemographicData models.DemographicData

	// Verificar que el rol existe en la base de datos
	if err := db.DB.First(&existingDemographicData, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DemographicData not found"})
		return
	}

	// Bind JSON al DTO para validar los datos
	var demographic_dataDTO request.CreateDemographicDataDTO
	if err := c.ShouldBindJSON(&demographic_dataDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Actualizar solo los campos permitidos
	existingDemographicData.RacialEthnicInformation = demographic_dataDTO.RacialEthnicInformation
	existingDemographicData.MaritalStatus = demographic_dataDTO.MaritalStatus
	existingDemographicData.Nationality = demographic_dataDTO.Nationality
	existingDemographicData.EmploymentInformation = demographic_dataDTO.EmploymentInformation

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingDemographicData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update demographic_data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"demographic_data": existingDemographicData})
}

// DeleteDemographicData godoc
// @Summary Delete a demographic_data
// @Description Remove a demographic_data by its ID
// @Tags demographics_data
// @Param id path string true "DemographicData ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /demographics_data/{id} [delete]
func DeleteDemographicData(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.DemographicData{}, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "DemographicData not found"})
		return
	}
	c.JSON(204, nil)
}
