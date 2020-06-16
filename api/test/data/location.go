package data

import (
	"github.com/andrewzulaybar/books/api/pkg/location"
)

// GetLocations reads the test data in `location.json` and returns it.
func GetLocations() location.Locations {
	var buffer struct {
		Locations location.Locations `json:"locations"`
	}
	loadBuffer(&buffer, "location.json")
	return buffer.Locations
}

// LoadLocations reads the test data in `location.json` and loads it into the database.
func LoadLocations(ls *location.Service) location.Locations {
	var locations location.Locations
	for _, loc := range GetLocations() {
		s, l := ls.PostLocation(&loc)
		if s.Err() == nil {
			locations = append(locations, *l)
		}
	}
	return locations
}
