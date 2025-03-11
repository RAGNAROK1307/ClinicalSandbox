package response

import "time"

type PatientResponseDTO struct {
	//IDRole               uint   `json:"id_rol" binding:"required"` // ID del rol asociado, requerido
	//IDIdentification     uint   `json:"id_identificacion" binding:"required"`
	//IDDemographicData    uint   `json:"id_datos_demograficos" binding:"required"`
	//IDUser               uint   `json:"id_usuarios" binding:"required"`
	IDPatient               uint      `json:"id_paciente"`
	FullName                string    `json:"nombre_completo"`
	BirthDate               time.Time `json:"fecha_nacimiento" example:"2025-01-20"`
	Gender                  string    `json:"genero"`
	Address                 string    `json:"direccion"`
	Phone                   string    `json:"telefono"`
	SocialSecurityNumber    string    `json:"numero_seguro_social"`
	RoleName                string    `json:"nombre_rol"`
	DocumentType            string    `json:"tipo_documento"`
	DocumentNumber          string    `json:"numero_documento"`
	RacialEthnicInformation string    `json:"informacion_etnica_racial"`
	MaritalStatus           string    `json:"estado_civil"`
	Nationality             string    `json:"nacionalidad"`
	EmploymentInformation   string    `json:"informacion_laboral"`
	UserName                string    `json:"nombre_usuario"`
}
