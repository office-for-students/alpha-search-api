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
	ErrLimitWrongType           = errors.New("limit value needs to be a number")
	ErrNegativeLimit            = errors.New("limit needs to be a positive number, limit cannot be lower than 0")
	ErrOffsetWrongType          = errors.New("offset value needs to be a number")
	ErrNegativeOffset           = errors.New("offset needs to be a positive number, offset cannot be lower than 0")
	ErrMultipleModes            = errors.New("cannot have both part-time and full-time filters set")
	ErrInvalidFilter            = errors.New("invalid filters")
	ErrDuplicateFilters         = errors.New("use of the same filter option more than once")
	ErrInvalidCountry           = errors.New("invalid countries")
	ErrLengthOfCourseWrongType  = errors.New("length_of_course values needs to be a number")
	ErrLengthOfCourseOutOfRange = errors.New("length_of_course values needs to be numbers between the range of 1 and 7")
	ErrEmptySearchTerm          = errors.New("empty search term")

	ErrCourseNotFound         = errors.New("course not found")
	ErrIndexNotFound          = errors.New("search index not found")
	ErrInstitutionNotFound    = errors.New("institution not found")
	ErrInternalServer         = errors.New("internal server error")
	ErrMarshallingQuery       = errors.New("failed to marshal query to bytes for request body to send to elastic")
	ErrParsingQueryParameters = errors.New("failed to parse query parameters, values must be an integer")
	ErrUnmarshallingJSON      = errors.New("failed to parse json body")
	ErrUnexpectedStatusCode   = errors.New("unexpected status code from elastic api")

	NotFoundMap = map[error]bool{
		ErrCourseNotFound:      true,
		ErrInstitutionNotFound: true,
	}
)
