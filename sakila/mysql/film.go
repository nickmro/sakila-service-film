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
	var specialFeatures string

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
		&specialFeatures,
		&film.LastUpdate,
	)

	film.SpecialFeatures = strings.Split(specialFeatures, ",")

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
		var specialFeatures string

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
			&specialFeatures,
			&film.LastUpdate,
		); err != nil {
			return nil, err
		}

		film.SpecialFeatures = strings.Split(specialFeatures, ",")

		films = append(films, &film)
	}

	return films, nil
}

// QueryFilmActors returns a film's actors.
func (db *FilmDB) QueryFilmActors(params sakila.FilmActorParams) ([]*sakila.FilmActor, error) {
	actors := []*sakila.FilmActor{}

	query, args := filmActorsQueryForParams(params)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	} else if rows.Err() != nil {
		return nil, rows.Err()
	}

	defer rows.Close() //nolint:errcheck

	for rows.Next() {
		var actor sakila.FilmActor

		err := rows.Scan(
			&actor.FilmID,
			&actor.ActorID,
			&actor.FirstName,
			&actor.LastName,
			&actor.LastUpdate,
		)
		if err != nil {
			return nil, err
		}

		actors = append(actors, &actor)
	}

	return actors, nil
}

// QueryFilmCategories returns film categories.
func (db *FilmDB) QueryFilmCategories(params sakila.FilmCategoryParams) ([]*sakila.FilmCategory, error) {
	categories := []*sakila.FilmCategory{}

	query, args := filmCategoriesQueryForParams(params)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	} else if rows.Err() != nil {
		return nil, rows.Err()
	}

	defer rows.Close() //nolint:errcheck

	for rows.Next() {
		var category sakila.FilmCategory

		err := rows.Scan(
			&category.FilmID,
			&category.CategoryID,
			&category.Name,
			&category.LastUpdate,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil
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

func filmActorsQueryForParams(params sakila.FilmActorParams) (query string, args []interface{}) {
	query = `SELECT
			film_actor.film_id,
			actor.actor_id,
			actor.first_name,
			actor.last_name,
			actor.last_update
		FROM actor
	`
	joins := []string{}
	wheres := []string{}

	if len(params.FilmIDs) > 0 {
		joins = append(joins, `INNER JOIN film_actor ON film_actor.actor_id = actor.actor_id`)

		wheres = append(wheres, fmt.Sprintf(`film_actor.film_id IN (%s)`, argString(len(params.FilmIDs))))

		for i := range params.FilmIDs {
			args = append(args, params.FilmIDs[i])
		}
	}

	if len(joins) > 0 {
		query += ` ` + strings.Join(joins, ` `)
	}

	if len(wheres) > 0 {
		query += fmt.Sprintf(` WHERE %s`, strings.Join(wheres, ` AND `))
	}

	return query, args
}

func filmCategoriesQueryForParams(params sakila.FilmCategoryParams) (query string, args []interface{}) {
	columns := []string{
		"film_category.film_id",
		"category.category_id",
		"category.name",
		"category.last_update",
	}
	table := "category"
	joins := []string{}
	wheres := []string{}
	whereArgs := []interface{}{}

	if len(params.FilmIDs) > 0 {
		joins = append(joins, "INNER JOIN film_category ON film_category.category_id = category.category_id")
		wheres = append(wheres, fmt.Sprintf("film_category.film_id IN (%s)", argString(len(params.FilmIDs))))

		for i := range params.FilmIDs {
			whereArgs = append(whereArgs, i)
		}
	}

	query = "SELECT " + strings.Join(columns, ", ") + " FROM " + table

	if len(joins) > 0 {
		query += " " + strings.Join(joins, " ")
	}

	if len(wheres) > 0 {
		query += " WHERE " + strings.Join(wheres, " AND ")

		args = append(args, whereArgs...)
	}

	return query, args
}
