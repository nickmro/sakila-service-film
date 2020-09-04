package graphql

import (
	"sakila/sakila-film-service/sakila"
	"sync"

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
