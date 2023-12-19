package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/kamva/mgm/v3"
	db "github.com/kerem-ozt/GoodBlast_API/models/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateTournament create new tournament record
func CreateTournament(participants ...primitive.ObjectID) (*db.Tournament, error) {

	todayTournament, _ := FindTournamentByStartDateToday()

	if todayTournament != nil {
		return nil, errors.New("tournament already exists for today")
	}

	tournament := db.NewTournament(participants)
	err := mgm.Coll(tournament).Create(tournament)
	if err != nil {
		return nil, errors.New("cannot create new tournament")
	}

	return tournament, nil
}

// GetTournaments get paginated tournaments list
func GetTournaments(page int, limit int) ([]db.Tournament, error) {
	var tournaments []db.Tournament

	findOptions := options.Find().
		SetSkip(int64(page * limit)).
		SetLimit(int64(limit + 1))

	err := mgm.Coll(&db.Tournament{}).SimpleFind(
		&tournaments,
		bson.M{},
		findOptions,
	)

	if err != nil {
		return nil, errors.New("cannot find tournaments")
	}

	return tournaments, nil
}

// GetTournamentById get tournament by id
func GetTournamentById(tournamentId primitive.ObjectID) (*db.Tournament, error) {
	tournament := &db.Tournament{}
	err := mgm.Coll(tournament).FindByID(tournamentId, tournament)
	if err != nil {
		return nil, errors.New("cannot find tournament")
	}

	return tournament, nil
}

// FindTournamentByStartDateToday find tournament by start_date
func FindTournamentByStartDateToday() (*db.Tournament, error) {
	// Get the current date in UTC
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Find the tournament with the matching start_date
	tournament := &db.Tournament{}
	err := mgm.Coll(tournament).First(bson.M{"startTime": startOfDay}, tournament)

	if err != nil {
		return nil, errors.New("cannot find tournament")
	}

	return tournament, nil
}

// CreateTournamentGroups create tournament groups
func CreateTournamentGroups() ([]db.TournamentGroup, error) {
	groups := make([]db.TournamentGroup, 0)

	// Create a new group
	group := db.TournamentGroup{
		GroupID:      primitive.NewObjectID(),
		Participants: []primitive.ObjectID{}, // Empty participants for now
	}

	// Add the new group to the groups slice
	groups = append(groups, group)

	// Find the tournament for today
	todayTournament, err := FindTournamentByStartDateToday()
	if err != nil {
		return nil, err
	}

	// Check if the tournament exists
	if todayTournament == nil {
		return nil, errors.New("no tournament found for today")
	}

	// Add the new groups to the existing groups
	todayTournament.Groups = append(todayTournament.Groups, groups...)

	// Save the updated tournament to the database
	err = mgm.Coll(todayTournament).Update(todayTournament)
	if err != nil {
		return nil, errors.New("cannot update tournament with groups: " + err.Error())
	}

	// Check if groups is not empty before returning
	if len(groups) > 0 {
		return groups, nil
	}

	return nil, errors.New("no groups created")
}

// Define a struct to store the participant ID and rank
type Participant struct {
	ID   primitive.ObjectID `bson:"id"`
	Rank int
}

// EnterTournament enter user to tournament
func ProgressTournament(tournamentID primitive.ObjectID) ([]Participant, error) {
	tournament := &db.Tournament{}

	// Find the tournament by ID
	err := mgm.Coll(tournament).FindByID(tournamentID, tournament)
	if err != nil {
		return nil, errors.New("cannot find tournament")
	}

	var winners []Participant

	// Iterate through groups
	for _, group := range tournament.Groups { //
		fmt.Println("Group:", group)
		var participants []Participant

		// Iterate through participants within the group
		for _, objID := range group.Participants {
			id, err := primitive.ObjectIDFromHex(objID.Hex())
			if err != nil {
				return nil, errors.New("invalid participant ID")
			}
			participants = append(participants, Participant{ID: id, Rank: 0})
		}

		// Progress the tournament within the group
		for round := 1; len(group.Participants) > 1; round++ {
			for i := len(group.Participants) - 1; i > 0; i-- {
				j := rand.Intn(i + 1)
				group.Participants[i], group.Participants[j] = group.Participants[j], group.Participants[i]
			}

			winnerCount := len(group.Participants) / 2

			winnersSlice := group.Participants[:winnerCount]

			group.Participants = group.Participants[:winnerCount]

			// var winners []Participant
			for _, winner := range winnersSlice {
				for j := range participants {
					if participants[j].ID == winner { // Access the ID field of the Participant struct
						participants[j].Rank = round
						winners = append(winners, participants[j])
						break
					}
				}
			}

			fmt.Println("Group Round", round, "Winners:", winners)
		}

		// Update progress for each participant in the group
		for _, participants := range participants {
			err := UpdateProgress(participants.ID, participants.Rank*100, 0)
			if err != nil {
				return nil, errors.New("cannot update user progress")
			}
		}

		// Sort participants by rank
		sort.Slice(participants, func(i, j int) bool {
			return participants[i].Rank > participants[j].Rank
		})

		// top3winnerIDs := []primitive.ObjectID{}
		// for i := 0; i < 3 && i < len(participants); i++ {
		// 	top3winnerIDs = append(top3winnerIDs, participants[i].ID)
		// }

		// Update progress for top 3 winners within the group with rewards
		rewards := []int{5000, 3000, 2000, 1000, 1000, 1000, 1000, 1000, 1000, 1000}

		for i, reward := range rewards {
			if i < len(participants) {
				err := UpdateProgress(participants[i].ID, 0, reward)
				if err != nil {
					return nil, errors.New("cannot update user progress")
				}
			}
		}

		fmt.Println("Sorted:", participants)

		winners = participants
	}

	// Save the updated tournament to the cache
	_ = CacheOneTournament(tournament.ID, winners)

	return winners, nil
}

// Get tournament winners from cache
func GetTournamentWinnersFromCache(tournamentID primitive.ObjectID) ([]Participant, error) {

	// Get the tournament from the cache
	tournament, err := GetTournamentFromCache(tournamentID)
	if err != nil {
		return nil, errors.New("cannot get tournament from cache")
	}

	var participants []Participant
	err = json.Unmarshal([]byte(tournament), &participants)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return participants, nil
}
