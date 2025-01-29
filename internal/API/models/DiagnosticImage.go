package models

import "time"

type DiagnosticImage struct {
	IDDiagnosticImage   uint              `json:"id_imagen" gorm:"primaryKey;column:id_imagen" swaggerignore:"true"`
	IDConsultationVisit uint              `json:"id_consulta" gorm:"column:id_consulta"`
	ImageDate           time.Time         `json:"fecha_imagen" gorm:"type:date;column:fecha_imagen"`
	ImageType           string            `json:"tipo_imagen" gorm:"type:varchar(50);column:tipo_imagen"`
	Description         string            `json:"descripcion" gorm:"type:text;column:descripcion"`
	ImageInterpretation string            `json:"interpretacion_imagen" gorm:"type:text;column:interpretacion_imagen"`
	ExternalFilePath    string            `json:"ruta_archivo_externo" gorm:"type:text;column:ruta_archivo_externo"`
	ConsultationVisit   ConsultationVisit `gorm:"foreignKey:IDConsultationVisit;references:IDConsultationVisit"`
}

func (DiagnosticImage) TableName() string {
	return "imagenes_diagnosticas"
}
