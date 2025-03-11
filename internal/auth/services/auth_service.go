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
