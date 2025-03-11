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

// CreatePatient godoc
// @Summary Create a new patient
// @Description Adds a new patient to the system
// @Tags patients
// @Accept json
// @Produce json
// @Param patient body request.CreatePatientDTO true "Patient data"
// @Success 201 {object} models.Patient
// @Failure 400 {object} map[string]string
// @Router /patients [post]
func CreatePatient(c *gin.Context) {
	var patientDTO request.CreatePatientDTO

	// Bind JSON to patientDTO
	if err := c.ShouldBindJSON(&patientDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha en formato string a time.Time
	birthDate, err := time.Parse("2006-01-02", patientDTO.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Mapea a modelo Patient y usa la variable birthDate
	patient := models.Patient{
		IDRole:               patientDTO.IDRole,
		IDIdentification:     patientDTO.IDIdentification,
		IDDemographicData:    patientDTO.IDDemographicData,
		IDUser:               patientDTO.IDUser,
		FullName:             patientDTO.FullName,
		BirthDate:            birthDate, // Aquí estamos usando la variable birthDate
		Gender:               patientDTO.Gender,
		Phone:                patientDTO.Phone,
		Address:              patientDTO.Address,
		SocialSecurityNumber: patientDTO.SocialSecurityNumber,
	}

	// Guarda el nuevo patient en la base de datos
	if err := db.DB.Create(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient"})
		return
	}

	//var role models.Role
	//db.DB.First(&role, patient.IDRole)

	// Generar la respuesta con UserResponseDTO
	patientResponse := response.PatientResponseDTO{
		IDPatient:            patient.IDPatient,
		FullName:             patient.FullName,
		BirthDate:            patient.BirthDate,
		Gender:               patient.Gender,
		Address:              patient.Address,
		Phone:                patient.Phone,
		SocialSecurityNumber: patient.SocialSecurityNumber,
	}

	c.JSON(http.StatusCreated, gin.H{"user": patientResponse})

	//c.JSON(http.StatusCreated, gin.H{"patient": patient})
}

// GetPatients godoc
// @Summary List all patients
// @Description Retrieve a list of all patients in the system
// @Tags patients
// @Produce json
// @Success 200 {array} models.Patient
// @Router /patients [get]
func GetPatients(c *gin.Context) {
	var patients []models.Patient
	db.DB.Preload("Role").Preload("Identification").Preload("DemographicData").Preload("User").Find(&patients)

	var patientsResponse []response.PatientResponseDTO
	for _, patient := range patients {
		patientsResponse = append(patientsResponse, response.PatientResponseDTO{
			IDPatient:               patient.IDPatient,
			FullName:                patient.FullName,
			BirthDate:               patient.BirthDate,
			Gender:                  patient.Gender,
			Address:                 patient.Address,
			Phone:                   patient.Phone,
			SocialSecurityNumber:    patient.SocialSecurityNumber,
			RoleName:                patient.Role.RoleName,
			UserName:                patient.User.UserName,
			DocumentType:            patient.Identification.DocumentType,
			DocumentNumber:          patient.Identification.DocumentNumber,
			RacialEthnicInformation: patient.DemographicData.RacialEthnicInformation,
			MaritalStatus:           patient.DemographicData.MaritalStatus,
			Nationality:             patient.DemographicData.Nationality,
			EmploymentInformation:   patient.DemographicData.EmploymentInformation,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": patientsResponse})
	//c.JSON(http.StatusOK, gin.H{"patients": patients})
}

// GetPatient godoc
// @Summary Get a patient by ID
// @Description Retrieve a single patient by its ID
// @Tags patients
// @Param id path string true "Patient ID"
// @Produce json
// @Success 200 {object} models.Patient
// @Failure 404 {object} map[string]string
// @Router /patients/{id} [get]
func GetPatient(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient
	if err := db.DB.Preload("Role").Preload("Identification").Preload("DemographicData").Preload("User").First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	patientResponse := response.PatientResponseDTO{
		IDPatient:               patient.IDPatient,
		FullName:                patient.FullName,
		BirthDate:               patient.BirthDate,
		Gender:                  patient.Gender,
		Address:                 patient.Address,
		Phone:                   patient.Phone,
		SocialSecurityNumber:    patient.SocialSecurityNumber,
		RoleName:                patient.Role.RoleName,
		UserName:                patient.User.UserName,
		DocumentType:            patient.Identification.DocumentType,
		DocumentNumber:          patient.Identification.DocumentNumber,
		RacialEthnicInformation: patient.DemographicData.RacialEthnicInformation,
		MaritalStatus:           patient.DemographicData.MaritalStatus,
		Nationality:             patient.DemographicData.Nationality,
		EmploymentInformation:   patient.DemographicData.EmploymentInformation,
	}

	c.JSON(http.StatusOK, gin.H{"user": patientResponse})
	//c.JSON(http.StatusOK, gin.H{"patient": patient})
}

// UpdatePatient godoc
// @Summary Update a patient
// @Description Update the information of an existing patient
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Param patient body request.CreatePatientDTO true "Updated patient data"
// @Success 200 {object} models.Patient
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /patients/{id} [put]
func UpdatePatient(c *gin.Context) {
	id := c.Param("id")
	var existingPatient models.Patient

	// Verificar que el usuario existe en la base de datos
	if err := db.DB.First(&existingPatient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// Bind JSON al DTO para validar los datos
	var patientDTO request.CreatePatientDTO
	if err := c.ShouldBindJSON(&patientDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha en formato string a time.Time
	birthDate, err := time.Parse("2006-01-02", patientDTO.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Actualizar solo los campos permitidos
	existingPatient.IDRole = patientDTO.IDRole
	existingPatient.IDIdentification = patientDTO.IDIdentification
	existingPatient.IDDemographicData = patientDTO.IDDemographicData
	existingPatient.IDUser = patientDTO.IDUser
	existingPatient.FullName = patientDTO.FullName
	existingPatient.BirthDate = birthDate // Usar la fecha convertida
	existingPatient.Gender = patientDTO.Gender
	existingPatient.Address = patientDTO.Address
	existingPatient.Phone = patientDTO.Phone
	existingPatient.SocialSecurityNumber = patientDTO.SocialSecurityNumber

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingPatient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient"})
		return
	}

	patientResponse := response.PatientResponseDTO{
		IDPatient:            existingPatient.IDPatient,
		FullName:             existingPatient.FullName,
		BirthDate:            existingPatient.BirthDate,
		Gender:               existingPatient.Gender,
		Address:              existingPatient.Address,
		Phone:                existingPatient.Phone,
		SocialSecurityNumber: existingPatient.SocialSecurityNumber,
	}

	c.JSON(http.StatusOK, gin.H{"user": patientResponse})
	//c.JSON(http.StatusOK, gin.H{"patient": existingPatient})
}

// DeletePatient godoc
// @Summary Delete a patient
// @Description Remove a patient by its ID
// @Tags patients
// @Param id path string true "Patient ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /patients/{id} [delete]
func DeletePatient(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Patient{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
