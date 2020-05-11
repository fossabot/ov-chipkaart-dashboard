package main

import (
	"github.com/pkg/errors"
)

// NSDalVrijCalculator calculates the price of journeys with the DalVrij discount
type NSDalVrijCalculator struct {
	priceFetcher   NSPriceFetcherService
	offPeakService NSOffPeakService
}

// NewNSDalVrijCalculator creates a new instance of an NSDalVrijCalculator
func NewNSDalVrijCalculator(priceFetcher NSPriceFetcherService, offPeakService NSOffPeakService) *NSDalVrijCalculator {
	return &NSDalVrijCalculator{
		priceFetcher:   priceFetcher,
		offPeakService: offPeakService,
	}
}

//NSDalVrijCalculatorResult represents the calculation result of an NS Journey
type NSDalVrijCalculatorResult struct {
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

func (result *NSDalVrijCalculatorResult) init() {
	result.OffPeakSecondClassPrice = NewEUR(0)
	result.OffPeakFirstClassPrice = NewEUR(0)
	result.PeakFirstClassPrice = NewEUR(0)
	result.PeakSecondClassPrice = NewEUR(0)
	result.PeakSupplementPrice = NewEUR(0)
	result.OffPeakSupplementPrice = NewEUR(0)
}

// addOffPeakJourneyPrice adds the price of an NSJourney when not in peak period
func (result *NSDalVrijCalculatorResult) addOffPeakJourneyPrice(journey NSJourneyPrice) {
	result.OffPeakJourneyCount++
}

// addPeakJourneyPrice adds the price of an NS Journey during the peak period
func (result *NSDalVrijCalculatorResult) addPeakJourneyPrice(journey NSJourneyPrice) {
	result.PeakFirstClassPrice = result.PeakFirstClassPrice.AddAmount(NewEUR(journey.FirstClassSingleFarePrice).Value())
	result.PeakSecondClassPrice = result.PeakSecondClassPrice.AddAmount(NewEUR(journey.SecondClassSingleFarePrice).Value())
	result.PeakJourneyCount++
}

// incrementPeakSupplement adds the peak supplement price
func (result *NSDalVrijCalculatorResult) incrementPeakSupplement() {
	result.PeakSupplementCount++
	result.PeakSupplementPrice = result.PeakSupplementPrice.AddAmount(supplementPricePeak)
}

// incrementOffPeakSupplement adds the off peak supplement price
func (result *NSDalVrijCalculatorResult) incrementOffPeakSupplement() {
	result.OffPeakSupplementCount++
	result.OffPeakSupplementPrice = result.OffPeakSupplementPrice.AddAmount(supplementPriceOffPeak)
}

// SupplementPrice returns the price of both off peak and peak supplement
func (result NSDalVrijCalculatorResult) SupplementPrice() Money {
	return result.OffPeakSupplementPrice.AddAmount(result.PeakSupplementPrice.Value())
}

// SupplementCount returns the total count of all supplements.
func (result NSDalVrijCalculatorResult) SupplementCount() int {
	return result.OffPeakSupplementCount + result.PeakSupplementCount
}

// Calculate calculates the total price
func (calculator NSDalVrijCalculator) Calculate(records []EnrichedRecord) (result NSDalVrijCalculatorResult) {
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
