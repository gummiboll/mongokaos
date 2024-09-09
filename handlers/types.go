package handlers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SingleResult struct {
	Document primitive.M `json:"document,omitempty"`
}

type MultipleResults struct {
	Documents []primitive.M `json:"documents"`
}
