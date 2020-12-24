package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sakila/sakila-film-service/sakila"
	"time"

	"github.com/go-redis/redis/v8"
)

// FilmCache is a Redis film cache.
type FilmCache struct {
	Client *Client
}

const (
	cachekeyFilmFmt = "sakila:film.service:film:%d"
)

const (
	timeoutDuration    = time.Second * 5
	expirationDuration = time.Minute * 5
)

// GetFilm returns a film from the cache.
func (c *FilmCache) GetFilm(id int) (*sakila.Film, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := c.Client.Get(ctx, filmCacheKey(id)).Result()
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

// SetFilm sets the film in the cache.
func (c *FilmCache) SetFilm(film *sakila.Film) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := marshalFilm(film)
	if err != nil {
		return err
	}

	return c.Client.Set(ctx, filmCacheKey(film.FilmID), val, expirationDuration).Err()
}

func filmCacheKey(id int) string {
	return fmt.Sprintf(cachekeyFilmFmt, id)
}

func marshalFilm(film *sakila.Film) (string, error) {
	b, err := json.Marshal(film)
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
