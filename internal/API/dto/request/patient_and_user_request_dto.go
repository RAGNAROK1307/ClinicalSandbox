package request

type CreateUserAndPatientDTO struct {
	UserName string `json:"nombre_usuario" binding:"required"`
	Password string `json:"contraseña" binding:"required,min=8"`
	//IDRole                  uint   `json:"id_rol" binding:"required"`
	FullName                string `json:"nombre_completo" binding:"required"`
	BirthDate               string `json:"fecha_nacimiento" binding:"required" example:"2025-01-20"`
	Gender                  string `json:"genero" binding:"required"`
	Address                 string `json:"direccion" binding:"required"`
	Phone                   string `json:"telefono" binding:"required"`
	SocialSecurityNumber    string `json:"numero_seguro_social" binding:"required"`
	DocumentType            string `json:"tipo_documento" binding:"required"`
	DocumentNumber          string `json:"numero_documento" binding:"required"`
	RacialEthnicInformation string `json:"informacion_etnica_racial" binding:"required"`
	MaritalStatus           string `json:"estado_civil" binding:"required"`
	Nationality             string `json:"nacionalidad" binding:"required"`
	EmploymentInformation   string `json:"informacion_laboral" binding:"required"`
}

type UpdateUserAndPatientDTO struct {
	UserName string `json:"nombre_usuario" binding:"omitempty"`
	//Password                string `json:"contraseña" binding:"required,min=8"`
	//IDRole                  uint   `json:"id_rol" binding:"omitempty"`
	FullName                string `json:"nombre_completo" binding:"omitempty"`
	BirthDate               string `json:"fecha_nacimiento" binding:"omitempty" example:"2025-01-20"`
	Gender                  string `json:"genero" binding:"omitempty"`
	Address                 string `json:"direccion" binding:"omitempty"`
	Phone                   string `json:"telefono" binding:"omitempty"`
	SocialSecurityNumber    string `json:"numero_seguro_social" binding:"omitempty"`
	DocumentType            string `json:"tipo_documento" binding:"omitempty"`
	DocumentNumber          string `json:"numero_documento" binding:"omitempty"`
	RacialEthnicInformation string `json:"informacion_etnica_racial" binding:"omitempty"`
	MaritalStatus           string `json:"estado_civil" binding:"omitempty"`
	Nationality             string `json:"nacionalidad" binding:"omitempty"`
	EmploymentInformation   string `json:"informacion_laboral" binding:"omitempty"`
}
