package main

import (
	"net/http"

	"github.com/AchoArnold/homework/services/json"
	"github.com/pkg/errors"
)

// Api Endpoints
const (
	apiEndpointPrices      = "https://gateway.apiportal.ns.nl/public-prijsinformatie/prices"
	apiEndpointAllStations = "https://gateway.apiportal.ns.nl/public-reisinformatie/api/v2/stations"
)

//NSAPIClient is a data structure for the price fetcher
type NSAPIClient struct {
	httpClient             HTTPClient
	publicTravelInfoAPIKey string
}

// NewNSAPIClient creates a new price fetcher service
func NewNSAPIClient(client HTTPClient, priceAPIKey string) *NSAPIClient {
	return &NSAPIClient{
		httpClient:             client,
		publicTravelInfoAPIKey: priceAPIKey,
	}
}

type allStationsAPIResponse struct {
	Links   struct{} `json:"links"`
	Payload []struct {
		Sporen               []interface{} `json:"sporen"`
		Synoniemen           []string      `json:"synoniemen"`
		HeeftFaciliteiten    bool          `json:"heeftFaciliteiten"`
		HeeftVertrektijden   bool          `json:"heeftVertrektijden"`
		HeeftReisassistentie bool          `json:"heeftReisassistentie"`
		Code                 string        `json:"code"`
		Namen                struct {
			Lang   string `json:"lang"`
			Kort   string `json:"kort"`
			Middel string `json:"middel"`
		} `json:"namen"`
		StationType   string  `json:"stationType"`
		Land          string  `json:"land"`
		UICCode       string  `json:"UICCode"`
		Lat           float64 `json:"lat"`
		Lng           float64 `json:"lng"`
		Radius        int     `json:"radius"`
		NaderenRadius int     `json:"naderenRadius"`
		EVACode       string  `json:"EVACode"`
		IngangsDatum  string  `json:"ingangsDatum"`
	} `json:"payload"`
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
func (service NSAPIClient) FetchJourneyPrice(nsJourney NSJourney) (price NSJourneyPrice, err error) {
	apiRequest, err := service.createGetRequest(apiEndpointPrices, service.publicTravelInfoAPIKey, nsJourney.ToMap())
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
		return price, errors.Wrapf(err, "cannot decode response into price response: payload = %+#v", nsJourney.ToMap())
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
		Year:                          nsJourney.Year,
		FromStationCode:               nsJourney.FromStationCode,
		ToStationCode:                 nsJourney.ToStationCode,
		FirstClassSingleFarePrice:     priceAPIResponse.getPriceForProductClass(productSingleFare, classFirst),
		SecondClassSingleFarePrice:    priceAPIResponse.getPriceForProductClass(productSingleFare, classSecond),
		FirstClassRouteBusinessPrice:  priceAPIResponse.getPriceForProductClass(productRouteFreeBusiness, classFirst),
		SecondClassRouteBusinessPrice: priceAPIResponse.getPriceForProductClass(productRouteFreeBusiness, classSecond),
		FirstClassRoutePrice:          priceAPIResponse.getPriceForProductClass(productRouteFree, classFirst),
		SecondClassRoutePrice:         priceAPIResponse.getPriceForProductClass(productRouteFreeBusiness, classSecond),
		Hash:                          nsJourney.NSPriceHash(),
	}

	return price, err
}

// GetAllStations returns all NS train stations
func (service NSAPIClient) GetAllStations() (stations []NSStation, error error) {
	apiRequest, err := service.createGetRequest(apiEndpointAllStations, service.publicTravelInfoAPIKey, nil)
	if err != nil {
		return stations, errors.Wrap(err, "cannot create get request for all stations")
	}

	response, err := service.doHTTPRequest(apiRequest)
	if err != nil {
		return stations, errors.Wrap(err, "cannot do http request for all stations")
	}

	if response.StatusCode != responseCodeOk && response.StatusCode != 400 {
		return stations, errors.Wrapf(errors.New("invalid response code for all stations"), "%d", response.StatusCode)
	}

	var allStationsRaw allStationsAPIResponse
	err = json.JsonDecode(&allStationsRaw, response.Body)
	if err != nil {
		return stations, errors.Wrapf(err, "cannot decode response into all stations struct")
	}

	for _, station := range allStationsRaw.Payload {
		stations = append(stations, NSStation{
			Name:          station.Namen.Lang,
			CurrentName:   station.Namen.Lang,
			Code:          station.Code,
			Country:       station.Land,
			EVACode:       station.EVACode,
			Latitude:      station.Lat,
			Longitude:     station.Lng,
			StartIngDate:  station.IngangsDatum,
			UICCode:       station.UICCode,
			IsDepreciated: false,
		})

		for _, name := range station.Synoniemen {
			stations = append(stations, NSStation{
				Name:          name,
				CurrentName:   station.Namen.Lang,
				Code:          station.Code,
				Country:       station.Land,
				EVACode:       station.EVACode,
				Latitude:      station.Lat,
				Longitude:     station.Lng,
				StartIngDate:  station.IngangsDatum,
				UICCode:       station.UICCode,
				IsDepreciated: true,
			})
		}
	}

	return stations, nil
}

func (service NSAPIClient) doHTTPRequest(request *http.Request) (*http.Response, error) {
	apiResponse, err := service.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot execute %s request for %s: ", request.Method, request.URL.String())
	}

	return apiResponse, nil
}

func (service NSAPIClient) createGetRequest(endpoint string, apiKey string, payload map[string]string) (*http.Request, error) {
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
