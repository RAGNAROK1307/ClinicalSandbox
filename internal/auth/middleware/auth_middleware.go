package middleware

import (
	"ClinicalSandBox/configs/db"
	"ClinicalSandBox/internal/API/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var jwtKey = []byte("Simclec")

type Claims struct {
	UserID   uint   `json:"user_id"`
	RoleID   uint   `json:"role_id"`
	UserName string `json:"user_name"`
	RoleName string `json:"role_name"`
	jwt.StandardClaims
}

func GenerateToken(userID uint, roleID uint, userName string, roleName string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		RoleID:   roleID,
		UserName: userName,
		RoleName: roleName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Expira en 24 horas
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		var user models.User
		if err := db.DB.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Establecer los nuevos campos en el contexto
		c.Set("user", user)
		c.Set("user_name", claims.UserName)
		c.Set("role_name", claims.RoleName)

		c.Next()
	}
}

func RoleMiddleware(roleID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		if user.(models.User).IDRole != roleID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
