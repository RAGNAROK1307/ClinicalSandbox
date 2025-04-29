package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/dto/response"
	"ClinicalSandBox/internal/API/models"
	"ClinicalSandBox/internal/auth/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Adds a new user to the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body request.CreateUserDTO true "User data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var userDTO request.CreateUserDTO

	// Bind JSON to userDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Hashear la contraseña
	hashedPassword, err := services.HashPassword(userDTO.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la contraseña"})
		return
	}

	// Mapea a modelo User con la contraseña hasheada
	user := models.User{
		IDRole:   userDTO.IDRole,
		UserName: userDTO.UserName,
		Password: hashedPassword,
	}

	// Guarda el nuevo usuario en la base de datos
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	var role models.Role
	db.DB.First(&role, user.IDRole)

	// Generar la respuesta con UserResponseDTO
	userResponse := response.UserResponseDTO{
		IDUser:   user.IDUser,
		UserName: user.UserName,
		RoleName: role.RoleName,
	}

	c.JSON(http.StatusCreated, gin.H{"user": userResponse})
}

/*func CreateUser(c *gin.Context) {
	var userDTO request.CreateUserDTO

	// Bind JSON to userDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Mapea a modelo User
	user := models.User{
		IDRole:   userDTO.IDRole,
		UserName: userDTO.UserName,
		Password: userDTO.Password,
	}

	// Guarda el nuevo usuario en la base de datos
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	var role models.Role
	db.DB.First(&role, user.IDRole)

	// Generar la respuesta con UserResponseDTO
	userResponse := response.UserResponseDTO{
		IDUser:   user.IDUser,
		UserName: user.UserName,
		RoleName: role.RoleName,
	}

	c.JSON(http.StatusCreated, gin.H{"user": userResponse})
	//c.JSON(http.StatusCreated, gin.H{"user": user})
}*/

// GetUsers godoc
// @Summary List all users
// @Description Retrieve a list of all users in the system
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := db.DB.Preload("Role").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var usersResponse []response.UserResponseDTO
	for _, user := range users {
		usersResponse = append(usersResponse, response.UserResponseDTO{
			IDUser:   user.IDUser,
			UserName: user.UserName,
			RoleName: user.Role.RoleName, // Aquí accedemos correctamente al nombre del rol
			//Password: user.Password,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": usersResponse})
	//c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Retrieve a single user by its ID
// @Tags users
// @Param id path string true "User ID"
// @Produce json
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := db.DB.Preload("Role").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userResponse := response.UserResponseDTO{
		IDUser:   user.IDUser,
		UserName: user.UserName,
		RoleName: user.Role.RoleName, // Asegúrate de que el campo en Role sea `NombreRol`
	}

	c.JSON(http.StatusOK, gin.H{"user": userResponse})
	//c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update the information of an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body request.CreateUserDTO true "Updated user data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var existingUser models.User

	// Verificar que el usuario existe en la base de datos
	if err := db.DB.First(&existingUser, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Bind JSON al DTO para validar los datos
	var userDTO request.CreateUserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Actualizar solo los campos permitidos
	existingUser.IDRole = userDTO.IDRole
	existingUser.UserName = userDTO.UserName
	existingUser.Password = userDTO.Password

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	var role models.Role
	db.DB.First(&role, existingUser.IDRole)

	userResponse := response.UserResponseDTO{
		IDUser:   existingUser.IDUser,
		UserName: existingUser.UserName,
		RoleName: role.RoleName,
	}

	c.JSON(http.StatusOK, gin.H{"user": userResponse})
	//c.JSON(http.StatusOK, gin.H{"user": existingUser})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Remove a user by its ID
// @Tags users
// @Param id path string true "User ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// UpdatePassword godoc
// @Summary Update user password
// @Description Updates user password after verifying current password
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Param password body request.UpdatePasswordDTO true "Password data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id}/password [put]
func UpdatePassword(c *gin.Context) {
	// Obtener ID del usuario objetivo
	userID := c.Param("id")

	// Bind del JSON
	var passwordDTO request.UpdatePasswordDTO
	if err := c.ShouldBindJSON(&passwordDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Buscar usuario
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Verificar contraseña actual
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordDTO.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Contraseña actual incorrecta"})
		return
	}

	// Hashear nueva contraseña
	hashedPassword, err := services.HashPassword(passwordDTO.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar contraseña"})
		return
	}

	// Actualizar contraseña
	user.Password = hashedPassword
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar contraseña"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada correctamente"})
}

// AdminUpdatePassword godoc
// @Summary Admin updates user password
// @Description Admin updates any user password without current password (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Param password body request.AdminUpdatePasswordDTO true "New password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id}/passwords [put]
func AdminUpdatePassword(c *gin.Context) {
	// Obtener ID del usuario objetivo
	userID := c.Param("id")

	// Bind del JSON
	var passwordDTO request.AdminUpdatePasswordDTO
	if err := c.ShouldBindJSON(&passwordDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Buscar usuario
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Hashear nueva contraseña
	hashedPassword, err := services.HashPassword(passwordDTO.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar contraseña"})
		return
	}

	// Actualizar contraseña
	user.Password = hashedPassword
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar contraseña"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Contraseña actualizada correctamente",
		"user_id": userID,
	})
}
