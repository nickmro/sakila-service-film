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
	actorType *graphql.Object
	//nolint:gochecknoglobals
	actorTypeSync sync.Once
)

// ActorType returns the graphQL actor type.
func ActorType() *graphql.Object {
	actorTypeSync.Do(func() {
		actorType = graphql.NewObject(
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
					"firstName": &graphql.Field{
						Type:        graphql.String,
						Description: "The actor first name.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if actor, ok := p.Source.(*sakila.Actor); ok {
								return actor.FirstName, nil
							}

							return nil, nil
						},
					},
					"lastName": &graphql.Field{
						Type:        graphql.String,
						Description: "The actor last name.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if actor, ok := p.Source.(*sakila.Actor); ok {
								return actor.LastName, nil
							}

							return nil, nil
						},
					},
					"lastUpdate": &graphql.Field{
						Type:        graphql.String,
						Description: "The actor last updated at time.",
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if actor, ok := p.Source.(*sakila.Actor); ok {
								return actor.LastUpdate, nil
							}

							return nil, nil
						},
					},
				},
			},
		)
	})

	return actorType
}

// FilmActorsResolver returns a film actors resolver.
func FilmActorsResolver(loader *dataloader.Loader) graphql.FieldResolveFn {
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
