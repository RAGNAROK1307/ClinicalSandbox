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

// CreateDoctorAndUser godoc
// @Summary Create a new doctor and user
// @Description Adds a new doctor and user to the system in a single transaction
// @Tags doctors
// @Accept json
// @Produce json
// @Param doctorAndUser body request.CreateHospitalEmployeeAndUserDTO true "Doctor and User data"
// @Success 201 {object} response.HospitalEmployeeAndUserResponseDTO
// @Failure 400 {object} map[string]string
// @Router /doctor-and-user [post]
func CreateDoctorAndUser(c *gin.Context) {
	var doctorAndUserDTO request.CreateHospitalEmployeeAndUserDTO

	// Bind JSON to doctorAndUserDTO
	if err := c.ShouldBindJSON(&doctorAndUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha en formato string a time.Time
	birthDate, err := time.Parse("2006-01-02", doctorAndUserDTO.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Obtener el ID del rol de Médico desde `SeedRoles`
	doctorRoleID := db.MedicoID
	if doctorRoleID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Doctor role ID not found"})
		return
	}

	// Iniciar una transacción
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Asignar siempre el rol de Médico (ID = 1)
	//const doctorRoleID = 1

	// Crear el usuario
	user := models.User{
		IDRole:   doctorRoleID,
		UserName: doctorAndUserDTO.UserName,
		Password: doctorAndUserDTO.Password,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Crear el identification
	identification := models.Identification{
		DocumentType:   doctorAndUserDTO.DocumentType,
		DocumentNumber: doctorAndUserDTO.DocumentNumber,
	}

	if err := tx.Create(&identification).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create identification"})
		return
	}

	// Crear el doctor (hospital employee)
	doctor := models.HospitalEmployee{
		IDRole:           doctorRoleID,
		IDIdentification: identification.IDIdentification, // Asignar el ID del identification creado
		IDUser:           user.IDUser,                     // Asignar el ID del usuario creado
		FullName:         doctorAndUserDTO.FullName,
		BirthDate:        birthDate,
		Gender:           doctorAndUserDTO.Gender,
		Phone:            doctorAndUserDTO.Phone,
		Address:          doctorAndUserDTO.Address,
	}

	if err := tx.Create(&doctor).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create doctor"})
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
		RoleName: "Médico",
	}

	doctorResponse := response.HospitalEmployeeResponseDTO{
		IDHospitalEmployee: doctor.IDHospitalEmployee,
		UserName:           user.UserName,
		RoleName:           "Médico",
		FullName:           doctor.FullName,
		BirthDate:          doctor.BirthDate,
		Gender:             doctor.Gender,
		Address:            doctor.Address,
		Phone:              doctor.Phone,
	}

	identificationResponse := response.IdentificationResponseDTO{
		IDIdentification: identification.IDIdentification,
		DocumentType:     identification.DocumentType,
		DocumentNumber:   identification.DocumentNumber,
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":           userResponse,
		"doctor":         doctorResponse,
		"identification": identificationResponse,
	})
}

// GetDoctorsAndUsers godoc
// @Summary List all doctors and their associated users
// @Description Retrieve a list of all doctors and their associated users in the system
// @Tags doctors
// @Produce json
// @Success 200 {array} response.HospitalEmployeeAndUserResponseDTO
// @Router /doctors-and-users [get]
func GetDoctorsAndUsers(c *gin.Context) {
	var doctors []models.HospitalEmployee

	// Obtener todos los doctores con sus relaciones
	if err := db.DB.Preload("User").Preload("Role").Preload("Identification").Where("id_rol = ?", db.MedicoID).Find(&doctors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctors"})
		return
	}

	var doctorsResponse []response.HospitalEmployeeAndUserResponseDTO
	for _, doctor := range doctors {
		doctorsResponse = append(doctorsResponse, response.HospitalEmployeeAndUserResponseDTO{
			User: response.UserResponseDTO{
				IDUser:   doctor.User.IDUser,
				UserName: doctor.User.UserName,
				RoleName: doctor.Role.RoleName,
			},
			HospitalEmployee: response.HospitalEmployeeResponseDTO{
				IDHospitalEmployee: doctor.IDHospitalEmployee,
				UserName:           doctor.User.UserName,
				RoleName:           doctor.Role.RoleName,
				FullName:           doctor.FullName,
				BirthDate:          doctor.BirthDate,
				Gender:             doctor.Gender,
				Address:            doctor.Address,
				Phone:              doctor.Phone,
			},
			Identification: response.IdentificationResponseDTO{
				IDIdentification: doctor.Identification.IDIdentification,
				DocumentType:     doctor.Identification.DocumentType,
				DocumentNumber:   doctor.Identification.DocumentNumber,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{"doctors": doctorsResponse})
}

// GetDoctorAndUserByID godoc
// @Summary Get a doctor and user by ID
// @Description Retrieve a single doctor and their associated user by ID
// @Tags doctors
// @Param id path string true "Doctor ID"
// @Produce json
// @Success 200 {object} response.HospitalEmployeeAndUserResponseDTO
// @Failure 404 {object} map[string]string
// @Router /doctors-and-users/{id} [get]
func GetDoctorAndUserByID(c *gin.Context) {
	id := c.Param("id")
	var doctor models.HospitalEmployee

	// Obtener el doctor con sus relaciones
	if err := db.DB.Preload("User").Preload("Role").Preload("Identification").First(&doctor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	// Verificar que el doctor tiene el rol de médico
	if doctor.IDRole != db.MedicoID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	doctorResponse := response.HospitalEmployeeAndUserResponseDTO{
		User: response.UserResponseDTO{
			IDUser:   doctor.User.IDUser,
			UserName: doctor.User.UserName,
			RoleName: doctor.Role.RoleName,
		},
		HospitalEmployee: response.HospitalEmployeeResponseDTO{
			IDHospitalEmployee: doctor.IDHospitalEmployee,
			UserName:           doctor.User.UserName,
			RoleName:           doctor.Role.RoleName,
			FullName:           doctor.FullName,
			BirthDate:          doctor.BirthDate,
			Gender:             doctor.Gender,
			Address:            doctor.Address,
			Phone:              doctor.Phone,
		},
		Identification: response.IdentificationResponseDTO{
			IDIdentification: doctor.Identification.IDIdentification,
			DocumentType:     doctor.Identification.DocumentType,
			DocumentNumber:   doctor.Identification.DocumentNumber,
		},
	}

	c.JSON(http.StatusOK, gin.H{"doctor": doctorResponse})
}

// UpdateDoctorAndUser godoc
// @Summary Update a doctor and user
// @Description Update the information of an existing doctor and their associated user
// @Tags doctors
// @Accept json
// @Produce json
// @Param id path string true "Doctor ID"
// @Param doctorAndUser body request.UpdateHospitalEmployeeAndUserDTO true "Updated doctor and user data"
// @Success 200 {object} response.HospitalEmployeeAndUserResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /doctors-and-users/{id} [put]
func UpdateDoctorAndUser(c *gin.Context) {
	id := c.Param("id")
	var existingDoctor models.HospitalEmployee

	// Verificar que el doctor existe
	if err := db.DB.Preload("User").Preload("Role").Preload("Identification").First(&existingDoctor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	// Verificar que el doctor tiene el rol de médico (id_rol = 1)
	if existingDoctor.IDRole != db.MedicoID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	// Bind JSON al DTO
	var doctorAndUserDTO request.UpdateHospitalEmployeeAndUserDTO
	if err := c.ShouldBindJSON(&doctorAndUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Mapa para actualizar solo los campos enviados
	updateFields := map[string]interface{}{}

	if doctorAndUserDTO.FullName != "" {
		updateFields["nombre_completo"] = doctorAndUserDTO.FullName
	}
	if doctorAndUserDTO.BirthDate != "" {
		birthDate, err := time.Parse("2006-01-02", doctorAndUserDTO.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		updateFields["fecha_nacimiento"] = birthDate
	}
	if doctorAndUserDTO.Gender != "" {
		updateFields["genero"] = doctorAndUserDTO.Gender
	}
	if doctorAndUserDTO.Phone != "" {
		updateFields["telefono"] = doctorAndUserDTO.Phone
	}
	if doctorAndUserDTO.Address != "" {
		updateFields["direccion"] = doctorAndUserDTO.Address
	}

	// Iniciar transacción
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Actualizar usuario si se envió algún dato
	if doctorAndUserDTO.UserName != "" {
		if err := tx.Model(&existingDoctor.User).Update("nombre_usuario", doctorAndUserDTO.UserName).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}

	// Actualizar identificación si se envió algún dato
	identificationFields := map[string]interface{}{}
	if doctorAndUserDTO.DocumentType != "" {
		identificationFields["tipo_documento"] = doctorAndUserDTO.DocumentType
	}
	if doctorAndUserDTO.DocumentNumber != "" {
		identificationFields["numero_documento"] = doctorAndUserDTO.DocumentNumber
	}

	if len(identificationFields) > 0 {
		if err := tx.Model(&models.Identification{}).
			Where("id_identificacion = ?", existingDoctor.IDIdentification).
			Updates(identificationFields).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update identification"})
			return
		}
	}

	// Actualizar el doctor solo con los valores enviados
	if len(updateFields) > 0 {
		if err := tx.Model(&models.HospitalEmployee{}).
			Where("id_personal_hospital = ?", existingDoctor.IDHospitalEmployee).
			Updates(updateFields).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update doctor"})
			return
		}
	}

	// Commit de la transacción
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Generar la respuesta actualizada
	doctorResponse := response.HospitalEmployeeAndUserResponseDTO{
		User: response.UserResponseDTO{
			IDUser:   existingDoctor.User.IDUser,
			UserName: existingDoctor.User.UserName,
			RoleName: "Médico",
		},
		HospitalEmployee: response.HospitalEmployeeResponseDTO{
			IDHospitalEmployee: existingDoctor.IDHospitalEmployee,
			UserName:           existingDoctor.User.UserName,
			RoleName:           "Médico",
			FullName:           existingDoctor.FullName,
			BirthDate:          existingDoctor.BirthDate,
			Gender:             existingDoctor.Gender,
			Address:            existingDoctor.Address,
			Phone:              existingDoctor.Phone,
		},
		Identification: response.IdentificationResponseDTO{
			IDIdentification: existingDoctor.Identification.IDIdentification,
			DocumentType:     existingDoctor.Identification.DocumentType,
			DocumentNumber:   existingDoctor.Identification.DocumentNumber,
		},
	}

	c.JSON(http.StatusOK, gin.H{"doctor": doctorResponse})
}

// DeleteDoctorAndUser godoc
// @Summary Delete a doctor and user
// @Description Remove a doctor and their associated user by ID
// @Tags doctors
// @Param id path string true "Doctor ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /doctors-and-users/{id} [delete]
func DeleteDoctorAndUser(c *gin.Context) {
	id := c.Param("id")
	var doctor models.HospitalEmployee

	// Verificar que el doctor existe
	if err := db.DB.Preload("User").Preload("Identification").First(&doctor, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	// Verificar que el doctor tiene el rol de médico
	if doctor.IDRole != db.MedicoID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	// Iniciar una transacción
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Eliminar el doctor
	if err := tx.Delete(&doctor).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete doctor"})
		return
	}

	// Eliminar el usuario
	if err := tx.Delete(&doctor.User).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Eliminar el identification
	if err := tx.Delete(&doctor.Identification).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete identification"})
		return
	}

	// Commit de la transacción
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
