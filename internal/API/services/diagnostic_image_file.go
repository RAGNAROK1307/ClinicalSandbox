package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/response"
	"ClinicalSandBox/internal/API/models"
	"bytes"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	_ "time"
)

// UploadDiagnosticImage godoc
// @Summary Upload diagnostic image
// @Description Upload a diagnostic image file (stored as base64 in database)
// @Tags diagnostic_images
// @Accept multipart/form-data
// @Param id path string true "DiagnosticImage ID"
// @Param file formData file true "Diagnostic image file"
// @Success 200 {object} response.DiagnosticImageFileResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /diagnostic_images/{id}/file [post]
func UploadDiagnosticImage(c *gin.Context) {
	imageID := c.Param("id")

	var diagnosticImage models.DiagnosticImage
	if err := db.DB.First(&diagnosticImage, imageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diagnostic image not found"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// ❌ Eliminado: Validación de tamaño de archivo
	// if file.Size > 10<<20 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB limit"})
	// 	return
	// }

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer src.Close()

	// ❌ Eliminado: Validación de que es una imagen
	// buff := make([]byte, 512)
	// if _, err := src.Read(buff); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image file"})
	// 	return
	// }

	// mimeType := http.DetectContentType(buff)
	// if !strings.HasPrefix(mimeType, "image/") {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is not an image"})
	// 	return
	// }

	// if _, err := src.Seek(0, 0); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
	// 	return
	// }

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}
	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Guardar directamente sin verificar tipo ni tamaño
	tx := db.DB.Begin()
	if err := tx.Model(&models.DiagnosticImage{}).
		Where("id_imagen = ?", imageID).
		Update("ruta_archivo_externo", imageBase64).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	tx.Commit()

	// ❌ mimeType puede estar vacío si no se valida antes
	c.JSON(http.StatusOK, response.DiagnosticImageFileResponseDTO{
		Image:    imageBase64,
		MimeType: "application/octet-stream", // tipo genérico
	})
}

// GetDiagnosticImageFile godoc
// @Summary Get diagnostic image
// @Description Get the diagnostic image (base64 encoded)
// @Tags diagnostic_images
// @Param id path string true "DiagnosticImage ID"
// @Produce json
// @Success 200 {object} response.DiagnosticImageFileResponseDTO
// @Failure 404 {object} map[string]string
// @Router /diagnostic_images/{id}/file [get]
func GetDiagnosticImageFile(c *gin.Context) {
	imageID := c.Param("id")

	var diagnosticImage models.DiagnosticImage
	if err := db.DB.First(&diagnosticImage, imageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diagnostic image not found"})
		return
	}

	if diagnosticImage.ExternalFilePath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No image found for this diagnostic"})
		return
	}

	// Intentar detectar el tipo MIME
	mimeType := "image/jpeg" // valor por defecto
	if decoded, err := base64.StdEncoding.DecodeString(diagnosticImage.ExternalFilePath); err == nil {
		if len(decoded) > 512 {
			mimeType = http.DetectContentType(decoded[:512])
		} else {
			mimeType = http.DetectContentType(decoded)
		}
	}

	c.JSON(http.StatusOK, response.DiagnosticImageFileResponseDTO{
		Image:    diagnosticImage.ExternalFilePath,
		MimeType: mimeType,
	})
}

// UpdateDiagnosticImageFile godoc
// @Summary Update diagnostic image
// @Description Update the diagnostic image (base64 encoded)
// @Tags diagnostic_images
// @Accept multipart/form-data
// @Param id path string true "DiagnosticImage ID"
// @Param file formData file true "Diagnostic image file"
// @Success 200 {object} response.DiagnosticImageFileResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /diagnostic_images/{id}/file [put]
func UpdateDiagnosticImageFile(c *gin.Context) {
	imageID := c.Param("id")

	// Verificar si existe la imagen diagnóstica
	var diagnosticImage models.DiagnosticImage
	if err := db.DB.First(&diagnosticImage, imageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Diagnostic image not found"})
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

	// ❌ Eliminado: Validación de tamaño
	// if file.Size > 10<<20 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB limit"})
	// 	return
	// }

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer src.Close()

	// ❌ Eliminado: Validación del tipo de archivo
	// buff := make([]byte, 512)
	// if _, err := src.Read(buff); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image file"})
	// 	return
	// }

	// mimeType := http.DetectContentType(buff)
	// if !strings.HasPrefix(mimeType, "image/") {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is not an image"})
	// 	return
	// }

	// if _, err := src.Seek(0, 0); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
	// 	return
	// }

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}
	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Guardar sin validaciones
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&models.DiagnosticImage{}).
		Where("id_imagen = ?", imageID).
		Update("ruta_archivo_externo", imageBase64).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update file"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// ❌ Tipo MIME genérico
	c.JSON(http.StatusOK, response.DiagnosticImageFileResponseDTO{
		Image:    imageBase64,
		MimeType: "application/octet-stream",
	})
}

// DeleteDiagnosticImageFile godoc
// @Summary Delete diagnostic image
// @Description Delete the diagnostic image from database
// @Tags diagnostic_images
// @Param id path string true "DiagnosticImage ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /diagnostic_images/{id}/file [delete]
func DeleteDiagnosticImageFile(c *gin.Context) {
	imageID := c.Param("id")

	// Verificar si existe la imagen diagnóstica
	var count int64
	if err := db.DB.Model(&models.DiagnosticImage{}).
		Where("id_imagen = ?", imageID).
		Count(&count).Error; err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diagnostic image not found"})
		return
	}

	// Actualizar en la base de datos (establecer imagen como vacía)
	if err := db.DB.Model(&models.DiagnosticImage{}).
		Where("id_imagen = ?", imageID).
		Update("ruta_archivo_externo", "").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
		return
	}

	c.Status(http.StatusNoContent)
}
