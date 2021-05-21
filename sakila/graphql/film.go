package graphql

import (
	"github.com/nickmro/sakila-service-film/sakila"

	"github.com/graphql-go/graphql"
)

// FilmResolver returns the film with the given ID.
func FilmResolver(service sakila.FilmService) graphql.FieldResolveFn {
	return func(params graphql.ResolveParams) (i interface{}, e error) {
		if filmID, ok := params.Args["filmId"].(int); ok {
			return service.GetFilm(params.Context, filmID)
		}

		return nil, nil
	}
}

// FilmsResolver returns films for the given parameters.
func FilmsResolver(service sakila.FilmService) graphql.FieldResolveFn {
	return func(params graphql.ResolveParams) (i interface{}, e error) {
		filmParams := sakila.FilmParams{}

		if limit, ok := params.Args["limit"].(int); ok {
			filmParams.Limit = limit
		}

		if offset, ok := params.Args["offset"].(int); ok {
			filmParams.Offset = offset
		}

		return service.GetFilms(params.Context, filmParams)
	}
}
