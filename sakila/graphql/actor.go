package graphql

import (
	"context"
	"strconv"

	"github.com/nickmro/sakila-service-film/sakila"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
)

// FilmActorsDataLoader loads data for film actors.
func FilmActorsDataLoader(service sakila.FilmService) *dataloader.Loader {
	options := []dataloader.Option{
		dataloader.WithCache(&dataloader.NoCache{}),
		dataloader.WithBatchCapacity(20),
	}

	return dataloader.NewBatchedLoader(func(
		ctx context.Context,
		keys dataloader.Keys,
	) []*dataloader.Result {
		filmIDs := make([]int, len(keys))
		for i := range keys {
			id, err := strconv.ParseInt(keys[i].String(), 10, 32)
			if err != nil {
				return []*dataloader.Result{{Error: err}}
			}
			filmIDs[i] = int(id)
		}

		actors, err := service.GetFilmActors(ctx, filmIDs...)
		if err != nil {
			return []*dataloader.Result{{Error: err}}
		}

		filmsMap := map[int][]*sakila.Actor{}
		for _, actor := range actors {
			if actors := filmsMap[actor.FilmID]; actors != nil {
				filmsMap[actor.FilmID] = append(actors, &actor.Actor)
			} else {
				filmsMap[actor.FilmID] = []*sakila.Actor{&actor.Actor}
			}
		}

		results := make([]*dataloader.Result, len(filmIDs))
		for i := range filmIDs {
			if actors, ok := filmsMap[filmIDs[i]]; ok {
				results[i] = &dataloader.Result{Data: actors}
			} else {
				results[i] = &dataloader.Result{Data: []*sakila.Actor{}}
			}
		}

		return results
	}, options...)
}

// FilmActorsResolver returns actors for the given films.
func FilmActorsResolver(service sakila.FilmService) graphql.FieldResolveFn {
	loader := FilmActorsDataLoader(service)

	return func(params graphql.ResolveParams) (interface{}, error) {
		if film, ok := params.Source.(*sakila.Film); ok {
			key := strconv.Itoa(film.FilmID)
			thunk := loader.Load(params.Context, dataloader.StringKey(key))

			return func() (interface{}, error) {
				return thunk()
			}, nil
		}

		return nil, nil
	}
}
