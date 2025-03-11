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

// CreateExecutiveAndUser godoc
// @Summary Create a new executive and user
// @Description Adds a new executive and user to the system in a single transaction
// @Tags executives
// @Accept json
// @Produce json
// @Param executiveAndUser body request.CreateHospitalEmployeeAndUserDTO true "Executive and User data"
// @Success 201 {object} response.HospitalEmployeeAndUserResponseDTO
// @Failure 400 {object} map[string]string
// @Router /executive-and-user [post]
func CreateExecutiveAndUser(c *gin.Context) {
	var executiveAndUserDTO request.CreateHospitalEmployeeAndUserDTO

	// Bind JSON to executiveAndUserDTO
	if err := c.ShouldBindJSON(&executiveAndUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha en formato string a time.Time
	birthDate, err := time.Parse("2006-01-02", executiveAndUserDTO.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Obtener el ID del rol de Directivo desde `SeedRoles`
	executiveRoleID := db.DirectivoID
	if executiveRoleID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Executive role ID not found"})
		return
	}

	// Iniciar una transacción
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Asignar siempre el rol de Directivo (ID = 2)
	//const executiveRoleID = 2

	// Crear el usuario
	user := models.User{
		IDRole:   executiveRoleID,
		UserName: executiveAndUserDTO.UserName,
		Password: executiveAndUserDTO.Password,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Crear el identification
	identification := models.Identification{
		DocumentType:   executiveAndUserDTO.DocumentType,
		DocumentNumber: executiveAndUserDTO.DocumentNumber,
	}

	if err := tx.Create(&identification).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create identification"})
		return
	}

	// Crear el executive (hospital employee)
	executive := models.HospitalEmployee{
		IDRole:           executiveRoleID,
		IDIdentification: identification.IDIdentification, // Asignar el ID del identification creado
		IDUser:           user.IDUser,                     // Asignar el ID del usuario creado
		FullName:         executiveAndUserDTO.FullName,
		BirthDate:        birthDate,
		Gender:           executiveAndUserDTO.Gender,
		Phone:            executiveAndUserDTO.Phone,
		Address:          executiveAndUserDTO.Address,
	}

	if err := tx.Create(&executive).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create executive"})
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
		RoleName: "Directivo",
	}

	executiveResponse := response.HospitalEmployeeResponseDTO{
		IDHospitalEmployee: executive.IDHospitalEmployee,
		UserName:           user.UserName,
		RoleName:           "Directivo",
		FullName:           executive.FullName,
		BirthDate:          executive.BirthDate,
		Gender:             executive.Gender,
		Address:            executive.Address,
		Phone:              executive.Phone,
	}

	identificationResponse := response.IdentificationResponseDTO{
		IDIdentification: identification.IDIdentification,
		DocumentType:     identification.DocumentType,
		DocumentNumber:   identification.DocumentNumber,
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":           userResponse,
		"executive":      executiveResponse,
		"identification": identificationResponse,
	})
}

// GetExecutivesAndUsers godoc
// @Summary List all executives and their associated users
// @Description Retrieve a list of all executives and their associated users in the system
// @Tags executives
// @Produce json
// @Success 200 {array} response.HospitalEmployeeAndUserResponseDTO
// @Router /executives-and-users [get]
func GetExecutivesAndUsers(c *gin.Context) {
	var executives []models.HospitalEmployee

	// Obtener todos los executivees con sus relaciones
	if err := db.DB.Preload("User").Preload("Role").Preload("Identification").Where("id_rol = ?", db.DirectivoID).Find(&executives).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch executives"})
		return
	}

	var executivesResponse []response.HospitalEmployeeAndUserResponseDTO
	for _, executive := range executives {
		executivesResponse = append(executivesResponse, response.HospitalEmployeeAndUserResponseDTO{
			User: response.UserResponseDTO{
				IDUser:   executive.User.IDUser,
				UserName: executive.User.UserName,
				RoleName: executive.Role.RoleName,
			},
			HospitalEmployee: response.HospitalEmployeeResponseDTO{
				IDHospitalEmployee: executive.IDHospitalEmployee,
				UserName:           executive.User.UserName,
				RoleName:           executive.Role.RoleName,
				FullName:           executive.FullName,
				BirthDate:          executive.BirthDate,
				Gender:             executive.Gender,
				Address:            executive.Address,
				Phone:              executive.Phone,
			},
			Identification: response.IdentificationResponseDTO{
				IDIdentification: executive.Identification.IDIdentification,
				DocumentType:     executive.Identification.DocumentType,
				DocumentNumber:   executive.Identification.DocumentNumber,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{"executives": executivesResponse})
}

// GetExecutiveAndUserByID godoc
// @Summary Get an executive and user by ID
// @Description Retrieve a single executive and their associated user by ID
// @Tags executives
// @Param id path string true "Executive ID"
// @Produce json
// @Success 200 {object} response.HospitalEmployeeAndUserResponseDTO
// @Failure 404 {object} map[string]string
// @Router /executives-and-users/{id} [get]
func GetExecutiveAndUserByID(c *gin.Context) {
	id := c.Param("id")
	var executive models.HospitalEmployee

	// Obtener el executive con sus relaciones
	if err := db.DB.Preload("User").Preload("Role").Preload("Identification").First(&executive, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Executive not found"})
		return
	}

	// Verificar que el executive tiene el rol de directivo
	if executive.IDRole != db.DirectivoID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Executive not found"})
		return
	}

	executiveResponse := response.HospitalEmployeeAndUserResponseDTO{
		User: response.UserResponseDTO{
			IDUser:   executive.User.IDUser,
			UserName: executive.User.UserName,
			RoleName: executive.Role.RoleName,
		},
		HospitalEmployee: response.HospitalEmployeeResponseDTO{
			IDHospitalEmployee: executive.IDHospitalEmployee,
			UserName:           executive.User.UserName,
			RoleName:           executive.Role.RoleName,
			FullName:           executive.FullName,
			BirthDate:          executive.BirthDate,
			Gender:             executive.Gender,
			Address:            executive.Address,
			Phone:              executive.Phone,
		},
		Identification: response.IdentificationResponseDTO{
			IDIdentification: executive.Identification.IDIdentification,
			DocumentType:     executive.Identification.DocumentType,
			DocumentNumber:   executive.Identification.DocumentNumber,
		},
	}

	c.JSON(http.StatusOK, gin.H{"executive": executiveResponse})
}

// UpdateExecutiveAndUser godoc
// @Summary Update an executive and user
// @Description Update the information of an existing executive and their associated user
// @Tags executives
// @Accept json
// @Produce json
// @Param id path string true "Executive ID"
// @Param executiveAndUser body request.UpdateHospitalEmployeeAndUserDTO true "Updated executive and user data"
// @Success 200 {object} response.HospitalEmployeeAndUserResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /executives-and-users/{id} [put]
func UpdateExecutiveAndUser(c *gin.Context) {
	id := c.Param("id")
	var existingExecutive models.HospitalEmployee

	// Verificar que el executive existe
	if err := db.DB.Preload("User").Preload("Role").Preload("Identification").First(&existingExecutive, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Executive not found"})
		return
	}

	// Verificar que el executive tiene el rol de médico (id_rol = 4)
	if existingExecutive.IDRole != db.DirectivoID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Executive not found"})
		return
	}

	// Bind JSON al DTO
	var executiveAndUserDTO request.UpdateHospitalEmployeeAndUserDTO
	if err := c.ShouldBindJSON(&executiveAndUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Mapa para actualizar solo los campos enviados
	updateFields := map[string]interface{}{}

	if executiveAndUserDTO.FullName != "" {
		updateFields["nombre_completo"] = executiveAndUserDTO.FullName
	}
	if executiveAndUserDTO.BirthDate != "" {
		birthDate, err := time.Parse("2006-01-02", executiveAndUserDTO.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		updateFields["fecha_nacimiento"] = birthDate
	}
	if executiveAndUserDTO.Gender != "" {
		updateFields["genero"] = executiveAndUserDTO.Gender
	}
	if executiveAndUserDTO.Phone != "" {
		updateFields["telefono"] = executiveAndUserDTO.Phone
	}
	if executiveAndUserDTO.Address != "" {
		updateFields["direccion"] = executiveAndUserDTO.Address
	}

	// Iniciar transacción
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Actualizar usuario si se envió algún dato
	if executiveAndUserDTO.UserName != "" {
		if err := tx.Model(&existingExecutive.User).Update("nombre_usuario", executiveAndUserDTO.UserName).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}

	// Actualizar identificación si se envió algún dato
	identificationFields := map[string]interface{}{}
	if executiveAndUserDTO.DocumentType != "" {
		identificationFields["tipo_documento"] = executiveAndUserDTO.DocumentType
	}
	if executiveAndUserDTO.DocumentNumber != "" {
		identificationFields["numero_documento"] = executiveAndUserDTO.DocumentNumber
	}

	if len(identificationFields) > 0 {
		if err := tx.Model(&models.Identification{}).
			Where("id_identificacion = ?", existingExecutive.IDIdentification).
			Updates(identificationFields).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update identification"})
			return
		}
	}

	// Actualizar el executive solo con los valores enviados
	if len(updateFields) > 0 {
		if err := tx.Model(&models.HospitalEmployee{}).
			Where("id_personal_hospital = ?", existingExecutive.IDHospitalEmployee).
			Updates(updateFields).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update executive"})
			return
		}
	}

	// Commit de la transacción
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Generar la respuesta actualizada
	executiveResponse := response.HospitalEmployeeAndUserResponseDTO{
		User: response.UserResponseDTO{
			IDUser:   existingExecutive.User.IDUser,
			UserName: existingExecutive.User.UserName,
			RoleName: "Directivo",
		},
		HospitalEmployee: response.HospitalEmployeeResponseDTO{
			IDHospitalEmployee: existingExecutive.IDHospitalEmployee,
			UserName:           existingExecutive.User.UserName,
			RoleName:           "Directivo",
			FullName:           existingExecutive.FullName,
			BirthDate:          existingExecutive.BirthDate,
			Gender:             existingExecutive.Gender,
			Address:            existingExecutive.Address,
			Phone:              existingExecutive.Phone,
		},
		Identification: response.IdentificationResponseDTO{
			IDIdentification: existingExecutive.Identification.IDIdentification,
			DocumentType:     existingExecutive.Identification.DocumentType,
			DocumentNumber:   existingExecutive.Identification.DocumentNumber,
		},
	}

	c.JSON(http.StatusOK, gin.H{"executive": executiveResponse})
}

// DeleteExecutiveAndUser godoc
// @Summary Delete an executive and user
// @Description Remove an executive and their associated user by ID
// @Tags executives
// @Param id path string true "Executive ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /executives-and-users/{id} [delete]
func DeleteExecutiveAndUser(c *gin.Context) {
	id := c.Param("id")
	var executive models.HospitalEmployee

	// Verificar que el executive existe
	if err := db.DB.Preload("User").Preload("Identification").First(&executive, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Executive not found"})
		return
	}

	// Verificar que el executive tiene el rol de directivo
	if executive.IDRole != db.DirectivoID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Executive not found"})
		return
	}

	// Iniciar una transacción
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Eliminar el executive
	if err := tx.Delete(&executive).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete executive"})
		return
	}

	// Eliminar el usuario
	if err := tx.Delete(&executive.User).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Eliminar el identification
	if err := tx.Delete(&executive.Identification).Error; err != nil {
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
