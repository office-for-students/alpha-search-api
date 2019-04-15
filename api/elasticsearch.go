package api

import (
	"context"

	"github.com/ofs/alpha-search-api/models"
)

// Elasticsearcher - An interface used to access elasticsearch
type Elasticsearcher interface {
	QueryCoursesSearch(ctx context.Context, index, term string, limit, offset int) (*models.SearchResponse, int, error)
	QueryInstitutionCoursesSearch(ctx context.Context, index, term string) (*models.SearchResponse, int, error)
}
