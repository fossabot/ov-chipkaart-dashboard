package main

import (
	"bytes"
	"github.com/AchoArnold/homework/services/json"
	"github.com/pkg/errors"
	"go.uber.org/ratelimit"
	"math"
	"net/http"
	"time"
)

const endpointAuthentication = "https://login.ov-chipkaart.nl/oauth2/token"
const endpointAuthorisation = "https://api2.ov-chipkaart.nl/femobilegateway/v1/api/authorize"
const endpointTransactions = "https://api2.ov-chipkaart.nl/femobilegateway/v1/transactions"

const contentTypeJson = "application/json"

const responseCodeOk = 200

const transactionRequestsPerSecond = 5

const timeoutApiRequest = 100 * time.Millisecond

type AuthenticationTokenResponse struct {
	IDToken          string `json:"id_token"`
	ErrorDescription string `json:"error_description"`
	Error            string `json:"error"`
}

type AuthorisationTokenResponse struct {
	ResponseCode int         `json:"c"`
	Value        string      `json:"o"`
	Error        interface{} `json:"e"`
}

type TransactionsResponse struct {
	ResponseCode int `json:"c"`
	Response     struct {
		TotalSize              int       `json:"totalSize"`
		NextOffset             int       `json:"nextOffset"`
		PreviousOffset         int       `json:"previousOffset"`
		Records                []Records `json:"records"`
		TransactionsRestricted bool      `json:"transactionsRestricted"`
		NextRequestContext     struct {
			StartDate string `json:"startDate"`
			EndDate   string `json:"endDate"`
			Offset    int    `json:"offset"`
		} `json:"nextRequestContext"`
	} `json:"o"`
	Error interface{} `json:"e"`
}

type Records struct {
	CheckInInfo            string  `json:"checkInInfo"`
	CheckInText            string  `json:"checkInText"`
	Fare                   float64 `json:"fare"`
	FareCalculation        string  `json:"fareCalculation"`
	FareText               string  `json:"fareText"`
	ModalType              string  `json:"modalType"`
	ProductInfo            string  `json:"productInfo"`
	ProductText            string  `json:"productText"`
	Pto                    string  `json:"pto"`
	TransactionDateTime    int64   `json:"transactionDateTime"`
	TransactionInfo        string  `json:"transactionInfo"`
	TransactionName        string  `json:"transactionName"`
	EPurseMut              float64 `json:"ePurseMut"`
	EPurseMutInfo          string  `json:"ePurseMutInfo"`
	TransactionExplanation string  `json:"transactionExplanation"`
	TransactionPriority    string  `json:"transactionPriority"`
}

type TransactionsPayload struct {
	AuthorisationToken string  `json:"authorizationToken"`
	MediumId           string  `json:"mediumId"`
	Offset             *int    `json:"omitempty,offset"`
	Locale             *string `json:"omitempty,locale"`
	StartDate          *string `json:"omitempty,startDate"`
	EndDate            *string `json:"omitempty,endDate"`
}

type ApiService struct {
	clientId     string
	clientSecret string
	httpClient   *http.Client
	locale       *string
}

func NewApiService(clientId, clientSecret, locale string, client *http.Client) *ApiService {
	return &ApiService{
		clientId:     clientId,
		clientSecret: clientSecret,
		httpClient:   client,
		locale:       &locale,
	}
}

func (service ApiService) FetchTransactions(username, password, cardNumber string) (records *[]Records, err error) {
	authenticationToken, err := service.getAuthenticationToken(username, password)
	if err != nil {
		return records, errors.Wrap(err, "could not fetch authentication token")
	}

	authorisationToken, err := service.getAuthorisationToken(authenticationToken)
	if err != nil {
		return records, errors.Wrap(err, "could not fetch authorisation token")
	}

	records, err = service.getTransactions(authorisationToken, cardNumber)
	if err != nil {
		return records, errors.Wrap(err, "could not fetch transactions")
	}

	return records, nil
}

func (service ApiService) getTransactions(authorisationToken AuthorisationTokenResponse, cardNumber string) (*[]Records, error) {
	payload := TransactionsPayload{
		AuthorisationToken: authorisationToken.Value,
		MediumId:           cardNumber,
		Locale:             service.locale,
		StartDate:          nil,
		EndDate:            nil,
	}

	transactionsResponse, err := service.getTransaction(payload)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot perform transactions request: payload = %+#v", payload)
	}

	records := transactionsResponse.Response.Records

	payload.StartDate = &transactionsResponse.Response.NextRequestContext.StartDate
	payload.EndDate = &transactionsResponse.Response.NextRequestContext.EndDate

	requestLimit := len(records)
	numberOfRequests := int(math.Ceil(float64(transactionsResponse.Response.TotalSize) / float64(requestLimit)))

	rateLimiter := ratelimit.New(transactionRequestsPerSecond)
	for i := 1; i < numberOfRequests; i++ {
		payload.Offset = &transactionsResponse.Response.NextRequestContext.Offset

		rateLimiter.Take()

		transactions, err := service.getTransaction(payload)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot perform transactions request: payload = %+#v", payload)
		}

		records = append(records, transactions.Response.Records...)
	}

	return &records, nil
}

func (service ApiService) getTransaction(payload TransactionsPayload) (transactionsResponse *TransactionsResponse, err error) {
	request, err := service.createPostRequest(endpointTransactions, payload)
	if err != nil {
		return transactionsResponse, errors.Wrapf(err, "cannot create transactions request: payload = %+#v", payload)
	}

	response, err := service.doHTTPRequest(request)
	if err != nil {
		return transactionsResponse, errors.Wrapf(err, "cannot perform transactions request: payload = %+#v", payload)
	}

	err = json.JsonDecode(&transactionsResponse, response.Body)
	if err != nil {
		return transactionsResponse, errors.Wrapf(err, "cannot decode response into transactions response: payload = %+#v", payload)
	}

	if transactionsResponse != nil && transactionsResponse.ResponseCode != responseCodeOk {
		return transactionsResponse, errors.Errorf("Invalid response code %d: payload = %+#v", transactionsResponse.ResponseCode, payload)
	}

	return transactionsResponse, nil
}

func (service ApiService) getAuthorisationToken(authenticationTokenResponse AuthenticationTokenResponse) (authorisationToken AuthorisationTokenResponse, err error) {
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

	if authorisationToken.ResponseCode != 200 {
		return authorisationToken, errors.Errorf("Response Code: %d, Error: %s", authorisationToken.ResponseCode, authorisationToken.Value)
	}

	return authorisationToken, nil
}

func (service ApiService) getAuthenticationToken(username, password string) (authenticationToken AuthenticationTokenResponse, err error) {
	payload := map[string]string{
		"username":      username,
		"password":      password,
		"client_id":     service.clientId,
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

func (service ApiService) doHTTPRequest(request *http.Request) (*http.Response, error) {
	apiResponse, err := service.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot execute %s request for %s: ", request.Method, request.URL.String())
	}

	return apiResponse, nil
}

func (service ApiService) createPostRequest(url string, request interface{}) (*http.Request, error) {
	requestBytes, err := json.JsonEncode(request)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot serialize request to json %#+v", request)
	}

	apiRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, errors.Wrap(err, "cannot create request for URL: "+url)
	}

	apiRequest.Header.Set("Accept", contentTypeJson)
	apiRequest.Header.Set("Content-Type", contentTypeJson)

	return apiRequest, nil
}
