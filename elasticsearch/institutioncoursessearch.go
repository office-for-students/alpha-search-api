package elasticsearch

import (
	"context"
	"encoding/json"

	"github.com/methods/go-methods-lib/log"
	errs "github.com/ofs/alpha-search-api/apierrors"
	"github.com/ofs/alpha-search-api/models"
	"github.com/pkg/errors"
)

// QueryInstitutionCoursesSearch builds query as a json body to call an elasticsearch index with
func (api *API) QueryInstitutionCoursesSearch(ctx context.Context, index, term string, filters map[string]string, countries, lengthOfCourse, institutions, subjects []string) (*models.SearchResponse, int, error) {
	response := &models.SearchResponse{}

	path := api.url + "/" + index + "/_search"

	logData := log.Data{"term": term, "path": path}

	log.InfoCtx(ctx, "searching index", logData)

	body := buildInstitutionSearchQuery(term, filters, countries, lengthOfCourse, institutions, subjects)

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

func buildInstitutionSearchQuery(term string, filters map[string]string, countries, lengthOfCourse, institutions, subjects []string) *Body {

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
		Size: 3500,
		Query: Query{
			Bool: Bool{
				Should: []Match{
					englishTitleMatch,
					welshTitleMatch,
				},
				MimimumShouldMatch: 1,
			},
		},
		Sort: []Criteria{{
			InstitutionName: "asc",
			Score:           "desc",
		}},
	}

	query = addQueryFilters(query, filters, countries, lengthOfCourse, institutions, subjects)

	return query
}
