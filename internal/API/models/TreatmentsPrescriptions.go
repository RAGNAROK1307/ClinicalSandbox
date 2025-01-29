package models

type TreatmentPrescription struct {
	IDTreatmentPrescription uint    `json:"id_tratamiento" gorm:"primaryKey;column:id_tratamiento" swaggerignore:"true"`
	IDPatient               uint    `json:"id_paciente" gorm:"column:id_paciente"`
	PrescribedMedication    string  `json:"medicamento_prescrito" gorm:"type:varchar(50);column:medicamento_prescrito"`
	Amount                  string  `json:"dosis" gorm:"type:varchar(50);column:dosis"`
	Frequency               string  `json:"frecuencia" gorm:"type:varchar(50);column:frecuencia"`
	Instructions            string  `json:"instrucciones" gorm:"type:text;column:instrucciones"`
	DurationTreatment       string  `json:"duracion_tratamiento" gorm:"type:text;column:duracion_tratamiento"`
	Patient                 Patient `gorm:"foreignKey:IDPatient ;references:IDPatient "`
}

func (TreatmentPrescription) TableName() string {
	return "tratamientos_prescripciones"
}
