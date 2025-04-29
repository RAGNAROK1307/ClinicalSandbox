package middleware

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	jwtKey          = []byte("Simclec")
	activeSessions  = make(map[uint]time.Time)
	sessionMutex    = &sync.Mutex{}
	sessionDuration = 5 * time.Minute
)

/*type Claims struct {
	UserID   uint   `json:"user_id"`
	RoleID   uint   `json:"role_id"`
	UserName string `json:"user_name"`
	RoleName string `json:"role_name"`
	jwt.StandardClaims
}*/

type Claims struct {
	UserID           uint   `json:"user_id"`
	RoleID           uint   `json:"role_id"`
	UserName         string `json:"user_name"`
	RoleName         string `json:"role_name"`
	HospitalEmployee uint   `json:"hospital_employee"` // ID del personal hospitalario
	Patient          uint   `json:"patient"`           // ID del paciente
	jwt.StandardClaims
}

/*func GenerateToken(userID uint, roleID uint, userName string, roleName string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		RoleID:   roleID,
		UserName: userName,
		RoleName: roleName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Expira en 24 horas
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}*/

func GenerateToken(userID uint, roleID uint, userName string, roleName string, hospitalEmployeeID uint, patientID uint) (string, error) {
	claims := &Claims{
		UserID:           userID,
		RoleID:           roleID,
		UserName:         userName,
		RoleName:         roleName,
		HospitalEmployee: hospitalEmployeeID,
		Patient:          patientID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(sessionDuration).Unix(), // Expira en 5 minutos
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Registrar la sesión activa
	sessionMutex.Lock()
	activeSessions[userID] = time.Now()
	sessionMutex.Unlock()

	return token.SignedString(jwtKey)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Verificar actividad reciente
		sessionMutex.Lock()
		lastActivity, exists := activeSessions[claims.UserID]
		sessionMutex.Unlock()

		if !exists || time.Since(lastActivity) > sessionDuration {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired due to inactivity"})
			c.Abort()
			return
		}

		// Actualizar tiempo de última actividad
		sessionMutex.Lock()
		activeSessions[claims.UserID] = time.Now()
		sessionMutex.Unlock()

		// Resto del middleware (sin cambios)
		var user models.User
		if err := db.DB.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("user", user)
		c.Set("user_id", claims.UserID)
		c.Set("user_name", claims.UserName)
		c.Set("role_name", claims.RoleName)
		c.Set("hospital_employee_id", claims.HospitalEmployee)
		c.Set("patient_id", claims.Patient)

		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Verificar si el rol del usuario está en la lista de roles permitidos
		userRole := user.(models.User).IDRole
		allowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

/*
	func ValidateUserAccess() gin.HandlerFunc {
		return func(c *gin.Context) {
			// Obtener los claims del contexto
			claimsInterface, exists := c.Get("claims")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
				c.Abort()
				return
			}

			// Convertir los claims a la estructura Claims
			claims, ok := claimsInterface.(*Claims)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims format"})
				c.Abort()
				return
			}

			// Obtener el ID de la ruta y convertirlo a uint
			requestedIDStr := c.Param("id")
			requestedID, err := strconv.ParseUint(requestedIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
				c.Abort()
				return
			}

			// Verificar el acceso según el rol
			switch claims.RoleName {
			case "Médico", "Directivo":
				// Para médicos y directivos, verificar el ID del personal hospitalario
				if uint(requestedID) != claims.HospitalEmployee {
					c.JSON(http.StatusForbidden, gin.H{"error": "You can only access your own information"})
					c.Abort()
					return
				}
			case "Paciente":
				// Para pacientes, verificar el ID del paciente
				if uint(requestedID) != claims.Patient {
					c.JSON(http.StatusForbidden, gin.H{"error": "You can only access your own information"})
					c.Abort()
					return
				}
			default:
				c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
				c.Abort()
				return
			}

			c.Next()
		}
	}
*/
func ValidateUserAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener los claims del contexto
		claimsInterface, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
			c.Abort()
			return
		}

		// Convertir los claims a la estructura Claims
		claims, ok := claimsInterface.(*Claims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims format"})
			c.Abort()
			return
		}

		// Si el usuario es administrador, permitir el acceso sin verificación adicional
		if claims.RoleName == "Administrador" {
			c.Next()
			return
		}

		// Obtener el ID de la ruta y convertirlo a uint
		requestedIDStr := c.Param("id")
		requestedID, err := strconv.ParseUint(requestedIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			c.Abort()
			return
		}

		// Verificar el acceso según el rol
		switch claims.RoleName {
		case "Médico", "Directivo":
			// Para médicos y directivos, verificar el ID del personal hospitalario
			if uint(requestedID) != claims.HospitalEmployee {
				c.JSON(http.StatusForbidden, gin.H{"error": "You can only access your own information"})
				c.Abort()
				return
			}
		case "Paciente":
			// Para pacientes, verificar el ID del paciente
			if uint(requestedID) != claims.Patient {
				c.JSON(http.StatusForbidden, gin.H{"error": "You can only access your own information"})
				c.Abort()
				return
			}
		default:
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func ValidateMedicalRecordAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsInterface, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
			c.Abort()
			return
		}

		claims, ok := claimsInterface.(*Claims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims format"})
			c.Abort()
			return
		}

		// Obtener el ID del paciente de la URL
		patientIDStr := c.Param("id")
		patientID, err := strconv.ParseUint(patientIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID format"})
			c.Abort()
			return
		}

		switch claims.RoleName {
		case "Médico":
			// Médicos pueden acceder a cualquier registro (o podrías añadir lógica para verificar
			// si este paciente está asignado a este médico)
			c.Next()
			return

		case "Paciente":
			// Pacientes solo pueden ver sus propios registros
			if uint(patientID) != claims.Patient {
				c.JSON(http.StatusForbidden, gin.H{"error": "You can only access your own medical records"})
				c.Abort()
				return
			}
			c.Next()
			return

		default:
			// Ningún otro rol tiene acceso (incluyendo administradores)
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
	}
}

func ValidatePasswordAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener claims del usuario autenticado
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontró información de autenticación"})
			c.Abort()
			return
		}

		authClaims := claims.(*Claims)
		targetUserID := c.Param("id")

		// Permitir siempre a administradores
		if authClaims.RoleName == "Administrador" {
			c.Next()
			return
		}

		// Para otros roles, verificar que coincida el ID de usuario
		if strconv.FormatUint(uint64(authClaims.UserID), 10) != targetUserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Solo puedes cambiar tu propia contraseña"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Función para cerrar sesión
func Logout(c *gin.Context) {
	claimsInterface, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No active session"})
		return
	}

	claims := claimsInterface.(*Claims)

	sessionMutex.Lock()
	delete(activeSessions, claims.UserID)
	sessionMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
