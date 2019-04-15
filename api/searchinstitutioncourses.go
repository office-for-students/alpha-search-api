package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ONSdigital/go-ns/log"
	errs "github.com/ofs/alpha-search-api/apierrors"
	"github.com/ofs/alpha-search-api/models"
	"github.com/pkg/errors"
)

// SearchInstitutionCourses retrieves a list of relevant results from search term
func (api *SearchAPI) SearchInstitutionCourses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, contextServiceName, searchAPI)
	defer drainBody(ctx, r)

	var err error

	term := r.FormValue("q")
	requestedLimit := r.FormValue("limit")
	requestedOffset := r.FormValue("offset")

	logData := log.Data{"api_config": api, "limit": requestedLimit, "offset": requestedOffset, "search_term": term}

	log.InfoCtx(ctx, "SearchInstitutionCourses handler: attempting to get list of courses relevant to search term", logData)

	limit := defaultLimit
	if requestedLimit != "" {
		limit, err = strconv.Atoi(requestedLimit)
		if err != nil {
			log.ErrorCtx(ctx, errors.WithMessage(err, "search Institution courses endpoint: request limit parameter error"), logData)

			Error(ctx, w, errs.ErrParsingQueryParameters)
			return
		}
	}

	offset := defaultOffset
	if requestedOffset != "" {
		offset, err = strconv.Atoi(requestedOffset)
		if err != nil {
			log.ErrorCtx(ctx, errors.WithMessage(err, "search Institution courses endpoint: request offset parameter error"), logData)

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
		log.ErrorCtx(ctx, errors.WithMessage(err, "search Institution courses endpoint: failed query parameter validation"), logData)

		Error(ctx, w, err)
		return
	}

	logData["limit"] = page.Limit
	logData["offset"] = page.Offset

	log.InfoCtx(ctx, "search Institution courses endpoint: just before querying search index", logData)
	// Search for courses in elasticsearch
	response, _, err := api.Elasticsearch.QueryInstitutionCoursesSearch(ctx, api.Index, term)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "search Institution courses endpoint: failed to query elastic search index"), logData)

		Error(ctx, w, err)
		return
	}

	institutions := groupCoursesByInstitution(response)

	searchResults := &models.InstitutionCoursesSearchResult{
		Count:  response.Hits.Total,
		Limit:  page.Limit,
		Offset: page.Offset,
		Items:  institutions,
	}

	searchResults.Count = len(searchResults.Items)

	b, err := json.Marshal(searchResults)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "search Institution courses endpoint: failed to marshal search resource into bytes"), logData)

		Error(ctx, w, errs.ErrInternalServer)
		return
	}

	log.InfoCtx(ctx, "SearchInstitutionCourses handler: successfully got list of course resources", logData)
	writeBody(ctx, w, b)
}

func groupCoursesByInstitution(response *models.SearchResponse) (institutions []models.Institution) {
	institutionMap := make(map[string][]models.Document)
	institutionOrder := make(map[int]string)

	count := 0
	for _, result := range response.Hits.HitList {
		doc := result.Source.Doc

		if _, ok := institutionMap[doc.Institution.UKPRNName]; !ok {
			institutionOrder[count] = doc.Institution.UKPRNName
			count++
		}

		// Add course to set
		var courses []models.Document

		courses = append(courses, institutionMap[doc.Institution.UKPRNName]...)
		courses = append(courses, doc)

		institutionMap[doc.Institution.UKPRNName] = courses
	}

	for _, institutionName := range institutionOrder {

		institution := models.Institution{
			UKPRNName: institutionName,
			Courses:   institutionMap[institutionName],
		}

		institutions = append(institutions, institution)
	}

	return
}
