package response

type UserResponseDTO struct {
	IDUser uint `json:"id_usuario"` // ID del usuario
	//IDRole   uint   `json:"id_rol"`         // ID del rol asociado
	UserName string `json:"nombre_usuario"`
	RoleName string `json:"nombre_rol"`
	//Password string `json:"contraseña"`
}
