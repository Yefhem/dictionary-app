package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Word struct {
	ID           primitive.ObjectID `bson:"_id"`
	CreatedAt    time.Time          `bson:"created_at"`
	Slug         string             `bson:"slug"`
	English      string             `bson:"english"`
	Turkish      string             `bson:"turkish"`
	Abbreviation string             `bson:"abbreviation"`
	Description  string             `bson:"description"`
}
