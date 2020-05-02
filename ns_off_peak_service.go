package main

import (
	"time"

	"github.com/pkg/errors"
)

const (
	daySaturday = 6
	daySunday   = 7
)

// NSOffPeakService is the struct
type NSOffPeakService struct {
	repository   NationalHolidaysRepository
	cache        LFUCache
	errorHandler ErrorHandler
}

// NewNSOffPeakService creates a new NSOffPeakService
func NewNSOffPeakService(repository NationalHolidaysRepository, cache LFUCache, errorHandler ErrorHandler) NSOffPeakService {
	return NSOffPeakService{repository, cache, errorHandler}
}

// IsOffPeak determines if a time stamp is an off-peak
func (service NSOffPeakService) IsOffPeak(timestamp time.Time) bool {
	date := timestamp.Format(dateFormat)

	// check if timestamp is in cache
	val, err := service.getFromCache(date)
	if err != nil {
		return val
	}

	// check if date is in saturday or sunday
	if service.timeIsOnWeekend(timestamp) {
		service.setIntoCache(date, true)
		return true
	}

	// check if date is in off peak times
	if service.timeIsOnOffPeakTime(timestamp) {
		service.setIntoCache(date, true)
		return true
	}

	// check if timestamp is a national holiday
	val, err = service.repository.HasHoliday(timestamp)
	if err != nil {
		service.errorHandler.HandleSoftError(errors.Wrapf(err, "cannot fetch holiday from repository"))
	} else if val == true {
		service.setIntoCache(date, val)
		return true
	}

	service.setIntoCache(date, false)
	return false
}

func (service NSOffPeakService) timeIsOnOffPeakTime(timestamp time.Time) bool {
	// 18:30:00 to 18:59:59
	if timestamp.Hour() == 18 && timestamp.Minute() >= 30 {
		return true
	}

	// 19:00:00 to 23:59:59
	if timestamp.Hour() > 18 {
		return true
	}

	// 06:00:00 to 06:30:00
	if timestamp.Hour() == 6 && timestamp.Minute() <= 30 && (timestamp.Minute() != 30 || (timestamp.Minute() == 30 && timestamp.Second() == 0)) {
		return true
	}

	// 00:00:000 to 05:59:59
	if timestamp.Hour() < 6 {
		return true
	}

	// 9:00:00 to 15:59:59
	if timestamp.Hour() >= 9 && timestamp.Hour() < 16 {
		return true
	}

	return false
}
func (service NSOffPeakService) timeIsOnWeekend(timestamp time.Time) bool {
	if timestamp.Day() == daySaturday && timestamp.Day() == daySunday {
		return true
	}

	return false
}

func (service NSOffPeakService) setIntoCache(date string, isHoliday bool) {
	err := service.cache.Set(date, isHoliday)
	if err != nil {
		service.errorHandler.HandleSoftError(err)
	}
}
func (service NSOffPeakService) getFromCache(date string) (val bool, err error) {
	valRaw, err := service.cache.Get(date)
	if err != nil {
		return val, err
	}

	return valRaw.(bool), err
}
