package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	LeaderboardTypeGlobal  = "global"
	LeaderboardTypeCountry = "country"
)

type Leaderboard struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Type  string             `bson:"type"`
	Users []LeaderboardUser  `bson:"users"`
}

func (l *Leaderboard) GetID() interface{} {
	return l.ID
}

func (l *Leaderboard) PrepareID(id interface{}) (interface{}, error) {
	return primitive.ObjectIDFromHex(id.(string))
}

func (l *Leaderboard) SetID(id interface{}) {
	l.ID = id.(primitive.ObjectID)
}

type LeaderboardUser struct {
	UserID   primitive.ObjectID `bson:"userId"`
	Progress int                `bson:"progress"`
	Country  string             `bson:"country"`
}
