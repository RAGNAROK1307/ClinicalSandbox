package request

type CreateRoleDTO struct {
	RoleName    string `json:"nombre_rol" binding:"required"`
	Description string `json:"descripcion" binding:"required"`
}
