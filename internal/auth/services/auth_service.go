package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/dto/response"
	"ClinicalSandBox/internal/API/models"
	"ClinicalSandBox/internal/auth/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"sync"
	"time"
)

/*
Este módulo del paquete `services` implementa las funciones de autenticación del sistema.

Función Login:
- Valida las credenciales del usuario enviadas en el cuerpo de la solicitud.
- Aplica un sistema de bloqueo temporal tras múltiples intentos fallidos de autenticación, con tiempos crecientes.
- Detecta patrones comunes de inyección SQL en el campo de usuario y omite la validación de contraseña en esos casos como mecanismo de mitigación básica.
- Si las credenciales son válidas, genera un token JWT que incluye información adicional según el rol del usuario:
  - Para médicos o directivos: se adjunta el ID del personal hospitalario.
  - Para pacientes: se adjunta el ID del paciente.
- Responde con un objeto JSON que indica el estado del inicio de sesión y el token de autenticación.

Función Logout:
- Llama al middleware para invalidar el token y cerrar la sesión del usuario.

Notas:
- Existe una vulnerabilidad de inyección SQL en el uso directo de `fmt.Sprintf()` al construir la consulta de búsqueda del usuario. Se recomienda reemplazar esto por una consulta parametrizada con GORM para mayor seguridad.
*/

type LoginAttempt struct {
	FailedAttempts int
	LockUntil      time.Time
}

var (
	loginAttempts = make(map[string]*LoginAttempt)
	mutex         sync.Mutex
)

func registerFailedAttempt(username string) {
	mutex.Lock()
	defer mutex.Unlock()

	attempt, exists := loginAttempts[username]
	if !exists {
		loginAttempts[username] = &LoginAttempt{FailedAttempts: 1}
		return
	}

	attempt.FailedAttempts++
	if attempt.FailedAttempts >= 3 {
		// Tiempo de bloqueo = 1 min, luego se duplica (2, 4, 8, etc.)
		lockDuration := time.Minute * time.Duration(1<<uint(attempt.FailedAttempts-3))
		attempt.LockUntil = time.Now().Add(lockDuration)
	}
}

func resetFailedAttempts(username string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(loginAttempts, username)
}

func Login(c *gin.Context) {
	var loginRequest request.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":      "Cuerpo de solicitud inválido",
			"blocked":      false,
			"remaining_ms": 0,
		})
		return
	}

	username := loginRequest.Username

	// Buscar el usuario (vulnerabilidad aquí)
	var user models.User
	query := fmt.Sprintf("SELECT * FROM usuarios WHERE nombre_usuario = '%s'", username)
	fmt.Println(">>> Ejecutando SQL:", query)

	if err := db.DB.Raw(query).Scan(&user).Error; err != nil || user.IDUser == 0 {
		registerFailedAttempt(username)

		mutex.Lock()
		defer mutex.Unlock()
		attempt := loginAttempts[username]
		var remaining int64
		var blocked bool
		if attempt != nil && time.Now().Before(attempt.LockUntil) {
			blocked = true
			remaining = time.Until(attempt.LockUntil).Milliseconds()
		}

		c.JSON(http.StatusUnauthorized, response.LoginResponse{
			Blocked:     blocked,
			RemainingMS: remaining,
			Message:     "Usuario no encontrado",
		})
		return
	}

	// Verificar la contraseña (solo si no es un intento de inyección)
	if !(strings.Contains(username, "'") || strings.Contains(username, "--") || strings.Contains(username, "1=1")) {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
			registerFailedAttempt(username)

			mutex.Lock()
			defer mutex.Unlock()
			attempt := loginAttempts[username]
			var remaining int64
			var blocked bool
			if attempt != nil && time.Now().Before(attempt.LockUntil) {
				blocked = true
				remaining = time.Until(attempt.LockUntil).Milliseconds()
			}

			c.JSON(http.StatusUnauthorized, response.LoginResponse{
				Blocked:     blocked,
				RemainingMS: remaining,
				Message:     "Contraseña incorrecta",
			})
			return
		}
	} else {
		fmt.Println(">>> Posible inyección SQL detectada: se omite la validación de contraseña")
	}

	// Restablecer intentos fallidos
	resetFailedAttempts(username)

	// Cargar el rol
	var role models.Role
	if err := db.DB.First(&role, user.IDRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":      "Error al cargar el rol del usuario",
			"blocked":      false,
			"remaining_ms": 0,
		})
		return
	}

	// Obtener IDs adicionales
	var hospitalEmployeeID uint
	var patientID uint
	switch role.RoleName {
	case "Médico", "Directivo":
		var hospitalEmployee models.HospitalEmployee
		if err := db.DB.Where("id_usuarios = ?", user.IDUser).First(&hospitalEmployee).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":      "Error al cargar el personal hospitalario",
				"blocked":      false,
				"remaining_ms": 0,
			})
			return
		}
		hospitalEmployeeID = hospitalEmployee.IDHospitalEmployee
	case "Paciente":
		var patient models.Patient
		if err := db.DB.Where("id_usuarios = ?", user.IDUser).First(&patient).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":      "Error al cargar el paciente",
				"blocked":      false,
				"remaining_ms": 0,
			})
			return
		}
		patientID = patient.IDPatient
	}

	// Generar token
	token, err := middleware.GenerateToken(user.IDUser, user.IDRole, user.UserName, role.RoleName, hospitalEmployeeID, patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":      "Error al generar el token",
			"blocked":      false,
			"remaining_ms": 0,
		})
		return
	}

	// Respuesta exitosa
	c.JSON(http.StatusOK, response.LoginResponse{
		Token:         token,
		Blocked:       false,
		RemainingMS:   0,
		ActiveSession: false,
		Message:       "Inicio de sesión exitoso",
	})
}

func Logout(c *gin.Context) {
	middleware.Logout(c)
}
