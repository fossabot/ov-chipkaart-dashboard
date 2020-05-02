package main

import "github.com/pkg/errors"

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
	FirstClassPrice  int
	SecondClassPrice int
	SupplementPrice  int
	SupplementCount  int
	JourneyCount     int
	Error            EnrichedRecordsError
}

// AddJourneyPrice adds the price of an NSJourney to the result
func (result *NSNoDiscountCalculatorResult) AddJourneyPrice(journey NSJourneyPrice) {
	result.FirstClassPrice += journey.FirstClassSingleFarePrice
	result.SecondClassPrice += journey.SecondClassSingleFarePrice
	result.JourneyCount++
}

// IncrementSupplement increments tne result by adding a supplement
func (result *NSNoDiscountCalculatorResult) IncrementSupplement(isOffPeak bool) {
	result.SupplementCount++
	if isOffPeak {
		result.SupplementPrice += supplementPriceOffPeak
	} else {
		result.SupplementPrice += supplementPricePeak
	}
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
			}

			result.AddJourneyPrice(journeyPrice)
		} else if record.IsSupplement() {
			result.IncrementSupplement(calculator.offPeakService.IsOffPeak(record.StartTime.ToTime()))
		}
	}

	return result
}
