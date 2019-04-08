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

// SearchResults represents a structure for a list of returned objects
type SearchResults struct {
	Count  int        `json:"count"`
	Items  []Document `json:"items"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

// SearchResult represents data on a single item of search results
type SearchResult struct {
	Doc Document `json:"doc"`
}

// Document is the nested document in single search result
type Document struct {
	KISCourseID      string          `json:"kis_course_id"`
	EnglishTitle     string          `json:"english_title"`
	DistanceLearning string          `json:"distance_learning,omitempty"`
	FoundationYear   string          `json:"foundation_year"`
	Institution      *Institution    `json:"institution"`
	Link             string          `json:"link"`
	Location         *LocationObject `json:"location,omitempty"`
	Matches          Matches         `json:"matches,omitempty"`
	Mode             string          `json:"mode"`
	NHSFunded        string          `json:"nhs_funded,omitempty"`
	Qualification    *Qualification  `json:"qualification,omitempty"`
	NumberOfChildren int             `json:"number_of_children"`
	SandwichYear     string          `json:"sandwich_year,omitempty"`
	YearAbroad       string          `json:"year_abroad,omitempty"`
}

// Institution represents institution data of a single item in returned list
type Institution struct {
	PublicUKPRN     string `json:"public_ukprn"`
	PublicUKPRNName string `json:"public_ukprn_name"`
	UKPRN           string `json:"ukprn"`
	UKPRNName       string `json:"ukprn_name"`
}

// LocationObject represents location data of a single item in returned list
type LocationObject struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// Matches represents a list of members and their arrays of character offsets that matched the search term
type Matches struct {
	KISCourseID     []Snippet `json:"kis_course_id,omitempty"`
	EnglishTitle    []Snippet `json:"english_title,omitempty"`
	WelshTitle      []Snippet `json:"welsh_title,omitempty"`
	InstitutionName []Snippet `json:"institution.public_ukprn_name,omitempty"`
}

// Qualification represents the qualification data of a single item in returned list
type Qualification struct {
	Code  string `json:"code,omitempty"`
	Label string `json:"label,omitempty"`
	Level string `json:"level,omitempty"`
	Name  string `json:"name,omitempty"`
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
func (page *PageVariables) ValidateQueryParameters(term string) error {
	if term == "" {
		return errs.ErrEmptySearchTerm
	}

	if page.Offset >= page.DefaultMaxResults {
		return ErrorMaximumOffsetReached(page.DefaultMaxResults)
	}

	if page.Offset+page.Limit > page.DefaultMaxResults {
		page.Limit = page.DefaultMaxResults - page.Offset
	}

	return nil
}