package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var (
	// ErrorInvalidStationName is the error that is returned when the station name does not exist in the database.
	ErrorInvalidStationName = errors.New("invalid station name")
)

// NSStationsCodeService is the data structure for a station code
type NSStationsCodeService struct {
	cache        LFUCache
	repository   NSStationsRepository
	errorHandler ErrorHandler
}

// NewNSStationsCodeService is a service for fetching the station code based on a station name
func NewNSStationsCodeService(repository NSStationsRepository, errorHandler ErrorHandler, cache LFUCache) NSStationsCodeService {
	return NSStationsCodeService{cache, repository, errorHandler}
}

// GetCodeForStationName gets the station code for a corresponding station name.
// It's fault tolerant if you pass the station code instead of the station name it won't error
func (service *NSStationsCodeService) GetCodeForStationName(stationName string) (nsStation NSStation, error error) {
	// converting string to lowercase for consistency
	stationName = strings.ToLower(stationName)

	// Search the cache for the code
	val, err := service.cache.Get(stationName)
	if err == nil {
		return val.(NSStation), error
	}

	// Search the database for the code
	nsStation, err = service.repository.GetByName(stationName)
	if err != nil {
		// stationName does not exist find by code instead
		nsStation, err = service.repository.GetByCode(stationName)
		if err != nil {
			return nsStation, ErrorInvalidStationName
		}

		// log this error for debugging
		service.errorHandler.HandleSoftError(
			errors.New(fmt.Sprintf("GetCodeForStationName() called with short code '%s' instead of stationName name", stationName)),
		)
	}

	// the stationName code exists so update the cache
	err = service.cache.Set(nsStation.Name, nsStation)
	if err != nil {
		// log this error for debugging
		service.errorHandler.HandleSoftError(err)
	}

	return nsStation, nil
}
