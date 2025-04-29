package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/response"
	"ClinicalSandBox/internal/API/models"
	"bytes"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	_ "time"
)

// UploadLaboratoryFile godoc
// @Summary Upload laboratory file
// @Description Upload a file (PDF, image, etc.) for a laboratory result
// @Tags laboratories
// @Accept multipart/form-data
// @Param id path string true "Laboratory ID"
// @Param file formData file true "Laboratory file"
// @Success 200 {object} response.LaboratoryFileResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /laboratories/{id}/file [post]
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

	// Validar tamaño del archivo (max 10MB)
	if file.Size > 10<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB limit"})
		return
	}

	// Leer el archivo
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer src.Close()

	// Validar tipo de archivo (permite PDF e imágenes)
	buff := make([]byte, 512)
	if _, err := src.Read(buff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	mimeType := http.DetectContentType(buff)
	if !strings.HasPrefix(mimeType, "image/") && !strings.HasPrefix(mimeType, "application/pdf") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image or PDF"})
		return
	}

	// Volver al inicio del archivo
	if _, err := src.Seek(0, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}

	// Convertir a Base64
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}
	fileBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Actualizar en la base de datos
	tx := db.DB.Begin()
	if err := tx.Model(&models.Laboratory{}).
		Where("id_laboratorio = ?", laboratoryID).
		Update("ruta_archivo_externo", fileBase64).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, response.LaboratoryFileResponseDTO{
		FileContent: fileBase64,
		MimeType:    mimeType,
		Message:     "File uploaded successfully",
	})
}

// GetLaboratoryFile godoc
// @Summary Get laboratory file
// @Description Get the laboratory file (base64 encoded)
// @Tags laboratories
// @Param id path string true "Laboratory ID"
// @Produce json
// @Success 200 {object} response.LaboratoryFileResponseDTO
// @Failure 404 {object} map[string]string
// @Router /laboratories/{id}/file [get]
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

	// Intentar detectar el tipo MIME
	mimeType := "application/octet-stream" // valor por defecto
	if decoded, err := base64.StdEncoding.DecodeString(laboratory.ExternalFilePath); err == nil {
		if len(decoded) > 512 {
			mimeType = http.DetectContentType(decoded[:512])
		} else {
			mimeType = http.DetectContentType(decoded)
		}
	}

	c.JSON(http.StatusOK, response.LaboratoryFileResponseDTO{
		FileContent: laboratory.ExternalFilePath,
		MimeType:    mimeType,
	})
}

// UpdateLaboratoryFile godoc
// @Summary Update laboratory file
// @Description Update the laboratory file (base64 encoded)
// @Tags laboratories
// @Accept multipart/form-data
// @Param id path string true "Laboratory ID"
// @Param file formData file true "Laboratory file"
// @Success 200 {object} response.LaboratoryFileResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /laboratories/{id}/file [put]
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

	// Validar tamaño del archivo (max 10MB)
	if file.Size > 10<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB limit"})
		return
	}

	// Leer y validar el archivo
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer src.Close()

	// Validar tipo de archivo
	buff := make([]byte, 512)
	if _, err := src.Read(buff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	mimeType := http.DetectContentType(buff)
	if !strings.HasPrefix(mimeType, "image/") && !strings.HasPrefix(mimeType, "application/pdf") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image or PDF"})
		return
	}

	// Volver al inicio del archivo para leerlo completo
	if _, err := src.Seek(0, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}

	// Convertir a Base64
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}
	fileBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Actualizar en la base de datos con transacción
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&models.Laboratory{}).
		Where("id_laboratorio = ?", laboratoryID).
		Update("ruta_archivo_externo", fileBase64).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update file"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Registrar la actualización
	log.Printf("File updated for laboratory ID: %s", laboratoryID)

	c.JSON(http.StatusOK, response.LaboratoryFileResponseDTO{
		FileContent: fileBase64,
		MimeType:    mimeType,
		Message:     "File updated successfully",
	})
}

// DeleteLaboratoryFile godoc
// @Summary Delete laboratory file
// @Description Delete the laboratory file from database
// @Tags laboratories
// @Param id path string true "Laboratory ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /laboratories/{id}/file [delete]
func DeleteLaboratoryFile(c *gin.Context) {
	laboratoryID := c.Param("id")

	// Verificar si existe el laboratorio
	var count int64
	if err := db.DB.Model(&models.Laboratory{}).
		Where("id_laboratorio = ?", laboratoryID).
		Count(&count).Error; err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Laboratory not found"})
		return
	}

	// Actualizar en la base de datos (establecer archivo como vacío)
	if err := db.DB.Model(&models.Laboratory{}).
		Where("id_laboratorio = ?", laboratoryID).
		Update("ruta_archivo_externo", "").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.Status(http.StatusNoContent)
}
