package sakila

import (
	"context"
	"time"
)

// Film represents a Sakila film.
type Film struct {
	FilmID             int       `json:"filmId"`
	Title              string    `json:"title"`
	Description        *string   `json:"description,omitempty"`
	ReleaseYear        *int      `json:"releaseYear,omitempty"`
	LanguageID         int       `json:"languageId"`
	OriginalLanguageID *int      `json:"originalLanguageId,omitempty"`
	RentalDuration     int       `json:"rentalDuration"`
	RentalRate         float64   `json:"rentalRate"`
	Length             *int      `json:"length,omitempty"`
	ReplacementCost    float64   `json:"replacementCost"`
	Rating             *string   `json:"rating,omitempty"`
	SpecialFeatures    []string  `json:"specialFeatures,omitempty"`
	Actors             []*Actor  `json:"actors"`
	LastUpdate         time.Time `json:"lastUpdate"`
}

// FilmParams are film query params.
type FilmParams struct {
	FilmIDs []int
	Limit   int
	Offset  int
}

// FilmService defines the interface for a film service.
type FilmService interface {
	GetFilm(ctx context.Context, filmID int) (*Film, error)
	GetFilms(ctx context.Context, params FilmParams) ([]*Film, error)
	GetFilmActors(ctx context.Context, filmIDs ...int) ([]*FilmActor, error)
}
