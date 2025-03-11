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

// CreateUserAndPatient godoc
// @Summary Create a new user, patient, demographic data, and identification
// @Description Adds a new user, patient, demographic data, and identification to the system in a single transaction
// @Tags patients
// @Accept json
// @Produce json
// @Param userAndPatient body request.CreateUserAndPatientDTO true "User, Patient, Demographic Data, and Identification"
// @Success 201 {object} response.UserAndPatientResponseDTO
// @Failure 400 {object} map[string]string
// @Router /user-and-patients [post]
func CreateUserAndPatient(c *gin.Context) {
	var userAndPatientDTO request.CreateUserAndPatientDTO

	// Bind JSON to userAndPatientDTO
	if err := c.ShouldBindJSON(&userAndPatientDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha en formato string a time.Time
	birthDate, err := time.Parse("2006-01-02", userAndPatientDTO.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Obtener el ID del rol de Paciente desde `SeedRoles`
	patientRoleID := db.PacienteID
	if patientRoleID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Patient role ID not found"})
		return
	}

	// Iniciar una transacción
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Asignar manualmente el ID del rol de Paciente (3)
	//const patientRoleID = 3

	// Crear el usuario
	user := models.User{
		IDRole:   patientRoleID, // Siempre asignar el rol de Paciente (3)
		UserName: userAndPatientDTO.UserName,
		Password: userAndPatientDTO.Password,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Crear el identification
	identification := models.Identification{
		DocumentType:   userAndPatientDTO.DocumentType,
		DocumentNumber: userAndPatientDTO.DocumentNumber,
	}

	if err := tx.Create(&identification).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create identification"})
		return
	}

	// Crear el demographic_data
	demographicData := models.DemographicData{
		RacialEthnicInformation: userAndPatientDTO.RacialEthnicInformation,
		MaritalStatus:           userAndPatientDTO.MaritalStatus,
		Nationality:             userAndPatientDTO.Nationality,
		EmploymentInformation:   userAndPatientDTO.EmploymentInformation,
	}

	if err := tx.Create(&demographicData).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create demographic data"})
		return
	}

	// Crear el paciente
	patient := models.Patient{
		IDRole:               patientRoleID,                     // Siempre asignar el rol de Paciente (3)
		IDIdentification:     identification.IDIdentification,   // Asignar el ID del identification creado
		IDDemographicData:    demographicData.IDDemographicData, // Asignar el ID del demographic_data creado
		IDUser:               user.IDUser,                       // Asignar el ID del usuario creado
		FullName:             userAndPatientDTO.FullName,
		BirthDate:            birthDate,
		Gender:               userAndPatientDTO.Gender,
		Phone:                userAndPatientDTO.Phone,
		Address:              userAndPatientDTO.Address,
		SocialSecurityNumber: userAndPatientDTO.SocialSecurityNumber,
	}

	if err := tx.Create(&patient).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient"})
		return
	}

	// Commit de la transacción
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Generar la respuesta
	userResponse := response.UserResponseDTO{
		IDUser:   user.IDUser,
		UserName: user.UserName,
		RoleName: "Paciente", // Asignar manualmente el nombre del rol
	}

	patientResponse := response.PatientResponseDTO{
		IDPatient:            patient.IDPatient,
		FullName:             patient.FullName,
		BirthDate:            patient.BirthDate,
		Gender:               patient.Gender,
		Address:              patient.Address,
		Phone:                patient.Phone,
		SocialSecurityNumber: patient.SocialSecurityNumber,
	}

	demographicDataResponse := response.DemographicDataResponseDTO{
		IDDemographicData:       demographicData.IDDemographicData,
		RacialEthnicInformation: demographicData.RacialEthnicInformation,
		MaritalStatus:           demographicData.MaritalStatus,
		Nationality:             demographicData.Nationality,
		EmploymentInformation:   demographicData.EmploymentInformation,
	}

	identificationResponse := response.IdentificationResponseDTO{
		IDIdentification: identification.IDIdentification,
		DocumentType:     identification.DocumentType,
		DocumentNumber:   identification.DocumentNumber,
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":             userResponse,
		"patient":          patientResponse,
		"demographic_data": demographicDataResponse,
		"identification":   identificationResponse,
	})
}

// GetUserAndPatients godoc
// @Summary List all users and patients
// @Description Retrieve a list of all users and patients with their demographic data and identification
// @Tags patients
// @Produce json
// @Success 200 {array} response.UserAndPatientResponseDTO
// @Router /user-and-patients [get]
func GetUserAndPatients(c *gin.Context) {
	var patients []models.Patient

	// Obtener todos los pacientes con sus relaciones
	if err := db.DB.Preload("User").Preload("Role").Preload("Identification").Preload("DemographicData").Find(&patients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch patients"})
		return
	}

	var patientsResponse []response.UserAndPatientResponseDTO
	for _, patient := range patients {
		patientsResponse = append(patientsResponse, response.UserAndPatientResponseDTO{
			User: response.UserResponseDTO{
				IDUser:   patient.User.IDUser,
				UserName: patient.User.UserName,
				RoleName: patient.Role.RoleName,
			},
			Patient: response.PatientResponseDTO{
				IDPatient:            patient.IDPatient,
				FullName:             patient.FullName,
				BirthDate:            patient.BirthDate,
				Gender:               patient.Gender,
				Address:              patient.Address,
				Phone:                patient.Phone,
				SocialSecurityNumber: patient.SocialSecurityNumber,
				RoleName:             patient.Role.RoleName,
			},
			DemographicData: response.DemographicDataResponseDTO{
				IDDemographicData:       patient.DemographicData.IDDemographicData,
				RacialEthnicInformation: patient.DemographicData.RacialEthnicInformation,
				MaritalStatus:           patient.DemographicData.MaritalStatus,
				Nationality:             patient.DemographicData.Nationality,
				EmploymentInformation:   patient.DemographicData.EmploymentInformation,
			},
			Identification: response.IdentificationResponseDTO{
				IDIdentification: patient.Identification.IDIdentification,
				DocumentType:     patient.Identification.DocumentType,
				DocumentNumber:   patient.Identification.DocumentNumber,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{"patients": patientsResponse})
}

// GetUserAndPatient godoc
// @Summary Get a user and patient by ID
// @Description Retrieve a single user and patient with their demographic data and identification by ID
// @Tags patients
// @Param id path string true "Patient ID"
// @Produce json
// @Success 200 {object} response.UserAndPatientResponseDTO
// @Failure 404 {object} map[string]string
// @Router /user-and-patients/{id} [get]
func GetUserAndPatient(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient

	// Obtener el paciente con sus relaciones
	if err := db.DB.Preload("User").Preload("Role").Preload("Identification").Preload("DemographicData").First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// Generar la respuesta
	patientResponse := response.UserAndPatientResponseDTO{
		User: response.UserResponseDTO{
			IDUser:   patient.User.IDUser,
			UserName: patient.User.UserName,
			RoleName: patient.Role.RoleName,
		},
		Patient: response.PatientResponseDTO{
			IDPatient:            patient.IDPatient,
			FullName:             patient.FullName,
			BirthDate:            patient.BirthDate,
			Gender:               patient.Gender,
			Address:              patient.Address,
			Phone:                patient.Phone,
			SocialSecurityNumber: patient.SocialSecurityNumber,
		},
		DemographicData: response.DemographicDataResponseDTO{
			IDDemographicData:       patient.DemographicData.IDDemographicData,
			RacialEthnicInformation: patient.DemographicData.RacialEthnicInformation,
			MaritalStatus:           patient.DemographicData.MaritalStatus,
			Nationality:             patient.DemographicData.Nationality,
			EmploymentInformation:   patient.DemographicData.EmploymentInformation,
		},
		Identification: response.IdentificationResponseDTO{
			IDIdentification: patient.Identification.IDIdentification,
			DocumentType:     patient.Identification.DocumentType,
			DocumentNumber:   patient.Identification.DocumentNumber,
		},
	}

	c.JSON(http.StatusOK, gin.H{"patient": patientResponse})
}

// UpdateUserAndPatient godoc
// @Summary Update a user and patient
// @Description Update the information of an existing user and patient with their demographic data and identification
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Param userAndPatient body request.UpdateUserAndPatientDTO true "Updated user and patient data"
// @Success 200 {object} response.UserAndPatientResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /user-and-patients/{id} [put]
func UpdateUserAndPatient(c *gin.Context) {
	id := c.Param("id")
	var existingPatient models.Patient

	// Verificar que el paciente existe
	if err := db.DB.Preload("User").Preload("Role").Preload("Identification").Preload("DemographicData").First(&existingPatient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// Verificar que el paciente tiene el rol de Paciente (id_rol = 3)
	if existingPatient.IDRole != db.PacienteID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// Bind JSON al DTO
	var userAndPatientDTO request.UpdateUserAndPatientDTO
	if err := c.ShouldBindJSON(&userAndPatientDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Mapa para actualizar solo los campos enviados
	updateFields := map[string]interface{}{}

	if userAndPatientDTO.FullName != "" {
		updateFields["nombre_completo"] = userAndPatientDTO.FullName
	}
	if userAndPatientDTO.BirthDate != "" {
		birthDate, err := time.Parse("2006-01-02", userAndPatientDTO.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		updateFields["fecha_nacimiento"] = birthDate
	}
	if userAndPatientDTO.Gender != "" {
		updateFields["genero"] = userAndPatientDTO.Gender
	}
	if userAndPatientDTO.Phone != "" {
		updateFields["telefono"] = userAndPatientDTO.Phone
	}
	if userAndPatientDTO.Address != "" {
		updateFields["direccion"] = userAndPatientDTO.Address
	}
	if userAndPatientDTO.SocialSecurityNumber != "" {
		updateFields["numero_seguro_social"] = userAndPatientDTO.SocialSecurityNumber
	}

	// Iniciar transacción
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Actualizar usuario si se envió algún dato
	if userAndPatientDTO.UserName != "" {
		if err := tx.Model(&existingPatient.User).Update("nombre_usuario", userAndPatientDTO.UserName).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}

	// Actualizar identificación si se envió algún dato
	identificationFields := map[string]interface{}{}
	if userAndPatientDTO.DocumentType != "" {
		identificationFields["tipo_documento"] = userAndPatientDTO.DocumentType
	}
	if userAndPatientDTO.DocumentNumber != "" {
		identificationFields["numero_documento"] = userAndPatientDTO.DocumentNumber
	}

	if len(identificationFields) > 0 {
		if err := tx.Model(&models.Identification{}).
			Where("id_identificacion = ?", existingPatient.IDIdentification).
			Updates(identificationFields).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update identification"})
			return
		}
	}

	// Actualizar datos demográficos si se envió algún dato
	demographicFields := map[string]interface{}{}
	if userAndPatientDTO.RacialEthnicInformation != "" {
		demographicFields["informacion_etnica_racial"] = userAndPatientDTO.RacialEthnicInformation
	}
	if userAndPatientDTO.MaritalStatus != "" {
		demographicFields["estado_civil"] = userAndPatientDTO.MaritalStatus
	}
	if userAndPatientDTO.Nationality != "" {
		demographicFields["nacionalidad"] = userAndPatientDTO.Nationality
	}
	if userAndPatientDTO.EmploymentInformation != "" {
		demographicFields["informacion_laboral"] = userAndPatientDTO.EmploymentInformation
	}

	if len(demographicFields) > 0 {
		if err := tx.Model(&models.DemographicData{}).
			Where("id_datos_demograficos = ?", existingPatient.IDDemographicData).
			Updates(demographicFields).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update demographic data"})
			return
		}
	}

	// Actualizar el paciente solo con los valores enviados
	if len(updateFields) > 0 {
		if err := tx.Model(&models.Patient{}).
			Where("id_paciente = ?", existingPatient.IDPatient).
			Updates(updateFields).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient"})
			return
		}
	}

	// Commit de la transacción
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Generar la respuesta actualizada
	patientResponse := response.UserAndPatientResponseDTO{
		User: response.UserResponseDTO{
			IDUser:   existingPatient.User.IDUser,
			UserName: existingPatient.User.UserName,
			RoleName: "Paciente",
		},
		Patient: response.PatientResponseDTO{
			IDPatient:            existingPatient.IDPatient,
			FullName:             existingPatient.FullName,
			BirthDate:            existingPatient.BirthDate,
			Gender:               existingPatient.Gender,
			Address:              existingPatient.Address,
			Phone:                existingPatient.Phone,
			SocialSecurityNumber: existingPatient.SocialSecurityNumber,
		},
		DemographicData: response.DemographicDataResponseDTO{
			IDDemographicData:       existingPatient.DemographicData.IDDemographicData,
			RacialEthnicInformation: existingPatient.DemographicData.RacialEthnicInformation,
			MaritalStatus:           existingPatient.DemographicData.MaritalStatus,
			Nationality:             existingPatient.DemographicData.Nationality,
			EmploymentInformation:   existingPatient.DemographicData.EmploymentInformation,
		},
		Identification: response.IdentificationResponseDTO{
			IDIdentification: existingPatient.Identification.IDIdentification,
			DocumentType:     existingPatient.Identification.DocumentType,
			DocumentNumber:   existingPatient.Identification.DocumentNumber,
		},
	}

	c.JSON(http.StatusOK, gin.H{"patient": patientResponse})
}

// DeleteUserAndPatient godoc
// @Summary Delete a user and patient
// @Description Remove a user and patient with their demographic data and identification by ID
// @Tags patients
// @Param id path string true "Patient ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /user-and-patients/{id} [delete]
func DeleteUserAndPatient(c *gin.Context) {
	id := c.Param("id")
	var patient models.Patient

	// Verificar que el paciente existe
	if err := db.DB.Preload("User").Preload("Identification").Preload("DemographicData").First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// Iniciar una transacción
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Eliminar el paciente
	if err := tx.Delete(&patient).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete patient"})
		return
	}

	// Eliminar el usuario
	if err := tx.Delete(&patient.User).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Eliminar la identificación
	if err := tx.Delete(&patient.Identification).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete identification"})
		return
	}

	// Eliminar los datos demográficos
	if err := tx.Delete(&patient.DemographicData).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete demographic data"})
		return
	}

	// Commit de la transacción
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
