package handling

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func HandleRemove() {
	if len(os.Args) != 3 {
		fmt.Println("Incorrect usage.")
		fmt.Println("Hint: Sport-Companion remove <teamname>")
		return
	}
	removeFromFile(os.Args[2])
}

func removeFromFile(aFavourite string) {
	f, err := os.Open(".favourite_teams")
	if err != nil {
		fmt.Println("No teams have been added to your watchlist")
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	detected := false
	var lines []string
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		if line[2] != aFavourite {
			lines = append(lines, strings.Join(line, ","))
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
	if err != nil {
		fmt.Println("Error creating favourite team file")
		return
	}
	defer g.Close()

	for _, line := range lines {
		_, err = g.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing to the file")
			return
		}
	}
}
