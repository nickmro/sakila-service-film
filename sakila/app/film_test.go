package app_test

import (
	"sakila/sakila-film-service/sakila"
	"sakila/sakila-film-service/sakila/app"
	"sakila/sakila-film-service/sakila/log"
	"sakila/sakila-film-service/sakila/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const unexpectedError = sakila.Error("unexpected error")

var _ = Describe("FilmService", func() {
	var service *app.FilmService
	var filmCache *mock.FilmCache
	var filmStore *mock.FilmStore

	BeforeEach(func() {
		logger, err := log.NewWriter("TEST")
		if err != nil {
			panic(err)
		}

		filmCache = &mock.FilmCache{}
		filmStore = &mock.FilmStore{}
		service = &app.FilmService{
			Cache:  filmCache,
			Store:  filmStore,
			Logger: logger,
		}
	})

	Describe("GetFilm", func() {
		It("returns the film from the cache", func() {
			var invoked bool

			filmCache.GetFilmFn = func(id int) (*sakila.Film, error) {
				invoked = true

				return &sakila.Film{FilmID: id}, nil
			}

			film, err := service.GetFilm(1)
			Expect(err).ToNot(HaveOccurred())
			Expect(invoked).To(BeTrue())
			Expect(film).ToNot(BeNil())
			Expect(film.FilmID).To(Equal(1))
		})

		Context("when the film is not cached", func() {
			BeforeEach(func() {
				filmCache.GetFilmFn = func(_ int) (*sakila.Film, error) {
					return nil, sakila.ErrorNotFound
				}
			})

			It("gets the film from the store", func() {
				var invoked bool

				filmStore.QueryFilmFn = func(id int) (*sakila.Film, error) {
					invoked = true

					return &sakila.Film{FilmID: id}, nil
				}

				film, err := service.GetFilm(1)
				Expect(err).ToNot(HaveOccurred())
				Expect(invoked).To(BeTrue())
				Expect(film).ToNot(BeNil())
				Expect(film.FilmID).To(Equal(1))
			})

			It("caches the film", func() {
				var invoked bool

				filmCache.SetFilmFn = func(film *sakila.Film) error {
					invoked = true

					return nil
				}

				_, err := service.GetFilm(1)
				Expect(err).ToNot(HaveOccurred())
				Eventually(func() bool {
					return invoked
				}).Should(BeTrue())
			})
		})

		Context("when getting the film from the cache fails", func() {
			BeforeEach(func() {
				filmCache.GetFilmFn = func(filmID int) (*sakila.Film, error) {
					return nil, unexpectedError
				}
			})

			It("gets the film from the store", func() {
				var invoked bool

				filmStore.QueryFilmFn = func(filmID int) (*sakila.Film, error) {
					invoked = true

					return &sakila.Film{FilmID: filmID}, nil
				}

				film, err := service.GetFilm(1)
				Expect(err).ToNot(HaveOccurred())
				Expect(invoked).To(BeTrue())
				Expect(film).ToNot(BeNil())
				Expect(film.FilmID).To(Equal(1))
			})
		})

		Context("when the film is not found", func() {
			BeforeEach(func() {
				filmCache.GetFilmFn = func(filmID int) (*sakila.Film, error) {
					return nil, sakila.ErrorNotFound
				}

				filmStore.QueryFilmFn = func(filmID int) (*sakila.Film, error) {
					return nil, sakila.ErrorNotFound
				}
			})

			It("returns a not found error", func() {
				_, err := service.GetFilm(1)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(sakila.ErrorNotFound.Error()))
			})
		})

		Context("when getting the film from the store fails", func() {
			BeforeEach(func() {
				filmCache.GetFilmFn = func(filmID int) (*sakila.Film, error) {
					return nil, sakila.ErrorNotFound
				}

				filmStore.QueryFilmFn = func(filmID int) (*sakila.Film, error) {
					return nil, unexpectedError
				}
			})

			It("returns an internal error", func() {
				_, err := service.GetFilm(1)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(sakila.ErrorInternal.Error()))
			})
		})
	})

	Describe("GetFilms", func() {
		It("returns the films from the cache", func() {
			var invoked bool

			filmCache.GetFilmsFn = func(params sakila.FilmQueryParams) ([]*sakila.Film, error) {
				invoked = true

				return []*sakila.Film{{}}, nil
			}

			films, err := service.GetFilms(sakila.FilmQueryParams{})
			Expect(err).ToNot(HaveOccurred())
			Expect(invoked).To(BeTrue())
			Expect(films).ToNot(BeNil())
			Expect(films).To(HaveLen(1))
		})

		Context("when the films are not cached", func() {
			BeforeEach(func() {
				filmCache.GetFilmsFn = func(params sakila.FilmQueryParams) ([]*sakila.Film, error) {
					return nil, sakila.ErrorNotFound
				}
			})

			It("returns the films from the store", func() {
				var invoked bool

				filmStore.QueryFilmsFn = func(_ sakila.FilmQueryParams) ([]*sakila.Film, error) {
					invoked = true

					return []*sakila.Film{{}, {}}, nil
				}

				films, err := service.GetFilms(sakila.FilmQueryParams{})
				Expect(invoked).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
				Expect(films).To(HaveLen(2))
			})
		})

		It("sends the 'limit' parameter", func() {
			var limit int

			filmCache.GetFilmsFn = func(params sakila.FilmQueryParams) ([]*sakila.Film, error) {
				if param := params[sakila.FilmQueryParamLimit]; param != nil {
					limit = param.(int)
				}

				return []*sakila.Film{{}}, nil
			}

			_, err := service.GetFilms(sakila.FilmQueryParams{
				sakila.FilmQueryParamLimit: 10,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(limit).To(Equal(10))
		})

		It("sends the 'offset' parameter", func() {
			var offset int

			filmCache.GetFilmsFn = func(params sakila.FilmQueryParams) ([]*sakila.Film, error) {
				offset = params[sakila.FilmQueryParamOffset].(int)

				return []*sakila.Film{{}}, nil
			}

			_, err := service.GetFilms(sakila.FilmQueryParams{
				sakila.FilmQueryParamOffset: 10,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(offset).To(Equal(10))
		})

		Context("when the 'limit' parameter is not provided", func() {
			It("defaults to 20", func() {
				var limit int

				filmCache.GetFilmsFn = func(params sakila.FilmQueryParams) ([]*sakila.Film, error) {
					if param := params[sakila.FilmQueryParamLimit]; param != nil {
						limit = param.(int)
					}

					return []*sakila.Film{{}}, nil
				}

				_, err := service.GetFilms(sakila.FilmQueryParams{})
				Expect(err).ToNot(HaveOccurred())
				Expect(limit).To(Equal(20))
			})
		})

		Context("when getting the films fails", func() {
			BeforeEach(func() {
				filmCache.GetFilmsFn = func(params sakila.FilmQueryParams) ([]*sakila.Film, error) {
					return nil, sakila.ErrorNotFound
				}

				filmStore.QueryFilmsFn = func(_ sakila.FilmQueryParams) ([]*sakila.Film, error) {
					return nil, unexpectedError
				}
			})

			It("returns an internal error", func() {
				_, err := service.GetFilms(sakila.FilmQueryParams{})
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(sakila.ErrorInternal.Error()))
			})
		})
	})

	Describe("GetFilmActors", func() {
		It("returns the film actors from the cache", func() {
			var invoked bool

			filmCache.GetFilmActorsFn = func(params sakila.FilmActorParams) ([]*sakila.FilmActor, error) {
				invoked = true

				return []*sakila.FilmActor{{}}, nil
			}

			actors, err := service.Cache.GetFilmActors(sakila.FilmActorParams{})
			Expect(err).ToNot(HaveOccurred())
			Expect(invoked).To(BeTrue())
			Expect(actors).ToNot(BeNil())
			Expect(actors).To(HaveLen(1))
		})

		Context("when the actors are not cached", func() {
			BeforeEach(func() {
				filmCache.GetFilmActorsFn = func(params sakila.FilmActorParams) ([]*sakila.FilmActor, error) {
					return nil, sakila.ErrorNotFound
				}
			})

			It("returns the actors from the store", func() {
				var invoked bool

				filmStore.QueryFilmActorsFn = func(_ sakila.FilmActorParams) ([]*sakila.FilmActor, error) {
					invoked = true

					return []*sakila.FilmActor{{}, {}}, nil
				}

				actors, err := service.GetFilmActors(sakila.FilmActorParams{})
				Expect(invoked).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
				Expect(actors).To(HaveLen(2))
			})
		})

		Context("when getting the actors fails", func() {
			BeforeEach(func() {
				filmCache.GetFilmActorsFn = func(params sakila.FilmActorParams) ([]*sakila.FilmActor, error) {
					return nil, sakila.ErrorNotFound
				}

				filmStore.QueryFilmActorsFn = func(_ sakila.FilmActorParams) ([]*sakila.FilmActor, error) {
					return nil, unexpectedError
				}
			})

			It("returns an internal error", func() {
				_, err := service.GetFilmActors(sakila.FilmActorParams{})
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(sakila.ErrorInternal.Error()))
			})
		})
	})

	Describe("GetFilmCategories", func() {
		It("returns the film categories from the cache", func() {
			var invoked bool

			filmCache.GetFilmCategoriesFn = func(params sakila.FilmCategoryParams) ([]*sakila.FilmCategory, error) {
				invoked = true

				return []*sakila.FilmCategory{{}}, nil
			}

			categories, err := service.Cache.GetFilmCategories(sakila.FilmCategoryParams{})
			Expect(err).ToNot(HaveOccurred())
			Expect(invoked).To(BeTrue())
			Expect(categories).ToNot(BeNil())
			Expect(categories).To(HaveLen(1))
		})

		Context("when the categories are not cached", func() {
			BeforeEach(func() {
				filmCache.GetFilmCategoriesFn = func(params sakila.FilmCategoryParams) ([]*sakila.FilmCategory, error) {
					return nil, sakila.ErrorNotFound
				}
			})

			It("returns the catgories from the store", func() {
				var invoked bool

				filmStore.QueryFilmCategoriesFn = func(_ sakila.FilmCategoryParams) ([]*sakila.FilmCategory, error) {
					invoked = true

					return []*sakila.FilmCategory{{}, {}}, nil
				}

				categories, err := service.GetFilmCategories(sakila.FilmCategoryParams{})
				Expect(invoked).To(BeTrue())
				Expect(err).ToNot(HaveOccurred())
				Expect(categories).To(HaveLen(2))
			})
		})

		Context("when getting the actors fails", func() {
			BeforeEach(func() {
				filmCache.GetFilmActorsFn = func(params sakila.FilmActorParams) ([]*sakila.FilmActor, error) {
					return nil, sakila.ErrorNotFound
				}

				filmStore.QueryFilmActorsFn = func(_ sakila.FilmActorParams) ([]*sakila.FilmActor, error) {
					return nil, unexpectedError
				}
			})

			It("returns an internal error", func() {
				_, err := service.GetFilmActors(sakila.FilmActorParams{})
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(sakila.ErrorInternal.Error()))
			})
		})
	})
})
