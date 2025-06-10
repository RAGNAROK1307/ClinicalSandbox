package response

type LoginResponse struct {
	Token         string `json:"token"`
	Blocked       bool   `json:"blocked"`
	RemainingMS   int64  `json:"remaining_ms"`
	Message       string `json:"message"`
	ActiveSession bool   `json:"active_session"`
}
