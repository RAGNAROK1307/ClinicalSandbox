package request

type CreatePatientDTO struct {
	IDRole               uint   `json:"id_rol" binding:"required"`
	IDIdentification     uint   `json:"id_identificacion" binding:"required"`
	IDDemographicData    uint   `json:"id_datos_demograficos" binding:"required"`
	IDUser               uint   `json:"id_usuarios" binding:"required"`
	FullName             string `json:"nombre_completo" binding:"required"`
	BirthDate            string `json:"fecha_nacimiento" binding:"required" example:"2025-01-20"`
	Gender               string `json:"genero" binding:"required"`
	Address              string `json:"direccion" binding:"required"`
	Phone                string `json:"telefono" binding:"required"`
	SocialSecurityNumber string `json:"numero_seguro_social" binding:"required"`
}
