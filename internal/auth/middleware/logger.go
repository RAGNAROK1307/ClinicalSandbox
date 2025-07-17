package middleware

import (
	"log"
	"os"
	"time"
)

/*
Este módulo implementa un sistema básico de registro de acciones de autenticación y autorización,
utilizado principalmente por el middleware de control de acceso del sistema clínico.

FUNCIONALIDADES PRINCIPALES:

- initLogger:
  Inicializa el sistema de logging:
    • Crea el directorio `logs/` si no existe.
    • Genera un archivo de log diario con nombre basado en la fecha actual (`auth_YYYY-MM-DD.log`).
    • Configura un `log.Logger` con prefijo `AUTH:` y formato que incluye timestamp y ubicación del log.

- logAuthAction:
  Registra acciones relevantes de autenticación o control de sesión.
  Guarda entradas con el siguiente formato:
      UserID: <id> | Username: <nombre> | Action: <descripción>
  Acciones típicas incluyen:
    • LOGIN
    • LOGOUT
    • SESSION_EXPIRED
    • TOKEN_RENEWED
    • LOGIN_BLOCKED (cuando ya hay una sesión activa)
    • BLACKLISTED_TOKEN (cuando un token ha sido invalidado)
    • INVALID_TOKEN (cuando el token es inválido o manipulado)
    • TOKEN_PROMOTED (cuando un token renovado pasa a ser el activo)

Este logger permite trazar las actividades clave relacionadas con la autenticación de usuarios,
sirviendo como mecanismo de auditoría y análisis ante posibles incidentes de seguridad.
*/

var (
	logFile *os.File
	logger  *log.Logger
)

func initLogger() {
	// Crear directorio si no existe
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	// Configurar archivo de log con fecha
	logFileName := "logs/auth_" + time.Now().Format("2006-01-02") + ".log"
	var err error
	logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	logger = log.New(logFile, "AUTH: ", log.LstdFlags|log.Lshortfile)
}

func logAuthAction(userID uint, username, action string) {
	if logger == nil {
		initLogger()
	}
	logger.Printf("UserID: %d | Username: %s | Action: %s\n", userID, username, action)
}
