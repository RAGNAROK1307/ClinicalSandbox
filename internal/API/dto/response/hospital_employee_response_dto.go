package response

import "time"

type HospitalEmployeeResponseDTO struct {
	//IDRole           uint      `json:"id_rol"` // ID del rol asociado, requerido
	//IDIdentification uint      `json:"id_identificacion"`
	//IDUser           uint      `json:"id_usuarios"`
	IDHospitalEmployee uint      `json:"id_personal_hospital"`
	UserName           string    `json:"nombre_usuario"`
	RoleName           string    `json:"nombre_rol"`
	FullName           string    `json:"nombre_completo"`
	BirthDate          time.Time `json:"fecha_nacimiento" example:"2025-01-20"`
	Gender             string    `json:"genero"`
	Phone              string    `json:"telefono"`
	Address            string    `json:"direccion"`
}
