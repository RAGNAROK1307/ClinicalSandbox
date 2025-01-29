package request

type CreateConsentAuthorizationDTO struct {
	IDPatient        uint   `json:"id_paciente" binding:"required"`
	ConsentType      string `json:"tipo_consentimiento" binding:"required"`
	ConsentDate      string `json:"fecha_consentimiento" binding:"required" example:"2025-01-20"`
	Details          string `json:"detalles" binding:"required"`
	ExternalFilePath string `json:"ruta_archivo_externo" binding:"required"`
}
