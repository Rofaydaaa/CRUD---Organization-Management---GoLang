package controller

import (
    "context"
    "net/http"
    "time"

    model "organization_management/pkg/database/mongodb/models"
    repository "organization_management/pkg/database/mongodb/repository"
	util "organization_management/pkg/utils"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"

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

        token, err := util.GenerateToken(userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Error generating token",
            })
            return
        }

        // Respond with token and message
        c.JSON(http.StatusOK, gin.H{
            "token":   token,
            "message": "Authentication successful",
        })
    }
}
