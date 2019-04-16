package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ONSdigital/go-ns/log"
	errs "github.com/ofs/alpha-search-api/apierrors"
	"github.com/ofs/alpha-search-api/models"
	"github.com/pkg/errors"
)

// SearchCourses retrieves a list of relevant results from search term
func (api *SearchAPI) SearchCourses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, contextServiceName, searchAPI)
	defer drainBody(ctx, r)

	var err error

	term := r.FormValue("q")
	filters := r.FormValue("filters")
	countries := r.FormValue("countries")
	lengthOfCourse := r.FormValue("length_of_course")
	institutions := r.FormValue("institutions")

	requestedLimit := r.FormValue("limit")
	requestedOffset := r.FormValue("offset")

	logData := log.Data{"api_config": api, "limit": requestedLimit, "offset": requestedOffset, "search_term": term}

	log.InfoCtx(ctx, "SearchCourses handler: attempting to get list of courses relevant to search term", logData)

	limit := defaultLimit
	if requestedLimit != "" {
		limit, err = strconv.Atoi(requestedLimit)
		if err != nil {
			log.ErrorCtx(ctx, errors.WithMessage(err, "search Courses endpoint: request limit parameter error"), logData)

			Error(ctx, w, errs.ErrParsingQueryParameters)
			return
		}
	}

	offset := defaultOffset
	if requestedOffset != "" {
		offset, err = strconv.Atoi(requestedOffset)
		if err != nil {
			log.ErrorCtx(ctx, errors.WithMessage(err, "search Courses endpoint: request offset parameter error"), logData)

			Error(ctx, w, errs.ErrParsingQueryParameters)
			return
		}
	}

	page := &models.PageVariables{
		DefaultMaxResults: api.DefaultMaxResults,
		Limit:             limit,
		Offset:            offset,
	}

	if err := page.ValidateQueryParameters(term); err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "search Courses endpoint: failed query parameter validation"), logData)

		Error(ctx, w, err)
		return
	}

	logData["limit"] = page.Limit
	logData["offset"] = page.Offset

	newFilters := make(map[string]string)
	if filters != "" {
		// Validate filters
		newFilters, err = models.ValidateFilters(filters)
		if err != nil {
			log.ErrorCtx(ctx, errors.WithMessage(err, "search Courses endpoint: failed filter validation"), logData)

			Error(ctx, w, err)
			return
		}
	}

	var newCountries []string
	if countries != "" {
		// Validate filter by countries
		newCountries, err = models.ValidateCountries(countries)
		if err != nil {
			log.ErrorCtx(ctx, errors.WithMessage(err, "search Courses endpoint: failed validation on country filters"), logData)

			Error(ctx, w, err)
			return
		}
	}

	var newLengthOfCourse []string
	if lengthOfCourse != "" {
		// Validate filter by length of course
		newLengthOfCourse, err = models.ValidateLengthOfCourse(lengthOfCourse)
		if err != nil {
			log.ErrorCtx(ctx, errors.WithMessage(err, "search Courses endpoint: failed validation on length of course filters"), logData)

			Error(ctx, w, err)
			return
		}
	}

	institutionList := strings.Split(institutions, ",")

	log.InfoCtx(ctx, "search Courses endpoint: just before querying search index", logData)
	// Search for courses in elasticsearch
	response, _, err := api.Elasticsearch.QueryCoursesSearch(ctx, api.Index, term, page.Limit, page.Offset, newFilters, newCountries, newLengthOfCourse, institutionList)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "search Courses endpoint: failed to query elastic search index"), logData)

		Error(ctx, w, err)
		return
	}

	searchResults := &models.CoursesSearchResults{
		TotalResults: response.Hits.Total,
		Limit:        page.Limit,
		Offset:       page.Offset,
	}

	for _, result := range response.Hits.HitList {

		result = getSnippets(ctx, result)

		doc := result.Source.Doc
		searchResults.Items = append(searchResults.Items, doc)
	}

	searchResults.Count = len(searchResults.Items)

	b, err := json.Marshal(searchResults)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "search Courses endpoint: failed to marshal search resource into bytes"), logData)

		Error(ctx, w, errs.ErrInternalServer)
		return
	}

	log.InfoCtx(ctx, "SearchCourses handler: successfully got list of course resources", logData)
	writeBody(ctx, w, b)
}

func getSnippets(ctx context.Context, result models.HitList) models.HitList {
	log.Debug("is highlight a thing", log.Data{"highlights?": result.Highlight})
	if len(result.Highlight.KISCourseID) > 0 {
		highlightedCode := result.Highlight.KISCourseID[0]
		var prevEnd int
		logData := log.Data{}
		for {
			start := prevEnd + strings.Index(highlightedCode, "\u0001S") + 1

			logData["start"] = start

			end := strings.Index(highlightedCode, "\u0001E")
			if end == -1 {
				break
			}
			logData["end"] = prevEnd + end - 2

			snippet := models.Snippet{
				Start: start,
				End:   prevEnd + end - 2,
			}

			prevEnd = snippet.End

			result.Source.Doc.Matches.KISCourseID = append(result.Source.Doc.Matches.KISCourseID, snippet)
			log.InfoCtx(ctx, "getSearch endpoint: added code snippet", logData)

			highlightedCode = string(highlightedCode[end+2:])
		}
	}

	if len(result.Highlight.EnglishTitle) > 0 {
		highlightedLabel := result.Highlight.EnglishTitle[0]
		var prevEnd int
		logData := log.Data{}
		for {
			start := prevEnd + strings.Index(highlightedLabel, "\u0001S") + 1

			logData["start"] = start

			end := strings.Index(highlightedLabel, "\u0001E")
			if end == -1 {
				break
			}
			logData["end"] = prevEnd + end - 2

			snippet := models.Snippet{
				Start: start,
				End:   prevEnd + end - 2,
			}

			prevEnd = snippet.End

			result.Source.Doc.Matches.EnglishTitle = append(result.Source.Doc.Matches.EnglishTitle, snippet)
			log.InfoCtx(ctx, "getSearch endpoint: added label snippet", logData)

			highlightedLabel = string(highlightedLabel[end+2:])
		}
	}

	return result
}
