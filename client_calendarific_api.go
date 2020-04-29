package main

import (
	"log"
	"net/http"
	"time"

	"github.com/AchoArnold/homework/services/json"
	"github.com/pkg/errors"
)

const (
	apiEndpointHolidays = "https://calendarific.com/api/v2/holidays"
	countryNL           = "NL"
	holidayTypeNational = "national"
)

// CalendarificAPIClient is the data structure for the api client
type CalendarificAPIClient struct {
	apiKey     string
	httpClient HTTPClient
}

// NewCalendarificAPIClient returns a new CalendarificAPIClient
func NewCalendarificAPIClient(apiKey string, httpClient HTTPClient) CalendarificAPIClient {
	return CalendarificAPIClient{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

type holidayAPIResponse struct {
	Meta struct {
		Code        int     `json:"code"`
		ErrorType   *string `json:"error_type"`
		ErrorDetail *string `json:"error_detail"`
	} `json:"meta"`
	Response struct {
		Holidays []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Country     struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"country"`
			Date struct {
				Iso      string `json:"iso"`
				Datetime struct {
					Year  int `json:"year"`
					Month int `json:"month"`
					Day   int `json:"day"`
				} `json:"datetime"`
			} `json:"date,omitempty"`
		} `json:"holidays"`
	} `json:"response"`
}

// FetchNationalHolidays fetches national holidays for the netherlands
func (apiClient CalendarificAPIClient) FetchNationalHolidays(timestamp time.Time) (holidays []Holiday, err error) {
	payload := map[string]string{
		"api_key": apiClient.apiKey,
		"year":    timestamp.Format(yearFormat),
		"country": countryNL,
		"type":    holidayTypeNational,
	}

	request, err := apiClient.createGetRequest(apiEndpointHolidays, payload)
	if err != nil {
		return holidays, errors.Wrap(err, "could not create request")
	}

	response, err := apiClient.doHTTPRequest(request)
	if err != nil {
		return holidays, errors.Wrap(err, "could not perform api request")
	}

	var apiResponse holidayAPIResponse
	err = json.JsonDecode(&apiResponse, response.Body)
	if err != nil {
		return holidays, errors.Wrap(err, "cannot decode api response into json")
	}

	if apiResponse.Meta.ErrorType != nil {
		return holidays, errors.Wrapf(err, "%s: %s", apiResponse.Meta.ErrorType, apiResponse.Meta.ErrorDetail)
	}

	for _, holiday := range apiResponse.Response.Holidays {
		log.Println(holiday.Date.Iso)
		val, err := time.Parse(dateFormat, holiday.Date.Iso)
		if err != nil {
			return holidays, errors.Wrap(err, "cannot convert time to date")
		}

		holidays = append(holidays, Holiday{
			ID:        NewTransactionID(),
			Timestamp: val,
			Name:      holiday.Name,
			Country:   countryNL,
			Date:      holiday.Date.Iso,
		})
	}

	return holidays, nil
}

func (apiClient CalendarificAPIClient) doHTTPRequest(request *http.Request) (*http.Response, error) {
	apiResponse, err := apiClient.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot execute %s request for %s: ", request.Method, request.URL.String())
	}

	return apiResponse, nil
}

func (apiClient CalendarificAPIClient) createGetRequest(endpoint string, payload map[string]string) (*http.Request, error) {
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

	return apiRequest, nil
}
