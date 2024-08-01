package handling

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/tremerj/Sport-Companion/handling/scheduling"
)

func HandleSchedule() {
	if len(os.Args) != 2 {
		fmt.Println("Incorrect usage.")
		fmt.Println("Correct usage: Sport-Companion schedule")
		return
	}
	printSchedule()
}

func printSchedule() {
	result, err := getSchedule()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range result {
		fmt.Printf("%s", v)
	}
	fmt.Print("\033[A")
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
		return scheduling.NextFiveNHL(arr[1], arr[2])
	case "NFL":
		return scheduling.NextFiveNFL(arr[1], arr[2])
	case "MLB":
		return scheduling.NextFiveMLB(arr[1], arr[2])
	case "NBA":
		return scheduling.NextFiveNBA(arr[1], arr[2])
	default:
		return "", errors.New("Unknown behavior")
	}
}
