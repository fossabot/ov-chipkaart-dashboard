package main

import (
	"bytes"
	"context"
	"github.com/AchoArnold/homework/services/json"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"github.com/pkg/errors"
)

const endpointAuthentication = "https://login.ov-chipkaart.nl/oauth2/token"
const endpointAuthorisation = "https://api2.ov-chipkaart.nl/femobilegateway/v1/api/authorize"
const endpointTransactions = "https://api2.ov-chipkaart.nl/femobilegateway/v1/transactions"

const contentTypeJson = "application/json"

type AuthenticationTokenResponse struct {
	IDToken          string `json:"id_token"`
	ErrorDescription string `json:"error_description"`
	Error            string `json:"error"`
}

type AuthorisationTokenResponse struct {
	ResponseCode int  `json:"c"`
	Value string      `json:"o"`
	E interface{} `json:"e"`
}


type ApiService struct {
	clientId string
	clientSecret string
}

func New(clientId, clientSecret string) ApiService {
	return ApiService{
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (service ApiService) getAuthorisationToken(ctx context.Context, authenticationTokenResponse *AuthenticationTokenResponse) (authorisationToken *AuthorisationTokenResponse, err error)  {
	payload := map[string]string{
		"authenticationToken": authenticationTokenResponse.IDToken,
	}

	request, err := createPostRequest(ctx, endpointAuthorisation, payload)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create authorisation request")
	}

	response, err := doHTTPRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "cannot perform authorisation request")
	}


	err = json.JsonDecode(&authorisationToken, response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot decode authorisation token response")
	}

	if authorisationToken.ResponseCode != 200 {
		return nil, errors.Errorf("Response Code: %d, Error: %s", authorisationToken.ResponseCode, authorisationToken.Value)
	}

	return authorisationToken, nil
}

func (service ApiService) getAuthenticationToken(ctx context.Context, username, password string) (authenticationToken *AuthenticationTokenResponse, err error) {
	payload := map[string]string{
		"username": username,
		"password": password,
		"client_id": service.clientId,
		"client_secret": service.clientSecret,
		"grant_type": "password",
		"scope" : "openid",
	}

	request, err := createPostRequest(ctx, endpointAuthentication, payload)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create authentication request")
	}

	response, err := doHTTPRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "cannot perform authentication request")
	}


	err = json.JsonDecode(&authenticationToken, response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot decode authentication token response")
	}

	if authenticationToken.Error != "" {
		return nil, errors.Wrap(errors.New(authenticationToken.Error), authenticationToken.ErrorDescription)
	}

	return authenticationToken, nil
}

func doHTTPRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}

	apiResponse, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot execute %s request for %s: ", request.Method, request.URL.String())
	}

	return apiResponse, nil
}

func createPostRequest(ctx context.Context, url string, request interface{}) (*http.Request, error) {
	requestBytes, err := json.JsonEncode(request)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot serialize request to json %#+v", request)
	}

	apiRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, errors.Wrap(err, "cannot create request for URL: "+url)
	}

	apiRequest.Header.Set("Accept", contentTypeJson)
	apiRequest.Header.Set("Content-Type", contentTypeJson)

	return apiRequest, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
