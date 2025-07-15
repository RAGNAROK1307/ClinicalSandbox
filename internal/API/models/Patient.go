package models

import "time"

type Patient struct {
	IDPatient            uint            `json:"id_paciente" gorm:"primaryKey;column:id_paciente" swaggerignore:"true"`
	IDRole               uint            `json:"id_rol" gorm:"column:id_rol"`
	IDIdentification     uint            `json:"id_identificacion" gorm:"column:id_identificacion"`
	IDDemographicData    uint            `json:"id_datos_demograficos" gorm:"column:id_datos_demograficos"`
	IDUser               uint            `json:"id_usuarios" gorm:"column:id_usuarios"`
	FullName             string          `json:"nombre_completo" gorm:"type:varchar(50);column:nombre_completo"`
	BirthDate            time.Time       `json:"fecha_nacimiento" gorm:"type:date;column:fecha_nacimiento"`
	Gender               string          `json:"genero" gorm:"type:varchar(20);column:genero"`
	Address              string          `json:"direccion" gorm:"type:varchar(50);column:direccion"`
	Phone                string          `json:"telefono" gorm:"type:varchar(20);column:telefono"`
	SocialSecurityNumber string          `json:"numero_seguro_social" gorm:"type:varchar(50);column:numero_seguro_social"`
	Role                 Role            `gorm:"foreignKey:IDRole;references:IDRole"`
	Identification       Identification  `gorm:"foreignKey:IDIdentification;references:IDIdentification"`
	DemographicData      DemographicData `gorm:"foreignKey:IDDemographicData;references:IDDemographicData"`
	User                 User            `gorm:"foreignKey:IDUser;references:IDUser"`
}

func (Patient) TableName() string {
	return "pacientes"
}
