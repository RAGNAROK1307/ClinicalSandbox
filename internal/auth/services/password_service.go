package services

import (
	"golang.org/x/crypto/bcrypt"
)

/*
Este módulo del paquete `services` proporciona funciones auxiliares para el manejo seguro de contraseñas mediante el algoritmo bcrypt.

Función HashPassword:
- Recibe una contraseña en texto plano y devuelve su representación en hash usando bcrypt con un costo por defecto (recomendado por la librería).
- Si ocurre un error durante la generación del hash, este se retorna junto con una cadena vacía.

Función CheckPassword:
- Compara una contraseña en texto plano con un hash previamente generado.
- Devuelve `true` si la contraseña coincide con el hash, o `false` si no coincide o si ocurre un error interno.

Estas funciones son esenciales para garantizar la seguridad de las credenciales almacenadas y validar inicios de sesión sin exponer contraseñas reales.
*/

// HashPassword genera un hash de la contraseña usando bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword compara una contraseña con su hash
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
