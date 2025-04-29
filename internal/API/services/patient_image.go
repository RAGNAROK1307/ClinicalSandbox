package services

import (
	"ClinicalSandBox/configs/db"
	_ "ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/dto/response"
	"ClinicalSandBox/internal/API/models"
	"bytes"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

// UploadPatientImage godoc
// @Summary Upload patient image
// @Description Upload a patient image for medical records
// @Tags medical_records
// @Accept multipart/form-data
// @Param id path string true "Patient ID"
// @Param file formData file true "Patient image"
// @Success 200 {object} response.MedicalRecordImageResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /medical-records/{id}/image [post]
func UploadPatientImage(c *gin.Context) {
	idPatient := c.Param("id")

	// Verificar si existe el historial médico
	var count int64
	if err := db.DB.Model(&models.MedicalRecord{}).
		Where("id_paciente = ?", idPatient).
		Count(&count).Error; err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found for the patient"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validar tamaño del archivo (max 5MB)
	if file.Size > 5<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 5MB limit"})
		return
	}

	// Leer el archivo
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer src.Close()

	// Validar que es una imagen
	buff := make([]byte, 512)
	if _, err := src.Read(buff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image file"})
		return
	}

	mimeType := http.DetectContentType(buff)
	if !strings.HasPrefix(mimeType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is not an image"})
		return
	}

	// Volver al inicio del archivo
	if _, err := src.Seek(0, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
		return
	}

	// Convertir a Base64
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
		return
	}
	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Actualizar en la base de datos
	if err := db.DB.Model(&models.MedicalRecord{}).
		Where("id_paciente = ?", idPatient).
		Update("imagen_paciente", imageBase64).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	c.JSON(http.StatusOK, response.MedicalRecordImageResponseDTO{
		Image:    imageBase64,
		MimeType: mimeType,
	})
}

// GetPatientImage godoc
// @Summary Get patient image
// @Description Get a patient image from medical records
// @Tags medical_records
// @Param id path string true "Patient ID"
// @Produce json
// @Success 200 {object} response.MedicalRecordImageResponseDTO
// @Failure 404 {object} map[string]string
// @Router /medical-records/{id}/image [get]
func GetPatientImage(c *gin.Context) {
	idPatient := c.Param("id")

	var medicalRecord models.MedicalRecord
	if err := db.DB.Select("imagen_paciente").
		Where("id_paciente = ?", idPatient).
		First(&medicalRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found"})
		return
	}

	if medicalRecord.PatientImage == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No image found for this patient"})
		return
	}

	// Intentar detectar el tipo MIME
	mimeType := "image/jpeg" // valor por defecto
	if decoded, err := base64.StdEncoding.DecodeString(medicalRecord.PatientImage); err == nil {
		if len(decoded) > 512 {
			mimeType = http.DetectContentType(decoded[:512])
		} else {
			mimeType = http.DetectContentType(decoded)
		}
	}

	c.JSON(http.StatusOK, response.MedicalRecordImageResponseDTO{
		Image:    medicalRecord.PatientImage,
		MimeType: mimeType,
	})
}

// UpdatePatientImage godoc
// @Summary Update patient image
// @Description Update a patient image in medical records
// @Tags medical_records
// @Accept multipart/form-data
// @Param id path string true "Patient ID"
// @Param file formData file true "Patient image"
// @Success 200 {object} response.MedicalRecordImageResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /medical-records/{id}/image [put]
func UpdatePatientImage(c *gin.Context) {
	idPatient := c.Param("id")

	// Verificar si existe el historial médico
	var medicalRecord models.MedicalRecord
	if err := db.DB.Where("id_paciente = ?", idPatient).First(&medicalRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found for the patient"})
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

	// Validar tamaño del archivo (max 5MB)
	if file.Size > 5<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 5MB limit"})
		return
	}

	// Leer y validar el archivo
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer src.Close()

	// Validar tipo de archivo (debe ser imagen)
	buff := make([]byte, 512)
	if _, err := src.Read(buff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image file"})
		return
	}

	mimeType := http.DetectContentType(buff)
	if !strings.HasPrefix(mimeType, "image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is not an image"})
		return
	}

	// Volver al inicio del archivo para leerlo completo
	if _, err := src.Seek(0, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
		return
	}

	// Convertir a Base64
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
		return
	}
	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Actualizar en la base de datos con transacción
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&models.MedicalRecord{}).
		Where("id_paciente = ?", idPatient).
		Updates(map[string]interface{}{
			"imagen_paciente":               imageBase64,
			"fecha_actualizacion_historial": time.Now(),
		}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update image"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Registrar la actualización
	log.Printf("Image updated for patient ID: %s", idPatient)

	c.JSON(http.StatusOK, response.MedicalRecordImageResponseDTO{
		Image:    imageBase64,
		MimeType: mimeType,
	})
}

// DeletePatientImage godoc
// @Summary Delete patient image
// @Description Delete a patient image from medical records
// @Tags medical_records
// @Param id path string true "Patient ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /medical-records/{id}/image [delete]
func DeletePatientImage(c *gin.Context) {
	idPatient := c.Param("id")

	// Verificar si existe el historial médico
	var count int64
	if err := db.DB.Model(&models.MedicalRecord{}).
		Where("id_paciente = ?", idPatient).
		Count(&count).Error; err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found for the patient"})
		return
	}

	// Actualizar en la base de datos (establecer imagen como vacía)
	if err := db.DB.Model(&models.MedicalRecord{}).
		Where("id_paciente = ?", idPatient).
		Update("imagen_paciente", "").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image"})
		return
	}

	c.Status(http.StatusNoContent)
}
