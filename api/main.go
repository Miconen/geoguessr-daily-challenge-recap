package api

import "fmt"

var Endpoint = "https://www.geoguessr.com/api/v3/"
var EndpointDaily = Endpoint + "challenges/daily-challenges/today"

func GetScoresEndpoint(id string) string {
	return fmt.Sprintf(Endpoint+"results/highscores/%s?friends=false&limit=26&minRounds=5&club=true", id)
}
