package sakila

import "time"

// Actor is a sakila film actor.
type Actor struct {
	ActorID    int       `json:"actor_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	LastUpdate time.Time `json:"last_update"`
}

// FilmActor is a sakila film actor.
type FilmActor struct {
	Actor
	FilmID int
}

// FilmActorParams are the params for film actors.
type FilmActorParams struct {
	FilmIDs []int
}
