package main

import (
	"net/http"
	"time"

	"github.com/AchoArnold/homework/services/json"
	"github.com/pkg/errors"
)

const endpointPricesAPI = "https://gateway.apiportal.ns.nl/public-prijsinformatie/prices"

//NSAPIService is a data structure for the price fetcher
type NSAPIService struct {
	httpClient  HTTPClient
	priceAPIKey string
}

// NewAPIPriceFetcherService creates a new price fetcher service
func NewAPIPriceFetcherService(client HTTPClient, priceAPIKey string) *NSAPIService {
	return &NSAPIService{
		httpClient:  client,
		priceAPIKey: priceAPIKey,
	}
}

// FetchOptions are options for fetching the price of a journey
type FetchOptions struct {
	toStation   string
	fromStation string
	date        time.Time
}

func (options FetchOptions) toMap() map[string]string {
	return map[string]string{
		"date":        options.date.Format(dateFormat),
		"toStation":   options.toStation,
		"fromStation": options.fromStation,
	}
}

type priceAPIResponse struct {
	PriceOptions []struct {
		Type           string `json:"type"`
		TariefEenheden int    `json:"tariefEenheden"`
		Prices         []struct {
			ClassType    string `json:"classType"`
			DiscountType string `json:"discountType"`
			ProductType  string `json:"productType"`
			Price        int    `json:"price"`
			Supplements  struct {
				Kaart *int `json:"kaart,omitempty"`
			} `json:"supplements"`
		} `json:"prices,omitempty"`
		Transporter string `json:"transporter,omitempty"`
		From        string `json:"from,omitempty"`
		To          string `json:"to,omitempty"`
		TotalPrices []struct {
			ClassType    string `json:"classType"`
			DiscountType string `json:"discountType"`
			ProductType  string `json:"productType"`
			Price        int    `json:"price"`
			Supplements  struct {
				Kaart *int `json:"kaart,omitempty"`
			} `json:"supplements"`
		} `json:"totalPrices,omitempty"`
		Trajecten []routePrices `json:"trajecten,omitempty"`
	} `json:"priceOptions"`
	FieldErrors *FieldErrors `json:"fieldErrors,omitempty"`
}

type routePrices struct {
	Transporter string `json:"transporter"`
	From        string `json:"from"`
	To          string `json:"to"`
	Prices      []struct {
		ClassType    string `json:"classType"`
		DiscountType string `json:"discountType"`
		ProductType  string `json:"productType"`
		Price        int    `json:"price"`
		Supplements  struct {
			Kaart *int `json:"kaart,omitempty"`
		} `json:"supplements"`
	} `json:"prices"`
}

func (response priceAPIResponse) routePrice() routePrices {
	for _, priceOptions := range response.PriceOptions {
		if len(priceOptions.Trajecten) == 1 {
			return priceOptions.Trajecten[0]
		}
	}
	return routePrices{}
}

func (response priceAPIResponse) getPriceForProductClass(product, class string) int {
	for _, price := range response.routePrice().Prices {
		if price.ProductType == product && price.ClassType == class && price.DiscountType == "NONE" {
			return price.Price
		}
	}
	return 0
}

// NSJourneyPrice is the price structure gotten from the API
type NSJourneyPrice struct {
	Input                         FetchOptions
	FromStationShortName          string
	toStationShortName            string
	firstClassSingleFarePrice     int
	secondClassSingleFarePrice    int
	firstClassRouteBusinessPrice  int
	secondClassRouteBusinessPrice int
	firstClassRoutePrice          int
	secondClassRoutePrice         int
}

// FieldErrors are errors from the API
type FieldErrors struct {
	FieldErrors []struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	} `json:"fieldErrors"`
}

// Error returns the validation error
func (fieldErrors FieldErrors) Error() string {
	return ErrorValidation
}

// Constant-like variables
var (
	ErrorValidation = errors.New("field_error").Error()
)

// FetchJourneyPrice fetches the price for a particular journey
func (service NSAPIService) FetchJourneyPrice(options FetchOptions) (price NSJourneyPrice, err error) {
	apiRequest, err := service.createGetRequest(endpointPricesAPI, service.priceAPIKey, options.toMap())
	if err != nil {
		return price, errors.Wrap(err, "cannot create get request")
	}

	response, err := service.doHTTPRequest(apiRequest)
	if err != nil {
		return price, errors.Wrap(err, "cannot do http request")
	}

	if response.StatusCode != responseCodeOk && response.StatusCode != 400 {
		return price, errors.Wrapf(errors.New("invalid response code"), "%d", response.StatusCode)
	}

	var priceAPIResponse priceAPIResponse
	err = json.JsonDecode(&priceAPIResponse, response.Body)
	if err != nil {
		return price, errors.Wrapf(err, "cannot decode response into price response: payload = %+#v", options.toMap())
	}

	if priceAPIResponse.FieldErrors != nil {
		return price, *priceAPIResponse.FieldErrors
	}

	var (
		classFirst  = "FIRST"
		classSecond = "SECOND"
	)

	var (
		productSingleFare        = "SINGLE_FARE"
		productRouteFree         = "TRAJECTVRIJ_MAAND"
		productRouteFreeBusiness = "TRAJECTVRIJ_NSBUSINESSKAART"
	)
	price = NSJourneyPrice{
		Input:                         FetchOptions{},
		FromStationShortName:          priceAPIResponse.routePrice().From,
		toStationShortName:            priceAPIResponse.routePrice().To,
		firstClassSingleFarePrice:     priceAPIResponse.getPriceForProductClass(productSingleFare, classFirst),
		secondClassSingleFarePrice:    priceAPIResponse.getPriceForProductClass(productSingleFare, classSecond),
		firstClassRouteBusinessPrice:  priceAPIResponse.getPriceForProductClass(productRouteFreeBusiness, classFirst),
		secondClassRouteBusinessPrice: priceAPIResponse.getPriceForProductClass(productRouteFreeBusiness, classSecond),
		firstClassRoutePrice:          priceAPIResponse.getPriceForProductClass(productRouteFree, classFirst),
		secondClassRoutePrice:         priceAPIResponse.getPriceForProductClass(productRouteFreeBusiness, classSecond),
	}

	return price, err
}

func (service NSAPIService) doHTTPRequest(request *http.Request) (*http.Response, error) {
	apiResponse, err := service.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot execute %s request for %s: ", request.Method, request.URL.String())
	}

	return apiResponse, nil
}

func (service NSAPIService) createGetRequest(endpoint string, apiKey string, payload map[string]string) (*http.Request, error) {
	apiRequest, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create request for URL: "+endpoint)
	}

	query := apiRequest.URL.Query()
	for key, value := range payload {
		query.Add(key, value)
	}

	apiRequest.URL.RawQuery = query.Encode()
	apiRequest.Header.Set("Accept", contentTypeJSON)
	apiRequest.Header.Set("Ocp-Apim-Subscription-Key", apiKey)

	return apiRequest, nil
}
