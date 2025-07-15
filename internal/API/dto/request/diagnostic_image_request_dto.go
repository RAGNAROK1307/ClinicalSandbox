package request

type CreateDiagnosticImageDTO struct {
	IDConsultationVisit uint   `json:"id_consulta" binding:"required"`
	ImageDate           string `json:"fecha_imagen" binding:"required" example:"2025-01-20"`
	ImageType           string `json:"tipo_imagen" binding:"required"`
	Description         string `json:"descripcion" binding:"required"`
	ImageInterpretation string `json:"interpretacion_imagen" binding:"required"`
}
