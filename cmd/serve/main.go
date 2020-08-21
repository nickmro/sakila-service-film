package main

import (
	"context"
	"database/sql"
	"fmt"

	"sakila/sakila-film-service/sakila/api"
	"sakila/sakila-film-service/sakila/config"
	"sakila/sakila-film-service/sakila/graphql"
	"sakila/sakila-film-service/sakila/http"
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
	db, err := sql.Open("mysql", env.GetMySQLURL())
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connecting to cache...")
	cache := redis.NewClient(env.GetRedisURL(), env.GetRedisPassword())
	defer cache.Close()

	err = cache.Ping(context.Background()).Err()
	if err != nil {
		logger.Fatal(err)
	}

	filmStore := &mysql.FilmDB{DB: db}
	filmCache := &redis.FilmCache{Cache: cache}
	actorStore := &mysql.ActorDB{DB: db}

	filmService := &api.FilmService{
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

	filmHandler := http.NewFilmHandler(filmService)
	graphqlHandler := http.NewGraphQLHandler(graphqlSchema)

	router := http.NewRouter(logger)
	router.Mount("/films", filmHandler)
	router.Mount("/graphql", graphqlHandler)

	logger.Info("Starting server...")

	addr := fmt.Sprintf(":%s", env.GetPort())
	logger.Info("Listening on", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal(err)
	}

	logger.Flush()
}
