package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sakila/sakila-film-service/sakila"
	"strconv"
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
	var film sakila.Film

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := c.Client.Get(ctx, filmCacheKey(c.CacheKeyPrefix, id)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, sakila.ErrorNotFound
	} else if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(val), &film); err != nil {
		return nil, err
	}

	return &film, nil
}

// GetFilms returns films from the cache.
func (c *FilmCache) GetFilms(params sakila.FilmQueryParams) ([]*sakila.Film, error) {
	var films []*sakila.Film

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := c.Client.Get(ctx, filmsCacheKey(c.CacheKeyPrefix, params)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, sakila.ErrorNotFound
	} else if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(val), &films); err != nil {
		return nil, err
	}

	return films, nil
}

// GetFilmActors returns film actors from the cache.
func (c *FilmCache) GetFilmActors(params sakila.FilmActorParams) ([]*sakila.FilmActor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := c.Client.Get(ctx, filmActorsCacheKey(c.CacheKeyPrefix, params)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, sakila.ErrorNotFound
	} else if err != nil {
		return nil, err
	}

	var actors []*sakila.FilmActor
	if err := json.Unmarshal([]byte(val), &actors); err != nil {
		return nil, err
	}

	return actors, nil
}

// GetFilmCategories returns film categories from the cache.
func (c *FilmCache) GetFilmCategories(params sakila.FilmCategoryParams) ([]*sakila.FilmCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := c.Client.Get(ctx, filmCategoriesCacheKey(c.CacheKeyPrefix, params)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, sakila.ErrorNotFound
	} else if err != nil {
		return nil, err
	}

	var categories []*sakila.FilmCategory
	if err := json.Unmarshal([]byte(val), &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// SetFilm sets the film in the cache.
func (c *FilmCache) SetFilm(film *sakila.Film) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := marshal(film)
	if err != nil {
		return err
	}

	return c.Client.Set(ctx, filmCacheKey(c.CacheKeyPrefix, film.FilmID), val, expirationDuration).Err()
}

// SetFilms sets the films in the cache.
func (c *FilmCache) SetFilms(films []*sakila.Film, params sakila.FilmQueryParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := marshal(films)
	if err != nil {
		return err
	}

	return c.Client.Set(ctx, filmsCacheKey(c.CacheKeyPrefix, params), val, expirationDuration).Err()
}

// SetFilmActors sets the film actors in the cache.
func (c *FilmCache) SetFilmActors(actors []*sakila.FilmActor, params sakila.FilmActorParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := marshal(actors)
	if err != nil {
		return err
	}

	return c.Client.Set(ctx, filmActorsCacheKey(c.CacheKeyPrefix, params), val, expirationDuration).Err()
}

// SetFilmCategories sets the film categories in the cache.
func (c *FilmCache) SetFilmCategories(categories []*sakila.FilmCategory, params sakila.FilmCategoryParams) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	val, err := marshal(categories)
	if err != nil {
		return err
	}

	return c.Client.Set(ctx, filmCategoriesCacheKey(c.CacheKeyPrefix, params), val, expirationDuration).Err()
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

func filmActorsCacheKey(prefix string, params sakila.FilmActorParams) string {
	b := strings.Builder{}
	b.WriteString(prefix)
	b.WriteString("::actors")

	if len(params.FilmIDs) > 0 {
		b.WriteString("::film_id:")

		ids := make([]string, len(params.FilmIDs))
		for i := range ids {
			ids[i] = strconv.Itoa(params.FilmIDs[i])
		}

		b.WriteString(strings.Join(ids, ","))
	}

	return b.String()
}

func filmCategoriesCacheKey(prefix string, params sakila.FilmCategoryParams) string {
	b := strings.Builder{}
	b.WriteString(prefix)
	b.WriteString("::categories")

	if len(params.FilmIDs) > 0 {
		b.WriteString("::film_id:")

		ids := make([]string, len(params.FilmIDs))
		for i := range ids {
			ids[i] = strconv.Itoa(params.FilmIDs[i])
		}

		b.WriteString(strings.Join(ids, ","))
	}

	return b.String()
}

func marshal(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(b), err
}
