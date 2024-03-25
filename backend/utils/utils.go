package utils

import (
	"github.com/kelvins/geocoder"
	_ "github.com/lib/pq"
)

func ReverseGeocodeToZipCode(lat, lng float64) string {
	location := geocoder.Location{
		Latitude:  lat,
		Longitude: lng,
	}
	addressList, _ := geocoder.GeocodingReverse(location)
	if len(addressList) > 0 {
		return addressList[0].PostalCode
	}
	return ""
}

func StringMissing(word string) bool {
	return word == "" || word == "Unknown" || word == "Null"
}
