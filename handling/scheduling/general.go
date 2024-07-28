package scheduling

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func makeGenericGameRequest(url string, host string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-rapidapi-key", os.Getenv("SPORT_API_KEY"))
	req.Header.Add("x-rapidapi-host", host)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getDescriptor(day int) string {
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func backFourHours(time string) string {
	hourMin := strings.Split(time, ":")
	intHour, _ := strconv.Atoi(hourMin[0])
	intHour = (intHour - 4) % 24
	return strconv.Itoa(intHour) + ":" + hourMin[1]
}
