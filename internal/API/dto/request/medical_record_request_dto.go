package request

import "time"

type CreateMedicalRecordDTO struct {
	IDPatient           uint      `json:"id_paciente" binding:"required"`
	PatientImage        string    `json:"imagen_paciente" binding:"required"`
	PreviousDiagnoses   string    `json:"diagnosticos_previos" binding:"required"`
	ChronicDiseases     string    `json:"enfermedades_cronicas" binding:"required"`
	Allergies           string    `json:"alergias" binding:"required"`
	CurrentMedications  string    `json:"medicamentos_actuales" binding:"required"`
	SurgeryHistory      string    `json:"historial_cirugias" binding:"required"`
	FamilyHistory       string    `json:"antecedentes_familiares" binding:"required"`
	CreationDateHistory time.Time `json:"fecha_creacion_historial" binding:"required"`
	HistoryUpdateDate   time.Time `json:"fecha_actualizacion_historial" binding:"required"`
}

type UpdateMedicalRecordDTO struct {
	PreviousDiagnoses   string    `json:"diagnosticos_previos"`
	ChronicDiseases     string    `json:"enfermedades_cronicas"`
	Allergies           string    `json:"alergias"`
	CurrentMedications  string    `json:"medicamentos_actuales"`
	SurgeryHistory      string    `json:"historial_cirugias"`
	FamilyHistory       string    `json:"antecedentes_familiares"`
	CreationDateHistory time.Time `json:"fecha_creacion_historial"`
	HistoryUpdateDate   time.Time `json:"fecha_actualizacion_historial"`
}
