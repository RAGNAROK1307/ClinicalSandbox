package models

type ClinicalNote struct {
	IDClinicalNote              uint              `json:"id_nota" gorm:"primaryKey;column:id_nota" swaggerignore:"true"`
	IDConsultationVisit         uint              `json:"id_consulta" gorm:"column:id_consulta"`
	ProgressNotes               string            `json:"notas_progreso" gorm:"type:text;column:notas_progreso"`
	ObservationsRecommendations string            `json:"observaciones_recomendaciones" gorm:"type:text;column:observaciones_recomendaciones"`
	ConsultationVisit           ConsultationVisit `gorm:"foreignKey:IDConsultationVisit ;references:IDConsultationVisit "`
}

func (ClinicalNote) TableName() string {
	return "notas_clinicas"
}
