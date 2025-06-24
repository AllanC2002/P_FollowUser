package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/AllanC2002/P_FollowUser/connection"
	"github.com/AllanC2002/P_FollowUser/functions"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var SECRET_KEY string

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	SECRET_KEY = os.Getenv("SECRET_KEY")
	if SECRET_KEY == "" {
		log.Fatal("SECRET_KEY not set in environment")
	}

	db, err := connection.ConnectDB()
	if err != nil {
		log.Fatal("Could not connect to DB:", err)
	}

	r := gin.Default()

	r.POST("/follow", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing or invalid"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
			return
		}
		idFollower := int(userIDFloat)

		var json struct {
			IdFollowing int `json:"id_following"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		message, code, err := functions.FollowUser(db, idFollower, json.IdFollowing)
		if err != nil {
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(code, gin.H{"message": message})
	})

	r.Run(":8080")
}
