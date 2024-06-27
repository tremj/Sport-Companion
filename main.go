package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type TeamSearch struct {
	Result int `json:"results"`
	Teams  []struct {
		ID int `json:"id"`
	} `json:"response"`
}

type LeagueSearch struct {
	Result  int `json:"results"`
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
	LreqURL := buildLeagueSearchURL(sport, os.Args[3])

	client := &http.Client{}

	Lreq, err := http.NewRequest("GET", LreqURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	apiKey := os.Getenv("SPORT_API_KEY")

	Lreq.Header.Add("x-rapidapi-key", apiKey)
	Lreq.Header.Add("x-rapidapi-host", LreqURL)

	fmt.Println(LreqURL)

	Lres, err := client.Do(Lreq)
	if err != nil {
		log.Fatal(err)
	}
	defer Lres.Body.Close()

	body, err := io.ReadAll(Lres.Body)
	if err != nil {
		log.Fatal(err)
	}

	var leagueSearch LeagueSearch
	err = json.Unmarshal(body, &leagueSearch)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(leagueSearch)
	if leagueSearch.Result == 0 {
		log.Fatal("No such league exist in the database.")
	}
	leagueID := strconv.Itoa(leagueSearch.Leagues[0].ID) // first result (better implementation later)
	team := strings.ReplaceAll(os.Args[4], "-", "%20")

	TreqURL := buildTeamSearchURL(sport, team, leagueID)

	Treq, err := http.NewRequest("GET", TreqURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	Treq.Header.Add("x-rapidapi-key", apiKey)
	Treq.Header.Add("x-rapidapi-host", TreqURL)

	Tres, err := client.Do(Treq)
	if err != nil {
		log.Fatal(err)
	}
	defer Tres.Body.Close()

	body, err = io.ReadAll(Tres.Body)
	if err != nil {
		log.Fatal(err)
	}

	var teamSearch TeamSearch
	err = json.Unmarshal(body, &teamSearch)
	if err != nil {
		log.Fatal(err)
	}
	if teamSearch.Result == 0 {
		log.Fatal("No such team exist in the database.")
	}
	teamID := strconv.Itoa(teamSearch.Teams[0].ID)
	favTeamString := sport + "-" + leagueID + "-" + teamID
	fmt.Println(favTeamString)

	_, set := os.LookupEnv("FAV_TEAM")
	if !set {
		err = os.Setenv("FAV_TEAM", favTeamString)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		newVal := os.Getenv("FAV_TEAM") + ":" + favTeamString
		err = os.Setenv("FAV_TEAM", newVal)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func addHeader() {

}

func buildLeagueSearchURL(sport string, league string) string {
	return "https://v1." + sport + ".api-sports.io/leagues?name=" + league
}

func buildTeamSearchURL(sport string, team string, leagueID string) string {
	return "https://v1." + sport + ".api-sports.io/teams?name=" + team + "&league=" + leagueID + "&season=2024"
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
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
