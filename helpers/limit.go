package helpers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/methods/go-methods-lib/log"
	errs "github.com/ofs/alpha-search-api/apierrors"
	"github.com/pkg/errors"
)

// CalculateLimit returns a valid limit for the number of items to be returned from query
func CalculateLimit(ctx context.Context, defaultLimit, maximumLimit int, requestedLimit string) (int, error) {
	if requestedLimit == "" {
		return defaultLimit, nil
	}

	var errorValues = make(map[string]string)
	errorValues["limit"] = requestedLimit

	requestedLimitNumber, err := strconv.Atoi(requestedLimit)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, errs.ErrLimitWrongType.Error()), log.Data{"requested_limit": requestedLimitNumber})
		return 0, errs.New(errs.ErrLimitWrongType, http.StatusBadRequest, errorValues)
	}

	if requestedLimitNumber < 0 {
		log.ErrorCtx(ctx, errs.ErrNegativeLimit, log.Data{"requested_limit": requestedLimitNumber})
		return 0, errs.New(errs.ErrNegativeLimit, http.StatusBadRequest, errorValues)
	}

	if requestedLimitNumber > maximumLimit {
		err := fmt.Errorf("limit exceeded maximum value, limit cannot be greater than [%d]", maximumLimit)

		log.ErrorCtx(ctx, err, log.Data{"requested_limit": requestedLimitNumber})
		return 0, errs.New(err, http.StatusBadRequest, errorValues)
	}

	return requestedLimitNumber, nil
}
