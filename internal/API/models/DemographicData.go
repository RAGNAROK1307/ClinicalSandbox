package models

type DemographicData struct {
	IDDemographicData       uint   `json:"id_datos_demograficos" gorm:"primaryKey;column:id_datos_demograficos" swaggerignore:"true"`
	RacialEthnicInformation string `json:"informacion_etnica_racial" gorm:"type:varchar(50);column:informacion_etnica_racial"`
	MaritalStatus           string `json:"estado_civil" gorm:"type:varchar(50);column:estado_civil"`
	Nationality             string `json:"nacionalidad" gorm:"type:varchar(50);column:nacionalidad"`
	EmploymentInformation   string `json:"informacion_laboral" gorm:"type:text;column:informacion_laboral"`
}

func (DemographicData) TableName() string {
	return "datos_demograficos"
}
