package main

import (
	"context"

	"github.com/gofrs/uuid"

	//"github.com/google/uuid"
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
	bsonService := BsonService{}
	nsStationsRepository := NewMongoNSStationsRepository(mongodb, collectionNSStations, bsonService)
	//
	//log.Printf("Fetching Stations")
	nsClient := NewNSAPIClient(&http.Client{}, os.Getenv("NS_API_KEY_PUBLIC_TRAVEL_INFORMATION"))
	//log.Printf("Stations fetch finished")
	//
	//stations, err := nsClient.GetAllStations()
	//if err != nil {
	//	log.Fatalf(err.Error())
	//}
	//
	//log.Printf("Storing stations in the database")
	//err = nsStationsRepository.Store(stations)
	//if err != nil {
	//	log.Fatalf(err.Error())
	//}
	//log.Printf("Finished storing stations")
	//
	rawRecordsRepository := NewRawRecordsRepository(mongodb, collectionRawRecords, bsonService)
	pricesRepository := NewMongoNSPricesRepository(mongodb, collectionNSPrices, bsonService)
	enrichedRecordsRepository := NewMongoNSEnrichedRecordsRepository(mongodb, collectionNSEnrichedRecords, bsonService)

	cache, err := lfucache.New(100)
	errorHandler := NewSentryErrorHandler()
	priceFetcher := NewNSPriceFetcher(nsClient, pricesRepository, errorHandler, cache)
	stationCodeService := NewNSStationsCodeService(nsStationsRepository, errorHandler, cache)
	enrichmentService := NewNSRawRecordsEnrichmentService(stationCodeService, priceFetcher)

	//
	//config := TransactionFetcherAPIServiceConfig{
	//	ClientID:     os.Getenv("CLIENT_ID"),
	//	ClientSecret: os.Getenv("CLIENT_SECRET"),
	//	Locale:       localeEnglish,
	//	Client:       &http.Client{},
	//}
	//apiService := NewAPIService(config)
	//
	//transactionConfig := TransactionFetchOptions{
	//	Username:   os.Getenv("OV_CHIPKAAT_USERNAME"),
	//	Password:   os.Getenv("OV_CHIPKAAT_PASSWORD"),
	//	CardNumber: os.Getenv("OV_CHIPKAAT_CARD_NUMBER"),
	//	StartDate:  time.Unix(0, 0),
	//	EndDate:    time.Now(),
	//}
	//

	globalTransactionID := NewTransactionID()

	//
	//log.Println("Fetching Transactions")
	//records, err := apiService.FetchTransactions(transactionConfig)
	//if err != nil {
	//	log.Panicf(errors.Wrapf(err, "%+v", err).Error())
	//}
	//log.Println("Finished Fetching transactions")
	//
	//// Enriching records
	//source := rawRecordSourceAPI
	//transactionID := globalTransactionID
	//for index := range records {
	//	recordID := NewTransactionID()
	//	records[index].TransactionID = &transactionID
	//	records[index].Source = &source
	//	records[index].ID = &recordID
	//}
	//
	//
	//log.Println(len(records))
	//log.Println("Inserting into database")
	//err = rawRecordsRepository.Store(records)
	//if err != nil {
	//	log.Panicf(errors.Wrapf(err, "%+v", err).Error())
	//}

	id, err := uuid.FromString("e7474312-00b4-4b49-ac1d-429d74111b85")
	if err != nil {
		errorHandler.HandleHardError(err)
	}

	globalTransactionID = TransactionID(id)

	getOptions := GetRawRecordsOptions{
		TransactionID: globalTransactionID,
		SortBy:        "transaction_timestamp",
		SortDirection: "DESC",
	}

	log.Println("fetching raw records from DB")
	records, err := rawRecordsRepository.GetByTransactionID(getOptions)
	if err != nil {
		errorHandler.HandleHardError(err)
	}

	log.Printf("%d raw records fetched\n", len(records))

	log.Println("Fetching enriched records")
	enrichmentResult := enrichmentService.Enrich(records)
	log.Println("Finished enriching records")

	log.Printf("%d enriched records and %d failed records\n", len(enrichmentResult.ValidRecords), len(enrichmentResult.Error.ErrorRecords))

	log.Println("Starting storing of enriched records")
	err = enrichedRecordsRepository.Store(enrichmentResult.ValidRecords)
	if err != nil {
		errorHandler.HandleHardError(err)
	}
	log.Println("Finished storing of enriched records")
}
