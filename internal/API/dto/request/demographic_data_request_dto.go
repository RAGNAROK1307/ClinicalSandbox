package request

type CreateDemographicDataDTO struct {
	RacialEthnicInformation string `json:"informacion_etnica_racial" binding:"required"`
	MaritalStatus           string `json:"estado_civil" binding:"required"`
	Nationality             string `json:"nacionalidad" binding:"required"`
	EmploymentInformation   string `json:"informacion_laboral" binding:"required"`
}
