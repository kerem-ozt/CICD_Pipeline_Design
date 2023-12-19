package services

import (
	"errors"
	"fmt"

	"github.com/kamva/mgm/v3"
	db "github.com/kerem-ozt/GoodBlast_API/models/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddUserToLeaderboard adds a user to the leaderboard
func AddUserToLeaderboard(userID primitive.ObjectID, progress int, leaderboardType string) error {
	// Check if the user is already in the leaderboard
	exists, err := isUserInLeaderboard(userID, leaderboardType)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user is already in the leaderboard")
	}

	// Create a new leaderboard user
	leaderboardUser := db.LeaderboardUser{
		UserID:   userID,
		Progress: progress,
	}

	// Get the corresponding leaderboard document or create a new one
	leaderboard, err := getOrCreateLeaderboard(leaderboardType)
	if err != nil {
		return err
	}

	// Add the user to the leaderboard
	leaderboard.Users = append(leaderboard.Users, leaderboardUser)

	// Save the updated leaderboard
	err = mgm.Coll(leaderboard).Update(leaderboard)
	if err != nil {
		return errors.New("failed to update leaderboard")
	}

	return nil
}

// GetGlobalLeaderboard returns the global leaderboard
func GetGlobalLeaderboard(leaderboardType string) (*db.Leaderboard, error) {
	leaderboard, err := getOrCreateLeaderboard(leaderboardType)
	if err != nil {
		return nil, err
	}

	return leaderboard, nil
}

// GetLeaderboard returns the leaderboard for a given type
func GetLeaderboardByCountry(leaderboardType string, country string) (*db.Leaderboard, error) {
	leaderboard, err := getOrCreateLeaderboard(leaderboardType)
	if err != nil {
		return nil, err
	}

	// Filter the leaderboard by country
	filteredLeaderboard := &db.Leaderboard{
		Type:  leaderboardType,
		Users: []db.LeaderboardUser{},
	}
	for _, user := range leaderboard.Users {
		if user.Country == country {
			filteredLeaderboard.Users = append(filteredLeaderboard.Users, user)
		}
	}

	return filteredLeaderboard, nil
}

// isUserInLeaderboard checks if a user is already in the leaderboard
func isUserInLeaderboard(userID primitive.ObjectID, leaderboardType string) (bool, error) {
	leaderboard, err := getOrCreateLeaderboard(leaderboardType)
	if err != nil {
		return false, err
	}

	// Check if the user exists in the leaderboard
	for _, user := range leaderboard.Users {
		if user.UserID == userID {
			return true, nil
		}
	}

	return false, nil
}

// getOrCreateLeaderboard retrieves an existing leaderboard or creates a new one
func getOrCreateLeaderboard(leaderboardType string) (*db.Leaderboard, error) {
	Leaderboard := &db.Leaderboard{}

	query := bson.M{"type": leaderboardType}
	err := mgm.Coll(&db.Leaderboard{}).First(query, Leaderboard)
	if err == nil {
		// Leaderboard already exists
		return Leaderboard, nil
	}

	// If not found, create a new leaderboard
	Leaderboard = &db.Leaderboard{
		Type:  leaderboardType,
		Users: []db.LeaderboardUser{},
	}

	err = mgm.Coll(Leaderboard).Create(Leaderboard)
	if err != nil {
		return nil, errors.New("failed to create leaderboard")
	}

	return Leaderboard, nil
}

// EnsureLeaderboardInitialized initializes the leaderboard with all users
func EnsureLeaderboardInitialized(leaderboardType string) (*db.Leaderboard, error) {

	leaderboard, err := getOrCreateLeaderboard(leaderboardType)
	if err != nil {
		return nil, err
	}

	// Get all users from the database
	users, err := getAllUsers()
	if err != nil {
		return nil, err
	}

	// Populate the leaderboard with all users
	for _, user := range users {
		// Check if the user is already in the leaderboard
		exists := false
		for _, leaderboardUser := range leaderboard.Users {
			if leaderboardUser.UserID == user.ID {
				exists = true
				break
			}
		}

		if !exists {
			// Add the user to the leaderboard
			leaderboardUser := db.LeaderboardUser{
				UserID:   user.ID,
				Progress: user.Progress, // You may customize this based on your requirements
				Country:  user.Country,
			}
			leaderboard.Users = append(leaderboard.Users, leaderboardUser)
		}
	}

	// Save the updated leaderboard
	err = mgm.Coll(leaderboard).Update(leaderboard)
	if err != nil {
		return nil, errors.New("failed to update leaderboard")
	}

	return leaderboard, nil
}

// getAllUsers retrieves all users from the database
func getAllUsers() ([]db.User, error) {
	var users []db.User
	// err := mgm.Coll(&db.User{}).SimpleFind(&users, nil)
	err := mgm.Coll(&db.User{}).SimpleFind(&users, bson.M{})
	if err != nil {
		// Add error logging here
		fmt.Println("Error getting users:", err)
		return nil, errors.New("failed to get all users")
	}

	return users, nil
}
