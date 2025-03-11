package response

type IdentificationResponseDTO struct {
	IDIdentification uint   `json:"id_identificacion"`
	DocumentType     string `json:"tipo_documento"`
	DocumentNumber   string `json:"numero_documento"`
}
