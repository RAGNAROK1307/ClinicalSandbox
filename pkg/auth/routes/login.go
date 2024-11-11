package routes

import (
	"ClinicalSandBox/pkg/auth/services"
	"encoding/json"
	"net/http"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if services.ValidateCredentials(creds.Username, creds.Password) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
