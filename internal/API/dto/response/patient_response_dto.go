package response

import "time"

type PatientResponseDTO struct {
	IDPatient            uint      `json:"id_paciente"`
	FullName             string    `json:"nombre_completo"`
	BirthDate            time.Time `json:"fecha_nacimiento" example:"2025-01-20"`
	Gender               string    `json:"genero"`
	Address              string    `json:"direccion"`
	Phone                string    `json:"telefono"`
	SocialSecurityNumber string    `json:"numero_seguro_social"`
}
