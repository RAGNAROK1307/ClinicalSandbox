package response

type UserResponseDTO struct {
	IDUser   uint   `json:"id_usuario"`
	UserName string `json:"nombre_usuario"`
	RoleName string `json:"nombre_rol"`
}
