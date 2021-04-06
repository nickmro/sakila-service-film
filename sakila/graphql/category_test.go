package graphql_test

// Category is a sakila category.
type Category struct {
	CategoryID int    `json:"categoryId"`
	Name       string `json:"name"`
	LastUpdate string `json:"lastUpdate"`
}
