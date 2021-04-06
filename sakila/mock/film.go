package mock

import "sakila/sakila-film-service/sakila"

// FilmService is a mock film service.
type FilmService struct {
	GetFilmFn           func(filmID int) (*sakila.Film, error)
	GetFilmsFn          func(params sakila.FilmQueryParams) ([]*sakila.Film, error)
	GetFilmActorsFn     func(params sakila.FilmActorParams) ([]*sakila.FilmActor, error)
	GetFilmCategoriesFn func(params sakila.FilmCategoryParams) ([]*sakila.FilmCategory, error)
}

// FilmStore is a mock film store.
type FilmStore struct {
	QueryFilmFn           func(filmID int) (*sakila.Film, error)
	QueryFilmsFn          func(params sakila.FilmQueryParams) ([]*sakila.Film, error)
	QueryFilmActorsFn     func(params sakila.FilmActorParams) ([]*sakila.FilmActor, error)
	QueryFilmCategoriesFn func(params sakila.FilmCategoryParams) ([]*sakila.FilmCategory, error)
}

// FilmCache is a mock film cache.
type FilmCache struct {
	GetFilmFn           func(filmID int) (*sakila.Film, error)
	SetFilmFn           func(film *sakila.Film) error
	GetFilmsFn          func(params sakila.FilmQueryParams) ([]*sakila.Film, error)
	SetFilmsFn          func(films []*sakila.Film, params sakila.FilmQueryParams) error
	GetFilmActorsFn     func(params sakila.FilmActorParams) ([]*sakila.FilmActor, error)
	SetFilmActorsFn     func(actors []*sakila.FilmActor, params sakila.FilmActorParams) error
	GetFilmCategoriesFn func(params sakila.FilmCategoryParams) ([]*sakila.FilmCategory, error)
	SetFilmCategoriesFn func(categories []*sakila.FilmCategory, params sakila.FilmCategoryParams) error
}

// GetFilm runs the mock function or returns an empty film.
func (s *FilmService) GetFilm(filmID int) (*sakila.Film, error) {
	if fn := s.GetFilmFn; fn != nil {
		return fn(filmID)
	}

	return &sakila.Film{}, nil
}

// GetFilms runs the mock function or returns an empty slice of films.
func (s *FilmService) GetFilms(params sakila.FilmQueryParams) ([]*sakila.Film, error) {
	if fn := s.GetFilmsFn; fn != nil {
		return fn(params)
	}

	return []*sakila.Film{}, nil
}

// GetFilmActors runs the mock function or returns an empty slice of film actors.
func (s *FilmService) GetFilmActors(params sakila.FilmActorParams) ([]*sakila.FilmActor, error) {
	if fn := s.GetFilmActorsFn; fn != nil {
		return fn(params)
	}

	return []*sakila.FilmActor{}, nil
}

// GetFilmCategories runs the mock function or returns an empty slice of film categories.
func (s *FilmService) GetFilmCategories(params sakila.FilmCategoryParams) ([]*sakila.FilmCategory, error) {
	if fn := s.GetFilmCategoriesFn; fn != nil {
		return fn(params)
	}

	return []*sakila.FilmCategory{}, nil
}

// QueryFilm runs the mock function or returns an empty film.
func (s *FilmStore) QueryFilm(filmID int) (*sakila.Film, error) {
	if fn := s.QueryFilmFn; fn != nil {
		return fn(filmID)
	}

	return &sakila.Film{}, nil
}

// QueryFilms runs the mock function or returns an empty slice of films.
func (s *FilmStore) QueryFilms(params sakila.FilmQueryParams) ([]*sakila.Film, error) {
	if fn := s.QueryFilmsFn; fn != nil {
		return fn(params)
	}

	return []*sakila.Film{}, nil
}

// QueryFilmActors runs the mock function or returns an empty slice of actors.
func (s *FilmStore) QueryFilmActors(params sakila.FilmActorParams) ([]*sakila.FilmActor, error) {
	if fn := s.QueryFilmActorsFn; fn != nil {
		return fn(params)
	}

	return []*sakila.FilmActor{}, nil
}

// QueryFilmCategories runs the mock function or returns an empty slice of categories.
func (s *FilmStore) QueryFilmCategories(params sakila.FilmCategoryParams) ([]*sakila.FilmCategory, error) {
	if fn := s.QueryFilmCategoriesFn; fn != nil {
		return fn(params)
	}

	return []*sakila.FilmCategory{}, nil
}

// GetFilm runs the mock function or returns an empty film.
func (c *FilmCache) GetFilm(filmID int) (*sakila.Film, error) {
	if fn := c.GetFilmFn; fn != nil {
		return fn(filmID)
	}

	return &sakila.Film{}, nil
}

// GetFilms runs the mock function or returns no films.
func (c *FilmCache) GetFilms(params sakila.FilmQueryParams) ([]*sakila.Film, error) {
	if fn := c.GetFilmsFn; fn != nil {
		return fn(params)
	}

	return []*sakila.Film{}, nil
}

// GetFilmActors runs the mock function or returns no actors.
func (c *FilmCache) GetFilmActors(params sakila.FilmActorParams) ([]*sakila.FilmActor, error) {
	if fn := c.GetFilmActorsFn; fn != nil {
		return fn(params)
	}

	return []*sakila.FilmActor{}, nil
}

// GetFilmCategories runs the mock function or returns no categories.
func (c *FilmCache) GetFilmCategories(params sakila.FilmCategoryParams) ([]*sakila.FilmCategory, error) {
	if fn := c.GetFilmCategoriesFn; fn != nil {
		return fn(params)
	}

	return []*sakila.FilmCategory{}, nil
}

// SetFilm runs the mock function or returns no error.
func (c *FilmCache) SetFilm(film *sakila.Film) error {
	if fn := c.SetFilmFn; fn != nil {
		return fn(film)
	}

	return nil
}

// SetFilms runs the mock function or returns no error.
func (c *FilmCache) SetFilms(films []*sakila.Film, params sakila.FilmQueryParams) error {
	if fn := c.SetFilmsFn; fn != nil {
		return fn(films, params)
	}

	return nil
}

// SetFilmActors runs the mock function or returns no error.
func (c *FilmCache) SetFilmActors(actors []*sakila.FilmActor, params sakila.FilmActorParams) error {
	if fn := c.SetFilmActorsFn; fn != nil {
		return fn(actors, params)
	}

	return nil
}

// SetFilmCategories runs the mock function or returns no error.
func (c *FilmCache) SetFilmCategories(categories []*sakila.FilmCategory, params sakila.FilmCategoryParams) error {
	if fn := c.SetFilmCategoriesFn; fn != nil {
		return fn(categories, params)
	}

	return nil
}
