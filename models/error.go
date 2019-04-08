package models

import errs "github.com/ofs/alpha-search-api/apierrors"

// ErrorResponse builds a list of errors for an unsuccessful request
type ErrorResponse struct {
	Errors []*ErrorObject `json:"errors"`
}

// ErrorObject contains an error message and error values
type ErrorObject struct {
	Error       string            `json:"error"`
	ErrorValues map[string]string `json:"error_values,omitempty"`
}

// CreateErrorObject formulates an error object from an error
func CreateErrorObject(err error) *ErrorObject {

	return &ErrorObject{Error: err.Error(), ErrorValues: err.(*errs.ErrorObject).Values()}
}
