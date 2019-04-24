package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ONSdigital/go-ns/log"
	errs "github.com/ofs/alpha-search-api/apierrors"
	"github.com/ofs/alpha-search-api/helpers"
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
	filters := r.FormValue("filters")
	countries := r.FormValue("countries")
	lengthOfCourse := r.FormValue("length_of_course")
	institutions := r.FormValue("institutions")

	requestedLimit := r.FormValue("limit")
	requestedOffset := r.FormValue("offset")

	logData := log.Data{"api_config": api, "limit": requestedLimit, "offset": requestedOffset, "search_term": term}

	log.InfoCtx(ctx, "SearchInstitutionCourses handler: attempting to get list of courses relevant to search term", logData)

	var errorObjects []*models.ErrorObject

	limit, err := helpers.CalculateLimit(ctx, defaultLimit, api.DefaultMaxResults, requestedLimit)
	if err != nil {
		errorObjects = append(errorObjects, &models.ErrorObject{Error: err.Error(), ErrorValues: err.(*errs.ErrorObject).Values()})
	}

	offset, err := helpers.CalculateOffset(ctx, requestedOffset)
	if err != nil {
		errorObjects = append(errorObjects, &models.ErrorObject{Error: err.Error(), ErrorValues: err.(*errs.ErrorObject).Values()})
	}

	page := &models.PageVariables{
		DefaultMaxResults: api.DefaultMaxResults,
		Limit:             limit,
		Offset:            offset,
	}

	if errorObject := page.ValidateQueryParameters(term); errorObject != nil {
		errorObjects = append(errorObjects, errorObject...)
	}

	logData["limit"] = page.Limit
	logData["offset"] = page.Offset

	newFilters := make(map[string]string)
	if filters != "" {
		var filterErrorObject []*models.ErrorObject

		// Validate filters
		newFilters, filterErrorObject = models.ValidateFilters(filters)
		if filterErrorObject != nil {
			errorObjects = append(errorObjects, filterErrorObject...)
		}
	}

	var newCountries []string
	if countries != "" {
		var countryErrorObject []*models.ErrorObject

		// Validate filter by countries
		newCountries, countryErrorObject = models.ValidateCountries(countries)
		if countryErrorObject != nil {
			errorObjects = append(errorObjects, countryErrorObject...)
		}
	}

	var newLengthOfCourse []string
	if lengthOfCourse != "" {
		var lengthOfCourseErrorObject []*models.ErrorObject

		// Validate filter by length of course
		newLengthOfCourse, lengthOfCourseErrorObject = models.ValidateLengthOfCourse(lengthOfCourse)
		if lengthOfCourseErrorObject != nil {
			errorObjects = append(errorObjects, lengthOfCourseErrorObject...)
		}
	}

	if errorObjects != nil {
		ErrorResponse(ctx, w, http.StatusBadRequest, &models.ErrorResponse{Errors: errorObjects})
		return
	}

	institutionList := strings.Split(strings.ToLower(institutions), ",")

	log.InfoCtx(ctx, "search Institution courses endpoint: just before querying search index", logData)
	// Search for courses in elasticsearch
	response, _, err := api.Elasticsearch.QueryInstitutionCoursesSearch(ctx, api.Index, term, newFilters, newCountries, newLengthOfCourse, institutionList)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "search Institution courses endpoint: failed to query elastic search index"), logData)

		Error(ctx, w, err)
		return
	}

	items := groupCoursesByInstitution(api.ShowScore, response)

	searchResults := &models.InstitutionCoursesSearchResult{
		Limit:        page.Limit,
		Offset:       page.Offset,
		TotalResults: len(items),
	}

	if len(items) > (page.Offset + page.Limit) {
		upper := page.Offset + page.Limit
		searchResults.Items = items[page.Offset:upper]
		searchResults.Count = len(searchResults.Items)
	} else {
		searchResults.Count = len(items)
		searchResults.Items = items
	}

	b, err := json.Marshal(searchResults)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "search Institution courses endpoint: failed to marshal search resource into bytes"), logData)

		Error(ctx, w, errs.ErrInternalServer)
		return
	}

	log.InfoCtx(ctx, "SearchInstitutionCourses handler: successfully got list of course resources", logData)
	writeBody(ctx, w, b)
}

func groupCoursesByInstitution(showScore bool, response *models.SearchResponse) (institutions []models.Institution) {
	institutionCourses := make(map[string][]models.Document)
	institutionOrder := make(map[int]models.Institution)

	count := 0
	for _, result := range response.Hits.HitList {
		doc := result.Source.Doc
		if showScore {
			doc.Score = result.Score
		}

		if _, ok := institutionCourses[doc.Institution.UKPRNName]; !ok {
			// store top level fields here
			institutionOrder[count] = models.Institution{
				PublicUKPRN:     doc.Institution.PublicUKPRN,
				PublicUKPRNName: doc.Institution.PublicUKPRNName,
				UKPRN:           doc.Institution.UKPRN,
				UKPRNName:       doc.Institution.UKPRNName,
			}
			count++
		}

		UKPRNName := doc.Institution.UKPRNName
		// TODO Remove nested institution doc from course object
		doc.Institution = nil

		// Add course to set
		var courses []models.Document

		courses = append(courses, institutionCourses[UKPRNName]...)
		courses = append(courses, doc)

		institutionCourses[UKPRNName] = courses
	}

	for i := 0; i < len(institutionOrder); i++ {
		institutionName := institutionOrder[i].UKPRNName

		institution := models.Institution{
			PublicUKPRN:     institutionOrder[i].PublicUKPRN,
			PublicUKPRNName: institutionOrder[i].PublicUKPRNName,
			UKPRN:           institutionOrder[i].UKPRN,
			UKPRNName:       institutionName,
			Count:           len(institutionCourses[institutionName]),
			Courses:         institutionCourses[institutionName],
		}

		if showScore {
			institution.Score = institutionCourses[institutionName][0].Score
		}

		institutions = append(institutions, institution)
	}

	return
}
