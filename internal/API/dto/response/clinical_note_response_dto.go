package response

type ClinicalResponseNoteDTO struct {
	IDConsultationVisit         uint   `json:"id_consulta"`
	ProgressNotes               string `json:"notas_progreso"`
	ObservationsRecommendations string `json:"observaciones_recomendaciones"`
}
