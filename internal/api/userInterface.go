package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tremerj/Sport-Companion/database"
	"gorm.io/gorm"
	"net/http"
)

/*
GetUsers is an API endpoint to get all Users in the application
*/
func GetUsers(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var users []database.Users
	db.Raw(`SELECT * FROM users`).Scan(&users)
	usersJSON, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(usersJSON)
	if err != nil {
		http.Error(w, "error writing response", http.StatusInternalServerError)
	}
}

/*
GetUserByID is an API endpoint to get specific User from the application
*/
func GetUserByID(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	vars := mux.Vars(r)
	id := vars["id"]
	var user database.Users
	db.Raw(`SELECT * FROM users WHERE id = $1`, id).Scan(&user)

	w.Header().Set("Content-Type", "application/json")
	userJSON, _ := json.Marshal(user)
	_, err := w.Write(userJSON)
	if err != nil {
		http.Error(w, "error writing response", http.StatusInternalServerError)
	}
}

/*
CreateUser is an API endpoint to create a User and add it to the database
*/
func CreateUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var userData, checkData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	db.Exec(`INSERT INTO users(username, password) VALUES ($1, $2)`, userData.Username, userData.Password).Rows()
	db.Raw(`SELECT * FROM users WHERE username = $1`, userData.Username).Scan(&checkData)
	if checkData.Username == "" {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}

/*
UpdateUser is an API endpoint to update the password of an existing User
*/
func UpdateUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var userData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}
	db.Exec(`UPDATE users SET password = $1 WHERE username = $2`, userData.Password, userData.Username)
}

/*
DeleteUser is an API endpoint that deletes the identified users in the endpoint URL
*/
func DeleteUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	username := mux.Vars(r)["username"]
	db.Exec(`DELETE FROM users WHERE username = $1`, username)
}
