package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tremerj/Sport-Companion/database"
	"gorm.io/gorm"
	"net/http"
)

/*
GetLeagues is an API endpoint to get all Teams in the application
*/
func GetLeagues(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var leagues []database.League
	db.Raw(`SELECT * FROM leagues`).Scan(&leagues)
	usersJSON, _ := json.Marshal(leagues)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(usersJSON)
	if err != nil {
		fmt.Println(err)
	}
}

/*
GetLeagueByID is an API endpoint to get specific League from the application
*/
func GetLeagueByID(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	vars := mux.Vars(r)
	id := vars["id"]
	var league database.SportTeam
	db.Raw(`SELECT * FROM leagues WHERE id = $1`, id).Scan(&league)

	w.Header().Set("Content-Type", "application/json")
	userJSON, _ := json.Marshal(league)
	_, err := w.Write(userJSON)
	if err != nil {
		fmt.Println(err)
	}
}

/*
CreateLeague is an API endpoint to create a League and add it to the database
*/
func CreateLeague(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var leagueData struct {
		Name  string `json:"name"`
		Sport string `json:"sport"`
	}

	if err := json.NewDecoder(r.Body).Decode(&leagueData); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	db.Exec(`INSERT INTO leagues(name, sport) VALUES ($1, $2)`, leagueData.Name, leagueData.Sport)
}

/*
UpdateLeague is an API endpoint to update the name of an existing League
*/
func UpdateLeague(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var leagueData struct {
		NewName string `json:"new_name"`
		OldName string `json:"old_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&leagueData); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}
	db.Exec(`UPDATE leagues SET name = $1 WHERE name = $2`, leagueData.NewName, leagueData.OldName)
}

/*
DeleteLeague is an API endpoint that deletes the identified League in the endpoint URL
*/
func DeleteLeague(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var leagueName struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&leagueName); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
	}
	db.Exec(`DELETE FROM leagues WHERE name = $1`, leagueName.Name)
}
