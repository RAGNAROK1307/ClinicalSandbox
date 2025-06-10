package middleware

import (
	"log"
	"os"
	"time"
)

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
