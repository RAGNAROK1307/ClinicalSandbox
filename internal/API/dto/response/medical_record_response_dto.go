package response

import "time"

type MedicalRecordResponseDTO struct {
	IDPatient uint `json:"id_paciente"` // ID del rol asociado, requerido
	//IDConsultationVisit     uint      `json:"id_consulta"`
	//IDTreatmentPrescription uint      `json:"id_tratamiento"`
	PatientImage        string    `json:"imagen_paciente"`
	PreviousDiagnoses   string    `json:"diagnosticos_previos"`
	ChronicDiseases     string    `json:"enfermedades_cronicas"`
	Allergies           string    `json:"alergias"`
	CurrentMedications  string    `json:"medicamentos_actuales"`
	SurgeryHistory      string    `json:"historial_cirugias" `
	FamilyHistory       string    `json:"antecedentes_familiares"`
	CreationDateHistory time.Time `json:"fecha_creacion_historial"`
	HistoryUpdateDate   time.Time `json:"fecha_actualizacion_historial"`
}

type MedicalRecordImageResponseDTO struct {
	Image    string `json:"image"`
	MimeType string `json:"mime_type,omitempty"`
}
