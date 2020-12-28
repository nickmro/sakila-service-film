package api_test

import (
	"encoding/json"
	"net/http/httptest"
	"sakila/sakila-film-service/sakila"
	"sakila/sakila-film-service/sakila/api"
	"sakila/sakila-film-service/sakila/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Film", func() {
	var filmHandler *api.FilmHandler
	var filmService *mock.FilmService

	BeforeEach(func() {
		filmService = &mock.FilmService{}
		filmHandler = api.NewFilmHandler(filmService)
	})

	Describe("GET films", func() {
		It("writes the films from the film service", func() {
			filmService.GetFilmsFn = func(_ map[sakila.FilmQueryParam]interface{}) ([]*sakila.Film, error) {
				return []*sakila.Film{{}, {}}, nil
			}

			r, _ := api.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()

			filmHandler.ServeHTTP(rr, r)
			films := filmsFromRecorder(rr)
			Expect(films).ToNot(BeNil())
			Expect(films).To(HaveLen(2))
		})

		Context("when the parameters are provided", func() {
			It("passes them to the service", func() {
				var first int
				var after int
				var category string

				filmService.GetFilmsFn = func(params map[sakila.FilmQueryParam]interface{}) ([]*sakila.Film, error) {
					first = params[sakila.FilmQueryParamFirst].(int)
					after = params[sakila.FilmQueryParamAfter].(int)
					category = params[sakila.FilmQueryParamCategory].(string)

					return []*sakila.Film{{}, {}}, nil
				}

				r, _ := api.NewRequest("GET", "/?category=Animation&first=20&after=100", nil)
				rr := httptest.NewRecorder()

				filmHandler.ServeHTTP(rr, r)
				Expect(first).To(Equal(20))
				Expect(after).To(Equal(100))
				Expect(category).To(Equal("Animation"))
			})
		})

		Context("when getting the films fails", func() {
			BeforeEach(func() {
				filmService.GetFilmsFn = func(_ map[sakila.FilmQueryParam]interface{}) ([]*sakila.Film, error) {
					return nil, sakila.ErrorInternal
				}
			})

			It("returns an error", func() {
				r, _ := api.NewRequest("GET", "/", nil)
				rr := httptest.NewRecorder()

				filmHandler.ServeHTTP(rr, r)
				Expect(rr.Code).To(Equal(500))
			})
		})
	})

	Describe("GET film", func() {
		It("writes the film from the film service", func() {
			filmService.GetFilmFn = func(filmID int) (*sakila.Film, error) {
				return &sakila.Film{FilmID: filmID}, nil
			}

			r, _ := api.NewRequest("GET", "/1", nil)
			rr := httptest.NewRecorder()

			filmHandler.ServeHTTP(rr, r)
			film := filmFromRecorder(rr)
			Expect(film).ToNot(BeNil())
			Expect(film.FilmID).To(Equal(1))
		})

		Context("when the id parameter is not an integer", func() {
			It("returns a bad request error", func() {
				r, _ := api.NewRequest("GET", "/title", nil)
				rr := httptest.NewRecorder()

				filmHandler.ServeHTTP(rr, r)
				Expect(rr.Code).To(Equal(400))
			})
		})

		Context("when getting the film fails", func() {
			BeforeEach(func() {
				filmService.GetFilmFn = func(filmID int) (*sakila.Film, error) {
					return nil, sakila.ErrorInternal
				}
			})

			It("writes an error", func() {
				r, _ := api.NewRequest("GET", "/1", nil)
				rr := httptest.NewRecorder()

				filmHandler.ServeHTTP(rr, r)
				Expect(rr.Code).To(Equal(500))
			})
		})
	})
})

func filmFromRecorder(rr *httptest.ResponseRecorder) *sakila.Film {
	var film sakila.Film

	b := rr.Body.Bytes()

	if err := json.Unmarshal(b, &film); err != nil {
		return nil
	}

	return &film
}

func filmsFromRecorder(rr *httptest.ResponseRecorder) []*sakila.Film {
	var films []*sakila.Film

	b := rr.Body.Bytes()

	if err := json.Unmarshal(b, &films); err != nil {
		return nil
	}

	return films
}
