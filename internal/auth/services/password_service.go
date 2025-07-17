package services

import (
	"golang.org/x/crypto/bcrypt"
)

/*
Este módulo proporciona utilidades para el manejo seguro de contraseñas en el sistema,
empleando el algoritmo bcrypt para hashing y verificación.

FUNCIONALIDADES:

- HashPassword:
  • Recibe una cadena de texto (contraseña en texto plano).
  • Genera un hash seguro usando bcrypt con el costo por defecto.
  • Retorna la cadena encriptada y un posible error.

- CheckPassword:
  • Recibe un hash y una contraseña en texto plano.
  • Verifica si la contraseña coincide con el hash usando bcrypt.
  • Retorna true si coinciden, false en caso contrario.

Estas funciones permiten:
  - Almacenar contraseñas de forma segura en la base de datos (nunca en texto plano).
  - Validar credenciales durante el inicio de sesión.
  - Cumplir buenas prácticas de ciberseguridad en el manejo de autenticación.

Este módulo es utilizado por el servicio de autenticación (`Login`) para verificar
credenciales, y por otros módulos donde se requiera cambiar o establecer contraseñas.
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
