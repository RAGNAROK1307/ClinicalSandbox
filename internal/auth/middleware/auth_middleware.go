package middleware

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/models"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
Este módulo del paquete `middleware` implementa un sistema avanzado de autenticación y control de sesiones
basado en JWT para un sistema de gestión clínica. Incluye seguridad mejorada con control de tokens activos,
renovación automática, blacklist, y registro de auditoría.

FUNCIONALIDADES PRINCIPALES:

- GenerateToken:
  Genera tokens JWT con datos extendidos del usuario (ID, rol, nombres, paciente o personal hospitalario).
  Controla si ya hay una sesión activa e impide múltiples inicios simultáneos. Soporta renovación forzada.

- AuthMiddleware:
  Middleware de autenticación que:
    • Valida el token recibido.
    • Verifica si el token pertenece al usuario autenticado y si no ha sido invalidado.
    • Permite renovar automáticamente tokens próximos a expirar.
    • Promueve un token renovado a sesión activa si corresponde.

- RoleMiddleware:
  Restringe acceso a rutas según los roles permitidos (por su ID).

- ValidateUserAccess:
  Permite que los usuarios accedan solo a sus propios recursos. Admin accede sin restricción.
  Verifica que el ID en la URL coincida con el asociado al usuario autenticado.

- ValidateUpdatePatient:
  Controla quién puede modificar datos del paciente. Solo pacientes, directivos y administradores
  tienen permisos según su relación.

- ValidateMedicalRecordAccess:
  Control de acceso a historias clínicas. Pacientes pueden ver sus propios registros, médicos pueden
  ver todos (puede extenderse a lógica adicional). Otros roles no tienen acceso.

- ValidatePasswordAccess:
  Asegura que los usuarios solo puedan modificar su propia contraseña, excepto administradores.

- Logout:
  Finaliza la sesión del usuario. Mueve el token a la blacklist, elimina las referencias activas,
  y registra la acción de cierre de sesión.

- checkExpiredSessions (goroutine):
  Revisa periódicamente (cada minuto) las sesiones activas para invalidar tokens por inactividad.
  Los tokens expirados se añaden a la blacklist.

- cleanExpiredBlacklistedTokens (goroutine):
  Limpia cada 24 horas los tokens expirados de la lista negra para mantener el sistema limpio.

Este middleware también registra cada acción de autenticación mediante `logAuthAction`, útil para
auditoría y trazabilidad. El diseño favorece la seguridad, control de acceso detallado por rol,
y manejo inteligente de sesiones concurrentes.
*/

var (
	jwtKey           = []byte("Simclec")
	activeSessions   = make(map[uint]time.Time)
	tokenBlacklist   = make(map[string]time.Time)
	blacklistMutex   = &sync.Mutex{}
	SessionMutex     = &sync.Mutex{}
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

func GenerateToken(userID uint, roleID uint, userName string, roleName string, hospitalEmployeeID uint, patientID uint, forceRenewal bool) (string, error) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	// Solo bloquear login si NO es renovación
	if _, exists := ActiveTokens[userID]; exists && !forceRenewal {
		logAuthAction(userID, userName, "LOGIN_BLOCKED (sesión activa)")
		return "", errors.New("ya existe una sesión activa para este usuario")
	}

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
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	if !forceRenewal {
		// Solo guardar como activa si es login original
		ActiveTokens[userID] = tokenString
		logAuthAction(userID, userName, "LOGIN")
	} else {
		// Guardar como token pendiente de promoción
		PendingTokens[userID] = tokenString
		logAuthAction(userID, userName, "TOKEN_RENEWED (pendiente de promoción)")
	}

	return tokenString, nil
}

func init() {
	// Goroutine para limpiar tokens vencidos periódicamente
	initLogger()
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
			logAuthAction(0, "", "MISSING_TOKEN")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]

		if isTokenBlacklisted(tokenString) {
			logAuthAction(0, "", "BLACKLISTED_TOKEN")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesión terminada"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			logAuthAction(0, "", "INVALID_TOKEN")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Verificar si el token coincide con el activo registrado
		tokenMutex.Lock()
		activeToken, exists := ActiveTokens[claims.UserID]
		pendingToken, pendingExists := PendingTokens[claims.UserID]
		tokenMutex.Unlock()

		// Si el token no es el activo ni el pendiente, es inválido
		if !exists || (tokenString != activeToken && (!pendingExists || tokenString != pendingToken)) {
			logAuthAction(claims.UserID, claims.UserName, "SESSION_REJECTED (token no coincide)")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Ya hay una sesión activa. Cierra la otra sesión primero.",
			})
			c.Abort()
			return
		}

		// Si es el token pendiente, promoverlo a activo
		if tokenString == pendingToken {
			tokenMutex.Lock()
			ActiveTokens[claims.UserID] = pendingToken
			delete(PendingTokens, claims.UserID)
			tokenMutex.Unlock()
			logAuthAction(claims.UserID, claims.UserName, "TOKEN_PROMOTED")
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
				true, // forzar renovación
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al renovar token"})
				c.Abort()
				return
			}

			c.Header("X-Renewed-Token", newToken)
		}

		// Exponer el header X-Renewed-Token al cliente
		c.Header("Access-Control-Expose-Headers", "X-Renewed-Token")

		// Actualizar tiempo de última actividad
		SessionMutex.Lock()
		activeSessions[claims.UserID] = time.Now()
		SessionMutex.Unlock()

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

	logAuthAction(claims.UserID, claims.UserName, "LOGOUT")
	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada correctamente"})
}

// Función para verificar tokens expirados periódicamente
func checkExpiredSessions() {
	for {
		time.Sleep(1 * time.Minute) // Verificar cada minuto
		now := time.Now()

		tokenMutex.Lock()
		for userID, lastActive := range activeSessions {
			if now.Sub(lastActive) > sessionDuration {
				if token, exists := ActiveTokens[userID]; exists {
					addToBlacklist(token, now)
					delete(ActiveTokens, userID)
					delete(activeSessions, userID)

					// Obtener nombre de usuario para el log (requeriría una consulta a la DB)
					logAuthAction(userID, "SYSTEM", "SESSION_EXPIRED")
				}
			}
		}
		tokenMutex.Unlock()
	}
}
