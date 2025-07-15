package response

type LaboratoryResponseDTO struct {
	IDLaboratory        uint   `json:"id_laboratorio"`
	IDConsultationVisit uint   `json:"id_consulta"`
	TestDate            string `json:"fecha_prueba" example:"2025-01-20"`
	TestType            string `json:"tipo_prueba" `
	TestResults         string `json:"resultados_prueba" `
	ExternalFilePath    string `json:"ruta_archivo_externo"`
}

type LaboratoryFileResponseDTO struct {
	FilePath string `json:"file_path"` // Ruta del archivo en el sistema
	MimeType string `json:"mime_type"` // Tipo MIME del archivo (ej: "image/jpeg",)
	Message  string `json:"message"`   // Mensaje opcional
}
