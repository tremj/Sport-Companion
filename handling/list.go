package handling

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func HandleList() {
	if len(os.Args) != 2 {
		fmt.Println("Incorrect usage.")
		fmt.Println("Correct usage: Sport-Companion list")
		return
	}
	printTeams()
}

func printTeams() {
	f, err := os.Open(".favourite_teams")
	if err != nil {
		fmt.Println("No teams have been added to your watchlist")
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		err := printTeam(line)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func printTeam(line string) error {
	split := strings.Split(line, ",")
	fmt.Println(split[2])
	return nil
}
