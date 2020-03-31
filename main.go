package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

const localeEnglish = "en-EN"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//func NewApiService(clientId, clientSecret, locale string, client *http.Client) *ApiService {
	apiService := NewApiService(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), localeEnglish, &http.Client{Timeout: timeoutApiRequest})

	_, err = apiService.FetchTransactions(os.Getenv("OV_CHIPKAAT_USERNAME"), os.Getenv("OV_CHIPKAAT_PASSWORD"), os.Getenv("OV_CHIPKAAT_CARD_NUMBER"))
	if err != nil {

	}
}
