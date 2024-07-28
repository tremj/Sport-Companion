package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/tremerj/Sport-Companion/handling"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please use one of the valid commands")
		return
	}

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading env variables")
		return
	}

	switch os.Args[1] {
	case "add":
		handling.HandleAdd()
	case "clear":
		handling.HandleClear()
	case "help":
		handling.HandleHelp()
	case "list":
		handling.HandleList()
	case "schedule":
		handling.HandleSchedule()
	case "remove":
		handling.HandleRemove()
	default:
		fmt.Println("Unknown command")
	}
}
