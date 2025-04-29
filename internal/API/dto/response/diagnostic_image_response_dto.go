package response

type DiagnosticImageResponseDTO struct {
	IDDiagnosticImage   uint   `json:"id_imagen"`
	IDConsultationVisit uint   `json:"id_consulta"`
	ImageDate           string `json:"fecha_imagen"  example:"2025-01-20"`
	ImageType           string `json:"tipo_imagen"`
	Description         string `json:"descripcion"`
	ImageInterpretation string `json:"interpretacion_imagen"`
	ExternalFilePath    string `json:"ruta_archivo_externo"`
}

type DiagnosticImageFileResponseDTO struct {
	Image    string `json:"image"`     // Imagen en base64
	MimeType string `json:"mime_type"` // Tipo MIME de la imagen
}
