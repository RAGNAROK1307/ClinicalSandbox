package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/dto/request"
	"ClinicalSandBox/internal/API/models"
	"github.com/gin-gonic/gin"
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

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// GetUsers godoc
// @Summary List all users
// @Description Retrieve a list of all users in the system
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(c *gin.Context) {
	var users []models.User
	db.DB.Preload("Role").Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
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
	c.JSON(http.StatusOK, gin.H{"user": user})
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

	c.JSON(http.StatusOK, gin.H{"user": existingUser})
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
