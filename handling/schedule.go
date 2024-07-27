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

func HandleSchedule() {
	printSchedule()
}

func printSchedule() {
	result, err := getSchedule()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range result {
		fmt.Println(v)
	}
}

func getSchedule() ([]string, error) {
	f, err := os.Open(".favourite_teams")
	defer f.Close()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	results := make(chan string)
	scanner := bufio.NewScanner(f)

	collectedResults := []string{}
	done := make(chan bool)
	go func() {
		for result := range results {
			collectedResults = append(collectedResults, result)
		}
		close(done)
	}()

	for scanner.Scan() {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			result, err := getNextFiveGames(line)
			if err != nil {
				results <- err.Error()
			} else {
				results <- result
			}
		}(scanner.Text())
	}

	wg.Wait()
	close(results)

	<-done

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return collectedResults, nil
}

func getNextFiveGames(line string) (string, error) {
	arr := strings.Split(line, ",")
	switch arr[0] {
	case "NHL":
		return nextFiveNHL(arr[1])
	case "NFL":
		return nextFiveNFL(arr[1])
	case "MLB":
		return nextFiveMLB(arr[1])
	case "NBA":
		return nextFiveNBA(arr[1])
	default:
		return "", errors.New("Unknown behavior")
	}
}

func nextFiveNHL(id string) (string, error) {
	url := gameNHLURL(id, 0)
	body, err := makeGenericGameRequest(url, "https://v1.hockey.api-sports.io")
	if err != nil {
		return "", err
	}

	var games objects.NHLGameSearch
	err = json.Unmarshal(body, &games)
	if err != nil {
		return "", err
	}

	if games.Results == 0 {
		url = gameNHLURL(id, 1)
		body, err = makeGenericGameRequest(url, "https://v1.hockey.api-sports.io")
		if err != nil {
			return "", err
		}

		err = json.Unmarshal(body, &games)
		if err != nil {
			return "", err
		}

		if games.Results == 0 {
			return "", errors.New("No games found.")
		}
	}

	formattedOutput := formatNHLGames(games.Games)
	return formattedOutput, nil
}

func gameNHLURL(id string, back int) string {
	return "https://v1.hockey.api-sports.io/games?league=" + objects.NHLLeagueID + "&id=" + id + "?season=" + strconv.Itoa(time.Now().Year()-back)
}

func formatNHLGames(games []objects.NHLGame) string {
	upcoming := getNextFiveGamesNHL(games)
	return formatNHL(upcoming)
}

func formatNHL(games []objects.NHLGame) string {
	var res string
	for _, g := range games {
		res = res + getNHLFormat(g)
	}
	return res
}

func getNHLFormat(g objects.NHLGame) string {
	// Format: "Monday September 25th @ HH:MM >>> <Home Team> vs. <Away Team>
	date := strings.Split(g.Date, "T")[0]
	dateObj, _ := time.Parse("2006-01-02", date)
	weekDay := dateObj.Weekday().String()
	month := dateObj.Month().String()
	day := dateObj.Day()
	dateDescriptor := getDescriptor(day)
	timeBackFourHours := backFourHours(g.Time)
	return weekDay + " " + month + " " + strconv.Itoa(day) + dateDescriptor + " @ " + timeBackFourHours + " >>> " + g.Teams.Home.Name + " vs. " + g.Teams.Away.Name
}

func getNextFiveGamesNHL(games []objects.NHLGame) []objects.NHLGame {
	startIndex := findStartNHL(games)
	var nextFive []objects.NHLGame
	if len(games)-startIndex < 5 {
		for i := startIndex; i < len(games); i++ {
			nextFive = append(nextFive, games[i])
		}
		return nextFive
	}

	nextFive = games[startIndex : startIndex+5]
	return nextFive
}

func findStartNHL(games []objects.NHLGame) int {
	if isBeforeTodayNHL(games[0]) {
		return 0
	}

	l, r := 0, len(games)-1
	var nextIndex int
	i := -1
	for l <= r {
		mid := l + (r-l)/2
		if isTodayNHL(games[mid]) {
			i = mid
			break
		} else if isBeforeTodayNHL(games[mid]) { // arr[i] < target
			l = mid + 1
		} else {
			nextIndex = mid
			r = mid - 1
		}
	}

	if i != -1 {
		return i
	}
	return nextIndex
}

func isBeforeTodayNHL(game objects.NHLGame) bool {
	date := strings.Split(game.Date, "T")[0]
	dateObj, _ := time.Parse("2006-01-02", date)
	todayString := strings.Split(time.Now().String(), " ")[0]
	formatToday, _ := time.Parse("2006-01-02", todayString)
	return dateObj.Before(formatToday)
}

func isTodayNHL(game objects.NHLGame) bool {
	date := strings.Split(game.Date, "T")[0]
	dateObj, _ := time.Parse("2006-01-02", date)
	todayString := strings.Split(time.Now().String(), " ")[0]
	formatToday, _ := time.Parse("2006-01-02", todayString)
	return dateObj.Equal(formatToday)
}

func nextFiveNFL(id string) (string, error)

func nextFiveMLB(id string) (string, error)

func nextFiveNBA(id string) (string, error)

func makeGenericGameRequest(url string, host string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-key", os.Getenv("SPORT_API_KEY"))
	req.Header.Add("x-rapidapi-host", host)

	client := &http.Client{}

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

func getDescriptor(day int) string {
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func backFourHours(time string) string {
	hourMin := strings.Split(time, ":")
	intHour, _ := strconv.Atoi(hourMin[0])
	intHour = (intHour - 4) % 24
	return strconv.Itoa(intHour) + ":" + hourMin[1]
}
