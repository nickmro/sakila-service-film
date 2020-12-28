package app

import (
	"errors"
	"sakila/sakila-film-service/sakila"
)

// FilmService is the film service.
type FilmService struct {
	ActorStore sakila.ActorStore
	FilmCache  sakila.FilmCache
	FilmStore  sakila.FilmStore
	Logger     sakila.Logger
}

const (
	defaultLimit = 20
)

// GetFilm returns the requested film.
func (s *FilmService) GetFilm(id int) (*sakila.Film, error) {
	film, err := s.FilmCache.GetFilm(id)
	if err != nil && !errors.Is(err, sakila.ErrorNotFound) {
		s.Logger.Error(err)
	}

	if film != nil {
		return film, nil
	}

	film, err = s.FilmStore.QueryFilm(id)
	if errors.Is(err, sakila.ErrorNotFound) {
		return nil, NewError(sakila.ErrorNotFound, "Film not found.")
	} else if err != nil {
		s.Logger.Error(err)

		return nil, NewError(sakila.ErrorInternal, "Internal error.")
	}

	actors, err := s.ActorStore.QueryFilmActors(id)
	if err != nil {
		s.Logger.Error(err)

		return nil, NewError(sakila.ErrorInternal, "Internal error.")
	}

	film.Actors = actors

	//nolint:errcheck
	go s.FilmCache.SetFilm(film)

	return film, nil
}

// GetFilms returns a list of films.
func (s *FilmService) GetFilms(params map[sakila.FilmQueryParam]interface{}) ([]*sakila.Film, error) {
	if _, ok := params[sakila.FilmQueryParamFirst].(int); !ok {
		params[sakila.FilmQueryParamFirst] = defaultLimit
	}

	films, err := s.FilmStore.QueryFilms(params)
	if err != nil {
		s.Logger.Error(err)

		return nil, NewError(sakila.ErrorInternal, "Internal error.")
	}

	for _, film := range films {
		actors, err := s.ActorStore.QueryFilmActors(film.FilmID)
		if err != nil {
			s.Logger.Error(err)

			return nil, NewError(sakila.ErrorInternal, "Internal error.")
		}

		film.Actors = actors
	}

	return films, nil
}
