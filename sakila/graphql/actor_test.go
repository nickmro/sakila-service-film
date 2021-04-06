package graphql_test

// Actor is a sakila actor.
type Actor struct {
	ActorID    int    `json:"actorId"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	LastUpdate string `json:"lastUpdate"`
}
