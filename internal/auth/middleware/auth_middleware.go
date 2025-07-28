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

/*
Este módulo del paquete `middleware` gestiona la seguridad de acceso en la API mediante JWT,
control de sesiones activas, lista negra de tokens (blacklist), y validaciones de autorización.

Principales funcionalidades:

- GenerateToken:
  Genera un token JWT con datos personalizados como ID de usuario, rol, nombre de usuario, y
  posibles IDs relacionados (personal hospitalario o paciente). El token tiene una duración
  de sesión limitada y se registra como sesión activa.

- AuthMiddleware:
  Middleware que valida tokens JWT en las peticiones protegidas. Si el token está próximo a expirar,
  se renueva automáticamente y se incluye en el header `X-Renewed-Token`.

- RoleMiddleware:
  Permite restringir el acceso a rutas según roles específicos (por su ID numérico).

- ValidateUserAccess:
  Valida que el usuario autenticado acceda únicamente a sus propios datos, según su rol.
  (Administrador puede acceder libremente; otros roles se validan por ID).

- ValidateUpdatePatient:
  Permite actualizar datos del paciente únicamente al propio paciente, a administradores
  o directivos.

- ValidateMedicalRecordAccess:
  Controla el acceso a registros médicos:
    - Pacientes solo acceden a su historial.
    - Médicos pueden acceder a todos (o se puede extender para filtrar por asignación).

- ValidatePasswordAccess:
  Solo permite cambiar la contraseña al propio usuario o a un administrador.

- Logout:
  Finaliza la sesión actual del usuario, agregando el token a la lista negra y
  eliminándolo de los registros de sesión activa.

- checkExpiredSessions (goroutine):
  Ejecuta limpieza periódica de sesiones inactivas y tokens expirados para mantener
  la seguridad del sistema.

- cleanExpiredBlacklistedTokens (goroutine):
  Elimina tokens expirados de la lista negra cada 24 horas.

Este middleware proporciona un modelo robusto de autenticación y control de acceso,
fundamental para garantizar la integridad y privacidad de los datos en la aplicación.
*/

var (
	jwtKey           = []byte("Vulcilab")
	activeSessions   = make(map[uint]time.Time)
	tokenBlacklist   = make(map[string]time.Time)
	blacklistMutex   = &sync.Mutex{}
	sessionMutex     = &sync.Mutex{}
	sessionDuration  = 5 * time.Minute
	blacklistCleanup = 24 * time.Hour        // Limpiar tokens vencidos cada 24h
	ActiveTokens     = make(map[uint]string) // userID -> token
	PendingTokens    = make(map[uint]string)
	tokenMutex       = &sync.Mutex{}
)

type Claims struct {
	UserID           uint   `json:"user_id"`
	RoleID           uint   `json:"role_id"`
	UserName         string `json:"user_name"`
	RoleName         string `json:"role_name"`
	HospitalEmployee uint   `json:"hospital_employee"` // ID del personal hospitalario
	Patient          uint   `json:"patient"`           // ID del paciente
	jwt.StandardClaims
}

func GenerateToken(userID uint, roleID uint, userName string, roleName string, hospitalEmployeeID uint, patientID uint) (string, error) {
	claims := &Claims{
		UserID:           userID,
		RoleID:           roleID,
		UserName:         userName,
		RoleName:         roleName,
		HospitalEmployee: hospitalEmployeeID,
		Patient:          patientID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(sessionDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Registrar la sesión activa
	sessionMutex.Lock()
	activeSessions[userID] = time.Now()
	sessionMutex.Unlock()

	return token.SignedString(jwtKey)
}

func init() {
	// Goroutine para limpiar tokens vencidos periódicamente
	go checkExpiredSessions()
	go func() {
		for {
			time.Sleep(blacklistCleanup)
			cleanExpiredBlacklistedTokens()
		}
	}()
}

func cleanExpiredBlacklistedTokens() {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	now := time.Now()
	for token, exp := range tokenBlacklist {
		if now.After(exp) {
			delete(tokenBlacklist, token)
		}
	}
}

func isTokenBlacklisted(tokenString string) bool {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	exp, exists := tokenBlacklist[tokenString]
	if !exists {
		return false
	}

	// Si el token ya expiró, no está realmente blacklisted
	return time.Now().Before(exp)
}

func addToBlacklist(tokenString string, exp time.Time) {
	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	tokenBlacklist[tokenString] = exp
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

		// Verificar si el token está a punto de expirar (por ejemplo, en menos de 1 minuto)
		timeLeft := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
		shouldRenew := timeLeft < time.Minute
		// Si necesita renovación, generar nuevo token
		if shouldRenew {
			newToken, err := GenerateToken(
				claims.UserID,
				claims.RoleID,
				claims.UserName,
				claims.RoleName,
				claims.HospitalEmployee,
				claims.Patient,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to renew token"})
				c.Abort()
				return
			}

			// Agregar el nuevo token a la respuesta
			c.Header("X-Renewed-Token", newToken)
		}

		// Exponer el header X-Renewed-Token al cliente
		c.Header("Access-Control-Expose-Headers", "X-Renewed-Token")

		// Actualizar tiempo de última actividad
		sessionMutex.Lock()
		activeSessions[claims.UserID] = time.Now()
		sessionMutex.Unlock()

		// Resto del middleware
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

func ValidateUpdatePatient() gin.HandlerFunc {
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
		if claims.RoleName == "Administrador" || claims.RoleName == "Directivo" {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "No hay sesión activa"})
		return
	}

	claims := claimsInterface.(*Claims)
	authHeader := c.GetHeader("Authorization")

	if authHeader != "" {
		tokenString := strings.Split(authHeader, " ")[1]
		exp := time.Unix(claims.ExpiresAt, 0)
		addToBlacklist(tokenString, exp)
	}

	tokenMutex.Lock()
	delete(ActiveTokens, claims.UserID)
	delete(activeSessions, claims.UserID)
	tokenMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada correctamente"})
}

// Función para verificar tokens expirados periódicamente
func checkExpiredSessions() {
	for {
		time.Sleep(1 * time.Minute)
		now := time.Now()

		tokenMutex.Lock()
		for userID, lastActive := range activeSessions {
			if now.Sub(lastActive) > sessionDuration {
				if token, exists := ActiveTokens[userID]; exists {
					addToBlacklist(token, now)
					delete(ActiveTokens, userID)
					delete(activeSessions, userID)
				}
			}
		}
		tokenMutex.Unlock()
	}
}
