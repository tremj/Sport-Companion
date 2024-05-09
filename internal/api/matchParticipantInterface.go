package api

import (
	"encoding/json"
	"github.com/tremerj/Sport-Companion/database"
	"gorm.io/gorm"
	"net/http"
)

/*
AddMatchParticipants adds an array of participants to a match.
Usually 2 but handles situations where F1 might be added
*/
func AddMatchParticipants(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var data struct {
		Participants []string `json:"participants"`
		Title        string   `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
	}

	var matchID uint
	db.Raw(`SELECT id FROM matches WHERE title = ?`, data.Title).Scan(&matchID)

	for _, participant := range data.Participants {
		var participantID uint
		db.Raw(`SELECT id FROM sport_teams WHERE name - $1`, participant).Scan(&participantID)
		if participantID != 0 {
			http.Error(w, "Team "+participant+" does not exist", http.StatusBadRequest)
			continue
		}
		db.Exec(`INSERT INTO match_teams (match_id, team_id) VALUES ($1, $2)`, matchID, participantID)
	}
}

/*
RemoveMatchParticipant removes a team from a match
*/
func RemoveMatchParticipant(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var data struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
	}

	var teamID, matchID uint
	db.Raw(`SELECT id FROM sport_teams WHERE name = $1`, data.Name).Scan(&teamID)
	db.Raw(`SELECT id FROM matches WHERE title = $1`, data.Title)
	db.Exec(`DELETE FROM match_teams WHERE team_id = $1 AND match_id = $2`, teamID, matchID)
}

/*
GetSchedule fetches all matches a team participates in
*/
func GetSchedule(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var data struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
		return
	}

	var teamID uint
	db.Raw(`SELECT id FROM sport_teams WHERE name = $1`, data.Name).Scan(&teamID)

	var matchIDs []uint
	db.Raw(`SELECT DISTINCT match_id FROM match_teams WHERE team_id = $1`, teamID).Scan(&matchIDs)

	var schedule []database.Match
	for _, matchID := range matchIDs {
		var match database.Match
		db.Raw(`SELECT * FROM matches WHERE id = $1`, matchID).Scan(&match)
		schedule = append(schedule, match)
	}
	jsonBody, err := json.Marshal(schedule)
	if err != nil {
		http.Error(w, "Error creating JSON", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBody)
	if err != nil {
		http.Error(w, "Error writing JSON", http.StatusInternalServerError)
	}
}
