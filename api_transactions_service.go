package main

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/AchoArnold/homework/services/json"
	"github.com/pkg/errors"
	"go.uber.org/ratelimit"
)

const endpointAuthentication = "https://login.ov-chipkaart.nl/oauth2/token"
const endpointAuthorisation = "https://api2.ov-chipkaart.nl/femobilegateway/v1/api/authorize"
const endpointTransactions = "https://api2.ov-chipkaart.nl/femobilegateway/v1/transactions"

const contentTypeJSON = "application/json"
const contentTypeFormURLEncoded = "application/x-www-form-urlencoded"

const responseCodeOk = 200

const transactionRequestsPerSecond = 5

type authenticationTokenResponse struct {
	IDToken          string `json:"id_token"`
	ErrorDescription string `json:"error_description"`
	Error            string `json:"error"`
}

type authorisationTokenResponse struct {
	ResponseCode int         `json:"c"`
	Value        string      `json:"o"`
	Error        interface{} `json:"e"`
}

type transactionsResponse struct {
	ResponseCode int `json:"c"`
	Response     struct {
		TotalSize              int      `json:"totalSize"`
		NextOffset             int      `json:"nextOffset"`
		PreviousOffset         int      `json:"previousOffset"`
		Records                []Record `json:"records"`
		TransactionsRestricted bool     `json:"transactionsRestricted"`
		NextRequestContext     struct {
			StartDate string `json:"startDate"`
			EndDate   string `json:"endDate"`
			Offset    int    `json:"offset"`
		} `json:"nextRequestContext"`
	} `json:"o"`
	Error interface{} `json:"e"`
}

// Record represents a transaction record
type Record struct {
	CheckInInfo            string   `json:"checkInInfo"`
	CheckInText            string   `json:"checkInText"`
	Fare                   *float64 `json:"fare"`
	FareCalculation        string   `json:"fareCalculation"`
	FareText               string   `json:"fareText"`
	ModalType              string   `json:"modalType"`
	ProductInfo            string   `json:"productInfo"`
	ProductText            string   `json:"productText"`
	Pto                    string   `json:"pto"`
	TransactionDateTime    int64    `json:"transactionDateTime"`
	TransactionInfo        string   `json:"transactionInfo"`
	TransactionName        string   `json:"transactionName"`
	EPurseMut              *float64 `json:"ePurseMut"`
	EPurseMutInfo          string   `json:"ePurseMutInfo"`
	TransactionExplanation string   `json:"transactionExplanation"`
	TransactionPriority    string   `json:"transactionPriority"`
}

type transactionsPayload struct {
	AuthorisationToken string `json:"authorizationToken"`
	MediumID           string `json:"mediumId"`
	Locale             string `json:"locale"`
	Offset             string `json:"offset"`
	StartDate          string `json:"startDate"`
	EndDate            string `json:"endDate"`
}

// TransactionFetcherAPIService is responsible for the fetching transactions using the ov-chipkaart API
type TransactionFetcherAPIService struct {
	clientID     string
	clientSecret string
	httpClient   HTTPClient
	locale       string
}

// TransactionFetcherAPIServiceConfig is the configuration for this service
type TransactionFetcherAPIServiceConfig struct {
	ClientID     string
	ClientSecret string
	Locale       string
	Client       HTTPClient
}

// NewAPIService Initializes the API service.
func NewAPIService(config TransactionFetcherAPIServiceConfig) *TransactionFetcherAPIService {
	return &TransactionFetcherAPIService{
		clientID:     config.ClientID,
		clientSecret: config.ClientSecret,
		httpClient:   config.Client,
		locale:       config.Locale,
	}
}

// FetchTransactions returns the transaction records based on the parameter provided.
func (service TransactionFetcherAPIService) FetchTransactions(options TransactionFetchOptions) (records *[]Record, err error) {
	authenticationToken, err := service.getAuthenticationToken(options.Username, options.Password)
	if err != nil {
		return records, errors.Wrap(err, "could not fetch authentication token")
	}

	authorisationToken, err := service.getAuthorisationToken(authenticationToken)
	if err != nil {
		return records, errors.Wrap(err, "could not fetch authorisation token")
	}

	records, err = service.getTransactions(authorisationToken, options)
	if err != nil {
		return records, errors.Wrap(err, "could not fetch transactions")
	}

	return records, nil
}

func (service TransactionFetcherAPIService) getTransactions(authorisationToken authorisationTokenResponse, options TransactionFetchOptions) (*[]Record, error) {
	payload := transactionsPayload{
		AuthorisationToken: authorisationToken.Value,
		MediumID:           options.CardNumber,
		Locale:             service.locale,
		StartDate:          options.StartDate.Format(dateFormat),
		EndDate:            options.EndDate.Format(dateFormat),
	}

	transactionsResponse, err := service.getTransaction(payload)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot perform transactions request: payload = %+v", payload)
	}

	records := transactionsResponse.Response.Records

	payload.StartDate = transactionsResponse.Response.NextRequestContext.StartDate
	payload.EndDate = transactionsResponse.Response.NextRequestContext.EndDate

	requestLimit := len(records)
	numberOfRequests := int(math.Ceil(float64(transactionsResponse.Response.TotalSize) / float64(requestLimit)))

	rateLimiter := ratelimit.New(transactionRequestsPerSecond)
	for i := 1; i < numberOfRequests; i++ {
		payload.Offset = strconv.Itoa(transactionsResponse.Response.NextRequestContext.Offset)

		rateLimiter.Take()

		transactions, err := service.getTransaction(payload)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot perform transactions request: payload = %+v", payload)
		}

		records = append(records, transactions.Response.Records...)
	}

	return &records, nil
}

func (service TransactionFetcherAPIService) getTransaction(payload transactionsPayload) (transactionsResponse *transactionsResponse, err error) {
	payloadAsMap, err := json.JsonToStringMap(payload)
	if err != nil {
		return transactionsResponse, errors.Wrapf(err, "cannot serialize request to map %#+v", payload)
	}

	request, err := service.createPostRequest(endpointTransactions, payloadAsMap)
	if err != nil {
		return transactionsResponse, errors.Wrapf(err, "cannot create transaction request: payload = %+#v", payloadAsMap)
	}

	response, err := service.doHTTPRequest(request)
	if err != nil {
		return transactionsResponse, errors.Wrapf(err, "cannot perform transaction request: payload = %+#v", request)
	}

	err = json.JsonDecode(&transactionsResponse, response.Body)
	if err != nil {
		return transactionsResponse, errors.Wrapf(err, "cannot decode response into transactions response: payload = %+#v", response)
	}

	if transactionsResponse != nil && transactionsResponse.ResponseCode != responseCodeOk {
		return transactionsResponse, errors.Errorf("Invalid response code %d: payload = %+#v", transactionsResponse.ResponseCode, payload)
	}

	return transactionsResponse, nil
}

func (service TransactionFetcherAPIService) getAuthorisationToken(authenticationTokenResponse authenticationTokenResponse) (authorisationToken authorisationTokenResponse, err error) {
	payload := map[string]string{
		"authenticationToken": authenticationTokenResponse.IDToken,
	}

	request, err := service.createPostRequest(endpointAuthorisation, payload)
	if err != nil {
		return authorisationToken, errors.Wrap(err, "cannot create authorisation request")
	}

	response, err := service.doHTTPRequest(request)
	if err != nil {
		return authorisationToken, errors.Wrap(err, "cannot perform authorisation request")
	}

	err = json.JsonDecode(&authorisationToken, response.Body)
	if err != nil {
		return authorisationToken, errors.Wrap(err, "cannot decode authorisation token response")
	}

	if authorisationToken.ResponseCode != responseCodeOk {
		return authorisationToken, errors.Errorf("Response Code: %d, Error: %s", authorisationToken.ResponseCode, authorisationToken.Value)
	}

	return authorisationToken, nil
}

func (service TransactionFetcherAPIService) getAuthenticationToken(username, password string) (authenticationToken authenticationTokenResponse, err error) {
	payload := map[string]string{
		"username":      username,
		"password":      password,
		"client_id":     service.clientID,
		"client_secret": service.clientSecret,
		"grant_type":    "password",
		"scope":         "openid",
	}

	request, err := service.createPostRequest(endpointAuthentication, payload)
	if err != nil {
		return authenticationToken, errors.Wrap(err, "cannot create authentication request")
	}

	response, err := service.doHTTPRequest(request)
	if err != nil {
		return authenticationToken, errors.Wrap(err, "cannot perform authentication request")
	}

	err = json.JsonDecode(&authenticationToken, response.Body)
	if err != nil {
		return authenticationToken, errors.Wrap(err, "cannot decode authentication token response")
	}

	if authenticationToken.Error != "" {
		return authenticationToken, errors.Wrap(errors.New(authenticationToken.Error), authenticationToken.ErrorDescription)
	}

	return authenticationToken, nil
}

func (service TransactionFetcherAPIService) doHTTPRequest(request *http.Request) (*http.Response, error) {
	apiResponse, err := service.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot execute %s request for %s: ", request.Method, request.URL.String())
	}

	return apiResponse, nil
}

func (service TransactionFetcherAPIService) createPostRequest(endpoint string, payload map[string]string) (*http.Request, error) {
	data := url.Values{}
	for key, val := range payload {
		data.Set(key, val)
	}

	apiRequest, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "cannot create request for URL: "+endpoint)
	}

	apiRequest.Header.Set("Accept", contentTypeJSON)
	apiRequest.Header.Set("Content-Type", contentTypeFormURLEncoded)

	return apiRequest, nil
}
