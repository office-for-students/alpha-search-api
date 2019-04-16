package elasticsearch

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ONSdigital/go-ns/log"
	errs "github.com/ofs/alpha-search-api/apierrors"
	"github.com/pkg/errors"
	awsauth "github.com/smartystreets/go-aws-auth"
)

// API aggregates a client and URL and other common data for accessing the API
type API struct {
	client       http.Client
	url          string
	signRequests bool
}

// NewElasticSearchAPI creates an API object
func NewElasticSearchAPI(client http.Client, elasticSearchAPIURL string, signRequests bool) *API {
	return &API{
		client:       client,
		url:          elasticSearchAPIURL,
		signRequests: signRequests,
	}
}

// Body represents the request body to elasticsearch
type Body struct {
	From      int              `json:"from"`
	Size      int              `json:"size"`
	Highlight *Highlight       `json:"highlight,omitempty"`
	Query     Query            `json:"query"`
	Sort      map[string]Order `json:"sort"`
}

// Highlight represents parts of the fields that matched
type Highlight struct {
	PreTags  []string          `json:"pre_tags,omitempty"`
	PostTags []string          `json:"post_tags,omitempty"`
	Fields   map[string]Object `json:"fields,omitempty"`
	Order    string            `json:"score,omitempty"`
}

// Object represents an empty object (as expected by elasticsearch)
type Object struct{}

// Query represents the request query details
type Query struct {
	Bool Bool `json:"bool"`
}

// Bool represents the desirable goals for query
type Bool struct {
	Must   []Match   `json:"must,omitempty"`
	Should []Match   `json:"should,omitempty"`
	Filter []Filters `json:"filter,omitempty"`
}

// Filters represents a list of items that are filterable
type Filters struct {
	Terms Terms `json:"terms,omitempty"`
}

// Terms represents a list of terms that are filterable
type Terms struct {
	Country                 []string `json:"doc.country_code,omitempty"`
	DistanceLearning        []string `json:"doc.distance_learning,omitempty"`
	FoundationYearAvailable []string `json:"doc.foundation_year,omitempty"`
	HonoursAward            []string `json:"doc.honours_award,omitempty"`
	Institutions            []string `json:"doc.institution.ukprn_name,omitempty"`
	LengthOfCourse          []string `json:"doc.length_of_course,omitempty"`
	Mode                    []string `json:"doc.mode,omitempty"`
	SandwichYear            []string `json:"doc.sandwich_year,omitempty"`
	YearAbroad              []string `json:"doc.year_abroad,omitempty"`
}

// Match represents the fields that the term should or must match within query
type Match struct {
	Match map[string]string `json:"match,omitempty"`
}

// Order contains the ordering (ascending or descending) on a particular field
type Order struct {
	Order string `json:"order"`
}

// CallElastic builds a request to elastic search based on the method, path and payload
func (api *API) CallElastic(ctx context.Context, path, method string, payload interface{}) ([]byte, int, error) {
	logData := log.Data{"url": path, "method": method}

	URL, err := url.Parse(path)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to create url for elastic call"), logData)
		return nil, 0, err
	}
	path = URL.String()
	logData["url"] = path

	var req *http.Request

	if payload != nil {
		req, err = http.NewRequest(method, path, bytes.NewReader(payload.([]byte)))
		req.Header.Add("Content-type", "application/json")
		logData["payload"] = string(payload.([]byte))
	} else {
		req, err = http.NewRequest(method, path, nil)
	}
	// check req, above, didn't error
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to create request for call to elastic"), logData)
		return nil, 0, err
	}

	if api.signRequests {
		awsauth.Sign(req)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to call elastic"), logData)
		return nil, 0, err
	}
	defer resp.Body.Close()

	logData["status_code"] = resp.StatusCode

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to read response body from call to elastic"), logData)
		return nil, resp.StatusCode, err
	}
	logData["json_body"] = string(jsonBody)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		log.ErrorCtx(ctx, errs.ErrUnexpectedStatusCode, logData)
		return nil, resp.StatusCode, errs.ErrUnexpectedStatusCode
	}

	return jsonBody, resp.StatusCode, nil
}
