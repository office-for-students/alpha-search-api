package models

import (
	"errors"
	"strconv"

	errs "github.com/ofs/alpha-search-api/apierrors"
)

func ErrorMaximumOffsetReached(m int) error {
	err := errors.New("the maximum offset has been reached, the offset cannot be more than " + strconv.Itoa(m))
	return err
}

type SearchResponse struct {
	Hits Hits `json:"hits"`
}

type Hits struct {
	Total   int       `json:"total"`
	HitList []HitList `json:"hits"`
}

type HitList struct {
	Highlight Highlight    `json:"highlight"`
	Score     float64      `json:"_score"`
	Source    SearchResult `json:"_source"`
}

type Highlight struct {
	KISCourseID     []string `json:"kis_course_id,omitempty"`
	EnglishTitle    []string `json:"english_title,omitempty"`
	WelshTitle      []string `json:"welsh_title,omitempty"`
	InstitutionName []string `json:"institution_public_name,omitempty"`
}

// CoursesSearchResults represents a structure for a list of returned objects
type CoursesSearchResults struct {
	TotalResults int        `json:"total_results"`
	Count        int        `json:"number_of_items"`
	Items        []Document `json:"items"`
	Limit        int        `json:"limit"`
	Offset       int        `json:"offset"`
}

// SearchResult represents data on a single item of search results
type SearchResult struct {
	Doc Document `json:"doc"`
}

// Document is the nested document in single search result
type Document struct {
	Score            float64         `json:"score,omitempty"`
	SortName         string          `json:"institution_name,omitempty"`
	KISCourseID      string          `json:"kis_course_id"`
	EnglishTitle     string          `json:"english_title"`
	Country          string          `json:"country"`
	DistanceLearning string          `json:"distance_learning,omitempty"`
	FoundationYear   string          `json:"foundation_year"`
	HonoursAward     string          `json:"honours_award"`
	Institution      *Institution    `json:"institution"`
	LengthOfCourse   string          `json:"length_of_course"`
	Link             string          `json:"link"`
	Location         *LocationObject `json:"location,omitempty"`
	Matches          Matches         `json:"matches,omitempty"`
	Mode             string          `json:"mode"`
	NHSFunded        string          `json:"nhs_funded,omitempty"`
	Qualification    *Qualification  `json:"qualification,omitempty"`
	SandwichYear     string          `json:"sandwich_year,omitempty"`
	YearAbroad       string          `json:"year_abroad,omitempty"`
}

// Matches represents a list of members and their arrays of character offsets that matched the search term
type Matches struct {
	KISCourseID     []Snippet `json:"kis_course_id,omitempty"`
	EnglishTitle    []Snippet `json:"english_title,omitempty"`
	WelshTitle      []Snippet `json:"welsh_title,omitempty"`
	InstitutionName []Snippet `json:"institution.public_ukprn_name,omitempty"`
}

// Snippet represents a pair of integers defining the start and end of a substring in the member that matched the search term
type Snippet struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// PageVariables are the necessary fields to determine paging
type PageVariables struct {
	DefaultMaxResults int
	Limit             int
	Offset            int
}

// ValidateQueryParameters represents a model for validating query parameters
func (page *PageVariables) ValidateQueryParameters(term string) []*ErrorObject {
	var errorObjects []*ErrorObject

	if term == "" {
		termErrorValue := make(map[string](string))
		termErrorValue["q"] = term
		errorObjects = append(errorObjects, &ErrorObject{Error: errs.ErrEmptySearchTerm.Error(), ErrorValues: termErrorValue})
	}

	if page.Offset >= page.DefaultMaxResults {
		pagingErrorValue := make(map[string](string))
		pagingErrorValue["offset"] = strconv.Itoa(page.Offset)
		errorObjects = append(errorObjects, &ErrorObject{Error: ErrorMaximumOffsetReached(page.DefaultMaxResults).Error(), ErrorValues: pagingErrorValue})
	}

	if errorObjects != nil {
		return errorObjects
	}

	if page.Offset+page.Limit > page.DefaultMaxResults {
		page.Limit = page.DefaultMaxResults - page.Offset
	}

	return nil
}
