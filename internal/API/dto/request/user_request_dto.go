package request

type CreateUserDTO struct {
	IDRole   uint   `json:"id_rol" binding:"required"`
	UserName string `json:"nombre_usuario" binding:"required"`
	Password string `json:"contraseña" binding:"required,min=8"`
}
