package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/dgrijalva/jwt-go"
)

// GenerateToken generates a JWT token for the given user ID.
func GenerateToken(user_id uint) (string, error) {
	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", err
	}

	// Create JWT claims
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()

	// Create token with claims and sign with secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

// TokenValid checks if the JWT token provided in the request is valid.
func TokenValid(c *gin.Context) error {
	// Extract token string from request
	tokenString := ExtractToken(c)

	// Parse and validate token
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	return err
}

// ExtractToken extracts the JWT token from the request.
func ExtractToken(c *gin.Context) string {
	// Check if token is passed as a query parameter
	token := c.Query("token")
	if token != "" {
		return token
	}

	// Check if token is passed in the Authorization header
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	// No token found
	return ""
}