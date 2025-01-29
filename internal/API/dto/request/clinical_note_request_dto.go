package request

type CreateClinicalNoteDTO struct {
	IDConsultationVisit         uint   `json:"id_consulta" binding:"required"`
	ProgressNotes               string `json:"notas_progreso" binding:"required"`
	ObservationsRecommendations string `json:"observaciones_recomendaciones" binding:"required"`
}
