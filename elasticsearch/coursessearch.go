package elasticsearch

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ONSdigital/go-ns/log"
	errs "github.com/ofs/alpha-search-api/apierrors"
	"github.com/ofs/alpha-search-api/models"
)

// QueryCoursesSearch builds query as a json body to call an elasticsearch index with
func (api *API) QueryCoursesSearch(ctx context.Context, index, term string, limit, offset int, filters map[string]string, countries, lengthOfCourse, institutions []string) (*models.SearchResponse, int, error) {
	response := &models.SearchResponse{}

	path := api.url + "/" + index + "/_search"

	logData := log.Data{"term": term, "path": path, "filters": filters}

	log.InfoCtx(ctx, "searching index", logData)

	body := buildSearchQuery(term, limit, offset, filters, countries, lengthOfCourse, institutions)

	log.InfoCtx(ctx, "searching index", log.Data{"query": body})

	bytes, err := json.Marshal(body)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "unable to marshal elastic search query to bytes"), logData)
		return nil, 0, errs.ErrMarshallingQuery
	}

	logData["request_body"] = string(bytes)

	responseBody, status, err := api.CallElastic(ctx, path, "GET", bytes)
	logData["status"] = status
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to call elasticsearch"), logData)
		return nil, status, errs.ErrIndexNotFound
	}

	logData["response_body"] = string(responseBody)

	if err = json.Unmarshal(responseBody, response); err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "unable to unmarshal json body"), logData)
		return nil, status, errs.ErrUnmarshallingJSON
	}

	log.InfoCtx(ctx, "search results", logData)

	return response, status, nil
}

func buildSearchQuery(term string, limit, offset int, filters map[string]string, countries, lengthOfCourse, institutions []string) *Body {
	var object Object
	highlight := make(map[string]Object)

	highlight["doc.english_title"] = object
	highlight["doc.welsh_title"] = object

	englishTitle := make(map[string]string)
	welshTitle := make(map[string]string)

	englishTitle["doc.english_title"] = term
	welshTitle["doc.welsh_title"] = term

	englishTitleMatch := Match{
		Match: englishTitle,
	}

	welshTitleMatch := Match{
		Match: welshTitle,
	}

	query := &Body{
		From: offset,
		Size: limit,
		Highlight: &Highlight{
			PreTags:  []string{"\u0001S"},
			PostTags: []string{"\u0001E"},
			Fields:   highlight,
		},
		Query: Query{
			Bool: Bool{
				Should: []Match{
					englishTitleMatch,
					welshTitleMatch,
				},
				MimimumShouldMatch: 1,
			},
		},
		Sort: []Criteria{
			{
				InstitutionName: "asc",
				Score:           "desc",
			},
		},
	}

	query = addQueryFilters(query, filters, countries, lengthOfCourse, institutions)

	return query
}

func addQueryFilters(query *Body, filters map[string]string, countries, lengthOfCourse, institutions []string) *Body {
	if len(filters) > 0 || len(countries) > 0 {
		query.Query.Bool.Filter = []Filters{}
	}

	for key, value := range filters {

		if key == "distance_learning" {
			if value == "true" {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							DistanceLearning: []string{"1", "2"},
						},
					},
				)
			} else {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							DistanceLearning: []string{"0", "2"},
						},
					},
				)
			}
		}

		if key == "foundation_year" {
			if value == "true" {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							FoundationYearAvailable: []string{"Optional", "Compulsory"},
						},
					},
				)
			} else {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							FoundationYearAvailable: []string{"Not available", "Optional"},
						},
					},
				)
			}
		}

		if key == "honours_award" {
			if value == "true" {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							HonoursAward: []string{"Available"},
						},
					},
				)
			} else {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							HonoursAward: []string{"Not available"},
						},
					},
				)
			}
		}

		if key == "sandwich_year" {
			if value == "true" {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							SandwichYear: []string{"Optional", "Compulsory"},
						},
					},
				)
			} else {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							SandwichYear: []string{"Not available", "Optional"},
						},
					},
				)
			}
		}

		if key == "year_abroad" {
			if value == "true" {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							YearAbroad: []string{"Optional", "Compulsory"},
						},
					},
				)
			} else {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							YearAbroad: []string{"Not available", "Optional"},
						},
					},
				)
			}
		}

		if key == "part_time" {
			if value == "true" {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							Mode: []string{"Part-time"},
						},
					},
				)
			} else {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							Mode: []string{"Full-time"},
						},
					},
				)
			}
		}

		if key == "full_time" {
			if value == "true" {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							Mode: []string{"Full-time"},
						},
					},
				)
			} else {
				query.Query.Bool.Filter = append(
					query.Query.Bool.Filter,
					Filters{
						Terms: Terms{
							Mode: []string{"Part-time"},
						},
					},
				)
			}
		}
	}

	if len(countries) > 0 {
		query.Query.Bool.Filter = append(
			query.Query.Bool.Filter,
			Filters{
				Terms: Terms{
					Country: countries,
				},
			},
		)
	}

	if len(lengthOfCourse) > 0 {
		query.Query.Bool.Filter = append(
			query.Query.Bool.Filter,
			Filters{
				Terms: Terms{
					LengthOfCourse: lengthOfCourse,
				},
			},
		)
	}

	if len(institutions) > 0 && institutions[0] != "" {
		query.Query.Bool.Filter = append(
			query.Query.Bool.Filter,
			Filters{
				Terms: Terms{
					Institutions: institutions,
				},
			},
		)
	}
	return query
}
