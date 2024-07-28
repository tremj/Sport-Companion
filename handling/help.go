package handling

import "fmt"

func HandleHelp() {
	fmt.Println("This CLI supports the MLB, NBA, NFL and NHL")
	fmt.Println("To add/remove teams from these leagues to your watchlist use the CLI as such: Sport-Companion <add|remove> <league> <team name>")
	fmt.Println("If you want to remove all your teams run this: Sport-Companion clear")
	fmt.Println("To see your teams run this command: Sport-Companion list")
	fmt.Println("To see the next 5 games of each of your teams run this: Sport-Companion schedule")
}
