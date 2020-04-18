package main

import (
	"context"

	"log"
	"net/http"
	"os"
	"time"

	lfucache "github.com/NdoleStudio/lfu-cache"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(errors.Wrapf(err, "cannot connect to mongoDB"))
	}

	mongodb := client.Database(os.Getenv("MONGODB_DB_NAME"))
	//loadNsStations(mongodb)
	//storeNSTransactions(mongodb)

	bsonService := NewBsonService()
	errorHandler := NewSentryErrorHandler()
	rawRecordsRepository := NewRawRecordsRepository(mongodb, collectionRawRecords, bsonService)
	enrichedRecordsRepository := NewMongoNSEnrichedRecordsRepository(mongodb, collectionNSEnrichedRecords, bsonService)
	cache, err := lfucache.New(100)
	nsClient := NewNSAPIClient(&http.Client{}, os.Getenv("NS_API_KEY_PUBLIC_TRAVEL_INFORMATION"))
	pricesRepository := NewMongoNSPricesRepository(mongodb, collectionNSPrices, bsonService)
	stationsRepository := NewMongoNSStationsRepository(mongodb, collectionNSStations, bsonService)
	priceFetcher := NewNSPriceFetcher(nsClient, pricesRepository, errorHandler, cache)
	stationCodeService := NewNSStationsCodeService(stationsRepository, errorHandler, cache)
	enrichmentService := NewNSRawRecordsEnrichmentService(stationCodeService, priceFetcher)

	log.Println("Fetchigng first transaction")
	id, err := rawRecordsRepository.First()
	if err != nil {
		errorHandler.HandleHardError(err)
	}
	log.Println("Finished fetching first transaction")

	globalTransactionID := *id.TransactionID
	getOptions := GetRawRecordsOptions{
		TransactionID: globalTransactionID,
		SortBy:        "transaction_timestamp",
		SortDirection: "ASC",
	}

	log.Println("fetching raw records from DB")
	rawRecords, err := rawRecordsRepository.GetByTransactionID(getOptions)
	if err != nil {
		errorHandler.HandleHardError(err)
	}

	log.Printf("%d raw records fetched\n", len(rawRecords))

	log.Println("Fetching enriched records")
	enrichmentResult := enrichmentService.Enrich(rawRecords)
	log.Println("Finished enriching records")

	log.Printf("%d enriched records and %d failed records\n", len(enrichmentResult.ValidRecords), len(enrichmentResult.Error.ErrorRecords))

	log.Println("Starting storing of enriched records")
	err = enrichedRecordsRepository.Store(enrichmentResult.ValidRecords)
	if err != nil {
		errorHandler.HandleHardError(err)
	}
	log.Println("Finished storing of enriched records")
}

func loadNsStations(mongodb *mongo.Database) {
	bsonService := BsonService{}
	nsStationsRepository := NewMongoNSStationsRepository(mongodb, collectionNSStations, bsonService)
	//
	log.Printf("Fetching Stations")
	nsClient := NewNSAPIClient(&http.Client{}, os.Getenv("NS_API_KEY_PUBLIC_TRAVEL_INFORMATION"))
	stations, err := nsClient.GetAllStations()
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Stations fetch finished")

	log.Printf("Storing stations in the database")
	err = nsStationsRepository.Store(stations)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Finished storing stations")
}

func storeNSTransactions(mongodb *mongo.Database) {
	bsonService := NewBsonService()
	rawRecordsRepository := NewRawRecordsRepository(mongodb, collectionRawRecords, bsonService)

	//
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

	globalTransactionID := NewTransactionID()

	//
	log.Println("Fetching Transactions")
	records, err := apiService.FetchTransactions(transactionConfig)
	if err != nil {
		log.Panicf(errors.Wrapf(err, "%+v", err).Error())
	}
	log.Println("Finished Fetching transactions")

	// Enriching records
	source := rawRecordSourceAPI
	transactionID := globalTransactionID
	for index := range records {
		recordID := NewTransactionID()
		records[index].TransactionID = &transactionID
		records[index].Source = &source
		records[index].ID = &recordID
	}

	log.Println(len(records))
	log.Println("Inserting into database")
	err = rawRecordsRepository.Store(records)
	if err != nil {
		log.Panicf(errors.Wrapf(err, "%+v", err).Error())
	}

}
