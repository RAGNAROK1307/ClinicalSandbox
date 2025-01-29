package services

func ValidateCredentials(username, password string) bool {
	// Aquí va la lógica de autenticación real
	return username == "user" && password == "pass"
}
