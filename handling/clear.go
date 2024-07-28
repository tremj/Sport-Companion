package handling

import (
	"fmt"
	"os"
)

func HandleClear() {
	err := Clear()
	if err != nil {
		fmt.Println(err)
	}
}

func Clear() error {
	f, err := os.OpenFile(".favourite_teams", os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
