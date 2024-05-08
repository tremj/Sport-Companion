package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tremerj/Sport-Companion/database"
	"github.com/tremerj/Sport-Companion/internal/api"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	db, err := gorm.Open(sqlite.Open("database/companionV2.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	err = db.AutoMigrate(&database.Users{}, &database.SportTeam{},
		&database.Match{}, &database.League{},
		&database.UserTeam{}, &database.MatchTeam{}, &database.LeagueTeam{})
	if err != nil {
		log.Fatal("Failed to migrate database", err)
	}

	router := mux.NewRouter()

	// user endpoints
	router.HandleFunc("/users", func(writer http.ResponseWriter, request *http.Request) {
		api.GetUsers(writer, request, db)
	}).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}", func(writer http.ResponseWriter, request *http.Request) {
		api.GetUserByID(writer, request, db)
	}).Methods("GET")
	router.HandleFunc("/users/create", func(writer http.ResponseWriter, request *http.Request) {
		api.CreateUser(writer, request, db)
	}).Methods("POST")
	router.HandleFunc("/users/update", func(writer http.ResponseWriter, request *http.Request) {
		api.UpdateUser(writer, request, db)
	}).Methods("PUT")
	router.HandleFunc("/users/delete/{username}", func(writer http.ResponseWriter, request *http.Request) {
		api.DeleteUser(writer, request, db)
	}).Methods("DELETE")
	router.HandleFunc("/users/addTeam", func(writer http.ResponseWriter, request *http.Request) {
		api.LinkTeamToUser(writer, request, db)
	}).Methods("POST")
	router.HandleFunc("/users/removeTeam", func(writer http.ResponseWriter, request *http.Request) {
		api.RemoveTeamFromUser(writer, request, db)
	}).Methods("DELETE")
	router.HandleFunc("/users/favourites", func(writer http.ResponseWriter, request *http.Request) {
		api.GetFavouriteTeams(writer, request, db)
	}).Methods("GET")

	// sport team endpoints
	router.HandleFunc("/teams", func(writer http.ResponseWriter, request *http.Request) {
		api.GetTeams(writer, request, db)
	}).Methods("GET")
	router.HandleFunc("/teams/{id}", func(writer http.ResponseWriter, request *http.Request) {
		api.GetTeamByID(writer, request, db)
	}).Methods("GET")
	router.HandleFunc("/teams/create", func(writer http.ResponseWriter, request *http.Request) {
		api.CreateTeam(writer, request, db)
	}).Methods("POST")
	router.HandleFunc("/teams/update", func(writer http.ResponseWriter, request *http.Request) {
		api.UpdateTeam(writer, request, db)
	}).Methods("PUT")
	router.HandleFunc("/teams/delete", func(writer http.ResponseWriter, request *http.Request) {
		api.DeleteTeam(writer, request, db)
	}).Methods("DELETE")
	//router.HandleFunc("/teams/fans", func(writer http.ResponseWriter, request *http.Request) {
	//	api.
	//})

	// league endpoints
	router.HandleFunc("/leagues", func(writer http.ResponseWriter, request *http.Request) {
		api.GetLeagues(writer, request, db)
	}).Methods("GET")
	router.HandleFunc("/leagues/{id}", func(writer http.ResponseWriter, request *http.Request) {
		api.GetLeagueByID(writer, request, db)
	}).Methods("GET")
	router.HandleFunc("/leagues/create", func(writer http.ResponseWriter, request *http.Request) {
		api.CreateLeague(writer, request, db)
	}).Methods("POST")
	router.HandleFunc("/leagues/update", func(writer http.ResponseWriter, request *http.Request) {
		api.UpdateLeague(writer, request, db)
	}).Methods("PUT")
	router.HandleFunc("/leagues/delete", func(writer http.ResponseWriter, request *http.Request) {
		api.DeleteLeague(writer, request, db)
	}).Methods("DELETE")

	// matches endpoints
	router.HandleFunc("/matches", func(writer http.ResponseWriter, request *http.Request) {
		api.GetMatches(writer, request, db)
	}).Methods("GET")
	router.HandleFunc("/matches/{id}", func(writer http.ResponseWriter, request *http.Request) {
		api.GetMatchByID(writer, request, db)
	}).Methods("GET")
	router.HandleFunc("/matches/create", func(writer http.ResponseWriter, request *http.Request) {
		api.CreateMatch(writer, request, db)
	}).Methods("POST")
	router.HandleFunc("/matches/update", func(writer http.ResponseWriter, request *http.Request) {
		api.UpdateMatch(writer, request, db)
	}).Methods("PUT")
	router.HandleFunc("/matches/delete", func(writer http.ResponseWriter, request *http.Request) {
		api.DeleteMatch(writer, request, db)
	}).Methods("DELETE")

	port := ":8080"
	fmt.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
