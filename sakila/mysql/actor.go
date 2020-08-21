package mysql

import (
	"database/sql"
	"sakila/sakila-film-service/sakila"
)

// ActorDB is a mysql actor store.
type ActorDB struct {
	*sql.DB
}

const (
	actorsQuery = `
		SELECT
			actor.actor_id,
			actor.first_name,
			actor.last_name,
			actor.last_update
		FROM actor
	`

	filmActorsQuery = actorsQuery + `
		INNER JOIN film_actor
			ON film_actor.actor_id = actor.actor_id
		WHERE film_actor.film_id = ?
	`
)

// QueryFilmActors returns a film's actors.
func (db *ActorDB) QueryFilmActors(filmID int) ([]*sakila.Actor, error) {
	actors := []*sakila.Actor{}

	rows, err := db.DB.Query(filmActorsQuery, filmID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var actor sakila.Actor

		err := rows.Scan(
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
