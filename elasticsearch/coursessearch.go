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
func (api *API) QueryCoursesSearch(ctx context.Context, index, term string, limit, offset int) (*models.SearchResponse, int, error) {
	response := &models.SearchResponse{}

	path := api.url + "/" + index + "/_search"

	logData := log.Data{"term": term, "path": path}

	log.InfoCtx(ctx, "searching index", logData)

	body := buildSearchQuery(term, limit, offset)

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

func buildSearchQuery(term string, limit, offset int) *Body {
	var object Object
	highlight := make(map[string]Object)

	highlight["doc.kis_course_id"] = object
	highlight["doc.english_title"] = object
	highlight["doc.welsh_title"] = object
	highlight["doc.institution.public_ukprn_name"] = object

	courseID := make(map[string]string)
	englishTitle := make(map[string]string)
	welshTitle := make(map[string]string)
	institutionName := make(map[string]string)

	courseID["doc.kis_course_id"] = term
	englishTitle["doc.english_title"] = term
	welshTitle["doc.welsh_title"] = term
	institutionName["doc.institution.public_ukprn_name"] = term

	courseIDMatch := Match{
		Match: courseID,
	}

	englishTitleMatch := Match{
		Match: englishTitle,
	}

	welshTitleMatch := Match{
		Match: welshTitle,
	}

	institutionNameMatch := Match{
		Match: institutionName,
	}

	sortbyScore := Order{
		Order: "desc",
	}

	sortOrders := make(map[string]Order)
	sortOrders["_score"] = sortbyScore

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
					courseIDMatch,
					englishTitleMatch,
					welshTitleMatch,
					institutionNameMatch,
				},
			},
		},
		Sort: sortOrders,
	}

	return query
}
