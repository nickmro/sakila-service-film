package sakila

import "time"

// Film represents a Sakila film.
type Film struct {
	FilmID             int       `json:"film_id"`
	Title              string    `json:"title"`
	Description        *string   `json:"description,omitempty"`
	ReleaseYear        *int      `json:"release_year,omitempty"`
	Actors             []*Actor  `json:"actors,omitempty"`
	LanguageID         int       `json:"language_id"`
	OriginalLanguageID *int      `json:"original_language_id,omitempty"`
	RentalDuration     int       `json:"rental_duration"`
	RentalRate         float64   `json:"rental_rate"`
	Length             *int      `json:"length,omitempty"`
	ReplacementCost    float64   `json:"replacement_cost"`
	Rating             *string   `json:"rating,omitempty"`
	SpecialFeatures    []uint8   `json:"special_features,omitempty"`
	LastUpdate         time.Time `json:"last_update"`
}

// FilmQueryParam is a film query parameter.
type FilmQueryParam string

// FilmQueryParams are film query parameters.
type FilmQueryParams map[FilmQueryParam]interface{}

// FilmService defines the interface for a film service.
type FilmService interface {
	GetFilm(id int) (*Film, error)
	GetFilms(params FilmQueryParams) ([]*Film, error)
}

// FilmStore defines the interface for film storage.
type FilmStore interface {
	QueryFilm(id int) (*Film, error)
	QueryFilms(params FilmQueryParams) ([]*Film, error)
}

// FilmCache defines the interface for the film cache.
type FilmCache interface {
	GetFilm(id int) (*Film, error)
	SetFilm(film *Film) error
}

const (
	// FilmQueryParamAfter indicates the ID of the record before the first one to return.
	FilmQueryParamAfter = FilmQueryParam("after")
	// FilmQueryParamFirst indicates how many records to return starting from the first matching record.
	FilmQueryParamFirst = FilmQueryParam("first")
	// FilmQueryParamCategory indicates the category of the films to return.
	FilmQueryParamCategory = FilmQueryParam("category")
)
