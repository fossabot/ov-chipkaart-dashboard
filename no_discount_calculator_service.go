package main

import (
	"github.com/pkg/errors"
)

const (
	supplementPriceOffPeak = 156
	supplementPricePeak    = 262
)

// NSNoDiscountCalculator calculates the price of a journey when there are no discounts.
type NSNoDiscountCalculator struct {
	priceFetcher   NSPriceFetcherService
	offPeakService NSOffPeakService
}

// NewNSNoDiscountCalculator creates a new instance of an NSNoDiscountCalculator
func NewNSNoDiscountCalculator(priceFetcher NSPriceFetcherService, offPeakService NSOffPeakService) *NSNoDiscountCalculator {
	return &NSNoDiscountCalculator{
		priceFetcher:   priceFetcher,
		offPeakService: offPeakService,
	}
}

//NSNoDiscountCalculatorResult represents the calculation result of an NS Journey
type NSNoDiscountCalculatorResult struct {
	FirstClassPrice        int
	SecondClassPrice       int
	PeakSupplementPrice    int
	PeakSupplementCount    int
	OffPeakSupplementPrice int
	OffPeakSupplementCount int
	JourneyCount           int
	Error                  EnrichedRecordsError
}

// addJourneyPrice adds the price of an NSJourney to the result
func (result *NSNoDiscountCalculatorResult) addJourneyPrice(journey NSJourneyPrice) {
	result.FirstClassPrice += journey.FirstClassSingleFarePrice
	result.SecondClassPrice += journey.SecondClassSingleFarePrice
	result.JourneyCount++
}

// incrementPeakSupplement adds the peak supplement price
func (result *NSNoDiscountCalculatorResult) incrementPeakSupplement() {
	result.PeakSupplementCount++
	result.PeakSupplementPrice += supplementPricePeak
}

// incrementOffPeakSupplement adds the off peak supplement price
func (result *NSNoDiscountCalculatorResult) incrementOffPeakSupplement() {
	result.OffPeakSupplementCount++
	result.OffPeakSupplementPrice += supplementPriceOffPeak
}

// SupplementPrice returns the price of both off peak and peak supplement
func (result NSNoDiscountCalculatorResult) SupplementPrice() int {
	return result.OffPeakSupplementPrice + result.PeakSupplementPrice
}

// SupplementCount returns the total count of all supplements.
func (result NSNoDiscountCalculatorResult) SupplementCount() int {
	return result.OffPeakSupplementCount + result.PeakSupplementCount
}

// Calculate calculates the total price
func (calculator NSNoDiscountCalculator) Calculate(records []EnrichedRecord) (result NSNoDiscountCalculatorResult) {
	var (
		errorRecords []ErrorEnrichedRecord
	)

	for _, record := range records {
		if record.IsNSJourney() {
			journeyPrice, err := calculator.priceFetcher.FetchPrice(record.NSJourney())
			if err != nil {
				errorRecords = append(errorRecords, ErrorEnrichedRecord{
					Record: record,
					Error:  errors.Wrapf(err, "cannot fetch price for record"),
				})
			} else {
				result.addJourneyPrice(journeyPrice)
			}
		} else if record.IsSupplement() {
			isOffPeak := calculator.offPeakService.IsOffPeak(record.StartTime.ToTime())
			if isOffPeak {
				result.incrementOffPeakSupplement()
			} else {
				result.incrementPeakSupplement()
			}
		}
	}

	return result
}
