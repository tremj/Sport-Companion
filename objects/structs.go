package objects

type TeamSearch struct {
	Results  int        `json:"results"`
	Response []Response `json:"response"`
}

type Response struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

const NHLLeagueID string = "57"

type NHLGameSearch struct {
	Results int       `json:"results"`
	Games   []NHLGame `json:"response"`
}

type NHLGame struct {
	Date  string `json:"date"`
	Time  string `json:"time"`
	Teams struct {
		Home struct {
			Name string `json:"home"`
		} `json:"home"`
		Away struct {
			Name string `json:"away"`
		}
	} `json:"teams"`
}

type NFLGameSearch struct {
	Results int       `json:"results"`
	Games   []NFLGame `json:"response"`
}

type NFLGame struct {
	Game struct {
		Date struct {
			Date string `json:"date"`
			Time string `json:"time"`
		} `json:"date"`
	} `json:"game"`
	Teams struct {
		Home struct {
			Name string `json:"name"`
		} `json:"home"`
		Away struct {
			Name string `json:"name"`
		} `json:"away"`
	} `json:"teams"`
}

const MLBLeagueID string = "1"

type MLBGameSearch struct {
	Results int       `json:"results"`
	Games   []MLBGame `json:"response"`
}

type MLBGame struct {
	Date  string `json:"date"`
	Time  string `json:"time"`
	Teams struct {
		Home struct {
			Name string `json:"name"`
		} `json:"home"`
		Away struct {
			Name string `json:"name"`
		} `json:"away"`
	} `json:"teams"`
}

type NBAGameSearch struct {
	Results int       `json:"results"`
	Games   []NBAGame `json:"response"`
}

type NBAGame struct {
	Date struct {
		Date string `json:"start"`
	} `json:"date"`
	Teams struct {
		Home struct {
			Name string `json:"name"`
		} `json:"home"`
		Away struct {
			Name string `json:"name"`
		} `json:"visitors"`
	} `json:"teams"`
}
