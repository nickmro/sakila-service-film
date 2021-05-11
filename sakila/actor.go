package sakila

// Actor is a sakila film actor.
type Actor struct {
	ActorID int `json:"actorId"`
}

// FilmActor is a sakila film actor.
type FilmActor struct {
	Actor
	FilmID int
}
