package models

// InstitutionCoursesSearchResponse represents a structure for a list of returned objects
type InstitutionCoursesSearchResult struct {
	Count  int           `json:"count"`
	Items  []Institution `json:"items"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

// Institution represents institution data of a single item in returned list
type Institution struct {
	PublicUKPRN     string     `json:"public_ukprn"`
	PublicUKPRNName string     `json:"public_ukprn_name"`
	UKPRN           string     `json:"ukprn"`
	UKPRNName       string     `json:"ukprn_name"`
	Courses         []Document `json:"courses,omitempty"`
}

// LocationObject represents location data of a single item in returned list
type LocationObject struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// Qualification represents the qualification data of a single item in returned list
type Qualification struct {
	Code  string `json:"code,omitempty"`
	Label string `json:"label,omitempty"`
	Level string `json:"level,omitempty"`
	Name  string `json:"name,omitempty"`
}
