package dto

type CreateRoleDTO struct {
	NombreRol   string `json:"nombre_rol" binding:"required"`
	Descripcion string `json:"descripcion" binding:"required"`
}
