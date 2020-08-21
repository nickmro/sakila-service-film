package graphql_test

import (
	"sakila/sakila-film-service/sakila"
	"time"
)

type Film struct {
	FilmID             int             `json:"filmId"`
	Title              string          `json:"title"`
	Description        *string         `json:"description,omitempty"`
	ReleaseYear        *int            `json:"releaseYear,omitempty"`
	Actors             []*sakila.Actor `json:"actors,omitempty"`
	LanguageID         int             `json:"languageId"`
	OriginalLanguageID *int            `json:"originalLanguageId,omitempty"`
	RentalDuration     int             `json:"rentalDuration"`
	RentalRate         float64         `json:"rentalRate"`
	Length             *int            `json:"length,omitempty"`
	ReplacementCost    float64         `json:"replacementCost"`
	Rating             *string         `json:"rating,omitempty"`
	SpecialFeatures    []uint8         `json:"specialFeatures,omitempty"`
	LastUpdate         time.Time       `json:"lastUpdate"`
}
