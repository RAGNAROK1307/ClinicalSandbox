package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/dto/response"
	"ClinicalSandBox/internal/API/models"
	"ClinicalSandBox/internal/auth/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	_ "time"
)

func Login(c *gin.Context) {
	var loginRequest request.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User
	if err := db.DB.Where("nombre_usuario = ?", loginRequest.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Contraseña incorrecta"})
		return
	}

	// Cargar la información del rol relacionado
	var role models.Role
	if err := db.DB.First(&role, user.IDRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cargar el rol del usuario"})
		return
	}

	// Obtener el ID del personal hospitalario o del paciente según el rol
	var hospitalEmployeeID uint
	var patientID uint

	switch role.RoleName {
	case "Médico", "Directivo":
		// Buscar el personal hospitalario asociado al usuario
		var hospitalEmployee models.HospitalEmployee
		if err := db.DB.Where("id_usuarios = ?", user.IDUser).First(&hospitalEmployee).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cargar la información del personal hospitalario"})
			return
		}
		hospitalEmployeeID = hospitalEmployee.IDHospitalEmployee
	case "Paciente":
		// Buscar el paciente asociado al usuario
		var patient models.Patient
		if err := db.DB.Where("id_usuarios = ?", user.IDUser).First(&patient).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cargar la información del paciente"})
			return
		}
		patientID = patient.IDPatient
	default:
		// Si el rol no es "Médico", "Directivo" o "Paciente", no se necesita un ID adicional
	}

	// Generar el token con los IDs correspondientes
	token, err := middleware.GenerateToken(user.IDUser, user.IDRole, user.UserName, role.RoleName, hospitalEmployeeID, patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token"})
		return
	}

	c.JSON(http.StatusOK, response.LoginResponse{
		Token: token,
	})
}

func Logout(c *gin.Context) {
	middleware.Logout(c)
}

/*func Login(c *gin.Context) {
	var loginRequest request.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("nombre_usuario = ?", loginRequest.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Contraseña incorrecta"})
		return
	}

	// Cargar la información del rol relacionado
	var role models.Role
	if err := db.DB.First(&role, user.IDRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al cargar el rol del usuario"})
		return
	}

	// Generar el token con user.IDUser, user.IDRole, user.UserName y role.RoleName
	token, err := middleware.GenerateToken(user.IDUser, user.IDRole, user.UserName, role.RoleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token"})
		return
	}

	c.JSON(http.StatusOK, response.LoginResponse{
		Token: token,
	})
}*/

/*func Login(c *gin.Context) {
	var loginRequest request.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("nombre_usuario = ?", loginRequest.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Comparar contraseñas directamente (sin bcrypt)
	if user.Password != loginRequest.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Contraseña incorrecta"})
		return
	}

	// Generar el token
	token, err := middleware.GenerateToken(user.IDUser, user.IDRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token"})
		return
	}

	c.JSON(http.StatusOK, response.LoginResponse{
		Token: token,
	})
}*/
