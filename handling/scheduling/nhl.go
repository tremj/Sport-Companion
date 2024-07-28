package scheduling

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/tremerj/Sport-Companion/objects"
)

func NextFiveNHL(id string) (string, error) {
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

	formattedOutput := formattedNHLGames(games.Games)
	return formattedOutput, nil
}

func gameNHLURL(id string, back int) string {
	return "https://v1.hockey.api-sports.io/games?league=" + objects.NHLLeagueID + "&team=" + id + "&season=" + strconv.Itoa(time.Now().Year()-back)
}

func formattedNHLGames(games []objects.NHLGame) string {
	upcoming := getNextFiveGamesNHL(games)
	return formatNHLGame(upcoming)
}

func formatNHLGame(games []objects.NHLGame) string {
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
	return weekDay + " " + month + " " + strconv.Itoa(day) + dateDescriptor + " @ " + timeBackFourHours + " >>> " + g.Teams.Home.Name + " vs. " + g.Teams.Away.Name + "\n"
}

func getNextFiveGamesNHL(games []objects.NHLGame) []objects.NHLGame {
	startIndex := findStartIndexNHL(games)
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

func findStartIndexNHL(games []objects.NHLGame) int {
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
