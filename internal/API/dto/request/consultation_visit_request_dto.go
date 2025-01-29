package request

import "time"

type CreateConsultationVisitDTO struct {
	IDPatient          uint      `json:"id_paciente" binding:"required"`
	IDHospitalEmployee uint      `json:"id_personal_hospital" binding:"required"`
	DateTimeVisit      time.Time `json:"fecha_hora_visita" binding:"required"`
	ReasonVisit        string    `json:"razon_visita" binding:"required"`
	MedicalNotes       string    `json:"notas_medicas" binding:"required"`
	ExamResults        string    `json:"resultados_examenes" binding:"required"`
}
