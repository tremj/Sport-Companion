package api

import (
	"encoding/json"
	"github.com/tremerj/Sport-Companion/database"
	"gorm.io/gorm"
	"net/http"
)

/*
LinkTeamToUser is an API endpoint that will link a User to a Team as a fan
*/
func LinkTeamToUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var data struct {
		Username string `json:"username"`
		TeamName string `json:"teamname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var userID, teamID uint

	db.Raw(`SELECT id FROM users WHERE username = ?`, data.Username).Scan(&userID)
	db.Raw(`SELECT id FROM sport_teams WHERE name = ?`, data.TeamName).Scan(&teamID)

	if userID == 0 || teamID == 0 {
		http.Error(w, "user or team not found", http.StatusNotFound)
		return
	}

	db.Exec(`INSERT INTO user_teams (user_id, team_id) VALUES ($1, $2)`, userID, teamID)
}

/*
RemoveTeamFromUser is an API endpoint that will remove a team
from a users favourite teams list and stop tracking their schedule
*/
func RemoveTeamFromUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var data struct {
		Username string `json:"username"`
		TeamName string `json:"teamname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var userID, teamID uint

	db.Raw(`SELECT id FROM users WHERE username = ?`, data.Username).Scan(&userID)
	db.Raw(`SELECT id FROM sport_teams WHERE name = ?`, data.TeamName).Scan(&teamID)

	if userID == 0 || teamID == 0 {
		http.Error(w, "user or team not found", http.StatusNotFound)
	}

	db.Exec(`DELETE FROM user_teams WHERE user_id = $1 AND team_id = $2`, userID, teamID)
}

/*
GetFavouriteTeams is an API endpoint to list all of a user's
favourite teams
*/
func GetFavouriteTeams(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var userData struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	var userID uint

	db.Raw(`SELECT id FROM users WHERE username = $1`, userData.Username).Scan(&userID)
	if userID == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	var teams []uint

	db.Raw(`SELECT DISTINCT team_id FROM user_teams WHERE user_id = $1`, userID).Scan(&teams)

	var teamsData []database.SportTeam
	for i := 0; i < len(teams); i++ {
		var data database.SportTeam
		db.Raw(`SELECT * FROM sport_teams WHERE id = $1`, teams[i]).Scan(&data)
		teamsData = append(teamsData, data)
	}

	jsonData, err := json.Marshal(teamsData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetFans(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var data struct {
		TeamName string `json:"teamname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
	}

	var teamID uint
	db.Raw(`SELECT id FROM sport_teams WHERE name = $1`, data.TeamName).Scan(&teamID)
	if teamID == 0 {
		http.Error(w, "team not found", http.StatusNotFound)
	}

	var users []uint
	db.Raw(`SELECT DISTINCT user_id FROM user_teams WHERE team_id = $1`, teamID).Scan(&users)

	var usernames []string
	for i := 0; i < len(users); i++ {
		var user string
		db.Raw(`SELECT username FROM users WHERE id = $1`, users[i]).Scan(&user)
		usernames = append(usernames, user)
	}
	jsonData, err := json.Marshal(usernames)
	if err != nil {
		http.Error(w, "Unable tp convert userData to JSON", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "Error writing JSON to response", http.StatusInternalServerError)
	}
}
