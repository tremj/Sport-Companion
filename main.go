package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

func getLeagueID(client *http.Client, requestURL string, apiKey string) (string, string) {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", "Error instanciating new request"
	}

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", requestURL)

	res, err := client.Do(req)
	if err != nil {
		return "", "Error fulfilling API request"
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", "Error reading JSON response body"
	}

	var leagueSearch LeagueSearch
	err = json.Unmarshal(body, &leagueSearch)
	if err != nil {
		return "", "Error parsing JSON reponse body"
	}

	if leagueSearch.Result == 0 {
		return "", "No such league exist in the database"
	}

	return strconv.Itoa(leagueSearch.Leagues[0].ID), ""
}

func getTeamID(client *http.Client, requestURL string, apiKey string) (string, string) {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", "Error instanciating new request"
	}

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", requestURL)

	res, err := client.Do(req)
	if err != nil {
		return "", "Error fulfilling API request"
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", "Error reading JSON response body"
	}

	var teamSearch TeamSearch
	err = json.Unmarshal(body, &teamSearch)
	if err != nil {
		return "", "Error parsing JSON response body"
	}

	if teamSearch.Result == 0 {
		return "", "No such team exists in the database"
	}

	return strconv.Itoa(teamSearch.Teams[0].ID), ""
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

func writeToFavourite(aFavourite string) string {
	if isAlreadySelected(aFavourite) {
		return "Team is already on your watchlist"
	} else {
		f, err := os.OpenFile(".favourite_teams", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "There was an error opening the watchlist file"
		}
		defer f.Close()

		_, err = f.WriteString(aFavourite + "\n")
		if err != nil {
			return "There was an error writing to the watchlist file"
		}
	}
	return ""
}

func addFavourite() {
	apiKey := os.Getenv("SPORT_API_KEY")
	sport := strings.ToLower(os.Args[2])
	LreqURL := buildLeagueSearchURL(sport, os.Args[3])
	client := &http.Client{}

	leagueID, errString := getLeagueID(client, LreqURL, apiKey)
	if leagueID == "" {
		fmt.Println(errString)
		return
	}

	team := strings.ReplaceAll(os.Args[4], " ", "%20")
	TreqURL := buildTeamSearchURL(sport, team, leagueID)

	teamID, errString := getTeamID(client, TreqURL, apiKey)
	if teamID == "" {
		fmt.Println(errString)
		return
	}

	code := writeToFavourite(sport + "," + leagueID + "," + teamID)
	if code != "" {
		fmt.Println(code)
	}
}

func addMMAFavourite() {
	code := writeToFavourite("mma")
	if code != "" {
		fmt.Println(code)
	}
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

	teamID, errString := getTeamID(client, reqURL, apiKey)
	if teamID == "" {
		fmt.Println(errString)
		return
	}

	code := writeToFavourite("formula-1," + choice + "," + teamID)
	if code != "" {
		fmt.Println(code)
	}
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

	teamID, errString := getTeamID(client, reqURL, apiKey)
	if teamID == "" {
		fmt.Println(errString)
		return
	}

	code := writeToFavourite(league + "," + teamID)
	if code != "" {
		fmt.Println(code)
	}
}

func removeFromFile(aFavourite string) {
	f, err := os.Open(".favourite_teams")
	if err != nil {
		fmt.Println("No teams have been added to your watchlist")
		return
	}

	scanner := bufio.NewScanner(f)
	detected := false
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != aFavourite {
			lines = append(lines, line)
		} else {
			detected = true
		}
	}

	if !detected {
		fmt.Println("No such team in your favourites")
		return
	}

	if err = scanner.Err(); err != nil {
		fmt.Println("Error reading the file")
		return
	}

	f.Close()

	g, err := os.OpenFile(".favourite_teams", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error creating favourite team file")
		return
	}
	defer g.Close()

	for _, line := range lines {
		_, err = g.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing to the file")
			return
		}
	}
}

func removeMMAFavourite() {
	removeFromFile("mma")
}

func removeF1Favourite() {
	apiKey := os.Getenv("SPORT_API_KEY")
	choice := os.Args[3]
	supporting := os.Args[4]

	reqURL := buildF1URL(choice, supporting)
	client := &http.Client{}

	teamID, errString := getTeamID(client, reqURL, apiKey)
	if teamID == "" {
		fmt.Println(errString)
		return
	}

	removeFromFile("formula-1," + choice + "," + teamID)
}

func removeOddFavourite() {
	apiKey := os.Getenv("SPORT_API_KEY")
	league := os.Args[2]
	team := os.Args[3]

	reqURL := buildOddURL(league, team)
	client := &http.Client{}

	teamID, errString := getTeamID(client, reqURL, apiKey)
	if teamID == "" {
		fmt.Println(errString)
		return
	}

	removeFromFile(league + "," + teamID)
}

func removeFavourite() {
	apiKey := os.Getenv("SPORT_API_KEY")
	sport := strings.ToLower(os.Args[2])
	LreqURL := buildLeagueSearchURL(sport, os.Args[3])
	client := &http.Client{}

	leagueID, errString := getLeagueID(client, LreqURL, apiKey)
	if leagueID == "" {
		fmt.Println(errString)
		return
	}

	team := strings.ReplaceAll(os.Args[4], " ", "%20")
	TreqURL := buildTeamSearchURL(sport, team, leagueID)

	teamID, errString := getTeamID(client, TreqURL, apiKey)
	if teamID == "" {
		fmt.Println(errString)
		return
	}

	removeFromFile(sport + "," + leagueID + "," + teamID)
}

func handleAdd() {
	if os.Args[2] == "MMA" && len(os.Args) == 3 {
		addMMAFavourite()
	} else if os.Args[2] == "Formula-1" && len(os.Args) == 5 {
		addF1Favourite()
	} else if (os.Args[2] == "NFL" || os.Args[2] == "NBA") && len(os.Args) == 4 {
		addOddFavourite()
	} else if len(os.Args) == 5 {
		addFavourite()
	} else {
		fmt.Println("Incorrect usage, please verify that you have the correct amount of command line arguments")
	}
}

func handleRemove() {
	if os.Args[2] == "MMA" && len(os.Args) == 3 {
		removeMMAFavourite()
	} else if os.Args[2] == "Formula-1" && len(os.Args) == 5 {
		removeF1Favourite()
	} else if (os.Args[2] == "NFL" || os.Args[2] == "NBA") && len(os.Args) == 4 {
		removeOddFavourite()
	} else if len(os.Args) == 5 {
		removeFavourite()
	} else {
		fmt.Println("Incorrect usage, please verify that you have the correct amount of command line arguments")
	}
}

func handleClear() {
	f, err := os.OpenFile("./.favourite_teams", os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error accessing favourite teams")
		return
	}
	defer f.Close()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading env variables")
		return
	}
	if len(os.Args) < 2 {
		showWeekSchedule()
	} else if os.Args[1] == "help" {
		fmt.Println("Usage: 'Sport-Companion schedule' to see next matches of your favourite teams")
		fmt.Println("General usage: 'Sport-Companion <add|remove> <sport> <league> <team>' to add a team to your watchlist")
		fmt.Println("Accepted sports are:")
		fmt.Printf("\tAFL\n\tBaseball\n\tBasketball\n\tFootball\n\tFormula-1\n\tHandball\n\tHockey\n\tMMA\n\tNBA\n\tNFL\n\tRugby\n\tVolleyball\n")
		fmt.Println("When trying to add/remove teams to your list from NBA or NFL no not specify the sport argument")
		fmt.Println("Formula 1 usage goes as such: Sport-Companion <add|remove> Formula-1 <team|driver> <team name|driver name>")
		fmt.Println("If you want to add MMA to your watchlist simply write MMA after add do not specify league or team arguments:w")
		fmt.Println("When writing team names and leagues please ensure that if there are spaces in the team or league name to surround it with double quotes -> \"")
		fmt.Println("Example: 'Montreal Canadiens' -> 'Montreal-Canadiens'")
	} else if os.Args[1] == "add" {
		handleAdd()
	} else if os.Args[1] == "remove" {
		handleRemove()
	} else if os.Args[1] == "clear" {
		handleClear()
	}
}
