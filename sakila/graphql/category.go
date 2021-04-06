package graphql

import (
	"sakila/sakila-film-service/sakila"
	"strconv"
	"sync"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
)

var (
	//nolint:gochecknoglobals
	categoryType *graphql.Object
	//nolint:gochecknoglobals
	categoryTypeSync sync.Once
)

// CategoryType returns the graphQL category type.
func CategoryType() *graphql.Object {
	categoryTypeSync.Do(func() {
		categoryType = graphql.NewObject(
			graphql.ObjectConfig{
				Name:        "Category",
				Description: "An Category is a Sakila film category.",
				Fields: graphql.Fields{
					"categoryId": &graphql.Field{
						Type:        graphql.Int,
						Description: "The category ID.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if category, ok := p.Source.(*sakila.Category); ok {
								return category.CategoryID, nil
							}

							return nil, nil
						},
					},
					"name": &graphql.Field{
						Type:        graphql.String,
						Description: "The category name.",
					},
					"lastUpdate": &graphql.Field{
						Type:        graphql.String,
						Description: "The category last updated at time.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if category, ok := p.Source.(*sakila.Category); ok {
								return category.LastUpdate, nil
							}

							return nil, nil
						},
					},
				},
			},
		)
	})

	return categoryType
}

// FilmCategoriesResolver returns a film categories resolver.
func FilmCategoriesResolver(loader *dataloader.Loader) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		if film, ok := p.Source.(*sakila.Film); ok {
			key := strconv.Itoa(film.FilmID)
			thunk := loader.Load(p.Context, dataloader.StringKey(key))

			return func() (interface{}, error) {
				return thunk()
			}, nil
		}

		return nil, nil
	}
}
