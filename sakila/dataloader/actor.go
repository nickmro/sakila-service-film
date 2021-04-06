package dataloader

import (
	"context"
	"sakila/sakila-film-service/sakila"
	"strconv"

	"github.com/graph-gophers/dataloader"
)

// FilmActorsDataLoader loads data for film actors.
func FilmActorsDataLoader(s sakila.FilmService) *dataloader.Loader {
	options := []dataloader.Option{
		dataloader.WithCache(&dataloader.NoCache{}),
	}

	return dataloader.NewBatchedLoader(func(
		ctx context.Context,
		keys dataloader.Keys,
	) []*dataloader.Result {
		ids := make([]int, len(keys))
		for i := range keys {
			id, err := strconv.ParseInt(keys[i].String(), 10, 32)
			if err != nil {
				return []*dataloader.Result{{Error: err}}
			}
			ids[i] = int(id)
		}

		actors, err := s.GetFilmActors(sakila.FilmActorParams{
			FilmIDs: ids,
		})
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

		results := make([]*dataloader.Result, len(ids))
		for i := range ids {
			if actors, ok := filmsMap[ids[i]]; ok {
				results[i] = &dataloader.Result{Data: actors}
			} else {
				results[i] = &dataloader.Result{Data: []*sakila.Actor{}}
			}
		}

		return results
	}, options...)
}
