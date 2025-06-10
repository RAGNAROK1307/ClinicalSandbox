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

	// Verificar si el usuario está bloqueado por intentos fallidos
	mutex.Lock()
	attempt, exists := loginAttempts[username]
	if exists && time.Now().Before(attempt.LockUntil) {
		remaining := time.Until(attempt.LockUntil)
		remainingMs := remaining.Milliseconds()
		mutex.Unlock()

		minutes := int(remaining.Minutes())
		seconds := int(remaining.Seconds()) % 60

		var timeParts []string
		if minutes > 0 {
			unit := "minuto"
			if minutes > 1 {
				unit += "s"
			}
			timeParts = append(timeParts, fmt.Sprintf("%d %s", minutes, unit))
		}
		if seconds > 0 {
			unit := "segundo"
			if seconds > 1 {
				unit += "s"
			}
			timeParts = append(timeParts, fmt.Sprintf("%d %s", seconds, unit))
		}

		timeMessage := strings.Join(timeParts, " y ")

		c.JSON(http.StatusTooManyRequests, response.LoginResponse{
			Blocked:     true,
			RemainingMS: remainingMs,
			Message:     fmt.Sprintf("Demasiados intentos fallidos. Intente de nuevo en %s.", timeMessage),
		})
		return
	}
	mutex.Unlock()

	// Buscar el usuario
	var user models.User
	if err := db.DB.Where("nombre_usuario = ?", username).First(&user).Error; err != nil {
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

	// Verificar la contraseña
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

	// Verificar si hay sesión activa (usando el middleware)
	middleware.SessionMutex.Lock()
	_, sessionActive := middleware.ActiveTokens[user.IDUser]
	middleware.SessionMutex.Unlock()

	if sessionActive {
		c.JSON(http.StatusConflict, response.LoginResponse{
			Blocked:       false,
			RemainingMS:   0,
			ActiveSession: true, // Indicar que hay sesión activa
			Message:       "Ya hay una sesión activa. Cierra la otra sesión primero.",
		})
		return
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
	token, err := middleware.GenerateToken(user.IDUser, user.IDRole, user.UserName, role.RoleName, hospitalEmployeeID, patientID, false)
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

func GetLockStatus(c *gin.Context) {
	username := c.Param("username")

	mutex.Lock()
	defer mutex.Unlock()

	attempt, exists := loginAttempts[username]
	if !exists || time.Now().After(attempt.LockUntil) {
		c.JSON(http.StatusOK, gin.H{
			"blocked":      false,
			"remaining_ms": 0,
			"message":      "El usuario no está bloqueado",
		})
		return
	}

	remaining := time.Until(attempt.LockUntil).Milliseconds()

	c.JSON(http.StatusOK, gin.H{
		"blocked":      true,
		"remaining_ms": remaining,
		"message":      "El usuario está temporalmente bloqueado",
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
}
*/
