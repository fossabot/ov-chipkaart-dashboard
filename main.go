package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
)

const localeEnglish = "en-EN"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = sentry.Init(sentry.ClientOptions{Dsn: os.Getenv("SENTRY_DSN")})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	sentry.CaptureMessage("It works!")

	dbURI := "mongodb+srv://" + os.Getenv("MONGODB_USERNAME") + ":" + os.Getenv("MONGODB_PASSWORD") + "@cluster0-iq6sq.mongodb.net/test?retryWrites=true&w=majority"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(errors.Wrapf(err, "cannot connect to mongoDB"))
	}

	rawRecordsRepository := NewRawRecordsRepository(client.Database(os.Getenv("MONGODB_DB_NAME")))

	config := TransactionFetcherAPIServiceConfig{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Locale:       localeEnglish,
		Client:       &http.Client{},
	}
	apiService := NewAPIService(config)

	transactionConfig := TransactionFetchOptions{
		Username:   os.Getenv("OV_CHIPKAAT_USERNAME"),
		Password:   os.Getenv("OV_CHIPKAAT_PASSWORD"),
		CardNumber: os.Getenv("OV_CHIPKAAT_CARD_NUMBER"),
		StartDate:  time.Unix(0, 0),
		EndDate:    time.Now(),
	}

	log.Println("Fetching Transactions")
	records, err := apiService.FetchTransactions(transactionConfig)
	if err != nil {
		log.Panicf(errors.Wrapf(err, "%+v", err).Error())
	}

	log.Println("Inserting into database")
	err = rawRecordsRepository.Store(*records, NewTransactionID())
	if err != nil {
		log.Panicf(errors.Wrapf(err, "%+v", err).Error())
	}

}
