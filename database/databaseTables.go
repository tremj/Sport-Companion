package database

type Users struct {
	ID       uint `gorm:"primary_key;auto_increment" json:"id"`
	Username string
	Password string
}

type SportTeam struct {
	ID       uint `gorm:"primary_key;auto_increment" json:"id"`
	Name     string
	Hometown string
}

type Match struct {
	ID   uint `gorm:"primary_key;auto_increment" json:"id"`
	Time string
}

type League struct {
	ID    uint `gorm:"primary_key;auto_increment" json:"id"`
	Name  string
	Sport string
}

type UserTeam struct {
	UserID uint
	TeamID uint
}

type MatchTeam struct {
	MatchID uint
	TeamID  uint
}

type LeagueTeam struct {
	LeagueID uint
	TeamID   uint
}
