package models

import (
	"strconv"
	"strings"

	errs "github.com/ofs/alpha-search-api/apierrors"
)

// ValidateFilters checks the filters set are valid
func ValidateFilters(filters string) (map[string]string, error) {
	var err error
	newFilters := make(map[string]string)
	fs := strings.Split(filters, ",")

	hasMode := false
	for _, filter := range fs {
		hasMode, err = checkFilterIsValid(filter, hasMode)
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(filter, "-") {
			f := strings.Split(filter, "-")
			newFilters[f[1]] = "false"
		} else {
			newFilters[filter] = "true"
		}
	}

	return newFilters, nil
}

func checkFilterIsValid(filter string, hasMode bool) (bool, error) {
	f := strings.TrimPrefix(filter, "-")

	switch f {
	case "distance_learning":
	case "honours_award":
	case "foundation_year":
	case "sandwich_year":
	case "year_abroad":
	case "full_time":
		if hasMode == true {
			return hasMode, errs.ErrBadFilter
		}
		hasMode = true
	case "part_time":
		if hasMode == true {
			return hasMode, errs.ErrBadFilter
		}
		hasMode = true
	default:
		return hasMode, errs.ErrBadFilter
	}

	return hasMode, nil
}

// ValidateCountries checks the filters set are valid
func ValidateCountries(countries string) ([]string, error) {
	var err error
	var mustHaveCountries, mustNotHaveCountries []string
	var countryCode string

	cs := strings.Split(countries, ",")

	hasCountry := make(map[string]bool)
	for _, country := range cs {
		countryCode, hasCountry, err = checkCountryIsValid(country, hasCountry)
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(country, "-") {
			mustNotHaveCountries = append(mustNotHaveCountries, countryCode)
		} else {
			mustHaveCountries = append(mustHaveCountries, countryCode)
		}
	}

	if len(mustHaveCountries) > 0 {
		return mustHaveCountries, nil
	}

	mustHaveCountries = convert(mustNotHaveCountries)

	return mustHaveCountries, nil
}

func checkCountryIsValid(country string, hasCountry map[string]bool) (string, map[string]bool, error) {
	c := strings.TrimPrefix(country, "-")

	var countryCode string
	switch c {
	case "england":
		countryCode = "1"
		if hasCountry[countryCode] == true {
			return countryCode, hasCountry, errs.ErrBadFilter
		}
		hasCountry[countryCode] = true
	case "northern_ireland":
		countryCode = "2"
		if hasCountry[countryCode] == true {
			return "", hasCountry, errs.ErrBadFilter
		}
		hasCountry[countryCode] = true
	case "scotland":
		countryCode = "3"
		if hasCountry[countryCode] == true {
			return "", hasCountry, errs.ErrBadFilter
		}
		hasCountry[countryCode] = true
	case "wales":
		countryCode = "4"
		if hasCountry[countryCode] == true {
			return "", hasCountry, errs.ErrBadFilter
		}
		hasCountry[countryCode] = true
	default:
		return "", hasCountry, errs.ErrBadFilter
	}

	return countryCode, hasCountry, nil
}

func convert(mustNotHaveCountries []string) []string {
	var mustHaveCountries []string

	countries := make(map[string]bool)
	countries["1"] = true
	countries["2"] = true
	countries["3"] = true
	countries["4"] = true

	for _, country := range mustNotHaveCountries {
		countries[country] = false
	}

	for country, ok := range countries {
		if ok {
			mustHaveCountries = append(mustHaveCountries, country)
		}
	}

	return mustHaveCountries
}

func ValidateLengthOfCourse(lengthOfCourse string) ([]string, error) {
	var newLengthOfCourse []string
	loc := strings.Split(lengthOfCourse, ",")

	for _, length := range loc {
		l, err := strconv.Atoi(length)
		if err != nil {
			return newLengthOfCourse, err
		}

		if l < 1 && l > 7 {
			return nil, errs.ErrBadFilter
		}

		newLengthOfCourse = append(newLengthOfCourse, length)
	}

	return newLengthOfCourse, nil
}
