package request

type CreateHospitalEmployeeDTO struct {
	IDRole           uint   `json:"id_rol" binding:"required"` // ID del rol asociado, requerido
	IDIdentification uint   `json:"id_identificacion" binding:"required"`
	IDUser           uint   `json:"id_usuarios" binding:"required"`
	FullName         string `json:"nombre_completo" binding:"required"`
	BirthDate        string `json:"fecha_nacimiento" binding:"required" example:"2025-01-20"`
	Gender           string `json:"genero" binding:"required"`
	Phone            string `json:"telefono" binding:"required"`
	Address          string `json:"direccion" binding:"required"`
}
