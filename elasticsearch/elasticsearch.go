package elasticsearch

import (
	"context"

	"github.com/ofs/alpha-search-api/models"
)

// Elasticsearcher - An interface used to access elasticsearch
type Elasticsearcher interface {
	QuerySearchIndex(ctx context.Context, index, term string, limit, offset int) (*models.SearchResponse, int, error)
}
