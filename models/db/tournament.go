package models

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var StartTime = time.Now().Truncate(24 * time.Hour)
var EndTime = StartTime.Add(24*time.Hour - time.Minute)

const (
	MaxParticipants = 35
	MinLevels       = 10
	EntryFee        = 500
)

type Tournament struct {
	mgm.DefaultModel `bson:",inline"`
	StartTime        time.Time            `bson:"startTime"`
	EndTime          time.Time            `bson:"endTime"`
	MinLevels        int                  `bson:"minLevels"`
	EntryFee         int                  `bson:"entryFee"`
	MaxParticipants  int                  `bson:"maxParticipants"`
	Participants     []primitive.ObjectID `json:"participants" binding:"required"`
	Scores           []TournamentScore    `bson:"scores"`
	Groups           []TournamentGroup    `bson:"groups"`
}

type TournamentScore struct {
	UserID primitive.ObjectID `bson:"userId"`
	Score  int                `bson:"score"`
}

type TournamentGroup struct {
	GroupID      primitive.ObjectID   `bson:"groupId"`
	Participants []primitive.ObjectID `json:"participants" binding:"required"`
}

func NewTournament(participants []primitive.ObjectID) *Tournament {
	return &Tournament{
		StartTime:       StartTime,
		EndTime:         EndTime,
		MinLevels:       10,
		EntryFee:        500,
		MaxParticipants: 35,
		Participants:    participants,
		Scores:          []TournamentScore{},
		Groups:          []TournamentGroup{},
	}
}

func (model *Tournament) CollectionName() string {
	return "tournaments"
}
