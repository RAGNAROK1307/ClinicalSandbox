package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/response"
	"ClinicalSandBox/internal/API/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	uploadDir = "Image_of_Laboratories"
)

// ensureUploadDir crea el directorio de uploads si no existe
func ensureUploadDir() error {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		return os.Mkdir(uploadDir, 0755)
	}
	return nil
}

// UploadLaboratoryFile guarda el archivo en el sistema de archivos
func UploadLaboratoryFile(c *gin.Context) {
	laboratoryID := c.Param("id")

	// Verificar si existe el laboratorio
	var laboratory models.Laboratory
	if err := db.DB.First(&laboratory, laboratoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Laboratory record not found"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Crear directorio si no existe
	if err := ensureUploadDir(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generar nombre único para el archivo (usamos el nombre original)
	newFilename := laboratoryID + filepath.Ext(file.Filename)
	filePath := filepath.Join(uploadDir, newFilename)

	// Guardar el archivo
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Actualizar en la base de datos (solo la ruta)
	tx := db.DB.Begin()
	if err := tx.Model(&models.Laboratory{}).
		Where("id_laboratorio = ?", laboratoryID).
		Update("ruta_archivo_externo", filePath).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, response.LaboratoryFileResponseDTO{
		FilePath: filePath,
		Message:  "File uploaded successfully",
	})
}

// GetLaboratoryFile devuelve el archivo desde el sistema de archivos
func GetLaboratoryFile(c *gin.Context) {
	laboratoryID := c.Param("id")

	var laboratory models.Laboratory
	if err := db.DB.First(&laboratory, laboratoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Laboratory not found"})
		return
	}

	if laboratory.ExternalFilePath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No file associated with this laboratory"})
		return
	}

	// Verificar que el archivo existe
	if _, err := os.Stat(laboratory.ExternalFilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found on server"})
		return
	}

	// Servir el archivo directamente
	c.File(laboratory.ExternalFilePath)
}

// UpdateLaboratoryFile actualiza el archivo en el sistema de archivos
func UpdateLaboratoryFile(c *gin.Context) {
	laboratoryID := c.Param("id")

	// Verificar si existe el laboratorio
	var laboratory models.Laboratory
	if err := db.DB.First(&laboratory, laboratoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Laboratory not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Eliminar archivo anterior si existe
	if laboratory.ExternalFilePath != "" {
		if err := os.Remove(laboratory.ExternalFilePath); err != nil && !os.IsNotExist(err) {
			log.Printf("Warning: Failed to remove old file: %v", err)
		}
	}

	// Crear directorio si no existe
	if err := ensureUploadDir(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generar nuevo nombre de archivo
	newFilename := laboratoryID + filepath.Ext(file.Filename)
	filePath := filepath.Join(uploadDir, newFilename)

	// Guardar el nuevo archivo
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Actualizar en la base de datos
	tx := db.DB.Begin()
	if err := tx.Model(&models.Laboratory{}).
		Where("id_laboratorio = ?", laboratoryID).
		Update("ruta_archivo_externo", filePath).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, response.LaboratoryFileResponseDTO{
		FilePath: filePath,
		Message:  "File updated successfully",
	})
}

// DeleteLaboratoryFile elimina el archivo del sistema de archivos
func DeleteLaboratoryFile(c *gin.Context) {
	laboratoryID := c.Param("id")

	// Obtener la ruta del archivo primero
	var laboratory models.Laboratory
	if err := db.DB.First(&laboratory, laboratoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Laboratory not found"})
		return
	}

	// Eliminar el archivo si existe
	if laboratory.ExternalFilePath != "" {
		if err := os.Remove(laboratory.ExternalFilePath); err != nil && !os.IsNotExist(err) {
			log.Printf("Warning: Failed to remove file: %v", err)
		}
	}

	// Actualizar la base de datos
	if err := db.DB.Model(&models.Laboratory{}).
		Where("id_laboratorio = ?", laboratoryID).
		Update("ruta_archivo_externo", "").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database"})
		return
	}

	c.Status(http.StatusNoContent)
}
