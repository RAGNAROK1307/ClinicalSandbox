package models

import "time"

type ConsultationVisit struct {
	IDConsultationVisit uint             `json:"id_consulta" gorm:"primaryKey;column:id_consulta" swaggerignore:"true"`
	IDPatient           uint             `json:"id_paciente" gorm:"column:id_paciente"`
	IDHospitalEmployee  uint             `json:"id_personal_hospital" gorm:"column:id_personal_hospital"`
	DateTimeVisit       time.Time        `json:"fecha_hora_visita" gorm:"type:date;column:fecha_hora_visita"`
	ReasonVisit         string           `json:"razon_visita" gorm:"type:text;column:razon_visita"`
	MedicalNotes        string           `json:"notas_medicas" gorm:"type:text;column:notas_medicas"`
	ExamResults         string           `json:"resultados_examenes" gorm:"type:text;column:resultados_examenes"`
	Patient             Patient          `gorm:"foreignKey:IDPatient;references:IDPatient"`
	HospitalEmployee    HospitalEmployee `gorm:"foreignKey:IDHospitalEmployee;references:IDHospitalEmployee"`
}

func (ConsultationVisit) TableName() string {
	return "consultas_visitas"
}
