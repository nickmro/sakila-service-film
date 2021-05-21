package graphql

import (
	"encoding/json"

	"github.com/nickmro/sakila-service-film/sakila"

	"github.com/graphql-go/graphql"
)

// Schema is a sakila graphQL schema.
type Schema struct {
	*graphql.Schema
}

// NewSchema returns a new graphQL schema.
func NewSchema(service sakila.FilmService) (*Schema, error) { //nolint:gocyclo
	actorType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Actor",
			Description: "An Actor is a Sakila film actor.",
			Fields: graphql.Fields{
				"actorId": &graphql.Field{
					Type:        graphql.Int,
					Description: "The actor ID.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if actor, ok := p.Source.(*sakila.Actor); ok {
							return actor.ActorID, nil
						}

						return nil, nil
					},
				},
			},
		},
	)

	filmType := graphql.NewObject(
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
					Type:        graphql.NewList(actorType),
					Description: "The film actors.",
					Resolve:     FilmActorsResolver(service),
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

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: graphql.NewObject(
				graphql.ObjectConfig{
					Name: "Query",
					Fields: graphql.Fields{
						"film": &graphql.Field{
							Description: "Returns the film with the given ID.",
							Type:        filmType,
							Args: graphql.FieldConfigArgument{
								"filmId": &graphql.ArgumentConfig{
									Type:        graphql.Int,
									Description: "The film ID.",
								},
							},
							Resolve: FilmResolver(service),
						},
						"films": &graphql.Field{
							Description: "Returns the films for the given parameters",
							Type:        graphql.NewList(filmType),
							Args: graphql.FieldConfigArgument{
								"limit": &graphql.ArgumentConfig{
									Type: graphql.Int,
								},
								"offset": &graphql.ArgumentConfig{
									Type: graphql.Int,
								},
							},
							Resolve: FilmsResolver(service),
						},
					},
				},
			),
		},
	)
	if err != nil {
		return nil, err
	}

	return &Schema{Schema: &schema}, nil
}

// Request takes a query to return data from the graphQL service.
func (s *Schema) Request(query string) ([]byte, error) {
	params := graphql.Params{Schema: *s.Schema, RequestString: query}

	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		return nil, r.Errors[0]
	}

	return json.Marshal(r.Data)
}
