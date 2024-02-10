package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organization struct {
    Id                   primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	OrganizationId		 string               `json:"organization_id,omitempty"`
    Name                 string               `json:"name" validate:"required"`
    Description          string               `json:"description" validate:"required"`
    OrganizationMembers  []OrganizationMember `json:"organization_members"`
}

type OrganizationMember struct {
    Name        string `json:"name" validate:"required"`
    UserEmail       string `json:"email" validate:"required"`
    AccessLevel string `json:"access_level" validate:"required"`
}
