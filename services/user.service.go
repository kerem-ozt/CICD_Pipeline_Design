package services

import (
	"errors"

	"github.com/kamva/mgm/v3"
	db "github.com/kerem-ozt/GoodBlast_API/models/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser create a user record
func CreateUser(name string, email string, plainPassword string, country string) (*db.User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("cannot generate hashed password")
	}

	user := db.NewUser(email, string(password), name, db.RoleUser, country, db.InitialLevel, db.InitialCoin, db.InitialProgress)
	err = mgm.Coll(user).Create(user)
	if err != nil {
		return nil, errors.New("cannot create new user")
	}

	return user, nil
}

// GetAllUsers get all users
func GetAllUsers() ([]*db.User, error) {
	users := []*db.User{}
	err := mgm.Coll(&db.User{}).SimpleFind(&users, bson.M{})
	if err != nil {
		return nil, errors.New("cannot find users")
	}

	return users, nil
}

// DeleteUser delete user by id
func DeleteUser(userId primitive.ObjectID) error {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(userId, user)
	if err != nil {
		return errors.New("cannot find user")
	}

	err = mgm.Coll(user).Delete(user)
	if err != nil {
		return errors.New("cannot delete user")
	}

	return nil
}

// FindUserById find user by id
func FindUserById(userId primitive.ObjectID) (*db.User, error) {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(userId, user)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return user, nil
}

// FindUserByEmail find user by email
func FindUserByEmail(email string) (*db.User, error) {
	user := &db.User{}
	err := mgm.Coll(user).First(bson.M{"email": email}, user)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return user, nil
}

// CheckUserMail search user by email, return error if someone uses
func CheckUserMail(email string) error {
	user := &db.User{}
	userCollection := mgm.Coll(user)
	err := userCollection.First(bson.M{"email": email}, user)
	if err == nil {
		return errors.New("email is already in use")
	}

	return nil
}

// UpdateProgress update user progress
func UpdateProgress(userId primitive.ObjectID, score int, coin int) error {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(userId, user)
	if err != nil {
		return errors.New("cannot find user")
	}

	newLevel := ((user.Progress + score) / 1000) + 1
	coinsEarned := (newLevel - user.Level) * 100

	user.Progress = (user.Progress + score) % 1000
	user.Level = newLevel

	user.Coin = user.Coin + coin + coinsEarned

	err = mgm.Coll(user).Update(user)
	if err != nil {
		return errors.New("cannot update user")
	}

	return nil
}

// EnterTournament enter user to tournament
func EnterTournament(userID primitive.ObjectID, tournamentID primitive.ObjectID) error {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(userID, user)
	if err != nil {
		return errors.New("cannot find user")
	}

	tournament := &db.Tournament{}
	err = mgm.Coll(tournament).FindByID(tournamentID, tournament)
	if err != nil {
		return errors.New("cannot find tournament")
	}

	if user.Coin < db.EntryFee {
		return errors.New("user does not have enough coin")
	}

	user.Coin -= db.EntryFee
	err = mgm.Coll(user).Update(user)
	if err != nil {
		return errors.New("cannot update user")
	}

	if user.Level < db.MinLevels {
		return errors.New("user level is not enough")
	}

	// Create a new group if the tournament has no groups
	if len(tournament.Groups) == 0 {
		newGroup := db.TournamentGroup{
			GroupID:      primitive.NewObjectID(),
			Participants: []primitive.ObjectID{userID},
		}
		tournament.Groups = append(tournament.Groups, newGroup)
	} else {
		// Find the first group in the tournament
		targetGroup := &tournament.Groups[0]

		// Check if the group is full
		if len(targetGroup.Participants) >= db.MaxParticipants {
			newGroup := db.TournamentGroup{
				GroupID:      primitive.NewObjectID(),
				Participants: []primitive.ObjectID{userID},
			}
			tournament.Groups = append(tournament.Groups, newGroup)
		}

		// Check if the user is already in the group
		for _, participantID := range targetGroup.Participants {
			if participantID == userID {
				return errors.New("user already in the tournament group")
			}
		}

		// Add the user to the group
		targetGroup.Participants = append(targetGroup.Participants, userID)
	}

	// Update the tournament with the modified group
	err = mgm.Coll(tournament).Update(tournament)
	if err != nil {
		return errors.New("cannot update tournament")
	}

	return nil
}
