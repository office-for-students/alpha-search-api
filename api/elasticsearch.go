package api

import (
	"context"

	"github.com/ofs/alpha-search-api/models"
)

// Elasticsearcher - An interface used to access elasticsearch
type Elasticsearcher interface {
	QueryCoursesSearch(ctx context.Context, index, term string, limit, offset int, listOfFilters map[string]string, listOfCountries, listOfLengthOfCourses, institutionList, subjects []string) (*models.SearchResponse, int, error)
	QueryInstitutionCoursesSearch(ctx context.Context, index, term string, filters map[string]string, countries, lengthOfCourse, institutions, subjects []string) (*models.SearchResponse, int, error)
}
