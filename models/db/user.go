package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/kamva/mgm/v3"
)

const (
	RoleUser        = "user"
	RoleAdmin       = "admin"
	InitialLevel    = 1
	InitialCoin     = 1000
	InitialProgress = 0
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string `json:"email" bson:"email"`
	Password         string `json:"-" bson:"password"`
	Name             string `json:"name" bson:"name"`
	Role             string `json:"role" bson:"role"`
	MailVerified     bool   `json:"mail_verified" bson:"mail_verified"`
	Level            int    `json:"level" bson:"level"`
	Coin             int    `json:"coin" bson:"coin"`
	Progress         int    `json:"progress" bson:"progress"`
	Country          string `json:"country" bson:"country"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Type  string `json:"type"`
}

func NewUser(email string, password string, name string, role string, country string, InitialLevel int, InitialCoin int, InitialProgress int) *User {
	return &User{
		Email:        email,
		Password:     password,
		Name:         name,
		Role:         role,
		Level:        InitialLevel,
		Coin:         InitialCoin,
		Progress:     InitialProgress,
		MailVerified: false,
		Country:      country,
	}
}

func (model *User) CollectionName() string {
	return "users"
}
