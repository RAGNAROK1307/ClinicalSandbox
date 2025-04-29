package request

type CreateMedicalRecordAndRelatedDTO struct {
	// Campos para MedicalRecord
	IDPatient uint `json:"id_paciente" binding:"required"`
	//PatientImage       string `json:"imagen_paciente" binding:"required"`
	PreviousDiagnoses  string `json:"diagnosticos_previos"`
	ChronicDiseases    string `json:"enfermedades_cronicas"`
	Allergies          string `json:"alergias"`
	CurrentMedications string `json:"medicamentos_actuales"`
	SurgeryHistory     string `json:"historial_cirugias"`
	FamilyHistory      string `json:"antecedentes_familiares"`

	// Campos para TreatmentPrescription
	PrescribedMedication string `json:"medicamento_prescrito"`
	Amount               string `json:"cantidad"`
	Frequency            string `json:"frecuencia"`
	Instructions         string `json:"instrucciones"`
	DurationTreatment    string `json:"duracion_tratamiento"`

	// Campos para ConsultationVisit
	IDHospitalEmployee uint   `json:"id_empleado_hospital" binding:"required"`
	DateTimeVisit      string `json:"fecha_hora_visita" binding:"required"`
	ReasonVisit        string `json:"motivo_visita" binding:"required"`
	MedicalNotes       string `json:"notas_medicas"`
	ExamResults        string `json:"resultados_examenes"`
}
