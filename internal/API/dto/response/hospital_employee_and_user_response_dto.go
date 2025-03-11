package response

type HospitalEmployeeAndUserResponseDTO struct {
	User             UserResponseDTO             `json:"user"`
	HospitalEmployee HospitalEmployeeResponseDTO `json:"hospital_employee"`
	Identification   IdentificationResponseDTO   `json:"identification"`
}
