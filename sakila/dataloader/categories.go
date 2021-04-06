package dataloader

import (
	"context"
	"sakila/sakila-film-service/sakila"
	"strconv"

	"github.com/graph-gophers/dataloader"
)

// FilmCategoriesDataLoader loads data for film categories.
func FilmCategoriesDataLoader(s sakila.FilmService) *dataloader.Loader {
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

		categories, err := s.GetFilmCategories(sakila.FilmCategoryParams{
			FilmIDs: ids,
		})
		if err != nil {
			return []*dataloader.Result{{Error: err}}
		}

		filmsMap := map[int][]*sakila.Category{}
		for _, category := range categories {
			if categories := filmsMap[category.FilmID]; categories != nil {
				filmsMap[category.FilmID] = append(categories, &category.Category)
			} else {
				filmsMap[category.FilmID] = []*sakila.Category{&category.Category}
			}
		}

		results := make([]*dataloader.Result, len(ids))
		for i := range ids {
			if categories, ok := filmsMap[ids[i]]; ok {
				results[i] = &dataloader.Result{Data: categories}
			} else {
				results[i] = &dataloader.Result{Data: []*sakila.Category{}}
			}
		}

		return results
	}, options...)
}
