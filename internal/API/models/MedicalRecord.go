package models

import "time"

type MedicalRecord struct {
	IDMedicalRecord     uint      `json:"id_historial" gorm:"primaryKey;column:id_historial" swaggerignore:"true"`
	IDPatient           uint      `json:"id_paciente" gorm:"column:id_paciente"` // Llave foránea que referencia la tabla roles
	PreviousDiagnoses   string    `json:"diagnosticos_previos" gorm:"type:text;column:diagnosticos_previos"`
	ChronicDiseases     string    `json:"enfermedades_cronicas" gorm:"type:text;column:enfermedades_cronicas"`
	Allergies           string    `json:"alergias" gorm:"type:text;column:alergias"`
	CurrentMedications  string    `json:"medicamentos_actuales" gorm:"type:text;column:medicamentos_actuales"`
	SurgeryHistory      string    `json:"historial_cirugias" gorm:"type:text;column:historial_cirugias"`
	FamilyHistory       string    `json:"antecedentes_familiares" gorm:"type:text;column:antecedentes_familiares"`
	CreationDateHistory time.Time `json:"fecha_creacion_historial" gorm:"type:date;column:fecha_creacion_historial"`
	HistoryUpdateDate   time.Time `json:"fecha_actualizacion_historial" gorm:"type:date;column:fecha_actualizacion_historial"`
	Patient             Patient   `gorm:"foreignKey:IDPatient;references:IDPatient"`
}

func (MedicalRecord) TableName() string {
	return "historial_medico"
}
