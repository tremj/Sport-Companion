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
	id, err := getTeamID(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = writeToFile("NBA," + id)
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
	id, err := getTeamID(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = writeToFile("NFL," + id)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func addNFLURL(season string) string {
	name := strings.ReplaceAll(os.Args[2], " ", "%20")
	return "https://v1.american-football.api-sports.io/teams?name=" + name + "&league=1&season=" + season
}

func addMLB() {
	id, err := getTeamID(2)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = writeToFile("MLB," + id)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func addMLBURL(season string) string {
	name := strings.ReplaceAll(os.Args[2], " ", "%20")
	return "https://v1.baseball.api-sports.io/teams?name=" + name + "&league=" + objects.MLBLeagueID + "&season=" + season
}

func addNBA() {
	id, err := getTeamID(3)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = writeToFile("NBA," + id)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func addNBAURL(season string) string {
	name := strings.ReplaceAll(os.Args[2], " ", "%20")
	return "https://v1.baseball.api-sports.io/teams?name=" + name + "&league=" + objects.MLBLeagueID + "&season=" + season
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

func getTeamID(id int) (string, error) {
	body, err := makeGenericRequest(0, id)
	if err != nil {
		return "", err
	}

	var teamSearch objects.TeamSearch
	err = json.Unmarshal(body, &teamSearch)
	if err != nil {
		return "", err
	}

	if teamSearch.GetResults() == 0 {
		body, err := makeGenericRequest(1, id)
		if err != nil {
			return "", err
		}

		err = json.Unmarshal(body, &teamSearch)
		if err != nil {
			return "", err
		}
		if teamSearch.Results == 0 {
			return "", errors.New("Inputted team is not an NBA team.")
		}
	}

	return strconv.Itoa(teamSearch.Response[0].ID), nil
}

func makeGenericRequest(back int, id int) ([]byte, error) {
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

func makeGenericURLHost(year string, id int) (string, string) {
	switch id {
	case 0: // NHL
		return addNHLURL(year), "https://v1.hockey.api-sports.io"
	case 1: // NFL
		return addNFLURL(year), "https://v1.american-football.api-sports.io"
	case 2: // MLB
		return addMLBURL(year), "https://v1.baseball.api-sport.io"
	case 3: //NBA
		return addNBAURL(year), "https://v2.basketball.api-sports.io"
	}
	return "", ""
}
