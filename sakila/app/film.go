package app

import (
	"errors"
	"sakila/sakila-film-service/sakila"
)

// FilmService is the film service.
type FilmService struct {
	Cache  sakila.FilmCache
	Store  sakila.FilmStore
	Logger sakila.Logger
}

const (
	defaultLimit = 20
)

// GetFilm returns the requested film.
func (s *FilmService) GetFilm(id int) (*sakila.Film, error) {
	film, err := s.Cache.GetFilm(id)
	if err != nil && !errors.Is(err, sakila.ErrorNotFound) {
		s.Logger.Error(err)
	}

	if film != nil {
		return film, nil
	}

	film, err = s.Store.QueryFilm(id)
	if errors.Is(err, sakila.ErrorNotFound) {
		return nil, NewError(sakila.ErrorNotFound, "Film not found.")
	} else if err != nil {
		s.Logger.Error(err)

		return nil, NewError(sakila.ErrorInternal, "Internal error.")
	}

	//nolint:errcheck
	go s.Cache.SetFilm(film)

	return film, nil
}

// GetFilms returns a list of films.
func (s *FilmService) GetFilms(params sakila.FilmQueryParams) ([]*sakila.Film, error) {
	if _, ok := params[sakila.FilmQueryParamLimit].(int); !ok {
		params[sakila.FilmQueryParamLimit] = defaultLimit
	}

	films, err := s.Cache.GetFilms(params)
	if err != nil && !errors.Is(err, sakila.ErrorNotFound) {
		s.Logger.Error(err)
	}

	if films != nil {
		return films, nil
	}

	films, err = s.Store.QueryFilms(params)
	if err != nil {
		s.Logger.Error(err)

		return nil, NewError(sakila.ErrorInternal, "Internal error.")
	}

	go s.cacheFilms(films, params)

	return films, nil
}

// GetFilmCategories returns a film's categories.
func (s *FilmService) GetFilmCategories(params sakila.FilmCategoryParams) ([]*sakila.FilmCategory, error) {
	categories, err := s.Cache.GetFilmCategories(params)
	if err != nil && !errors.Is(err, sakila.ErrorNotFound) {
		s.Logger.Error(err)
	}

	if categories != nil {
		return categories, nil
	}

	categories, err = s.Store.QueryFilmCategories(params)
	if err != nil {
		s.Logger.Error(err)

		return nil, sakila.ErrorInternal
	}

	go s.cacheFilmCategories(categories, params)

	return categories, nil
}

// GetFilmActors returns a film's actors.
func (s *FilmService) GetFilmActors(params sakila.FilmActorParams) ([]*sakila.FilmActor, error) {
	actors, err := s.Cache.GetFilmActors(params)
	if err != nil && !errors.Is(err, sakila.ErrorNotFound) {
		s.Logger.Error(err)
	}

	if actors != nil {
		return actors, nil
	}

	actors, err = s.Store.QueryFilmActors(params)
	if err != nil {
		s.Logger.Error(err)

		return nil, sakila.ErrorInternal
	}

	go s.cacheFilmActors(actors, params)

	return actors, nil
}

func (s *FilmService) cacheFilms(films []*sakila.Film, params sakila.FilmQueryParams) {
	if err := s.Cache.SetFilms(films, params); err != nil {
		s.Logger.Error(err)
	}
}

func (s *FilmService) cacheFilmActors(actors []*sakila.FilmActor, params sakila.FilmActorParams) {
	if err := s.Cache.SetFilmActors(actors, params); err != nil {
		s.Logger.Error(err)
	}
}

func (s *FilmService) cacheFilmCategories(categories []*sakila.FilmCategory, params sakila.FilmCategoryParams) {
	if err := s.Cache.SetFilmCategories(categories, params); err != nil {
		s.Logger.Error(err)
	}
}
