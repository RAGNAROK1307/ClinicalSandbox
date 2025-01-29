package models

type Identification struct {
	IDIdentification uint   `json:"id_identificacion" gorm:"primaryKey;column:id_identificacion" swaggerignore:"true"`
	DocumentType     string `json:"tipo_documento" gorm:"type:varchar(20);column:tipo_documento"`
	DocumentNumber   string `json:"numero_documento" gorm:"type:varchar(20);column:numero_documento"`
}

func (Identification) TableName() string {
	return "identificacion"
}
