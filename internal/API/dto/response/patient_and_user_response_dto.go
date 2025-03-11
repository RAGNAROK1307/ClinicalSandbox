package response

type UserAndPatientResponseDTO struct {
	User            UserResponseDTO            `json:"user"`
	Patient         PatientResponseDTO         `json:"patient"`
	DemographicData DemographicDataResponseDTO `json:"demographic_data"`
	Identification  IdentificationResponseDTO  `json:"identification"`
}
