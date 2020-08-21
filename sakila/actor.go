package sakila

import "time"

// Actor is a sakila film actor.
type Actor struct {
	ActorID    int       `json:"actor_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	LastUpdate time.Time `json:"last_update"`
}

// ActorStore defines the operations that may be performed on an actor store.
type ActorStore interface {
	QueryFilmActors(filmID int) ([]*Actor, error)
}
