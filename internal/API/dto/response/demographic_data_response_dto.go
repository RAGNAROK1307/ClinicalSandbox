package response

type DemographicDataResponseDTO struct {
	IDDemographicData       uint   `json:"id_datos_demograficos"`
	RacialEthnicInformation string `json:"informacion_etnica_racial"`
	MaritalStatus           string `json:"estado_civil"`
	Nationality             string `json:"nacionalidad"`
	EmploymentInformation   string `json:"informacion_laboral"`
}
