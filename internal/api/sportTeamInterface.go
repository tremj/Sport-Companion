package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tremerj/Sport-Companion/database"
	"gorm.io/gorm"
	"net/http"
)

/*
GetTeams is an API endpoint to get all Teams in the application
*/
func GetTeams(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var teams []database.SportTeam
	db.Raw(`SELECT * FROM sport_teams`).Scan(&teams)
	usersJSON, _ := json.Marshal(teams)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(usersJSON)
	if err != nil {
		http.Error(w, "error writing response", http.StatusInternalServerError)
	}
}

/*
GetTeamByID is an API endpoint to get specific Team from the application
*/
func GetTeamByID(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	vars := mux.Vars(r)
	id := vars["id"]
	var team database.SportTeam
	db.Raw(`SELECT * FROM sport_teams WHERE id = $1`, id).Scan(&team)

	w.Header().Set("Content-Type", "application/json")
	userJSON, _ := json.Marshal(team)
	_, err := w.Write(userJSON)
	if err != nil {

		http.Error(w, "error writing response", http.StatusInternalServerError)
	}
}

/*
CreateTeam is an API endpoint to create a Team and add it to the database
*/
func CreateTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var teamData struct {
		Name     string `json:"name"`
		Hometown string `json:"hometown"`
	}

	if err := json.NewDecoder(r.Body).Decode(&teamData); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	db.Exec(`INSERT INTO sport_teams(name, hometown) VALUES ($1, $2)`, teamData.Name, teamData.Hometown)
}

/*
UpdateTeam is an API endpoint to update the hometown of an existing Team
*/
func UpdateTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var teamData struct {
		Name     string `json:"name"`
		Hometown string `json:"hometown"`
	}

	if err := json.NewDecoder(r.Body).Decode(&teamData); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}
	db.Exec(`UPDATE sport_teams SET hometown = $1 WHERE name = $2`, teamData.Hometown, teamData.Name)
}

/*
DeleteTeam is an API endpoint that deletes the identified team in the endpoint URL
*/
func DeleteTeam(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var teamName struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&teamName); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
	}
	db.Exec(`DELETE FROM sport_teams WHERE name = $1`, teamName.Name)
}
