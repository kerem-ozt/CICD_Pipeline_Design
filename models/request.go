package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var passwordRule = []validation.Rule{
	validation.Required,
	validation.Length(8, 32),
	validation.Match(regexp.MustCompile(`^\S+$`)).Error("cannot contain whitespaces"),
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Country  string `json:"country"`
}

func (a RegisterRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required, validation.Length(3, 64)),
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, passwordRule...),
		validation.Field(&a.Country, validation.Required),
	)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a LoginRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, passwordRule...),
	)
}

type RefreshRequest struct {
	Token string `json:"token"`
}

func (a RefreshRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(
			&a.Token,
			validation.Required,
			// validation.Match(regexp.MustCompile("^\\S+$")).Error("cannot contain whitespaces"),
		),
	)
}

type TournamentRequest struct {
	MinLevels       string               `json:"MinLevels"`
	EntryFee        string               `json:"EntryFee"`
	MaxParticipants string               `json:"MaxParticipants"`
	Participants    []primitive.ObjectID `json:"participants" binding:"required"`
}

func (a TournamentRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.MinLevels, validation.Required),
		validation.Field(&a.EntryFee, validation.Required),
		validation.Field(&a.MaxParticipants, validation.Required),
		validation.Field(&a.Participants),
	)
}

type ProgressRequest struct {
	UserID primitive.ObjectID `json:"UserID"`
	Score  int                `json:"Score"`
	Coin   int                `json:"Coin"`
}

func (a ProgressRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.UserID, validation.Required),
		validation.Field(&a.Score),
		validation.Field(&a.Coin),
	)
}
