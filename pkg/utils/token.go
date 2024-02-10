package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"errors"

	"github.com/gin-gonic/gin"
	jwt "github.com/dgrijalva/jwt-go"
)

// GenerateToken generates a JWT access and refresh token for the given user ID and email.
func GenerateToken(userID uint, email string) (string, string, error) {
    tokenLifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
    if err != nil {
        return "", "", err
    }

    // Create JWT claims for access token
    accessClaims := jwt.MapClaims{}
    accessClaims["authorized"] = true
    accessClaims["user_id"] = userID
    accessClaims["email"] = email
    accessClaims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()

    // Create access token with claims and sign with secret
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("API_SECRET")))
    if err != nil {
        return "", "", err
    }

    // Create JWT claims for refresh token
    refreshClaims := jwt.MapClaims{}
    refreshClaims["authorized"] = true
    refreshClaims["user_id"] = userID
    refreshClaims["email"] = email
    refreshClaims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix() // Set refresh token expiration to 30 days

    // Create refresh token with claims and sign with secret
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("API_SECRET")))
    if err != nil {
        return "", "", err
    }

    return accessTokenString, refreshTokenString, nil
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

// ExtractUserEmail extracts the user email from the JWT token claims.
func ExtractUserEmail(c *gin.Context) (string, error) {
    tokenString := ExtractToken(c)
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("API_SECRET")), nil
    })
    if err != nil {
        return "", err
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return "", errors.New("Invalid token claims")
    }
    email, ok := claims["email"].(string)
    if !ok {
        return "", errors.New("Email claim not found")
    }
    return email, nil
}
