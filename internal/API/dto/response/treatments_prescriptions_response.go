package response

type TreatmentPrescriptionResponseDTO struct {
	IDPatient               uint   `json:"id_paciente"`
	IDTreatmentPrescription uint   `json:"id_tratamiento"`
	PrescribedMedication    string `json:"medicamento_prescrito"`
	Amount                  string `json:"dosis"`
	Frequency               string `json:"frecuencia"`
	Instructions            string `json:"instrucciones"`
	DurationTreatment       string `json:"duracion_tratamiento"`
}
