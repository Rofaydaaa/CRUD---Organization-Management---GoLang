package controller

import (
    "context"
    "net/http"
    "time"

    model "organization_management/pkg/database/mongodb/models"
    repository "organization_management/pkg/database/mongodb/repository"

    "github.com/gin-gonic/gin"
)

func CreateOrganization() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        var org model.Organization

        // Validate the request body
        if err := c.BindJSON(&org); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Use the validator library to validate required fields
        if validationErr := validate.Struct(&org); validationErr != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
            return
        }

        // Insert the organization into the database
        _, organizationId, err := repository.InsertOrganization(ctx, org)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, gin.H{"organization_id": organizationId})
    }
}

func ReadOrganization() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        // Extract organization ID from the request path parameters
        orgID := c.Param("organization_id")

        // Retrieve the organization from the database by its ID
        org, err := repository.GetOrganizationByID(ctx, orgID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organization"})
            return
        }

		// Check if organization members list is nil, and replace with an empty list if so
        var members []model.OrganizationMember
        if org.OrganizationMembers != nil {
            members = org.OrganizationMembers
        } else {
            members = []model.OrganizationMember{}
        }

        // Return the organization details in the specified JSON format
        c.JSON(http.StatusOK, gin.H{
            "organization_id":        orgID,
            "name":                   org.Name,
            "description":            org.Description,
            "organization_members":   members,
        })
    }
}

func ReadAllOrganizations() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Retrieve all organizations from the database
		orgs, err := repository.GetAllOrganizations(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organizations"})
			return
		}

		// Prepare the response JSON array
		var orgList []gin.H
		for _, org := range orgs {
			// Check if organization members list is nil, and replace with an empty list if so
			var members []model.OrganizationMember
			if org.OrganizationMembers != nil {
				members = org.OrganizationMembers
			} else {
				members = []model.OrganizationMember{}
			}

			// Append organization details to the response array
			orgList = append(orgList, gin.H{
				"organization_id":      org.OrganizationId,
				"name":                 org.Name,
				"description":          org.Description,
				"organization_members": members,
			})
		}

		// Return the response JSON array
		c.JSON(http.StatusOK, orgList)
	}
}

func UpdateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Extract organization ID from the request path parameters
		orgID := c.Param("organization_id")

		// Retrieve the organization from the database by its ID
		org, err := repository.GetOrganizationByID(ctx, orgID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organization"})
			return
		}

		// Bind the request body to a struct
		var updateData model.Organization
		if err := c.BindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update organization fields if they're not empty
		if updateData.Name != "" {
			org.Name = updateData.Name
		}
		if updateData.Description != "" {
			org.Description = updateData.Description
		}

		// Update the organization in the database
		err = repository.UpdateOrganization(ctx, org)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
			return
		}

		// Return the updated organization details
		c.JSON(http.StatusOK, gin.H{
			"organization_id": orgID,
			"name":            org.Name,
			"description":     org.Description,
		})
	}
}

// DeleteOrganization deletes an organization by its ID.
func DeleteOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Extract organization ID from the request path parameters
		orgID := c.Param("organization_id")

		// Delete the organization from the database
		err := repository.DeleteOrganization(ctx, orgID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return success message
		c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
	}
}

// InviteUserToOrganization invites a user to join an organization.
func InviteUserToOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Extract organization ID from the request path parameters
		orgID := c.Param("organization_id")

		// Retrieve the organization from the database by its ID
		org, err := repository.GetOrganizationByID(ctx, orgID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve organization"})
			return
		}

		// Check if the organization exists
		if org == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Organization not found"})
			return
		}

		// Bind the request body to a struct
		var inviteData struct {
			UserEmail string `json:"user_email" binding:"required"`
		}
		if err := c.BindJSON(&inviteData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the user email is already in the organization
		for _, member := range org.OrganizationMembers {
			if member.UserEmail == inviteData.UserEmail {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User is already a member of the organization"})
				return
			}
		}

		// Retrieve the user by email
		user, err := repository.GetUserByEmail(ctx, inviteData.UserEmail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		// Check if the user exists
		if user == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		// Add the user to the organization members
		orgMember := model.OrganizationMember{
			Name:        user.Name,
			UserEmail:       inviteData.UserEmail,
			AccessLevel: "member",
		}
		org.OrganizationMembers = append(org.OrganizationMembers, orgMember)

		// Update the organization in the database
		if err := repository.UpdateOrganization(ctx, org); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
			return
		}

		// Return success message
		c.JSON(http.StatusOK, gin.H{"message": "User invited successfully"})
	}
}