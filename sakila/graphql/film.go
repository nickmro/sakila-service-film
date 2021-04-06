package graphql

import (
	"sakila/sakila-film-service/sakila"
	"sakila/sakila-film-service/sakila/dataloader"
	"sync"

	"github.com/graphql-go/graphql"
)

var (
	//nolint:gochecknoglobals
	filmType *graphql.Object
	//nolint:gochecknoglobals
	filmTypeSync sync.Once
)

// FilmType returns the graphQL film type.
func FilmType(s sakila.FilmService) *graphql.Object {
	actorsLoader := dataloader.FilmActorsDataLoader(s)
	categoriesLoader := dataloader.FilmCategoriesDataLoader(s)

	filmTypeSync.Do(func() {
		filmType = graphql.NewObject(
			graphql.ObjectConfig{
				Name:        "Film",
				Description: "A Film is a Sakila film.",
				Fields: graphql.Fields{
					"filmId": &graphql.Field{
						Type:        graphql.Int,
						Description: "The film ID.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if film, ok := p.Source.(*sakila.Film); ok {
								return film.FilmID, nil
							}

							return nil, nil
						},
					},
					"title": &graphql.Field{
						Type:        graphql.String,
						Description: "The film title.",
					},
					"description": &graphql.Field{
						Type:        graphql.String,
						Description: "The film description.",
					},
					"releaseYear": &graphql.Field{
						Type:        graphql.Int,
						Description: "The film release year.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if film, ok := p.Source.(*sakila.Film); ok {
								return film.ReleaseYear, nil
							}

							return nil, nil
						},
					},
					"actors": &graphql.Field{
						Type:        graphql.NewList(ActorType()),
						Description: "The film actors.",
						Resolve:     FilmActorsResolver(actorsLoader),
					},
					"categories": &graphql.Field{
						Type:        graphql.NewList(CategoryType()),
						Description: "The film categories.",
						Resolve:     FilmCategoriesResolver(categoriesLoader),
					},
					"languageId": &graphql.Field{
						Type:        graphql.Int,
						Description: "The film language ID.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if film, ok := p.Source.(*sakila.Film); ok {
								return film.LanguageID, nil
							}

							return nil, nil
						},
					},
					"originalLanguageId": &graphql.Field{
						Type:        graphql.Int,
						Description: "The film original language ID.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if film, ok := p.Source.(*sakila.Film); ok {
								return film.OriginalLanguageID, nil
							}

							return nil, nil
						},
					},
					"rentalDuration": &graphql.Field{
						Type:        graphql.Int,
						Description: "The film rental duration.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if film, ok := p.Source.(*sakila.Film); ok {
								return film.RentalDuration, nil
							}

							return nil, nil
						},
					},
					"rentalRate": &graphql.Field{
						Type:        graphql.Float,
						Description: "The film rental rate.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if film, ok := p.Source.(*sakila.Film); ok {
								return film.RentalRate, nil
							}

							return nil, nil
						},
					},
					"length": &graphql.Field{
						Type:        graphql.Int,
						Description: "The film length.",
					},
					"replacementCost": &graphql.Field{
						Type:        graphql.Float,
						Description: "The film replacement cost.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if film, ok := p.Source.(*sakila.Film); ok {
								return film.ReplacementCost, nil
							}

							return nil, nil
						},
					},
					"rating": &graphql.Field{
						Type:        graphql.String,
						Description: "The film rating.",
					},
					"specialFeatures": &graphql.Field{
						Type:        graphql.NewList(graphql.String),
						Description: "The film special features.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if film, ok := p.Source.(*sakila.Film); ok {
								return film.SpecialFeatures, nil
							}

							return nil, nil
						},
					},
					"lastUpdate": &graphql.Field{
						Type:        graphql.DateTime,
						Description: "The film last update time.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if film, ok := p.Source.(*sakila.Film); ok {
								return film.LastUpdate, nil
							}

							return nil, nil
						},
					},
				},
			},
		)
	})

	return filmType
}
