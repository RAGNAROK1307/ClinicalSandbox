package request

type CreateLaboratoryDTO struct {
	IDConsultationVisit uint   `json:"id_consulta" binding:"required"`
	TestDate            string `json:"fecha_prueba" binding:"required" example:"2025-01-20"`
	TestType            string `json:"tipo_prueba" binding:"required"`
	TestResults         string `json:"resultados_prueba" binding:"required"`
	ExternalFilePath    string `json:"ruta_archivo_externo" binding:"required"`
}
