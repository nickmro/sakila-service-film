package mysql

import (
	"context"
	"database/sql"
	"errors"
	"sakila/sakila-film-service/sakila"
	"sakila/sakila-film-service/sakila/mysql/sqlbuilder"
	"strings"
)

// FilmService is a film service backed by a MySQL DB.
type FilmService struct {
	DB     *DB
	Logger sakila.Logger
}

// GetFilm returns a film.
func (service *FilmService) GetFilm(ctx context.Context, filmID int) (*sakila.Film, error) {
	var film sakila.Film
	var specialFeatures string

	query, args := filmQueryForParams(sakila.FilmParams{
		FilmIDs: []int{filmID},
		Limit:   1,
	})

	err := service.DB.QueryRow(query, args...).Scan(
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
		service.logError(err)
		return nil, sakila.ErrorInternal
	}

	return &film, nil
}

// GetFilms returns the films.
func (service *FilmService) GetFilms(ctx context.Context, params sakila.FilmParams) ([]*sakila.Film, error) {
	films := []*sakila.Film{}

	query, args := filmQueryForParams(params)

	rows, err := service.DB.Query(query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return films, nil
	} else if err != nil {
		service.logError(err)
		return nil, sakila.ErrorInternal
	} else if err := rows.Err(); err != nil {
		service.logError(err)
		return nil, sakila.ErrorInternal
	}

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
			service.logError(err)
			return nil, sakila.ErrorInternal
		}

		film.SpecialFeatures = strings.Split(specialFeatures, ",")

		films = append(films, &film)
	}

	if err := rows.Close(); err != nil {
		service.logError(err)
	}

	return films, nil
}

// GetFilmActors returns a film's actors.
func (service *FilmService) GetFilmActors(ctx context.Context, filmIDs ...int) ([]*sakila.FilmActor, error) {
	actors := []*sakila.FilmActor{}

	stmt := sqlbuilder.Select(
		"film_id",
		"actor_id",
	).
		From("film_actor")

	stmt.Where("film_actor.film_id IN (%v)", formattedIDs(filmIDs)...)

	query, args := stmt.Build()

	rows, err := service.DB.Query(query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return actors, nil
	} else if err != nil {
		service.logError(err)
		return nil, sakila.ErrorInternal
	} else if err := rows.Err(); err != nil {
		service.logError(err)
		return nil, sakila.ErrorInternal
	}

	for rows.Next() {
		var actor sakila.FilmActor

		err := rows.Scan(
			&actor.FilmID,
			&actor.ActorID,
		)
		if err != nil {
			return nil, err
		}

		actors = append(actors, &actor)
	}

	if err := rows.Close(); err != nil {
		service.logError(err)
	}

	return actors, nil
}

func (service *FilmService) logError(err error) {
	if logger := service.Logger; logger != nil {
		logger.Error(err)
	}
}

func filmQueryForParams(params sakila.FilmParams) (query string, args []interface{}) {
	stmt := sqlbuilder.Select(
		"film.film_id",
		"film.title",
		"film.description",
		"film.release_year",
		"film.language_id",
		"film.original_language_id",
		"film.rental_duration",
		"film.rental_rate",
		"film.length",
		"film.replacement_cost",
		"film.rating",
		"film.special_features",
		"film.last_update",
	).
		From("film")

	if ids := params.FilmIDs; len(ids) > 0 {
		stmt.Where("film.film_id IN (%v)", formattedIDs(ids)...)
	}

	if limit := params.Limit; limit > 0 {
		stmt.Limit(limit)
	}

	if offset := params.Offset; offset > 0 {
		stmt.Offset(offset)
	}

	return stmt.Build()
}

func formattedIDs(ids []int) []interface{} {
	formattedIDs := make([]interface{}, len(ids))
	for i := range ids {
		formattedIDs[i] = ids[i]
	}

	return formattedIDs
}
