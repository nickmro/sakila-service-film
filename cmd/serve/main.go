package main

import (
	"context"
	"fmt"
	"net/http"
	"sakila/sakila-film-service/sakila/api"
	"sakila/sakila-film-service/sakila/app"
	"sakila/sakila-film-service/sakila/config"
	"sakila/sakila-film-service/sakila/graphql"
	"sakila/sakila-film-service/sakila/health"
	"sakila/sakila-film-service/sakila/log"
	"sakila/sakila-film-service/sakila/mysql"
	"sakila/sakila-film-service/sakila/redis"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	env, err := config.GetEnv(".env")
	if err != nil {
		panic(err)
	}

	logger, err := log.NewWriter(log.Environment(env.GetLogger()))
	if err != nil {
		panic(err)
	}

	logger.Info("Connecting to DB...")

	db, err := mysql.Open(env.GetMySQLURL())
	if err != nil {
		logger.Fatal(err)
	}

	//nolint:errcheck
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connecting to cache...")

	cache := redis.NewClient(env.GetRedisURL(), env.GetRedisPassword())

	//nolint:errcheck
	defer cache.Close()

	err = cache.Ping(context.Background()).Err()
	if err != nil {
		logger.Fatal(err)
	}

	filmStore := &mysql.FilmDB{DB: db}
	filmCache := &redis.FilmCache{Client: cache, CacheKeyPrefix: env.GetRedisCacheKeyPrefix()}
	actorStore := &mysql.ActorDB{DB: db}

	filmService := &app.FilmService{
		ActorStore: actorStore,
		FilmCache:  filmCache,
		FilmStore:  filmStore,
		Logger:     logger,
	}

	logger.Info("Building GraphQL schema...")

	graphqlSchema, err := graphql.NewSchema(filmService)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Starting health checker...")

	checker, err := health.NewChecker(&health.Checks{
		DB: &health.Check{
			Name:    "mysql",
			Checker: db,
		},
		Cache: &health.Check{
			Name:    "redis",
			Checker: cache,
		},
	})
	if err != nil {
		logger.Fatal(err)
	}

	if err := checker.Start(); err != nil {
		logger.Fatal(err)
	}

	router := api.NewRouter(logger)
	router.Mount("/films", api.NewFilmHandler(filmService))
	router.Mount("/graphql", graphql.NewHandler(graphqlSchema))
	router.Mount("/healthz", health.NewHandler(checker))
	router.Mount("/readyz", health.NewHandler(checker))

	logger.Info("Starting server...")

	addr := fmt.Sprintf(":%s", env.GetPort())

	logger.Info("Listening on", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal(err)
	}

	logger.Flush()
}
