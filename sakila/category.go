package sakila

import "time"

// Category is a sakila film category.
type Category struct {
	CategoryID int       `json:"category_id"`
	Name       string    `json:"name"`
	LastUpdate time.Time `json:"last_update"`
}

// FilmCategory is a sakila film category.
type FilmCategory struct {
	Category
	FilmID int
}

// FilmCategoryParams are the params for film categories.
type FilmCategoryParams struct {
	FilmIDs []int
}
