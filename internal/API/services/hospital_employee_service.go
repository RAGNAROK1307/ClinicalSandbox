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

// CreateHospitalEmployee godoc
// @Summary Create a new hospital_employee
// @Description Adds a new hospital_employee to the system
// @Tags hospital_employees
// @Accept json
// @Produce json
// @Param hospital_employee body request.CreateHospitalEmployeeDTO true "HospitalEmployee data"
// @Success 201 {object} models.HospitalEmployee
// @Failure 400 {object} map[string]string
// @Router /hospital_employees [post]
func CreateHospitalEmployee(c *gin.Context) {
	var hospital_employeeDTO request.CreateHospitalEmployeeDTO

	// Bind JSON to hospital_employeeDTO
	if err := c.ShouldBindJSON(&hospital_employeeDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha en formato string a time.Time
	birthDate, err := time.Parse("2006-01-02", hospital_employeeDTO.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Mapea a modelo HospitalEmployee y usa la variable birthDate
	hospital_employee := models.HospitalEmployee{
		IDRole:           hospital_employeeDTO.IDRole,
		IDIdentification: hospital_employeeDTO.IDIdentification,
		IDUser:           hospital_employeeDTO.IDUser,
		FullName:         hospital_employeeDTO.FullName,
		BirthDate:        birthDate, // Aquí estamos usando la variable birthDate
		Gender:           hospital_employeeDTO.Gender,
		Phone:            hospital_employeeDTO.Phone,
		Address:          hospital_employeeDTO.Address,
	}

	// Guarda el nuevo hospital_employee en la base de datos
	if err := db.DB.Create(&hospital_employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hospital_employee"})
		return
	}

	hospitalEmployeeResponse := response.HospitalEmployeeResponseDTO{
		IDHospitalEmployee: hospital_employee.IDHospitalEmployee,
		UserName:           hospital_employee.User.UserName,
		RoleName:           hospital_employee.Role.RoleName,
		FullName:           hospital_employee.FullName,
		BirthDate:          hospital_employee.BirthDate,
		Gender:             hospital_employee.Gender,
		Address:            hospital_employee.Address,
		Phone:              hospital_employee.Phone,
	}

	c.JSON(http.StatusCreated, gin.H{"user": hospitalEmployeeResponse})
	//c.JSON(http.StatusCreated, gin.H{"hospital_employee": hospital_employee})
}

// GetHospitalEmployees godoc
// @Summary List all hospital_employees
// @Description Retrieve a list of all hospital_employees in the system
// @Tags hospital_employees
// @Produce json
// @Success 200 {array} models.HospitalEmployee
// @Router /hospital_employees [get]
func GetHospitalEmployees(c *gin.Context) {
	var hospital_employees []models.HospitalEmployee

	if err := db.DB.Preload("Role").Preload("Identification").Preload("User").Find(&hospital_employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var hospitalEmployeesResponse []response.HospitalEmployeeResponseDTO
	for _, hospital_employee := range hospital_employees {
		hospitalEmployeesResponse = append(hospitalEmployeesResponse, response.HospitalEmployeeResponseDTO{
			IDHospitalEmployee: hospital_employee.IDHospitalEmployee,
			UserName:           hospital_employee.User.UserName,
			RoleName:           hospital_employee.Role.RoleName,
			FullName:           hospital_employee.FullName,
			BirthDate:          hospital_employee.BirthDate,
			Gender:             hospital_employee.Gender,
			Address:            hospital_employee.Address,
			Phone:              hospital_employee.Phone,
		})
	}

	c.JSON(http.StatusOK, gin.H{"hospital_employees": hospitalEmployeesResponse})
	//c.JSON(http.StatusOK, gin.H{"hospital_employees": hospital_employees})
}

// GetHospitalEmployee godoc
// @Summary Get a hospital_employee by ID
// @Description Retrieve a single hospital_employee by its ID
// @Tags hospital_employees
// @Param id path string true "HospitalEmployee ID"
// @Produce json
// @Success 200 {object} models.HospitalEmployee
// @Failure 404 {object} map[string]string
// @Router /hospital_employees/{id} [get]
func GetHospitalEmployee(c *gin.Context) {
	id := c.Param("id")
	var hospital_employee models.HospitalEmployee
	if err := db.DB.Preload("Role").Preload("Identification").Preload("User").First(&hospital_employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "HospitalEmployee not found"})
		return
	}

	hospitalEmployeeResponse := response.HospitalEmployeeResponseDTO{
		IDHospitalEmployee: hospital_employee.IDHospitalEmployee,
		UserName:           hospital_employee.User.UserName,
		RoleName:           hospital_employee.Role.RoleName,
		FullName:           hospital_employee.FullName,
		BirthDate:          hospital_employee.BirthDate,
		Gender:             hospital_employee.Gender,
		Address:            hospital_employee.Address,
		Phone:              hospital_employee.Phone,
	}

	c.JSON(http.StatusOK, gin.H{"hospital_employee": hospitalEmployeeResponse})
	//c.JSON(http.StatusOK, gin.H{"hospital_employee": hospital_employee})
}

// UpdateHospitalEmployee godoc
// @Summary Update a hospital_employee
// @Description Update the information of an existing hospital_employee
// @Tags hospital_employees
// @Accept json
// @Produce json
// @Param id path string true "HospitalEmployee ID"
// @Param hospital_employee body request.CreateHospitalEmployeeDTO true "Updated hospital_employee data"
// @Success 200 {object} models.HospitalEmployee
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /hospital_employees/{id} [put]
func UpdateHospitalEmployee(c *gin.Context) {
	id := c.Param("id")
	var existingHospitalEmployee models.HospitalEmployee

	// Verificar que el usuario existe en la base de datos
	if err := db.DB.First(&existingHospitalEmployee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "HospitalEmployee not found"})
		return
	}

	// Bind JSON al DTO para validar los datos
	var hospital_employeeDTO request.CreateHospitalEmployeeDTO
	if err := c.ShouldBindJSON(&hospital_employeeDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convertir la fecha en formato string a time.Time
	birthDate, err := time.Parse("2006-01-02", hospital_employeeDTO.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Actualizar solo los campos permitidos
	existingHospitalEmployee.IDRole = hospital_employeeDTO.IDRole
	existingHospitalEmployee.IDIdentification = hospital_employeeDTO.IDIdentification
	existingHospitalEmployee.IDUser = hospital_employeeDTO.IDUser
	existingHospitalEmployee.FullName = hospital_employeeDTO.FullName
	existingHospitalEmployee.BirthDate = birthDate // Usar la fecha convertida
	existingHospitalEmployee.Gender = hospital_employeeDTO.Gender
	existingHospitalEmployee.Phone = hospital_employeeDTO.Phone
	existingHospitalEmployee.Address = hospital_employeeDTO.Address

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingHospitalEmployee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hospital_employee"})
		return
	}

	hospitalEmployeeResponse := response.HospitalEmployeeResponseDTO{
		IDHospitalEmployee: existingHospitalEmployee.IDHospitalEmployee,
		UserName:           existingHospitalEmployee.User.UserName,
		RoleName:           existingHospitalEmployee.Role.RoleName,
		FullName:           existingHospitalEmployee.FullName,
		BirthDate:          existingHospitalEmployee.BirthDate,
		Gender:             existingHospitalEmployee.Gender,
		Address:            existingHospitalEmployee.Address,
		Phone:              existingHospitalEmployee.Phone,
	}

	c.JSON(http.StatusCreated, gin.H{"user": hospitalEmployeeResponse})

	//c.JSON(http.StatusOK, gin.H{"hospital_employee": existingHospitalEmployee})
}

// DeleteHospitalEmployee godoc
// @Summary Delete a hospital_employee
// @Description Remove a hospital_employee by its ID
// @Tags hospital_employees
// @Param id path string true "HospitalEmployee ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /hospital_employees/{id} [delete]
func DeleteHospitalEmployee(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.HospitalEmployee{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "HospitalEmployee not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
