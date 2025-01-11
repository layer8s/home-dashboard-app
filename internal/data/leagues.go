package data

type League struct {
	//create league struct in json
	Id          int64 `json:"id"`
	LeagueId    int   `json:"league_id"`
	Year        int   `json:"year"`
	TeamCount   int   `json:"team_count"`
	CurrentWeek int   `json:"current_week"`
	NflWeek     int   `json:"nfl_week"`
}
