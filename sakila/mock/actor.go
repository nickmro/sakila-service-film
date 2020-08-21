package mock

import "sakila/sakila-film-service/sakila"

// ActorStore is a mock actor store.
type ActorStore struct {
	QueryFilmActorsFn func(filmID int) ([]*sakila.Actor, error)
}

// QueryFilmActors runs the mock function or an empty slice of actors.
func (s *ActorStore) QueryFilmActors(filmID int) ([]*sakila.Actor, error) {
	if fn := s.QueryFilmActorsFn; fn != nil {
		return fn(filmID)
	}

	return []*sakila.Actor{}, nil
}
