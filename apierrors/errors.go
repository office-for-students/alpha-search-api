package apierrors

import (
	"errors"
)

// New returns an error that formats as the given text.
func New(err error, status int, values map[string]string) error {
	return &ErrorObject{
		Code:    status,
		Keys:    values,
		Message: err.Error(),
	}
}

// ErrorObject is a trivial implementation of error.
type ErrorObject struct {
	Code    int
	Keys    map[string]string
	Message string
}

func (e *ErrorObject) Error() string {
	return e.Message
}

// Status represents the status code to return from error
func (e *ErrorObject) Status() int {
	return e.Code
}

// Values represents a list of key value pairs to return from error
func (e *ErrorObject) Values() map[string]string {
	return e.Keys
}

// A list of error messages for Dataset API
var (
	ErrCourseNotFound         = errors.New("course not found")
	ErrIndexNotFound          = errors.New("search index not found")
	ErrInstitutionNotFound    = errors.New("institution not found")
	ErrInternalServer         = errors.New("internal server error")
	ErrMarshallingQuery       = errors.New("failed to marshal query to bytes for request body to send to elastic")
	ErrParsingQueryParameters = errors.New("failed to parse query parameters, values must be an integer")
	ErrEmptySearchTerm        = errors.New("empty search term")
	ErrUnmarshallingJSON      = errors.New("failed to parse json body")
	ErrUnexpectedStatusCode   = errors.New("unexpected status code from elastic api")

	NotFoundMap = map[error]bool{
		ErrCourseNotFound:      true,
		ErrInstitutionNotFound: true,
	}
)
