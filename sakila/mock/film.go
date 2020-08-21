package mock

import "sakila/sakila-film-service/sakila"

// FilmService is a mock film service.
type FilmService struct {
	GetFilmFn  func(filmID int) (*sakila.Film, error)
	GetFilmsFn func(params map[sakila.FilmQueryParam]interface{}) ([]*sakila.Film, error)
}

// FilmStore is a mock film store.
type FilmStore struct {
	QueryFilmFn  func(filmID int) (*sakila.Film, error)
	QueryFilmsFn func(params map[sakila.FilmQueryParam]interface{}) ([]*sakila.Film, error)
}

// FilmCache is a mock film cache.
type FilmCache struct {
	GetFilmFn func(filmID int) (*sakila.Film, error)
	SetFilmFn func(film *sakila.Film) error
}

// GetFilm runs the mock function or returns an empty film.
func (s *FilmService) GetFilm(filmID int) (*sakila.Film, error) {
	if fn := s.GetFilmFn; fn != nil {
		return fn(filmID)
	}

	return &sakila.Film{}, nil
}

// GetFilms runs the mock function or returns an empty slice of films.
func (s *FilmService) GetFilms(params map[sakila.FilmQueryParam]interface{}) ([]*sakila.Film, error) {
	if fn := s.GetFilmsFn; fn != nil {
		return fn(params)
	}

	return []*sakila.Film{}, nil
}

// QueryFilm runs the mock function or returns an empty film.
func (s *FilmStore) QueryFilm(filmID int) (*sakila.Film, error) {
	if fn := s.QueryFilmFn; fn != nil {
		return fn(filmID)
	}

	return &sakila.Film{}, nil
}

// QueryFilms runs the mock function or returns an empty array of films.
func (s *FilmStore) QueryFilms(params map[sakila.FilmQueryParam]interface{}) ([]*sakila.Film, error) {
	if fn := s.QueryFilmsFn; fn != nil {
		return fn(params)
	}

	return []*sakila.Film{}, nil
}

// GetFilm runs the mock function or returns an empty film.
func (c *FilmCache) GetFilm(filmID int) (*sakila.Film, error) {
	if fn := c.GetFilmFn; fn != nil {
		return fn(filmID)
	}

	return &sakila.Film{}, nil
}

// SetFilm runs the mock function or returns no error.
func (c *FilmCache) SetFilm(film *sakila.Film) error {
	if fn := c.SetFilmFn; fn != nil {
		return fn(film)
	}

	return nil
}
