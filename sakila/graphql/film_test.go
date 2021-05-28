package graphql_test

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nickmro/sakila-service-film/sakila"
	"github.com/nickmro/sakila-service-film/sakila/graphql"
	"github.com/nickmro/sakila-service-film/sakila/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Data struct {
	Film  *sakila.Film   `json:"film,omitempty"`
	Films []*sakila.Film `json:"films,omitempty"`
}

var _ = Describe("Schema", func() {
	var schema *graphql.Schema
	var filmService *mock.FilmService

	BeforeEach(func() {
		filmService = &mock.FilmService{}
		s, err := graphql.NewSchema(filmService)
		if err != nil {
			panic(err)
		}
		schema = s
	})

	Describe("film", func() {
		BeforeEach(func() {
			filmService.GetFilmFn = func(ctx context.Context, filmIDs int) (*sakila.Film, error) {
				return &sakila.Film{
					FilmID: 1,
					Title:  "ACADEMY DINOSAUR",
					Description: stringP(
						"A Epic Drama of a Feminist And a Mad Scientist who must Battle a Teacher in The Canadian Rockies",
					),
					ReleaseYear:        intP(2006),
					LanguageID:         1,
					OriginalLanguageID: intP(1),
					RentalDuration:     6,
					RentalRate:         0.99,
					Length:             intP(86),
					ReplacementCost:    20.99,
					Rating:             stringP("PG"),
					SpecialFeatures:    []string{"Documentary"},
					LastUpdate:         time.Now(),
				}, nil
			}

			filmService.GetFilmActorsFn = func(ctx context.Context, filmIDs ...int) ([]*sakila.FilmActor, error) {
				return []*sakila.FilmActor{
					{
						FilmID: 1,
						Actor: sakila.Actor{
							ActorID: 1,
						},
					},
				}, nil
			}
		})

		It("returns the film", func() {
			query := `
				{
					film(filmId: 1) {
						filmId
						title
						description
						releaseYear
						languageId
						originalLanguageId
						rentalDuration
						rentalRate
						length
						replacementCost
						rating
						specialFeatures
						lastUpdate
						actors {
							actorId
						}
					}
				}
			`

			b, err := schema.Request(query)
			Expect(err).ToNot(HaveOccurred())

			data := dataFromBytes(b)
			Expect(data).ToNot(BeNil())
			Expect(data.Film).ToNot(BeNil())
			Expect(data.Film.FilmID).To(Equal(1))
			Expect(data.Film.Title).To(Equal("ACADEMY DINOSAUR"))
			Expect(data.Film.Description).ToNot(BeNil())
			Expect(*data.Film.Description).To(Equal(
				"A Epic Drama of a Feminist And a Mad Scientist who must Battle a Teacher in The Canadian Rockies",
			))
			Expect(data.Film.ReleaseYear).ToNot(BeNil())
			Expect(*data.Film.ReleaseYear).To(Equal(2006))
			Expect(data.Film.LanguageID).To(Equal(1))
			Expect(data.Film.OriginalLanguageID).ToNot(BeNil())
			Expect(*data.Film.OriginalLanguageID).To(Equal(1))
			Expect(data.Film.RentalDuration).To(Equal(6))
			Expect(data.Film.RentalRate).To(Equal(0.99))
			Expect(data.Film.Length).ToNot(BeNil())
			Expect(*data.Film.Length).To(Equal(86))
			Expect(data.Film.ReplacementCost).To(Equal(20.99))
			Expect(data.Film.Rating).ToNot(BeNil())
			Expect(*data.Film.Rating).To(Equal("PG"))
			Expect(data.Film.SpecialFeatures).ToNot(BeNil())
			Expect(data.Film.SpecialFeatures).To(HaveLen(1))
			Expect(data.Film.SpecialFeatures[0]).To(Equal("Documentary"))
			Expect(data.Film.LastUpdate).ToNot(BeZero())
			Expect(data.Film.Actors).To(HaveLen(1))
			Expect(data.Film.Actors[0].ActorID).To(Equal(1))
		})
	})

	Describe("films", func() {
		BeforeEach(func() {
			filmService.GetFilmsFn = func(ctx context.Context, params sakila.FilmParams) ([]*sakila.Film, error) {
				return []*sakila.Film{
					{
						FilmID: 1,
						Title:  "ACADEMY DINOSAUR",
						Description: stringP(
							"A Epic Drama of a Feminist And a Mad Scientist who must Battle a Teacher in The Canadian Rockies",
						),
						ReleaseYear:        intP(2006),
						LanguageID:         1,
						OriginalLanguageID: intP(1),
						RentalDuration:     6,
						RentalRate:         0.99,
						Length:             intP(86),
						ReplacementCost:    20.99,
						Rating:             stringP("PG"),
						SpecialFeatures:    []string{"Documentary"},
						LastUpdate:         time.Now(),
					},
				}, nil
			}

			filmService.GetFilmActorsFn = func(ctx context.Context, filmIDs ...int) ([]*sakila.FilmActor, error) {
				return []*sakila.FilmActor{
					{
						FilmID: 1,
						Actor: sakila.Actor{
							ActorID: 1,
						},
					},
				}, nil
			}
		})

		It("returns the films", func() {
			query := `
				{
					films {
						filmId
						title
						description
						releaseYear
						languageId
						originalLanguageId
						rentalDuration
						rentalRate
						length
						replacementCost
						rating
						specialFeatures
						lastUpdate
						actors {
							actorId
						}
					}
				}
			`

			b, err := schema.Request(query)
			Expect(err).ToNot(HaveOccurred())

			data := dataFromBytes(b)
			Expect(data).ToNot(BeNil())

			films := data.Films
			Expect(films).ToNot(BeNil())
			Expect(films).To(HaveLen(1))

			film := data.Films[0]
			Expect(film.Title).To(Equal("ACADEMY DINOSAUR"))
			Expect(film.Description).ToNot(BeNil())
			Expect(*film.Description).To(Equal(
				"A Epic Drama of a Feminist And a Mad Scientist who must Battle a Teacher in The Canadian Rockies",
			))
			Expect(film.ReleaseYear).ToNot(BeNil())
			Expect(*film.ReleaseYear).To(Equal(2006))
			Expect(film.LanguageID).To(Equal(1))
			Expect(film.OriginalLanguageID).ToNot(BeNil())
			Expect(*film.OriginalLanguageID).To(Equal(1))
			Expect(film.RentalDuration).To(Equal(6))
			Expect(film.RentalRate).To(Equal(0.99))
			Expect(film.Length).ToNot(BeNil())
			Expect(*film.Length).To(Equal(86))
			Expect(film.ReplacementCost).To(Equal(20.99))
			Expect(film.Rating).ToNot(BeNil())
			Expect(*film.Rating).To(Equal("PG"))
			Expect(film.SpecialFeatures).ToNot(BeNil())
			Expect(film.SpecialFeatures).To(HaveLen(1))
			Expect(film.SpecialFeatures[0]).To(Equal("Documentary"))
			Expect(film.LastUpdate).ToNot(BeZero())
			Expect(film.Actors).To(HaveLen(1))
			Expect(film.Actors[0].ActorID).To(Equal(1))
		})

		Context("when the 'limit', 'offset', and 'category' parameters are provided", func() {
			It("passes them to the film service", func() {
				var limit int
				var offset int

				filmService.GetFilmsFn = func(ctx context.Context, params sakila.FilmParams) ([]*sakila.Film, error) {
					limit = params.Limit
					offset = params.Offset

					return []*sakila.Film{{}}, nil
				}

				query := `
					{
						films(limit: 20, offset: 100) {
							filmId
						}
					}
				`

				_, err := schema.Request(query)
				Expect(err).NotTo(HaveOccurred())
				Expect(limit).To(Equal(20))
				Expect(offset).To(Equal(100))
			})
		})
	})
})

func stringP(s string) *string {
	return &s
}

func intP(i int) *int {
	return &i
}

func dataFromBytes(b []byte) *Data {
	var data Data

	if err := json.Unmarshal(b, &data); err != nil {
		panic(err)
	}

	return &data
}
