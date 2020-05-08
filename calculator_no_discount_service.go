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
	OffPeakFirstClassPrice  Money
	OffPeakSecondClassPrice Money
	OffPeakJourneyCount     int
	PeakFirstClassPrice     Money
	PeakSecondClassPrice    Money
	PeakJourneyCount        int
	PeakSupplementPrice     Money
	PeakSupplementCount     int
	OffPeakSupplementPrice  Money
	OffPeakSupplementCount  int
	Error                   EnrichedRecordsError
}

func (result *NSNoDiscountCalculatorResult) init() {
	result.OffPeakSecondClassPrice = NewEUR(0)
	result.OffPeakFirstClassPrice = NewEUR(0)
	result.PeakFirstClassPrice = NewEUR(0)
	result.PeakSecondClassPrice = NewEUR(0)
	result.PeakSupplementPrice = NewEUR(0)
	result.OffPeakSupplementPrice = NewEUR(0)
}

// addOffPeakJourneyPrice adds the price of an NSJourney when not in peak period
func (result *NSNoDiscountCalculatorResult) addOffPeakJourneyPrice(journey NSJourneyPrice) {
	result.OffPeakJourneyCount++
}

// addPeakJourneyPrice adds the price of an NS Journey during the peak period
func (result *NSNoDiscountCalculatorResult) addPeakJourneyPrice(journey NSJourneyPrice) {
	result.PeakFirstClassPrice = result.PeakFirstClassPrice.AddAmount(NewEUR(journey.FirstClassSingleFarePrice).Value())
	result.PeakSecondClassPrice = result.PeakSecondClassPrice.AddAmount(NewEUR(journey.SecondClassSingleFarePrice).Value())
	result.PeakJourneyCount++
}

// incrementPeakSupplement adds the peak supplement price
func (result *NSNoDiscountCalculatorResult) incrementPeakSupplement() {
	result.PeakSupplementCount++
	result.PeakSupplementPrice = result.PeakSupplementPrice.AddAmount(supplementPricePeak)
}

// incrementOffPeakSupplement adds the off peak supplement price
func (result *NSNoDiscountCalculatorResult) incrementOffPeakSupplement() {
	result.OffPeakSupplementCount++
	result.OffPeakSupplementPrice = result.OffPeakSupplementPrice.AddAmount(supplementPriceOffPeak)
}

// SupplementPrice returns the price of both off peak and peak supplement
func (result NSNoDiscountCalculatorResult) SupplementPrice() Money {
	return result.OffPeakSupplementPrice.AddAmount(result.PeakSupplementPrice.Value())
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
		isOffPeak := calculator.offPeakService.IsOffPeak(record.StartTime.ToTime())
		if record.IsNSJourney() {
			journeyPrice, err := calculator.priceFetcher.FetchPrice(record.NSJourney())
			if err != nil {
				errorRecords = append(errorRecords, ErrorEnrichedRecord{
					Record: record,
					Error:  errors.Wrapf(err, "cannot fetch price for record"),
				})
			} else {
				if isOffPeak {
					result.addOffPeakJourneyPrice(journeyPrice)
				} else {
					result.addPeakJourneyPrice(journeyPrice)
				}
			}
		} else if record.IsSupplement() {
			if isOffPeak {
				result.incrementOffPeakSupplement()
			} else {
				result.incrementPeakSupplement()
			}
		}
	}

	return result
}
