package models

import "time"

type HospitalEmployee struct {
	IDHospitalEmployee uint           `json:"id_personal_hospital" gorm:"primaryKey;column:id_personal_hospital" swaggerignore:"true"`
	IDRole             uint           `json:"id_rol" gorm:"column:id_rol"` // Llave foránea que referencia la tabla roles
	IDIdentification   uint           `json:"id_identificacion" gorm:"column:id_identificacion"`
	FullName           string         `json:"nombre_completo" gorm:"type:varchar(50);column:nombre_completo"`
	BirthDate          time.Time      `json:"fecha_nacimiento" gorm:"type:date;column:fecha_nacimiento"`
	Gender             string         `json:"genero" gorm:"type:varchar(20);column:genero"`
	Phone              string         `json:"telefono" gorm:"type:varchar(20);column:telefono"`
	Address            string         `json:"direccion" gorm:"type:varchar(50);column:direccion"`
	Role               Role           `gorm:"foreignKey:IDRole;references:IDRole"` // Carga la información del rol relacionado
	Identification     Identification `gorm:"foreignKey:IDIdentification;references:IDIdentification"`
}

func (HospitalEmployee) TableName() string {
	return "personal_hospital"
}
