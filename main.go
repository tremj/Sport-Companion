package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
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
	logErr(err)

	apiKey := os.Getenv("SPORT_API_KEY")

	Lreq.Header.Add("x-rapidapi-key", apiKey)
	Lreq.Header.Add("x-rapidapi-host", LreqURL)

	Lres, err := client.Do(Lreq)
	logErr(err)
	defer Lres.Body.Close()

	body, err := io.ReadAll(Lres.Body)
	logErr(err)

	var leagueSearch LeagueSearch
	err = json.Unmarshal(body, &leagueSearch)
	logErr(err)

	if leagueSearch.Result == 0 {
		log.Fatal("No such league exist in the database.")
	}

	leagueID := strconv.Itoa(leagueSearch.Leagues[0].ID) // first result (better implementation later)
	team := strings.ReplaceAll(os.Args[4], "-", "%20")

	TreqURL := buildTeamSearchURL(sport, team, leagueID)

	Treq, err := http.NewRequest("GET", TreqURL, nil)
	logErr(err)

	Treq.Header.Add("x-rapidapi-key", apiKey)
	Treq.Header.Add("x-rapidapi-host", TreqURL)

	Tres, err := client.Do(Treq)
	logErr(err)

	defer Tres.Body.Close()

	body, err = io.ReadAll(Tres.Body)
	logErr(err)

	var teamSearch TeamSearch
	err = json.Unmarshal(body, &teamSearch)
	logErr(err)

	if teamSearch.Result == 0 {
		log.Fatal("No such team exist in the database.")
	}
	teamID := strconv.Itoa(teamSearch.Teams[0].ID)
	favTeamString := sport + "-" + leagueID + "-" + teamID
	fmt.Println(favTeamString)
	f, err := os.Open(".favourite_teams")
	logErr(err)

}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
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
	} else if os.Args[1] == "add" {
		if len(os.Args) != 5 {
			fmt.Println("Usage: Sport-Companion [-a|--add] <sport> <league> <team>")
		} else {
			addFavourite()
		}
	}
}
