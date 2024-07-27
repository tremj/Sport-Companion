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
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/tremerj/Sport-Companion/handling"
)

type TeamSearch struct {
	Result int `json:"results"`
	Teams  []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"response"`
}

type LeagueSearch struct {
	Result  int `json:"results"`
	Leagues []struct {
		ID int `json:"id"`
	} `json:"response"`
}

type GameSearch struct {
	Results int    `json:"results"`
	Games   []Game `json:"response"`
}

type Game struct {
	Date  string `json:"date"`
	Teams struct {
		Home struct {
			Name string `json:"name"`
		} `json:"home"`
		Away struct {
			Name string `json:"name"`
		} `json:"away"`
	} `json:"teams"`
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
		defer f.Close()
		if err != nil {
			return "There was an error opening the watchlist file"
		}

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
	defer f.Close()
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

	g, err := os.OpenFile(".favourite_teams", os.O_WRONLY|os.O_TRUNC, 0644)
	defer g.Close()
	if err != nil {
		fmt.Println("Error creating favourite team file")
		return
	}

	for _, line := range lines {
		_, err = g.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing to the file")
			return
		}
	}
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
	if (os.Args[2] == "NFL" || os.Args[2] == "NBA") && len(os.Args) == 4 {
		addOddFavourite()
	} else if len(os.Args) == 5 {
		addFavourite()
	} else {
		fmt.Println("Incorrect usage, please verify that you have the correct amount of command line arguments")
	}
}

func handleRemove() {
	if os.Args[2] == "Formula-1" && len(os.Args) == 5 {
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
	f, err := os.OpenFile(".favourite_teams", os.O_TRUNC, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println("Error accessing favourite teams")
		return
	}
}

func handleHelp() {
	fmt.Println("Usage: 'Sport-Companion schedule' to see next matches of your favourite teams")
	fmt.Println("General usage: 'Sport-Companion <add|remove> <sport> <league> <team>' to add a team to your watchlist")
	fmt.Println("Accepted sports are:")
	fmt.Printf("\tAFL\n\tBaseball\n\tBasketball\n\tFootball\n\tFormula-1\n\tHandball\n\tHockey\n\tNBA\n\tNFL\n\tRugby\n\tVolleyball\n")
	fmt.Println("When trying to add/remove teams to your list from NBA or NFL no not specify the sport argument")
	fmt.Println("When writing team names and leagues please ensure that if there are spaces in the team or league name to surround it with double quotes -> \"")
	fmt.Println("Example: Montreal Canadiens -> \"Montreal Canadiens\"")
}

func getHostAndUrl(team string, year string, purpose int8) (string, string) {
	if year == "" {
		year = strings.Split(time.Now().String(), "-")[0]
	}

	var endpoint string
	var identifier string
	// getting teams
	if purpose == 0 {
		endpoint = "teams"
		identifier = "id"
	} else if purpose == 1 {
		endpoint = "games"
		identifier = "team"
	}
	url, host := "", ""
	args := strings.Split(team, ",")

	switch strings.ToLower(args[0]) {
	case "nba":
		url = "https://v2.nba.api-sports.io/" + endpoint + "?" + identifier + "=" + args[1]
		host = "v2.nba.api-sports.io"
	case "nfl":
		url = "https://v1.american-football.api-sports.io/" + endpoint + "?" + identifier + "=" + args[1] + "&season=" + year + "&league=1"
		host = "v1.american-football.api-sports.io"
	default:
		url = "https://v1." + args[0] + ".api-sports.io/" + endpoint + "?" + identifier + "=" + args[2] + "&league=" + args[1] + "&season=" + year
		host = "v1." + args[0] + "api-sports.io"
	}

	return url, host
}

func printTeam(team string) string {
	body, errS := makeRequest(team, "", 0)
	if errS != "" {
		return errS
	}

	var teamSearch TeamSearch
	err := json.Unmarshal(body, &teamSearch)
	if err != nil {
		return err.Error()
	}
	if teamSearch.Result == 0 {
		body, errS = makeRequest(team, strconv.Itoa(time.Now().Year()-1), 0) // go back 1 year
		if errS != "" {
			return errS
		}
		err = json.Unmarshal(body, &teamSearch)
	}
	fmt.Println(teamSearch.Teams[0].Name)
	return ""
}

func handleList() {
	f, err := os.Open(".favourite_teams")
	defer f.Close()
	if err != nil {
		fmt.Println("No teams have been added to your watchlist")
		return
	}

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			err := printTeam(line)
			if err != "" {
				fmt.Println(err)
			}
		}(line)
	}

	wg.Wait()
}

func makeRequest(team string, year string, purpose int8) ([]byte, string) {
	url, host := getHostAndUrl(team, year, purpose)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err.Error()
	}

	client := &http.Client{}

	req.Header.Add("x-rapidapi-key", os.Getenv("SPORT_API_KEY"))
	req.Header.Add("x-rapidapi-host", host)

	res, err := client.Do(req)
	if err != nil {
		return nil, err.Error()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err.Error()
	}

	return body, ""
}

func isToday(game Game) bool {
	date := strings.Split(game.Date, "T")[0]
	gameDate, _ := time.Parse("2006-01-02", date)
	today := time.Now()
	return gameDate.Equal(today)
}

func isBeforeToday(game Game) bool {
	date := strings.Split(game.Date, "T")[0]
	gameDate, _ := time.Parse("2006-01-02", date)
	today := time.Now()
	return gameDate.Before(today)
}

func isWeekAway(game Game) bool {
	date := strings.Split(game.Date, "T")[0]
	gameDate, _ := time.Parse("2006-01-02", date)
	today := time.Now()
	diff := gameDate.Sub(today).Hours()
	return diff < 168 // 168 hours in 1 week
}

func formattedWeeklySchedule(games []Game) string {
	// TODO
	return ""
}

func nextWeekSchedule(games []Game) []Game {
	var eligible []Game
	var mid, next_i int
	i := -1
	l, r := 0, len(games)-1
	for l <= r {
		mid = l + (r-l)/2
		fmt.Println(games[mid].Date)
		if isToday(games[mid]) {
			i = mid
			break
		} else if isBeforeToday(games[mid]) { // arr[i] < target
			l = mid + 1
		} else {
			next_i = mid
			r = mid - 1
		}
	}

	if i != mid {
		i = next_i
	}

	for isWeekAway(games[i]) {
		eligible = append(eligible, games[i])
		i++
	}

	// fOut := formattedWeeklySchedule(eligible)
	return eligible
}

func getWeeklySchedule(team string) string {
	body, errS := makeRequest(team, strconv.Itoa(time.Now().Year()), 1)
	if errS != "" {
		return errS
	}

	var games GameSearch
	err := json.Unmarshal(body, &games)
	if err != nil {
		return err.Error()
	}
	fmt.Println(team)
	fmt.Println(games.Results)
	// sometimes current year is not equivalent to season year in API architecture
	if games.Results == 0 {
		body, errS = makeRequest(team, strconv.Itoa(time.Now().Year()-1), 1) // go back 1 year
		if errS != "" {
			return errS
		}
		err = json.Unmarshal(body, &games)
		if err != nil {
			return err.Error()
		}
		if games.Results == 0 {
			return "No scheduled games"
		}
	}

	weekAway := nextWeekSchedule(games.Games)
	if len(weekAway) == 0 {
		fmt.Printf("There are no games this week, the next game is on \n")
	}
	for _, g := range weekAway {
		fmt.Println(g.Date)
	}
	return ""

}

func showWeekSchedule() {
	f, err := os.Open(".favourite_teams")
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	// var wg sync.WaitGroup
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// wg.Add(1)
		// go func(line string) {
		// 	defer wg.Done()
		err := getWeeklySchedule(line)
		if err != "" {
			fmt.Println(err)
		}
		// }(line)
	}

	// wg.Wait()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please use one of the valid commands")
		return
	}

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading env variables")
		return
	}

	switch os.Args[1] {
	case "add":
		handling.HandleAdd()
	case "clear":
		handling.HandleClear()
	case "help":
		handling.HandleHelp()
	case "list":
		handling.HandleList()
	case "schedule":
		// TODO
	case "remove":
		handling.HandleRemove()
	default:
		fmt.Println("Unknown command")
	}
}
