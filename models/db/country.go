package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Country struct {
	ID    primitive.ObjectID   `bson:"_id,omitempty"`
	Name  string               `bson:"name"`
	Users []primitive.ObjectID `bson:"users"`
}
