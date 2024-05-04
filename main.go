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

	router.HandleFunc("/users", func(writer http.ResponseWriter, request *http.Request) {
		api.GetUsers(writer, request, db)
	}).Methods("GET")
	router.HandleFunc("/users/{id}", func(writer http.ResponseWriter, request *http.Request) {
		api.GetUserByID(writer, request, db)
	}).Methods("GET")
	router.HandleFunc("/users/create", func(writer http.ResponseWriter, request *http.Request) {
		api.CreateUser(writer, request, db)
	}).Methods("POST")

	port := ":8080"
	fmt.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
