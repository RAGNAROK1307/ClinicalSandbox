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
	FileContent string `json:"file_content"` // Contenido del archivo en base64
	MimeType    string `json:"mime_type"`    // Tipo MIME del archivo (ej: "image/jpeg", "application/pdf")
	Message     string `json:"message"`      // Mensaje opcional
}
