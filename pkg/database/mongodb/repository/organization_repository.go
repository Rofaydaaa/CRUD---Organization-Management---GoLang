package repository

import (
    "context"
	"errors"

    model "organization_management/pkg/database/mongodb/models"
    database "organization_management/pkg/database/mongodb"

    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var orgCollection *mongo.Collection

func init() {
    orgCollection = database.GetCollection(database.DB, "organizations")
}

// InsertOrganization inserts a new organization into the database.
func InsertOrganization(ctx context.Context, org model.Organization) (*mongo.InsertOneResult, string, error) {
	// Generate a UUID for the organization ID
	organizationID := uuid.New().String()

    newOrg := model.Organization{
		OrganizationId: organizationID,
        Name:         org.Name,
        Description:  org.Description,
        OrganizationMembers:      org.OrganizationMembers,
    }

    // Insert the new organization into the database
	result, err := orgCollection.InsertOne(ctx, newOrg)
	if err != nil {
		return nil, "", err
	}

	return result, organizationID, nil
}

// GetOrganizationByID retrieves an organization by its ID from the database.
func GetOrganizationByID(ctx context.Context, id string) (*model.Organization, error) {
    // Define a filter to find the organization by ID
    filter := bson.M{"organizationid": id}

    // Query the database to find the organization
    var org model.Organization
    err := orgCollection.FindOne(ctx, filter).Decode(&org)
    if err != nil {
        return nil, err // Organization not found or error occurred
    }

    return &org, nil // Return the organization if found
}

// GetAllOrganizations retrieves all organizations from the database.
func GetAllOrganizations(ctx context.Context) ([]model.Organization, error) {
	// Initialize an empty filter to retrieve all documents
	filter := bson.M{}

	// Retrieve all organizations from the database
	cursor, err := orgCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode all documents into a slice of structs
	var orgs []model.Organization
	if err := cursor.All(ctx, &orgs); err != nil {
		return nil, err
	}

	return orgs, nil
}

func UpdateOrganization(ctx context.Context, org *model.Organization) error {
	// Check if organization exists
	_, err := GetOrganizationByID(ctx, org.OrganizationId)
	if err != nil {
		return errors.New("Organization not found")
	}

	// Define filter to find organization by ID
	filter := bson.M{"organizationid": org.OrganizationId}

	// Define update data
	update := primitive.M{
		"$set": primitive.M{
			"name":        org.Name,
			"description": org.Description,
			"organizationmembers": org.OrganizationMembers,
		},
	}

	// Perform update operation
	_, err = orgCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// DeleteOrganization deletes an organization from the database by its ID.
func DeleteOrganization(ctx context.Context, orgID string) error {
	// Check if organization exists
	_, err := GetOrganizationByID(ctx, orgID)
	if err != nil {
		return errors.New("Organization not found")
	}

	// Define filter to find organization by ID
	filter := bson.M{"organizationid": orgID}

	// Delete the organization from the database
	_, err = orgCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

