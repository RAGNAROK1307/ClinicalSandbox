package response

import "time"

type ConsultationVisitResponseDTO struct {
	IDPatient                uint      `json:"id_paciente"`
	IDConsultationVisit      uint      `json:"id_consulta"`
	IDHospitalEmployee       uint      `json:"id_personal_hospital"`
	HospitalEmployeeFullName string    `json:"hospital_employee_full_name"`
	DateTimeVisit            time.Time `json:"fecha_hora_visita"`
	ReasonVisit              string    `json:"razon_visita"`
	MedicalNotes             string    `json:"notas_medicas"`
	ExamResults              string    `json:"resultados_examenes"`
}
