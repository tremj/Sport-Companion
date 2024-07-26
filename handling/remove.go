package handling

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func HandleRemove() {
	id, err := getTeamID(os.Args[2])
	if err != nil {
		fmt.Println(err)
	}
	removeFromFile(strings.ToUpper(os.Args[2]) + "," + id)
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
