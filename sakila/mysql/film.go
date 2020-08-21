package mysql

import (
	"database/sql"

	"sakila/sakila-film-service/sakila"
)

// FilmDB is a connection to film database.
type FilmDB struct {
	*sql.DB
}

const (
	filmsQuery = `
		SELECT
			film.film_id,
			film.title,
			film.description,
			film.release_year,
			film.language_id,
			film.original_language_id,
			film.rental_duration,
			film.rental_rate,
			film.length,
			film.replacement_cost,
			film.rating,
			film.special_features,
			film.last_update
		FROM film
	`

	filmQuery = filmsQuery + `
		WHERE film_id = ?
	`
)

// QueryFilm returns the film with the given ID or an error.
func (db *FilmDB) QueryFilm(id int) (*sakila.Film, error) {
	var film sakila.Film

	err := db.DB.QueryRow(filmQuery, id).Scan(
		&film.FilmID,
		&film.Title,
		&film.Description,
		&film.ReleaseYear,
		&film.LanguageID,
		&film.OriginalLanguageID,
		&film.RentalDuration,
		&film.RentalRate,
		&film.Length,
		&film.ReplacementCost,
		&film.Rating,
		&film.SpecialFeatures,
		&film.LastUpdate,
	)

	if err == sql.ErrNoRows {
		return nil, sakila.ErrorNotFound
	} else if err != nil {
		return nil, err
	}

	return &film, nil
}

// QueryFilms returns the films to query or an error.
func (db *FilmDB) QueryFilms(params map[sakila.FilmQueryParam]interface{}) ([]*sakila.Film, error) {
	films := []*sakila.Film{}

	query, args := filmQueryForParams(params)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return films, err
	}

	for rows.Next() {
		var film sakila.Film

		if err := rows.Scan(
			&film.FilmID,
			&film.Title,
			&film.Description,
			&film.ReleaseYear,
			&film.LanguageID,
			&film.OriginalLanguageID,
			&film.RentalDuration,
			&film.RentalRate,
			&film.Length,
			&film.ReplacementCost,
			&film.Rating,
			&film.SpecialFeatures,
			&film.LastUpdate,
		); err != nil {
			return nil, err
		}

		films = append(films, &film)
	}

	return films, nil
}

func filmQueryForParams(params map[sakila.FilmQueryParam]interface{}) (query string, args []interface{}) {
	query = filmsQuery

	if category, ok := params[sakila.FilmQueryParamCategory].(string); ok {
		args = append(args, category)
		query += `
			INNER JOIN film_category ON film_category.film_id = film.film_id
			INNER JOIN category ON category.category_id = film_category.category_id
			WHERE category.name = ?`
	}

	if after, ok := params[sakila.FilmQueryParamAfter].(int); ok {
		args = append(args, after)
		query += ` WHERE film.film_id > ?`
	}

	query += ` ORDER BY film.film_id ASC`

	if first, ok := params[sakila.FilmQueryParamFirst].(int); ok {
		args = append(args, first)
		query += ` LIMIT ?`
	}

	return query, args
}
