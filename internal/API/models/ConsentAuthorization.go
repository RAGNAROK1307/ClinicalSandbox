package models

import "time"

type ConsentAuthorization struct {
	IDConsentAuthorization uint      `json:"id_consentimiento" gorm:"primaryKey;column:id_consentimiento" swaggerignore:"true"`
	IDPatient              uint      `json:"id_paciente" gorm:"column:id_paciente"`
	ConsentType            string    `json:"tipo_consentimiento" gorm:"type:varchar(50);column:tipo_consentimiento"`
	ConsentDate            time.Time `json:"fecha_consentimiento" gorm:"type:date;column:fecha_consentimiento"`
	Details                string    `json:"detalles" gorm:"type:text;column:detalles"`
	ExternalFilePath       string    `json:"ruta_archivo_externo" gorm:"type:text;column:ruta_archivo_externo"`
	Patient                Patient   `gorm:"foreignKey:IDPatient;references:IDPatient"`
}

func (ConsentAuthorization) TableName() string {
	return "consentimientos_autorizaciones"
}
