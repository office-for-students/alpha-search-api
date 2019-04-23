package models

import (
	"strconv"
	"strings"

	"github.com/ONSdigital/go-ns/log"
	errs "github.com/ofs/alpha-search-api/apierrors"
	"github.com/ofs/alpha-search-api/helpers"
)

// ValidateFilters checks the filters set are valid
func ValidateFilters(filters string) (map[string]string, []*ErrorObject) {
	var errorObjects []*ErrorObject
	var err error

	newFilters := make(map[string]string)
	fs := strings.Split(filters, ",")

	countFilters := make(map[string]int)
	var invalidFilters, duplicateFilters []string

	for _, filter := range fs {
		filterWithoutPrefix := strings.TrimPrefix(filter, "-")
		countFilters[filterWithoutPrefix]++

		// Check filter exists in whitelist
		// TODO could use a map instead of switch statment?
		err = checkFilterIsValid(filterWithoutPrefix)
		if err != nil {
			invalidFilters = append(invalidFilters, filter)
		}

		// Contextualise filter by using prefix to determine
		// whether filter is must exist or must not exist
		if strings.HasPrefix(filter, "-") {
			newFilters[filterWithoutPrefix] = "false"
		} else {
			newFilters[filterWithoutPrefix] = "true"
		}

		// Find duplicate filters
		if countFilters[filterWithoutPrefix] > 1 {
			duplicateFilters = append(duplicateFilters, filterWithoutPrefix)
		}
	}
	if len(invalidFilters) > 0 {
		invalifFilterList := map[string]string{"filters": helpers.StringifyWords(invalidFilters)}
		errorObjects = append(errorObjects, &ErrorObject{Error: errs.ErrInvalidFilter.Error(), ErrorValues: invalifFilterList})
	}

	if len(duplicateFilters) > 0 {
		duplicateFilterList := map[string]string{"filters": helpers.StringifyWords(duplicateFilters)}
		errorObjects = append(errorObjects, &ErrorObject{Error: errs.ErrDuplicateFilters.Error(), ErrorValues: duplicateFilterList})
	}

	// Check use of part_time and full_time filters
	_, ptFound := countFilters["part_time"]
	_, ftFound := countFilters["full_time"]
	if ptFound && ftFound {
		errorObjects = append(errorObjects, &ErrorObject{Error: errs.ErrMultipleModes.Error(), ErrorValues: map[string]string{"filters": "part_time,full_time"}})
	}

	if errorObjects != nil {
		return newFilters, errorObjects
	}

	return newFilters, nil
}

func checkFilterIsValid(filter string) error {

	switch filter {
	case "distance_learning":
	case "honours_award":
	case "foundation_year":
	case "sandwich_year":
	case "year_abroad":
	case "full_time":
	case "part_time":
	default:
		return errs.ErrInvalidFilter
	}

	return nil
}

// ValidateCountries checks the filters set are valid
func ValidateCountries(countries string) ([]string, []*ErrorObject) {
	var errorObjects []*ErrorObject
	var err error

	var mustHaveCountries, mustNotHaveCountries, invalidCountries []string
	var countryCode string

	cs := strings.Split(countries, ",")

	for _, country := range cs {
		countryCode, err = checkCountryIsValid(country)
		if err != nil {
			invalidCountries = append(invalidCountries, country)
		}

		if strings.HasPrefix(country, "-") {
			mustNotHaveCountries = append(mustNotHaveCountries, countryCode)
		} else {
			mustHaveCountries = append(mustHaveCountries, countryCode)
		}
	}

	if len(invalidCountries) > 0 {
		invalidCountryList := map[string]string{"countries": helpers.StringifyWords(invalidCountries)}
		errorObjects = append(errorObjects, &ErrorObject{Error: errs.ErrInvalidCountry.Error(), ErrorValues: invalidCountryList})
		return nil, errorObjects
	}

	if len(mustHaveCountries) > 0 {
		return mustHaveCountries, nil
	}

	mustHaveCountries = convert(mustNotHaveCountries)

	return mustHaveCountries, nil
}

func checkCountryIsValid(country string) (string, error) {
	c := strings.TrimPrefix(country, "-")

	var countryCode string
	switch c {
	case "england":
		countryCode = "XF"
	case "northern_ireland":
		countryCode = "XG"
	case "scotland":
		countryCode = "XH"
	case "wales":
		countryCode = "XI"
	default:
		return "", errs.ErrInvalidCountry
	}

	return countryCode, nil
}

func convert(mustNotHaveCountries []string) []string {
	var mustHaveCountries []string

	countries := make(map[string]bool)
	countries["XF"] = true
	countries["XG"] = true
	countries["XH"] = true
	countries["XI"] = true

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

func ValidateLengthOfCourse(lengthOfCourse string) ([]string, []*ErrorObject) {
	var errorObjects []*ErrorObject
	var newLengthOfCourse, invalidType, outOfRange []string

	loc := strings.Split(lengthOfCourse, ",")

	for _, length := range loc {
		log.Debug("length", log.Data{"length": length})
		l, err := strconv.Atoi(length)
		if err != nil {
			invalidType = append(invalidType, length)
			continue
		}

		if l < 1 || l > 7 {
			outOfRange = append(outOfRange, length)
		}

		newLengthOfCourse = append(newLengthOfCourse, length)
	}

	if len(invalidType) > 0 {
		invalidTypeList := map[string]string{"countries": helpers.StringifyWords(invalidType)}
		errorObjects = append(errorObjects, &ErrorObject{Error: errs.ErrLengthOfCourseWrongType.Error(), ErrorValues: invalidTypeList})
	}

	if len(outOfRange) > 0 {
		outOfRangeList := map[string]string{"countries": helpers.StringifyWords(outOfRange)}
		errorObjects = append(errorObjects, &ErrorObject{Error: errs.ErrLengthOfCourseOutOfRange.Error(), ErrorValues: outOfRangeList})
	}

	if errorObjects != nil {
		log.Debug("did we get here?", log.Data{"error": errorObjects})
		return newLengthOfCourse, errorObjects
	}

	return newLengthOfCourse, nil
}
