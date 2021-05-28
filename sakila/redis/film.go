package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nickmro/sakila-service-film/sakila"

	"github.com/go-redis/cache/v8"
)

// FilmService is a cached film service.
type FilmService struct {
	sakila.FilmService
	Cache          *Cache
	CacheKeyPrefix string
	TTL            time.Duration
	Logger         sakila.Logger
}

// GetFilm returns a film from the cache.
func (service *FilmService) GetFilm(ctx context.Context, id int) (*sakila.Film, error) {
	var film sakila.Film

	item := &cache.Item{
		Ctx:   ctx,
		Key:   service.filmCacheKey(id),
		Value: &film,
		Do: func(i *cache.Item) (interface{}, error) {
			return service.FilmService.GetFilm(ctx, id)
		},
		TTL: service.TTL,
	}

	err := service.Cache.Once(item)
	if err != nil && errors.Is(err, sakila.ErrorNotFound) {
		return nil, sakila.ErrorNotFound
	} else if err != nil {
		service.logError(err)
		return nil, sakila.ErrorInternal
	}

	return &film, err
}

// GetFilms returns films from the cache.
func (service *FilmService) GetFilms(ctx context.Context, params sakila.FilmParams) ([]*sakila.Film, error) {
	var films []*sakila.Film

	item := &cache.Item{
		Ctx:   ctx,
		Key:   service.filmsCacheKey(params),
		Value: &films,
		Do: func(i *cache.Item) (interface{}, error) {
			return service.FilmService.GetFilms(ctx, params)
		},
		TTL: service.TTL,
	}

	err := service.Cache.Once(item)
	if err != nil && errors.Is(err, sakila.ErrorNotFound) {
		return nil, sakila.ErrorNotFound
	} else if err != nil {
		service.logError(err)
		return nil, sakila.ErrorInternal
	}

	return films, err
}

// GetFilmActors returns film actors from the cache.
func (service *FilmService) GetFilmActors(ctx context.Context, filmIDs ...int) ([]*sakila.FilmActor, error) {
	var actors []*sakila.FilmActor

	item := &cache.Item{
		Ctx:   ctx,
		Key:   service.actorsCacheKey(filmIDs...),
		Value: &actors,
		Do: func(i *cache.Item) (interface{}, error) {
			return service.FilmService.GetFilmActors(ctx, filmIDs...)
		},
		TTL: service.TTL,
	}

	err := service.Cache.Once(item)
	if err != nil && errors.Is(err, sakila.ErrorNotFound) {
		return nil, sakila.ErrorNotFound
	} else if err != nil {
		service.logError(err)
		return nil, sakila.ErrorInternal
	}

	return actors, err
}

func (service *FilmService) logError(err error) {
	if logger := service.Logger; logger != nil {
		logger.Error(err)
	}
}

func (service *FilmService) cacheKey(key string) string {
	if prefix := service.CacheKeyPrefix; prefix != "" {
		return prefix + "::" + key
	}

	return key
}

func (service *FilmService) filmCacheKey(id int) string {
	key := hashedKey("film::id:" + strconv.Itoa(id))
	return service.cacheKey(key)
}

func (service *FilmService) filmsCacheKey(params sakila.FilmParams) string {
	b := strings.Builder{}

	b.WriteString("films::")

	if ids := params.FilmIDs; len(ids) > 0 {
		b.WriteString("::ids:")

		for i := range ids {
			if i > 0 {
				b.WriteString(",")
			}

			b.WriteString(strconv.Itoa(ids[i]))
		}
	}

	if limit := params.Limit; limit > 0 {
		b.WriteString("::limit:" + strconv.Itoa(limit))
	}

	if offset := params.Offset; offset > 0 {
		b.WriteString(fmt.Sprintf("::offset:" + strconv.Itoa(offset)))
	}

	return service.cacheKey(hashedKey(b.String()))
}

func (service *FilmService) actorsCacheKey(filmIDs ...int) string {
	b := strings.Builder{}

	for i := range filmIDs {
		if i > 0 {
			b.WriteString(",")
		}

		b.WriteString("::film_ids:" + strconv.Itoa(filmIDs[i]))
	}

	return service.cacheKey(hashedKey(b.String()))
}
