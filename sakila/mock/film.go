package mock

import (
	"context"
	"sakila/sakila-film-service/sakila"
)

// FilmService is a mock film service.
type FilmService struct {
	GetFilmFn       func(ctx context.Context, filmID int) (*sakila.Film, error)
	GetFilmsFn      func(ctx context.Context, params sakila.FilmParams) ([]*sakila.Film, error)
	GetFilmActorsFn func(ctx context.Context, filmIDs ...int) ([]*sakila.FilmActor, error)
}

// GetFilm runs the mock function or returns an empty film.
func (s *FilmService) GetFilm(ctx context.Context, filmID int) (*sakila.Film, error) {
	if fn := s.GetFilmFn; fn != nil {
		return fn(ctx, filmID)
	}

	return &sakila.Film{}, nil
}

// GetFilms runs the mock function or returns an empty slice of films.
func (s *FilmService) GetFilms(ctx context.Context, params sakila.FilmParams) ([]*sakila.Film, error) {
	if fn := s.GetFilmsFn; fn != nil {
		return fn(ctx, params)
	}

	return []*sakila.Film{}, nil
}

// GetFilmActors runs the mock function or returns an empty slice of film actors.
func (s *FilmService) GetFilmActors(ctx context.Context, filmIDs ...int) ([]*sakila.FilmActor, error) {
	if fn := s.GetFilmActorsFn; fn != nil {
		return fn(ctx, filmIDs...)
	}

	return []*sakila.FilmActor{}, nil
}
