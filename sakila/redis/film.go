package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sakila/sakila-film-service/sakila"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// FilmCache is a Redis film cache.
type FilmCache struct {
	Client         *Client
	CacheKeyPrefix string
}

const (
	timeoutDuration    = time.Second * 5
	expirationDuration = time.Minute * 5
)

// GetFilm returns a film from the cache.
func (c *FilmCache) GetFilm(id int) (*sakila.Film, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := c.Client.Get(ctx, filmCacheKey(c.CacheKeyPrefix, id)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, sakila.ErrorNotFound
	} else if err != nil {
		return nil, err
	}

	film, err := unmarshalFilm(val)
	if err != nil {
		return nil, err
	}

	return film, nil
}

// GetFilms returns films from the cache.
func (c *FilmCache) GetFilms(params sakila.FilmQueryParams) ([]*sakila.Film, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := c.Client.Get(ctx, filmsCacheKey(c.CacheKeyPrefix, params)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, sakila.ErrorNotFound
	} else if err != nil {
		return nil, err
	}

	films, err := unmarshalFilms(val)
	if err != nil {
		return nil, err
	}

	return films, nil
}

// SetFilm sets the film in the cache.
func (c *FilmCache) SetFilm(film *sakila.Film) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := marshalFilm(film)
	if err != nil {
		return err
	}

	return c.Client.Set(ctx, filmCacheKey(c.CacheKeyPrefix, film.FilmID), val, expirationDuration).Err()
}

// SetFilms sets the films in the cache.
func (c *FilmCache) SetFilms(films []*sakila.Film, params sakila.FilmQueryParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := marshalFilms(films)
	if err != nil {
		return err
	}

	return c.Client.Set(ctx, filmsCacheKey(c.CacheKeyPrefix, params), val, expirationDuration).Err()
}

func filmCacheKey(prefix string, id int) string {
	b := strings.Builder{}

	if prefix != "" {
		b.WriteString(prefix)
		b.WriteString("::")
	}

	b.WriteString(fmt.Sprintf("film::id:%d", id))

	return b.String()
}

func filmsCacheKey(prefix string, params sakila.FilmQueryParams) string {
	b := strings.Builder{}

	if prefix != "" {
		b.WriteString(prefix)
		b.WriteString("::")
	}

	b.WriteString("films")

	if limit, ok := params[sakila.FilmQueryParamLimit].(int); ok {
		b.WriteString(fmt.Sprintf("::limit:%d", limit))
	}

	if offset, ok := params[sakila.FilmQueryParamOffset].(int); ok {
		b.WriteString(fmt.Sprintf("::offset:%d", offset))
	}

	if category, ok := params[sakila.FilmQueryParamCategory].(string); ok {
		b.WriteString("::category:")
		b.WriteString(category)
	}

	return b.String()
}

func marshalFilm(film *sakila.Film) (string, error) {
	b, err := json.Marshal(film)
	if err != nil {
		return "", err
	}

	return string(b), err
}

func marshalFilms(films []*sakila.Film) (string, error) {
	b, err := json.Marshal(films)
	if err != nil {
		return "", err
	}

	return string(b), err
}

func unmarshalFilm(val string) (*sakila.Film, error) {
	var film sakila.Film

	if err := json.Unmarshal([]byte(val), &film); err != nil {
		return nil, err
	}

	return &film, nil
}

func unmarshalFilms(val string) ([]*sakila.Film, error) {
	var films []*sakila.Film

	if err := json.Unmarshal([]byte(val), &films); err != nil {
		return nil, err
	}

	return films, nil
}
