package types

import (
	"go.mongodb.org/mongo-driver/bson"
)

type RequestData struct {
	DataSource string   `json:"dataSource"`
	Database   string   `json:"database"`
	Collection string   `json:"collection"`
	Filter     bson.D   `json:"filter"`
	Pipeline   []bson.D `json:"pipeline"`
	Document   bson.D   `json:"document"`
	Update     bson.D   `json:"update"`
	Upsert     bool     `json:"upsert"`
}
