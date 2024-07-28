package scheduling

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/tremerj/Sport-Companion/objects"
)

func NextFiveMLB(id string) (string, error) {
	url := gameMLBURL(id, 0)
	body, err := makeGenericGameRequest(url, "https://v1.baseball.api-sports.io")
	if err != nil {
		return "", err
	}

	var games objects.MLBGameSearch
	err = json.Unmarshal(body, &games)
	if err != nil {
		return "", err
	}

	if games.Results == 0 {
		url = gameMLBURL(id, 1)
		body, err = makeGenericGameRequest(url, "https://v1.baseball.api-sports.io")
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

	formattedOutput := formattedMLBGames(games.Games)
	return formattedOutput, nil
}

func gameMLBURL(id string, back int) string {
	return "https://v1.baseball.api-sports.io/games?league=" + objects.MLBLeagueID + "&team=" + id + "&season=" + strconv.Itoa(time.Now().Year()-back)
}

func formattedMLBGames(games []objects.MLBGame) string {
	upcoming := getNextFiveGamesMLB(games)
	return formatMLBGame(upcoming)
}

func formatMLBGame(games []objects.MLBGame) string {
	var res string
	for _, g := range games {
		res = res + getMLBFormat(g)
	}
	return res
}

func getMLBFormat(g objects.MLBGame) string {
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

func getNextFiveGamesMLB(games []objects.MLBGame) []objects.MLBGame {
	startIndex := findStartIndexMLB(games)
	var nextFive []objects.MLBGame
	if len(games)-startIndex < 5 {
		for i := startIndex; i < len(games); i++ {
			nextFive = append(nextFive, games[i])
		}
		return nextFive
	}

	nextFive = games[startIndex : startIndex+5]
	return nextFive
}

func findStartIndexMLB(games []objects.MLBGame) int {
	l, r := 0, len(games)-1
	var nextIndex int
	i := -1
	for l <= r {
		mid := l + (r-l)/2
		if isTodayMLB(games[mid]) {
			i = mid
			break
		} else if isBeforeTodayMLB(games[mid]) { // arr[i] < target
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

func isBeforeTodayMLB(game objects.MLBGame) bool {
	date := strings.Split(game.Date, "T")[0]
	dateObj, _ := time.Parse("2006-01-02", date)
	todayString := strings.Split(time.Now().String(), " ")[0]
	formatToday, _ := time.Parse("2006-01-02", todayString)
	return dateObj.Compare(formatToday) == -1
}

func isTodayMLB(game objects.MLBGame) bool {
	date := strings.Split(game.Date, "T")[0]
	dateObj, _ := time.Parse("2006-01-02", date)
	todayString := strings.Split(time.Now().String(), " ")[0]
	formatToday, _ := time.Parse("2006-01-02", todayString)
	return dateObj.Equal(formatToday)
}
