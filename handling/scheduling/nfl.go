package scheduling

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/tremerj/Sport-Companion/objects"
)

func NextFiveNFL(id string) (string, error) {
	url := gameNFLURL(id, 0)
	body, err := makeGenericGameRequest(url, "https://v1.american-football.api-sports.io")
	if err != nil {
		return "", err
	}

	var games objects.NFLGameSearch
	err = json.Unmarshal(body, &games)
	if err != nil {
		return "", err
	}

	if games.Results == 0 {
		url = gameNFLURL(id, 1)
		body, err = makeGenericGameRequest(url, "https://v1.american-football.api-sports.io")
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

	formattedOutput := formattedNFLGames(games.Games)
	return formattedOutput, nil
}

func gameNFLURL(id string, back int) string {
	return "https://v1.american-football.api-sports.io/games?team=" + id + "&season=" + strconv.Itoa(time.Now().Year()-back)
}

func formattedNFLGames(games []objects.NFLGame) string {
	upcoming := getNextFiveGamesNFL(games)
	return formatNFLGame(upcoming)
}

func formatNFLGame(games []objects.NFLGame) string {
	var res string
	for _, g := range games {
		res = res + getNFLFormat(g)
	}
	return res
}

func getNFLFormat(g objects.NFLGame) string {
	// Format: "Monday September 25th @ HH:MM >>> <Home Team> vs. <Away Team>
	dateObj, _ := time.Parse("2006-01-02", g.Game.Date.Date)
	weekDay := dateObj.Weekday().String()
	month := dateObj.Month().String()
	day := dateObj.Day()
	dateDescriptor := getDescriptor(day)
	timeBackFourHours := backFourHours(g.Game.Date.Time)
	return weekDay + " " + month + " " + strconv.Itoa(day) + dateDescriptor + " @ " + timeBackFourHours + " >>> " + g.Teams.Home.Name + " vs. " + g.Teams.Away.Name + "\n"
}

func getNextFiveGamesNFL(games []objects.NFLGame) []objects.NFLGame {
	startIndex := findStartIndexNFL(games)
	var nextFive []objects.NFLGame
	if len(games)-startIndex < 5 {
		for i := startIndex; i < len(games); i++ {
			nextFive = append(nextFive, games[i])
		}
		return nextFive
	}

	nextFive = games[startIndex : startIndex+5]
	return nextFive
}

func findStartIndexNFL(games []objects.NFLGame) int {
	l, r := 0, len(games)-1
	var nextIndex int
	i := -1
	for l <= r {
		mid := l + (r-l)/2
		if isTodayNFL(games[mid]) {
			i = mid
			break
		} else if isBeforeTodayNFL(games[mid]) { // arr[i] < target
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

func isBeforeTodayNFL(game objects.NFLGame) bool {
	dateObj, _ := time.Parse("2006-01-02", game.Game.Date.Date)
	todayString := strings.Split(time.Now().String(), " ")[0]
	formatToday, _ := time.Parse("2006-01-02", todayString)
	return dateObj.Before(formatToday)
}

func isTodayNFL(game objects.NFLGame) bool {
	dateObj, _ := time.Parse("2006-01-02", game.Game.Date.Date)
	todayString := strings.Split(time.Now().String(), " ")[0]
	formatToday, _ := time.Parse("2006-01-02", todayString)
	return dateObj.Equal(formatToday)
}
