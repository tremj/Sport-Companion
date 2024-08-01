package scheduling

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/tremerj/Sport-Companion/objects"
)

func NextFiveNBA(id string, name string) (string, error) {
	url := gameNBAURL(id, 0)
	body, err := makeGenericGameRequest(url, "https://v1.hockey.api-sports.io")
	if err != nil {
		return "", err
	}

	var games objects.NBAGameSearch
	err = json.Unmarshal(body, &games)
	if err != nil {
		return "", err
	}

	if games.Results == 0 {
		url = gameNBAURL(id, 1)
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

	formattedOutput := formattedNBAGames(games.Games)
	finalResults := name + ":\n" + formattedOutput + "\n"
	return finalResults, nil
}

func gameNBAURL(id string, back int) string {
	return "https://v1.hockey.api-sports.io/games?id=" + id + "&season=" + strconv.Itoa(time.Now().Year()-back)
}

func formattedNBAGames(games []objects.NBAGame) string {
	upcoming := getNextFiveGamesNBA(games)
	return formatNBAGame(upcoming)
}

func formatNBAGame(games []objects.NBAGame) string {
	var res string
	for _, g := range games {
		res = res + getNBAFormat(g)
	}
	return res
}

func getNBAFormat(g objects.NBAGame) string {
	// Format: "Monday September 25th @ HH:MM >>> <Home Team> vs. <Away Team>
	date := strings.Split(g.Date.Date, "T")
	dateObj, _ := time.Parse("2006-01-02", date[0])
	weekDay := dateObj.Weekday().String()
	month := dateObj.Month().String()
	day := dateObj.Day()
	dateDescriptor := getDescriptor(day)
	timeSplit := strings.Split(date[1], ":")
	timeBackFourHours := backFourHours(timeSplit[0] + ":" + timeSplit[1])
	return weekDay + " " + month + " " + strconv.Itoa(day) + dateDescriptor + " @ " + timeBackFourHours + " >>> " + g.Teams.Home.Name + " vs. " + g.Teams.Away.Name + "\n"
}

func getNextFiveGamesNBA(games []objects.NBAGame) []objects.NBAGame {
	startIndex := findStartIndexNBA(games)
	var nextFive []objects.NBAGame
	if len(games)-startIndex < 5 {
		for i := startIndex; i < len(games); i++ {
			nextFive = append(nextFive, games[i])
		}
		return nextFive
	}

	nextFive = games[startIndex : startIndex+5]
	return nextFive
}

func findStartIndexNBA(games []objects.NBAGame) int {
	l, r := 0, len(games)-1
	var nextIndex int
	i := -1
	for l <= r {
		mid := l + (r-l)/2
		if isTodayNBA(games[mid]) {
			i = mid
			break
		} else if isBeforeTodayNBA(games[mid]) { // arr[i] < target
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

func isBeforeTodayNBA(game objects.NBAGame) bool {
	date := strings.Split(game.Date.Date, "T")[0]
	dateObj, _ := time.Parse("2006-01-02", date)
	todayString := strings.Split(time.Now().String(), " ")[0]
	formatToday, _ := time.Parse("2006-01-02", todayString)
	return dateObj.Before(formatToday)
}

func isTodayNBA(game objects.NBAGame) bool {
	date := strings.Split(game.Date.Date, "T")[0]
	dateObj, _ := time.Parse("2006-01-02", date)
	todayString := strings.Split(time.Now().String(), " ")[0]
	formatToday, _ := time.Parse("2006-01-02", todayString)
	return dateObj.Equal(formatToday)
}
