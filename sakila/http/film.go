package http

import (
	"encoding/json"
	"net/http"
	"sakila/sakila-film-service/sakila"
	"strconv"

	"github.com/go-chi/chi"
)

// FilmHandler handles film requests.
type FilmHandler struct {
	*chi.Mux
	FilmService sakila.FilmService
}

// NewFilmHandler returns a new film service handler.
func NewFilmHandler(service sakila.FilmService) *FilmHandler {
	mux := chi.NewMux()
	mux.Get("/", getFilmsHandlerFunc(service))
	mux.Get("/{id}", getFilmHandlerFunc(service))

	return &FilmHandler{Mux: mux}
}

func getFilmsHandlerFunc(service sakila.FilmService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := filmQueryParams(r)

		films, err := service.GetFilms(params)
		if err != nil {
			w.WriteHeader(statusForError(err))
			//nolint:errcheck
			json.NewEncoder(w).Encode(err)

			return
		}

		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		json.NewEncoder(w).Encode(films)
	}
}

func getFilmHandlerFunc(service sakila.FilmService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filmID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			//nolint:errcheck
			json.NewEncoder(w).Encode(http.StatusText(http.StatusBadRequest))

			return
		}

		film, err := service.GetFilm(filmID)
		if err != nil {
			w.WriteHeader(statusForError(err))
			//nolint:errcheck
			json.NewEncoder(w).Encode(err)

			return
		}

		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		json.NewEncoder(w).Encode(film)
	}
}

func filmQueryParams(r *http.Request) map[sakila.FilmQueryParam]interface{} {
	query := r.URL.Query()
	params := map[sakila.FilmQueryParam]interface{}{}

	if firstParam := query[string(sakila.FilmQueryParamFirst)]; len(firstParam) > 0 {
		if first, err := strconv.Atoi(firstParam[0]); err == nil {
			params[sakila.FilmQueryParamFirst] = first
		}
	}

	if afterParam := query[string(sakila.FilmQueryParamAfter)]; len(afterParam) > 0 {
		if after, err := strconv.Atoi(afterParam[0]); err == nil {
			params[sakila.FilmQueryParamAfter] = after
		}
	}

	if categoryParam := query[string(sakila.FilmQueryParamCategory)]; len(categoryParam) > 0 {
		params[sakila.FilmQueryParamCategory] = categoryParam[0]
	}

	return params
}
