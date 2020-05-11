package main

import (
	"github.com/pkg/errors"
)

// NSPriceFetcherService gets the price for an NS journey
type NSPriceFetcherService struct {
	apiClient        *NSAPIClient
	pricesRepository NSPricesRepository
	cache            LFUCache
	errorHandler     ErrorHandler
}

// NewNSPriceFetcher creates a new instance of the NSPriceFetcherService
func NewNSPriceFetcher(
	apiClient *NSAPIClient,
	pricesRepository NSPricesRepository,
	errorHandler ErrorHandler,
	cache LFUCache,
) NSPriceFetcherService {
	return NSPriceFetcherService{
		apiClient,
		pricesRepository,
		cache,
		errorHandler,
	}
}

// FetchPrice returns the NSJourneyPrice for an NSJourney
func (priceFetcher *NSPriceFetcherService) FetchPrice(nsJourney NSJourney) (price NSJourneyPrice, err error) {
	// Fetch price in Cache
	val, err := priceFetcher.cache.Get(nsJourney.NSPriceHash())
	if err == nil {
		return val.(NSJourneyPrice), err
	}

	// Fetch Price in DB
	price, err = priceFetcher.pricesRepository.GetByHash(nsJourney.NSPriceHash())
	if err == nil {
		// price is not in cache so store in cache
		err = priceFetcher.cache.Set(nsJourney.NSPriceHash(), price)
		return price, err
	}

	// handle error gracefully since we still have the API as a backup
	if err != ErrNotFound {
		priceFetcher.errorHandler.HandleSoftError(errors.Wrap(err, "could not fetch prices by hash value"))
	}

	// Fetch Price using the API
	journeyPrice, err := priceFetcher.apiClient.FetchJourneyPrice(nsJourney)
	if err != nil {
		return price, errors.Wrap(err, "cannot fetch price using API")
	}

	// Store the newly fetched price
	err = priceFetcher.pricesRepository.Store(journeyPrice)

	// No need to cause a panic
	if err != nil {
		priceFetcher.errorHandler.HandleSoftError(errors.Wrap(err, "cannot store price in mongodb"))
		return price, nil
	}

	return price, err
}
