package main

import (
	"bufio"
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

func buildLeagueSearchURL(sport string, league string) string {
	return "https://v1." + sport + ".api-sports.io/leagues?name=" + league
}

func buildTeamSearchURL(sport string, team string, leagueID string) string {
	return "https://v1." + sport + ".api-sports.io/teams?name=" + team + "&league=" + leagueID + "&season=2024"
}

func getLeagueID(client *http.Client, requestURL string, apiKey string) string {
	Lreq, err := http.NewRequest("GET", requestURL, nil)
	logErr(err)

	Lreq.Header.Add("x-rapidapi-key", apiKey)
	Lreq.Header.Add("x-rapidapi-host", requestURL)

	Lres, err := client.Do(Lreq)
	logErr(err)
	body, err := io.ReadAll(Lres.Body)
	Lres.Body.Close()
	logErr(err)

	var leagueSearch LeagueSearch
	err = json.Unmarshal(body, &leagueSearch)
	logErr(err)

	if leagueSearch.Result == 0 {
		log.Fatal("No such league exist in the database.")
	}

	return strconv.Itoa(leagueSearch.Leagues[0].ID) // first result (better implementation later)
}

func getTeamID(client *http.Client, requestURL string, apiKey string) string {
	req, err := http.NewRequest("GET", requestURL, nil)
	logErr(err)

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", requestURL)

	res, err := client.Do(req)
	logErr(err)

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	logErr(err)

	var teamSearch TeamSearch
	err = json.Unmarshal(body, &teamSearch)
	logErr(err)

	if teamSearch.Result == 0 {
		log.Fatal("No such team exist in the database.")
	}

	return strconv.Itoa(teamSearch.Teams[0].ID)
}

func isAlreadySelected(aFavourite string) bool {
	f, err := os.Open(".favourite_teams")
	defer f.Close()
	if err != nil {
		return false
	}
	fScanner := bufio.NewScanner(f)
	fScanner.Split(bufio.ScanLines)

	for fScanner.Scan() {
		if aFavourite == fScanner.Text() {
			return true
		}
	}
	return false
}

func writeToFavourite(aFavourite string) {
	if isAlreadySelected(aFavourite) {
		log.Fatal("Team has already been added to your watchlist")
	} else {
		f, err := os.OpenFile(".favourite_teams", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		logErr(err)

		_, err = f.WriteString(aFavourite + "\n")
		logErr(err)
	}
}

func addFavourite() {
	apiKey := os.Getenv("SPORT_API_KEY")
	sport := strings.ToLower(os.Args[2])
	LreqURL := buildLeagueSearchURL(sport, os.Args[3])
	client := &http.Client{}

	leagueID := getLeagueID(client, LreqURL, apiKey) // first result (better implementation later)

	team := strings.ReplaceAll(os.Args[4], " ", "%20")
	TreqURL := buildTeamSearchURL(sport, team, leagueID)

	teamID := getTeamID(client, TreqURL, apiKey)

	favTeamString := sport + "," + leagueID + "," + teamID

	writeToFavourite(favTeamString)
}

func logErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func addMMAFavourite() {
	writeToFavourite("mma")
}

func buildF1URL(choice string, supporting string) string {
	choice = choice + "s"
	return "https://v1.formula-1.api-sports.io/" + choice + "?name=" + strings.ReplaceAll(supporting, " ", "%20")
}

func addF1Favourite() {
	apiKey := os.Getenv("SPORT_API_KEY")
	choice := os.Args[3]
	supporting := os.Args[4]
	reqURL := buildF1URL(choice, supporting)
	client := &http.Client{}
	supportID := getTeamID(client, reqURL, apiKey)
	writeToFavourite("formula-1," + choice + "," + supportID)
}

func buildOddURL(league string, team string) string {
	prefix := "v1.american-football"
	if league == "NBA" {
		prefix = "v2.nba"
	}
	return "https://" + prefix + ".api-sports.io/teams?name=" + strings.ReplaceAll(team, " ", "%20")
}

func addOddFavourite() {
	apiKey := os.Getenv("SPORT_API_KEY")
	league := os.Args[2]
	team := os.Args[3]
	reqURL := buildOddURL(league, team)
	client := &http.Client{}
	teamID := getTeamID(client, reqURL, apiKey)
	writeToFavourite(league + "," + teamID)
}

func main() {
	err := godotenv.Load()
	logErr(err)
	if len(os.Args) < 2 {
		showWeekSchedule()
	} else if os.Args[1] == "help" {
		fmt.Println("Usage: 'Sport-Companion' to see next matches of your favourite teams")
		fmt.Println("Usage: 'Sport-Companion add <sport> <league> <team>' to add a team to your watchlist")
		fmt.Println("Accepted sports are:")
		fmt.Printf("\tAFL\n\tBaseball\n\tBasketball\n\tFootball\n\tFormula-1\n\tHandball\n\tHockey\n\tMMA\n\tNBA\n\tNFL\n\tRugby\n\tVolleyball\n")
		fmt.Println("When trying to add teams to your list from NBA or NFL no not specify the sport argument")
		fmt.Println("Formula 1 usage goes as such: Sport-Companion add Formula-1 <team|driver> <team name|driver name>")
		fmt.Println("If you want to add MMA to your watchlist simply write MMA after add do not specify league or team arguments:w")
		fmt.Println("When writing team names and leagues please ensure that if there are spaces in the team or league name to surround it with double quotes -> \"")
		fmt.Println("Example: 'Montreal Canadiens' -> 'Montreal-Canadiens'")
	} else if os.Args[1] == "add" {
		if os.Args[2] == "MMA" {
			addMMAFavourite()
		} else if os.Args[2] == "Formula-1" {
			addF1Favourite()
		} else if os.Args[2] == "NFL" || os.Args[2] == "NBA" {
			addOddFavourite()
		} else if len(os.Args) != 5 {
			fmt.Println("Usage: Sport-Companion add <sport> <league> <team>")
		} else {
			addFavourite()
		}
	}
}
