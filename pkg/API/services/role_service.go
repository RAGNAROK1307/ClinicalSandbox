package services

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/pkg/API/dto"
	"ClinicalSandBox/pkg/API/models" // Importa tu paquete de modelos donde tienes la estructura Role
	"github.com/gin-gonic/gin"       // Para el framework Gin
	"net/http"
	//"gorm.io/gorm"                 // Para trabajar con GORM, si no lo has hecho ya
	_ "net/http" // Para constantes HTTP como http.StatusNotFound, etc.
)

// CreateRole godoc
// @Summary Create a new role
// @Description Adds a new role to the system
// @Tags roles
// @Accept json
// @Produce json
// @Param role body dto.CreateRoleDTO true "Role data"
// @Failure 400 {object} map[string]string
// @Success 201 {object} models.Role
// @Router /roles [post]
func CreateRole(c *gin.Context) {
	var roleDTO dto.CreateRoleDTO

	// Bind JSON to roleDTO
	if err := c.ShouldBindJSON(&roleDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Mapea a modelo
	role := models.Role{
		NombreRol:   roleDTO.NombreRol,
		Descripcion: roleDTO.Descripcion,
	}

	// Guarda el nuevo rol en la base de datos
	if err := db.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"role": role})
}

// GetRoles godoc
// @Summary List all roles
// @Description Retrieve a list of all roles in the system
// @Tags roles
// @Produce json
// @Success 200 {array} models.Role
// @Router /roles [get]
func GetRoles(c *gin.Context) {
	var roles []models.Role
	db.DB.Find(&roles)
	c.JSON(200, gin.H{"roles": roles})
}

// GetRole godoc
// @Summary Get a role by ID
// @Description Retrieve a single role by its ID
// @Tags roles
// @Param id path string true "Role ID"
// @Produce json
// @Success 200 {object} models.Role
// @Failure 404 {object} map[string]string
// @Router /roles/{id} [get]
func GetRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role
	if err := db.DB.First(&role, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Role not found"})
		return
	}
	c.JSON(200, gin.H{"role": role})
}

// UpdateRole godoc
// @Summary Update a role
// @Description Update the information of an existing role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param role body dto.CreateRoleDTO true "Updated role data"
// @Success 200 {object} models.Role
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /roles/{id} [put]
func UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var existingRole models.Role

	// Verificar que el rol existe en la base de datos
	if err := db.DB.First(&existingRole, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// Bind JSON al DTO para validar los datos
	var roleDTO dto.CreateRoleDTO
	if err := c.ShouldBindJSON(&roleDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Actualizar solo los campos permitidos
	existingRole.NombreRol = roleDTO.NombreRol
	existingRole.Descripcion = roleDTO.Descripcion

	// Guardar los cambios en la base de datos
	if err := db.DB.Save(&existingRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": existingRole})
}

// DeleteRole godoc
// @Summary Delete a role
// @Description Remove a role by its ID
// @Tags roles
// @Param id path string true "Role ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /roles/{id} [delete]
func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Role{}, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Role not found"})
		return
	}
	c.JSON(204, nil)
}
