package response

import "time"

type HospitalEmployeeResponseDTO struct {
	IDHospitalEmployee uint      `json:"id_personal_hospital"`
	FullName           string    `json:"nombre_completo"`
	BirthDate          time.Time `json:"fecha_nacimiento" example:"2025-01-20"`
	Gender             string    `json:"genero"`
	Phone              string    `json:"telefono"`
	Address            string    `json:"direccion"`
}
