package main

import (
	"context"
	"fmt"
	"net/http"
	"sakila/sakila-film-service/sakila/config"
	"sakila/sakila-film-service/sakila/graphql"
	"sakila/sakila-film-service/sakila/health"
	"sakila/sakila-film-service/sakila/log"
	"sakila/sakila-film-service/sakila/mysql"
	"sakila/sakila-film-service/sakila/redis"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

func main() { // nolint:gocyclo
	env, err := config.GetEnv(".env")
	if err != nil {
		panic(err)
	}

	logger, err := log.NewWriter(log.Environment(env.GetLogger()))
	if err != nil {
		panic(err)
	}

	fmt.Println("Connecting to DB...")

	db, err := mysql.Open(env.GetMySQLURL())
	if err != nil {
		panic(err)
	}

	//nolint:errcheck
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connecting to cache...")

	redisClient := redis.NewClient(&redis.ClientParams{
		Host:     env.GetRedisHost(),
		Port:     env.GetRedisPort(),
		Password: env.GetRedisPassword(),
	})

	//nolint:errcheck
	defer redisClient.Close()

	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	cache, err := redis.NewCache(redisClient)
	if err != nil {
		panic(err)
	}

	filmService := &mysql.FilmService{
		DB:     db,
		Logger: logger,
	}
	filmCache := &redis.FilmService{
		FilmService:    filmService,
		Cache:          cache,
		CacheKeyPrefix: env.GetRedisKeyPrefix(),
	}

	fmt.Println("Building GraphQL schema...")

	graphqlSchema, err := graphql.NewSchema(filmCache)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting health checker...")

	checker, err := health.NewChecker(&health.Checks{
		DB: &health.Check{
			Name:    "mysql",
			Checker: db,
		},
		Cache: &health.Check{
			Name:    "redis",
			Checker: redisClient,
		},
	})
	if err != nil {
		panic(err)
	}

	if err := checker.Start(); err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	router.Mount("/graphql", graphql.NewHandler(graphqlSchema))
	router.Mount("/healthz", health.NewHandler(checker))
	router.Mount("/readyz", health.NewHandler(checker))

	fmt.Println("Starting server...")

	addr := fmt.Sprintf(":%s", env.GetPort())

	fmt.Println("Listening on", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		panic(err)
	}

	logger.Flush()
}
