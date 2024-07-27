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
	"sync"
	"time"

	"github.com/tremerj/Sport-Companion/objects"
)

func HandleList() {
	printTeams()
}

func printTeams() {
	f, err := os.Open(".favourite_teams")
	if err != nil {
		fmt.Println("No teams have been added to your watchlist")
		return
	}
	defer f.Close()

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			err := printTeam(line)
			if err != nil {
				fmt.Println(err)
			}
		}(line)
	}

	wg.Wait()
}

func printTeam(line string) error {
	split := strings.Split(line, ",")
	url, host := getURLHost(split[0], split[1], 0)
	body, err := makeRequest(url, host)
	if err != nil {
		return err
	}
	name, err := getTeamName(body, split[0], split[1])
	if err != nil {
		return err
	}

	fmt.Println(name)
	return nil
}

func getTeamName(body []byte, league string, id string) (string, error) {
	var teamSearch objects.TeamSearch
	err := json.Unmarshal(body, &teamSearch)
	if err != nil {
		return "", err
	}

	if teamSearch.Results == 0 {
		url, host := getURLHost(league, id, 1)
		body, err := makeRequest(url, host)
		if err != nil {
			return "", err
		}

		err = json.Unmarshal(body, &teamSearch)
		if err != nil {
			return "", err
		}
		if teamSearch.Results == 0 {
			return "", errors.New("Team does not exist anymore")
		}
	}

	return teamSearch.Response[0].Name, nil
}

func makeRequest(url string, host string) ([]byte, error) {
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

func getURLHost(league string, id string, back int) (string, string) {
	switch league {
	case "NHL":
		return listNHLURL(id, back), "https://v1.hockey.api-sports.io"
	case "NFL":
		return listNFLURL(id, back), "https://v1.american-football.api-sports.io"
	case "MLB":
		return listMLBURL(id, back), "https://v1.baseball.api-sports.io"
	case "NBA":
		return listNBAURL(id, back), "https://v2.nba.api-sports.io"
	default:
		return "", ""
	}
}

func listNHLURL(id string, back int) string {
	return "https://v1.hockey.api-sports.io/teams?league=" + objects.NHLLeagueID + "&id=" + id + "&season=" + strconv.Itoa(time.Now().Year()-back)
}

func listNFLURL(id string, back int) string {
	return "https://v1.american-football.api-sports.io/teams?id=" + id + "&season=" + strconv.Itoa(time.Now().Year()-back)
}

func listMLBURL(id string, back int) string {
	return "https://v1.baseball.api-sports.io/teams?id=" + id + "&league=" + objects.MLBLeagueID + "&season=" + strconv.Itoa(time.Now().Year()-back)
}

func listNBAURL(id string, back int) string {
	return "https://v2.nba.api-sports.io/teams?id=" + id + "&season=" + strconv.Itoa(time.Now().Year()-back)
}
