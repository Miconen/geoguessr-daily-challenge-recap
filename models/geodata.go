package models

type GeoData struct {
	PlaceID     int      `json:"place_id"`
	Licence     string   `json:"licence"`
	OsmType     string   `json:"osm_type"`
	OsmID       int      `json:"osm_id"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	Category    string   `json:"category"`
	Type        string   `json:"type"`
	PlaceRank   int      `json:"place_rank"`
	Importance  float64  `json:"importance"`
	Addresstype string   `json:"addresstype"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Address     Address  `json:"address"`
	Boundingbox []string `json:"boundingbox"`
}

type Address struct {
	Locality     string `json:"locality"`
	Village      string `json:"village"`
	County       string `json:"county"`
	Province     string `json:"province"`
	ISO31662Lvl6 string `json:"ISO3166-2-lvl6"`
	State        string `json:"state"`
	ISO31662Lvl4 string `json:"ISO3166-2-lvl4"`
	Postcode     string `json:"postcode"`
	Country      string `json:"country"`
	CountryCode  string `json:"country_code"`
}

// Guesses
type RoundGeoData struct {
	RoundNumber string
	Location    GeoData
}

type GuessGeoData struct {
	RoundNumber int
	Guess       GeoData
}

type PlayerGeoData struct {
	PlayerID string
	Rounds   []GuessGeoData
}

type GameGeoData struct {
	ActualLocations []RoundGeoData
	PlayerGuesses   []PlayerGeoData
}
