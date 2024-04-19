package endpoints

import (
	"fmt"
	"io"
	"net/http"
)

func makeSoccerRequest(userSpec string) (string, bool) {
	errorString := "Something went wrong."
	URL := "https://v1.football.api-sports.io//" + userSpec
	method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		fmt.Println(err)
		return errorString, false
	}

	req.Header.Add("x-rapidapi-key", "5c027afffa1013a9baac52513787e122")
	req.Header.Add("x-rapidapi-host", URL)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return errorString, false
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return errorString, false
	}
	return string(body), true
}
