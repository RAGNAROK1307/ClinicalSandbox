package response

type ConsentAuthorizationResponseDTO struct {
	IDConsentAuthorization uint   `json:"id_consentimiento"`
	IDPatient              uint   `json:"id_paciente"`
	ConsentType            string `json:"tipo_consentimiento"`
	ConsentDate            string `json:"fecha_consentimiento" example:"2025-01-20"`
	Details                string `json:"detalles"`
	FullName               string `json:"nombre_completo"`
}
