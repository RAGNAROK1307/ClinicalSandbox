package request

type CreateIdentificationDTO struct {
	DocumentType   string `json:"tipo_documento" binding:"required"`
	DocumentNumber string `json:"numero_documento" binding:"required"`
}
