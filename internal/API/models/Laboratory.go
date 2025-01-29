package models

import "time"

type Laboratory struct {
	IDLaboratory        uint              `json:"id_laboratorio" gorm:"primaryKey;column:id_laboratorio" swaggerignore:"true"`
	IDConsultationVisit uint              `json:"id_consulta" gorm:"column:id_consulta"`
	TestDate            time.Time         `json:"fecha_prueba" gorm:"type:date;column:fecha_prueba"`
	TestType            string            `json:"tipo_prueba" gorm:"type:varchar(50);column:tipo_prueba"`
	TestResults         string            `json:"resultados_prueba" gorm:"type:text;column:resultados_prueba"`
	ExternalFilePath    string            `json:"ruta_archivo_externo" gorm:"type:text;column:ruta_archivo_externo"`
	ConsultationVisit   ConsultationVisit `gorm:"foreignKey:IDConsultationVisit;references:IDConsultationVisit"`
}

func (Laboratory) TableName() string {
	return "laboratorio"
}
