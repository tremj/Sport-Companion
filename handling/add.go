package handling

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tremerj/Sport-Companion/objects"
)

func HandleAdd() {
	if len(os.Args) != 4 {
		fmt.Println("Incorrect usage.")
		fmt.Println("Hint: Sport-Companion add <league> <teamname>")
		return
	}
	switch strings.ToLower(os.Args[2]) {
	case "nhl":
		addNHL()
	case "nfl":
		addNFL()
	case "nba":
		addNBA()
	case "mlb":
		addMLB()
	default:
		fmt.Println("Unsupported league was inputted")
	}
}

func addNHL() {
	id, err := getTeamID("NHL")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = writeToFile("NHL," + id + "," + os.Args[3])
	if err != nil {
		fmt.Println(err)
		return
	}
}

func addNHLURL(season string) string {
	teamName := strings.ReplaceAll(os.Args[3], " ", "%20")
	return "https://v1.hockey.api-sports.io/teams?name=" + teamName + "&league=" + objects.NHLLeagueID + "&season=" + season
}

func addNFL() {
	id, err := getTeamID("NFL")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = writeToFile("NFL," + id + "," + os.Args[3])
	if err != nil {
		fmt.Println(err)
		return
	}
}

func addNFLURL(season string) string {
	name := strings.ReplaceAll(os.Args[3], " ", "%20")
	return "https://v1.american-football.api-sports.io/teams?name=" + name + "&league=" + objects.NFLLeagueID + "&season=" + season
}

func addMLB() {
	id, err := getTeamID("MLB")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = writeToFile("MLB," + id + "," + os.Args[3])
	if err != nil {
		fmt.Println(err)
		return
	}
}

func addMLBURL(season string) string {
	name := strings.ReplaceAll(os.Args[3], " ", "%20")
	return "https://v1.baseball.api-sports.io/teams?name=" + name + "&league=" + objects.MLBLeagueID + "&season=" + season
}

func addNBA() {
	id, err := getTeamID("NBA")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = writeToFile("NBA," + id + "," + os.Args[3])
	if err != nil {
		fmt.Println(err)
		return
	}
}

func addNBAURL() string {
	name := strings.ReplaceAll(os.Args[3], " ", "%20")
	return "https://v2.nba.api-sports.io/teams?name=" + name
}

func addToFile(line string) error {
	if isAlreadyAdded(line) {
		return errors.New("This team has already been added to your list")
	}

	err := writeToFile(line)
	if err != nil {
		return err
	}

	return nil
}

func writeToFile(line string) error {
	f, err := os.OpenFile(".favourite_teams", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = f.WriteString(line + "\n")
	if err != nil {
		return err
	}

	return nil
}

func isAlreadyAdded(line string) bool {
	f, err := os.Open(".favourite_teams")
	defer f.Close()
	if err != nil {
		return false
	}

	fScanner := bufio.NewScanner(f)
	fScanner.Split(bufio.ScanLines)

	for fScanner.Scan() {
		if line == fScanner.Text() {
			return true
		}
	}
	return false
}

func getTeamID(id string) (string, error) {
	body, err := makeGenericRequest(0, id)
	if err != nil {
		return "", err
	}

	var teamSearch objects.TeamSearch
	err = json.Unmarshal(body, &teamSearch)
	if err != nil {
		return "", err
	}

	if teamSearch.Results == 0 {
		body, err := makeGenericRequest(1, id)
		if err != nil {
			return "", err
		}

		err = json.Unmarshal(body, &teamSearch)
		if err != nil {
			return "", err
		}
		if teamSearch.Results == 0 {
			return "", errors.New("Inputted team is not an " + id + " team.")
		}
	}

	return strconv.Itoa(teamSearch.Response[0].ID), nil
}

func makeGenericRequest(back int, id string) ([]byte, error) {
	url, host := makeGenericURLHost(strconv.Itoa(time.Now().Year()-back), id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req.Header.Add("x-rapidapi-key", os.Getenv("SPORT_API_KEY"))
	req.Header.Add("x-rapidapi-host", host)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func makeGenericURLHost(year string, id string) (string, string) {
	switch strings.ToLower(id) {
	case "nhl":
		return addNHLURL(year), "https://v1.hockey.api-sports.io"
	case "nfl":
		return addNFLURL(year), "https://v1.american-football.api-sports.io"
	case "mlb":
		return addMLBURL(year), "https://v1.baseball.api-sport.io"
	case "nba":
		return addNBAURL(), "https://v2.basketball.api-sports.io"
	}
	return "", ""
}
