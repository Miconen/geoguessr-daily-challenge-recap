package utils

import "github.com/Miconen/geoguessr-daily-challenge-recap/models"

func GetCountryCode(data models.GeoData) string {
	return data.Address.CountryCode
}
