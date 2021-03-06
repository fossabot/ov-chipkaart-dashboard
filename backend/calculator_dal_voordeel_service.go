package main

import (
	"github.com/pkg/errors"
)

const (
	offPeakPriceDiscount = float64(0.6)
	peakPriceDiscount    = float64(0.80)
)

// NSDalVoordeelCalculator calculates the price of journeys with the DalVoordeel discount
type NSDalVoordeelCalculator struct {
	priceFetcher   NSPriceFetcherService
	offPeakService NSOffPeakService
}

// NewNSDalVoordeelCalculator creates a new instance of an NSDalVoordeelCalculator
func NewNSDalVoordeelCalculator(priceFetcher NSPriceFetcherService, offPeakService NSOffPeakService) *NSDalVoordeelCalculator {
	return &NSDalVoordeelCalculator{
		priceFetcher:   priceFetcher,
		offPeakService: offPeakService,
	}
}

//NSDalVoordeelCalculatorResult represents the calculation result of an NS Journey
type NSDalVoordeelCalculatorResult struct {
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

func (result *NSDalVoordeelCalculatorResult) init() {
	result.OffPeakSecondClassPrice = NewEUR(0)
	result.OffPeakFirstClassPrice = NewEUR(0)
	result.PeakFirstClassPrice = NewEUR(0)
	result.PeakSecondClassPrice = NewEUR(0)
	result.PeakSupplementPrice = NewEUR(0)
	result.OffPeakSupplementPrice = NewEUR(0)
}

// addOffPeakJourneyPrice adds the price of an NSJourney when not in peak period
func (result *NSDalVoordeelCalculatorResult) addOffPeakJourneyPrice(journey NSJourneyPrice) {
	result.OffPeakFirstClassPrice = result.OffPeakFirstClassPrice.AddAmount(NewEUR(journey.FirstClassSingleFarePrice).Multiply(offPeakPriceDiscount).Value())
	result.OffPeakSecondClassPrice = result.OffPeakSecondClassPrice.AddAmount(NewEUR(journey.SecondClassSingleFarePrice).Multiply(offPeakPriceDiscount).Value())
	result.OffPeakJourneyCount++
}

// addPeakJourneyPrice adds the price of an NS Journey during the peak period
func (result *NSDalVoordeelCalculatorResult) addPeakJourneyPrice(journey NSJourneyPrice) {
	result.PeakFirstClassPrice = result.PeakFirstClassPrice.AddAmount(NewEUR(journey.FirstClassSingleFarePrice).Value())
	result.PeakSecondClassPrice = result.PeakSecondClassPrice.AddAmount(NewEUR(journey.SecondClassSingleFarePrice).Value())
	result.PeakJourneyCount++
}

// incrementPeakSupplement adds the peak supplement price
func (result *NSDalVoordeelCalculatorResult) incrementPeakSupplement() {
	result.PeakSupplementCount++
	result.PeakSupplementPrice = result.PeakSupplementPrice.AddAmount(supplementPricePeak)
}

// incrementOffPeakSupplement adds the off peak supplement price
func (result *NSDalVoordeelCalculatorResult) incrementOffPeakSupplement() {
	result.OffPeakSupplementCount++
	result.OffPeakSupplementPrice = result.OffPeakSupplementPrice.AddAmount(supplementPriceOffPeak)
}

// SupplementPrice returns the price of both off peak and peak supplement
func (result NSDalVoordeelCalculatorResult) SupplementPrice() Money {
	return result.OffPeakSupplementPrice.AddAmount(result.PeakSupplementPrice.Value())
}

// SupplementCount returns the total count of all supplements.
func (result NSDalVoordeelCalculatorResult) SupplementCount() int {
	return result.OffPeakSupplementCount + result.PeakSupplementCount
}

// Calculate calculates the total price
func (calculator NSDalVoordeelCalculator) Calculate(records []EnrichedRecord) (result NSDalVoordeelCalculatorResult) {
	var (
		errorRecords []ErrorEnrichedRecord
	)

	result.init()
	for _, record := range records {
		isOffPeak := calculator.offPeakService.IsOffPeak(record.StartTime.ToTime())
		if record.IsNSJourney() {
			journeyPrice, err := calculator.priceFetcher.FetchPrice(record.NSJourney())
			if err != nil {
				errorRecords = append(errorRecords, ErrorEnrichedRecord{
					Record: record,
					Error:  errors.Wrap(err, "cannot fetch price for record"),
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
