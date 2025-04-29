package response

type MedicalRecordAndRelatedResponseDTO struct {
	MedicalRecord          MedicalRecordResponseDTO           `json:"medical_record"`
	TreatmentPrescriptions []TreatmentPrescriptionResponseDTO `json:"treatment_prescription"`
	ConsultationVisits     []ConsultationVisitResponseDTO     `json:"consultation_visit"`
	Patient                PatientResponseDTO                 `json:"patient"`
}
