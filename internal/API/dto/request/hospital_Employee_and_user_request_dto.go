package request

type CreateHospitalEmployeeAndUserDTO struct {
	UserName string `json:"nombre_usuario" binding:"required"`
	Password string `json:"contraseña" binding:"required,min=8"`
	//IDRole         uint   `json:"id_rol" binding:"required"`
	FullName       string `json:"nombre_completo" binding:"required"`
	BirthDate      string `json:"fecha_nacimiento" binding:"required" example:"2025-01-20"`
	Gender         string `json:"genero" binding:"required"`
	Phone          string `json:"telefono" binding:"required"`
	Address        string `json:"direccion" binding:"required"`
	DocumentType   string `json:"tipo_documento" binding:"required"`
	DocumentNumber string `json:"numero_documento" binding:"required"`
}

type UpdateHospitalEmployeeAndUserDTO struct {
	UserName string `json:"nombre_usuario" binding:"omitempty"`
	//Password       string `json:"contraseña" binding:"omitempty,min=8"` // No es requerida en PUT
	FullName       string `json:"nombre_completo" binding:"omitempty"`
	BirthDate      string `json:"fecha_nacimiento" binding:"omitempty" example:"2025-01-20"`
	Gender         string `json:"genero" binding:"omitempty"`
	Phone          string `json:"telefono" binding:"omitempty"`
	Address        string `json:"direccion" binding:"omitempty"`
	DocumentType   string `json:"tipo_documento" binding:"omitempty"`
	DocumentNumber string `json:"numero_documento" binding:"omitempty"`
}
