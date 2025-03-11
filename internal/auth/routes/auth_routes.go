package routes

// AuthRoutes configura las rutas de autenticación
/*func AuthRoutes(router *gin.Engine, authService services.AuthService) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", func(c *gin.Context) {
			var loginRequest struct {
				Username string `json:"nombre_usuario" binding:"required"`
				Password string `json:"contraseña" binding:"required"`
			}

			if err := c.ShouldBindJSON(&loginRequest); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de inicio de sesión inválidos"})
				return
			}

			token, err := authService.Login(loginRequest.Username, loginRequest.Password)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"token": token})
		})

		authGroup.GET("/validate", func(c *gin.Context) {
			tokenString := c.GetHeader("Authorization")
			if tokenString == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token no proporcionado"})
				return
			}

			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id_usuarios": claims.IDUser,
				"id_rol":      claims.IDRole,
			})
		})
	}
}*/
