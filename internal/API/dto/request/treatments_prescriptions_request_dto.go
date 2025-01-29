package request

type CreateTreatmentPrescriptionDTO struct {
	IDPatient            uint   `json:"id_paciente" binding:"required"`
	PrescribedMedication string `json:"medicamento_prescrito" binding:"required"`
	Amount               string `json:"dosis" binding:"required"`
	Frequency            string `json:"frecuencia" binding:"required"`
	Instructions         string `json:"instrucciones" binding:"required"`
	DurationTreatment    string `json:"duracion_tratamiento" binding:"required"`
}
