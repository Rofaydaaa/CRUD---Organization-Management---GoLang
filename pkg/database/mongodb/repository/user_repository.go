package repository

import (
    "context"

    model "organization_management/pkg/database/mongodb/models"
	database "organization_management/pkg/database/mongodb"

    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

var userCollection *mongo.Collection

func init() {
    userCollection = database.GetCollection(database.DB, "users")
}

// InsertUser inserts a new user into the database.
func InsertUser(ctx context.Context, user model.User) (*mongo.InsertOneResult, error) {
    newUser := model.User{
        Id:       primitive.NewObjectID(),
        Name:     user.Name,
        Email:    user.Email,
        Password: user.Password,
    }

    return userCollection.InsertOne(ctx, newUser)
}

// GetUserByEmail retrieves a user by email from the database.
func GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	// Define a filter to find the user by email
	filter := bson.M{"email": email}

	// Query the database to find the user
	var user model.User
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err // User not found or error occurred
	}

	return &user, nil // Return the user if found
}