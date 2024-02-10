package controller

import (
    "context"
    "net/http"
    "time"
    "os"
    "fmt"

    model "organization_management/pkg/database/mongodb/models"
    repository "organization_management/pkg/database/mongodb/repository"
	util "organization_management/pkg/utils"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/dgrijalva/jwt-go"

)

var validate = validator.New()

func RegisterUser() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        var user model.User

        // Validate the request body
        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
            return
        }

        // Use the validator library to validate required fields
        if validationErr := validate.Struct(&user); validationErr != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": validationErr.Error(),
            })
            return
        }

        // Validate the email address
        if !model.ValidateEmail(user.Email) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
            return
        }

        // Check if the user already exists
        existingUser, err := repository.GetUserByEmail(ctx, user.Email)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
            return
        }
        if existingUser != nil {
            c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
            return
        }

        // Hash the user's password before saving it
        if err := user.HashPassword(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": err.Error(),
            })
            return
        }

        // Insert the user into the database
        result, err := repository.InsertUser(ctx, user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": err.Error(),
            })
            return
        }

        c.JSON(http.StatusCreated, gin.H{
            "status":  http.StatusCreated,
            "message": "success",
            "data":    result,
        })
    }
}

// LoginInput represents the input data for login request
type LoginInput struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginUser handles the login request
func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        var input LoginInput

        // Bind JSON request body to input struct
        if err := c.BindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
            return
        }

        // Validate input
        if validationErr := validate.Struct(input); validationErr != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": validationErr.Error(),
            })
            return
        }

        // Check if user exists
        user, err := repository.GetUserByEmail(ctx, input.Email)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "User with this email does not exist",
            })
            return
        }

        // Verify password
        if err := user.VerifyPassword(input.Password, user.Password); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Password is incorrect",
            })
            return
        }

        // Convert ObjectID to timestamp
        timestamp := user.Id.Timestamp()
        // Get the seconds part of the timestamp
        seconds := timestamp.Unix()
        // Convert seconds to uint
        userID := uint(seconds)

        token, refreshToken, err := util.GenerateToken(userID, user.Email)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Error generating tokens",
            })
            return
        }

        // Respond with tokens and message
        c.JSON(http.StatusOK, gin.H{
            "access_token":  token,
            "refresh_token": refreshToken,
            "message":       "Authentication successful",
        })
    }
}

// RefreshToken handles the refresh token request
func RefreshToken() gin.HandlerFunc {
    return func(c *gin.Context) {
        var input struct {
            RefreshToken string `json:"refresh_token" binding:"required"`
        }
        if err := c.BindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Parse and validate the refresh token
        token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(os.Getenv("API_SECRET")), nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
            return
        }

        // Extract user ID and email from the refresh token claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token claims"})
            return
        }
        userID, ok := claims["user_id"].(float64)
        if !ok {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in refresh token"})
            return
        }
        email, ok := claims["email"].(string)
        if !ok {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email in refresh token"})
            return
        }

        // Generate a new access token
        accessToken, refreshToken, err := util.GenerateToken(uint(userID), email)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
            return
        }

        // Respond with new access token and message
        c.JSON(http.StatusOK, gin.H{
            "access_token": accessToken,
            "refresh_token": refreshToken,
            "message":      "Access token refreshed successfully",
        })
    }
}