package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type LeagueSearch struct {
	Result  int `json:"result"`
	Leagues []struct {
		ID int `json:"id"`
	} `json:"response"`
}

func showWeekSchedule() {
	favTeam := os.Getenv("FAV-TEAM")
	if favTeam == "" {
		fmt.Println("You have no favourite teams to check up on!")
	}
	// TODO
	// list all favourite team games for the next week (test with baseball season)
}

func addFavourite() {
	sport := strings.ToLower(os.Args[2])
	reqURL := buildLeagueSearchURL(sport, strings.ToLower(os.Args[3]))

	client := &http.Client{}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	apiKey := os.Getenv("SPORT-API-KEY")

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", reqURL)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var leagueSearch LeagueSearch
	err = json.Unmarshal(body, &leagueSearch)
	if err != nil {
		log.Fatal(err)
	}
	if leagueSearch.Result == 0 {
		log.Fatal("No such league exist in the database.")
	}
	leagueID := leagueSearch.Leagues[0].ID // first result (better implementation later)
}

func buildLeagueSearchURL(sport string, league string) string {
	return "https://v1." + sport + ".api-sports.io/teams?league=" + league
}

func main() {
	if len(os.Args) < 2 {
		showWeekSchedule()
	} else if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("Usage: 'Sport-Companion' to see next matches of your favourite teams")
		fmt.Println("Usage: 'Sport-Companion [-a|--add] <sport> <league> <team>' to add a team to your watchlist")
		fmt.Println("Accepted sports are:")
		fmt.Printf("\tAFL\n\tBaseball\n\tBasketball\n\tFootball\n\tFormula-1\n\tHandball\n\tHockey\n\tMMA\n\tNBA\n\tNFL\n\tRugby\n\tVolleyball\n")
		fmt.Println("When writing team names and leagues please ensure that spaces are replaced with '-'")
		fmt.Println("Example: 'Montreal Canadiens' -> 'Montreal-Canadiens'")
	} else if os.Args[1] == "-a" || os.Args[1] == "--add" {
		if len(os.Args) != 5 {
			fmt.Println("Usage: Sport-Companion [-a|--add] <sport> <league> <team>")
		} else {
			addFavourite()
		}
	}
}
