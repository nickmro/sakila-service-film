package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"sakila/sakila-film-service/sakila"
	"strings"
)

// FilmDB is a connection to film database.
type FilmDB struct {
	*DB
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

	if errors.Is(err, sql.ErrNoRows) {
		return nil, sakila.ErrorNotFound
	} else if err != nil {
		return nil, err
	}

	return &film, nil
}

// QueryFilms returns the films to query or an error.
func (db *FilmDB) QueryFilms(params sakila.FilmQueryParams) ([]*sakila.Film, error) {
	films := []*sakila.Film{}

	query, args := filmQueryForParams(params)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	} else if err := rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close() //nolint:errcheck

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

func filmQueryForParams(params sakila.FilmQueryParams) (query string, args []interface{}) {
	query = filmsQuery
	wheres := []string{}

	if category, ok := params[sakila.FilmQueryParamCategory].(string); ok {
		args = append(args, category)

		query += `
			INNER JOIN film_category ON film_category.film_id = film.film_id
			INNER JOIN category ON category.category_id = film_category.category_id`

		wheres = append(wheres, `category.name = ?`)
	}

	if len(wheres) > 0 {
		query += fmt.Sprintf(` WHERE %s`, strings.Join(wheres, ` AND `))
	}

	query += ` ORDER BY film.film_id ASC`

	if limit, ok := params[sakila.FilmQueryParamLimit].(int); ok {
		args = append(args, limit)

		query += ` LIMIT ?`
	}

	if offset, ok := params[sakila.FilmQueryParamOffset].(int); ok {
		args = append(args, offset)

		query += ` OFFSET ?`
	}

	return query, args
}
