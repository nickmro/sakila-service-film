package graphql

import (
	"encoding/json"
	"sakila/sakila-film-service/sakila"

	"github.com/graphql-go/graphql"
)

// Schema is a sakila graphQL schema.
type Schema struct {
	*graphql.Schema
}

// NewSchema returns a new graphQL schema.
func NewSchema(s sakila.FilmService) (*Schema, error) {
	filmType := FilmType(s)

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
							Resolve: FilmResolver(s),
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
								"category": &graphql.ArgumentConfig{
									Type: graphql.String,
								},
							},
							Resolve: FilmsResolver(s),
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

// FilmResolver returns the film with the given ID.
func FilmResolver(s sakila.FilmService) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (i interface{}, e error) {
		if id, ok := p.Args["filmId"].(int); ok {
			return s.GetFilm(id)
		}

		return nil, nil
	}
}

// FilmsResolver returns films for the given parameters.
func FilmsResolver(s sakila.FilmService) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (i interface{}, e error) {
		params := sakila.FilmQueryParams{}

		if limit, ok := p.Args["limit"].(int); ok {
			params[sakila.FilmQueryParamLimit] = limit
		}

		if offset, ok := p.Args["offset"].(int); ok {
			params[sakila.FilmQueryParamOffset] = offset
		}

		if category, ok := p.Args["category"].(string); ok {
			params[sakila.FilmQueryParamCategory] = category
		}

		return s.GetFilms(params)
	}
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
