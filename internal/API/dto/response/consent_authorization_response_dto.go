package response

type ConsentAuthorizationResponseDTO struct {
	IDConsentAuthorization uint   `json:"id_consentimiento"`
	ConsentType            string `json:"tipo_consentimiento"`
	ConsentDate            string `json:"fecha_consentimiento" example:"2025-01-20"`
	Details                string `json:"detalles"`
}
