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
GetMatches is an API endpoint to get all Matches in the application
*/
func GetMatches(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	fmt.Println("GetMatches")
	var matches []database.Match
	db.Raw(`SELECT * FROM matches`).Scan(&matches)
	usersJSON, _ := json.Marshal(matches)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(usersJSON)
	if err != nil {
		http.Error(w, "error writing response", http.StatusInternalServerError)
	}
}

/*
GetMatchByID is an API endpoint to get specific Match from the application
*/
func GetMatchByID(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	fmt.Println("GetMatchByID")
	vars := mux.Vars(r)
	id := vars["id"]
	var match database.Match
	db.Raw(`SELECT * FROM matches WHERE id = $1`, id).Scan(&match)

	w.Header().Set("Content-Type", "application/json")
	userJSON, _ := json.Marshal(match)
	_, err := w.Write(userJSON)
	if err != nil {
		http.Error(w, "error writing response", http.StatusInternalServerError)
	}
}

/*
CreateMatch is an API endpoint to create a Match and add it to the database
*/
func CreateMatch(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	fmt.Println("CreateMatch")
	var matchData struct {
		Title string `json:"title"`
		Time  string `json:"time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&matchData); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	db.Exec(`INSERT INTO matches(title, time) VALUES ($1, $2)`, matchData.Title, matchData.Time)
}

/*
UpdateMatch is an API endpoint to update the Time of an existing Match
*/
func UpdateMatch(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	fmt.Println("UpdateMatch")
	var matchData struct {
		Title string `json:"title"`
		Time  string `json:"time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&matchData); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}
	db.Exec(`UPDATE matches SET time = $1 WHERE title = $2`, matchData.Time, matchData.Title)
}

/*
DeleteMatch is an API endpoint that deletes the identified Match in the endpoint URL
*/
func DeleteMatch(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	fmt.Println("DeleteMatch")
	var matchTitle struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&matchTitle); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
	}
	db.Exec(`DELETE FROM matches WHERE title = $1`, matchTitle.Title)
}
