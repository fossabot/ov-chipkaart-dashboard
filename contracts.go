package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

// DBTimestamp stores the timestamp when persisting objects in the DB
type DBTimestamp struct {
	CreatedAt *time.Time `bson:"created_at,omitempty"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
}

// TransactionFetchOptions are the options needed when fetching a list of transactions
type TransactionFetchOptions struct {
	Username   string
	Password   string
	CardNumber string
	StartDate  time.Time
	EndDate    time.Time
}

// RawRecord represents a transaction record
type RawRecord struct {
	DBTimestamp
	ID                     *TransactionID     `bson:"id"`
	TransactionID          *TransactionID     `bson:"transaction_id"`
	CheckInInfo            string             `json:"checkInInfo" bson:"check_in_info"`
	CheckInText            string             `json:"checkInText" bson:"check_in_text"`
	Fare                   *float64           `json:"fare" bson:"fare"`
	FareCalculation        string             `json:"fareCalculation" bson:"fare_calculation"`
	FareText               string             `json:"fareText" bson:"fare_text"`
	ModalType              string             `json:"modalType" bson:"modal_type"`
	ProductInfo            string             `json:"productInfo" bson:"product_info"`
	ProductText            string             `json:"productText" bson:"product_text"`
	Pto                    string             `json:"pto" bson:"pto"`
	TransactionDateTime    TimeInMilliSeconds `json:"transactionDateTime" bson:"transaction_timestamp"`
	TransactionInfo        string             `json:"transactionInfo" bson:"transaction_info"`
	TransactionName        TransactionName    `json:"transactionName" bson:"transaction_name"`
	EPurseMut              *float64           `json:"ePurseMut" bson:"e_purse_mut"`
	EPurseMutInfo          string             `json:"ePurseMutInfo" bson:"e_purse_mut_info"`
	TransactionExplanation string             `json:"transactionExplanation" bson:"transaction_explanation"`
	TransactionPriority    string             `json:"transactionPriority" bson:"transaction_priority"`
	Source                 *RawRecordSource   `bson:"source"`
}

// IsCheckIn determines if a record is a check in record
func (record RawRecord) IsCheckIn() bool {
	return record.TransactionName.IsTheSameAs(transactionNameCheckIn)
}

// IsNSSupplement determines if a records is a surcharge
func (record RawRecord) IsNSSupplement() bool {
	return record.TransactionName.IsTheSameAs(transactionNameIntercityDirectSurcharge)
}

// IsCheckOut determines if a record is checkout transaction.
func (record RawRecord) IsCheckOut() bool {
	return record.TransactionName.IsTheSameAs(transactionNameCheckOut)
}

// IsRET is used to determine if a raw record is from the RET company
func (record RawRecord) IsRET() bool {
	return record.Pto == companyNameRET.String()
}

// IsNS is used to determine if a raw record is from the NS company
func (record RawRecord) IsNS() bool {
	return record.Pto == companyNameNS.String() && record.ModalType == "Trein"
}

// HTTPClient is the class used to perform http requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RawRecordsRepository is used to persist raw transaction records
type RawRecordsRepository interface {
	Store(records []RawRecord) (err error)
}

// NSJourneyPrice represents the price for an NS journey
type NSJourneyPrice struct {
	DBTimestamp                   `bson:"db_timestamp,omitempty"`
	Year                          string `bson:"year"`
	FromStationCode               string `bson:"from_station_code"`
	ToStationCode                 string `bson:"to_station_code"`
	FirstClassSingleFarePrice     int    `bson:"first_class_single_fare_price"`
	SecondClassSingleFarePrice    int    `bson:"second_class_single_fare_price"`
	FirstClassRouteBusinessPrice  int    `bson:"first_class_route_business_price"`
	SecondClassRouteBusinessPrice int    `bson:"second_class_route_business_price"`
	FirstClassRoutePrice          int    `bson:"fist_class_route_price"`
	SecondClassRoutePrice         int    `bson:"second_class_route_price"`
	Hash                          string `bson:"hash"`
}

// EstimatedDurationInMilliSeconds gives an estimate of the duration of a journey based on the price
func (price NSJourneyPrice) EstimatedDurationInMilliSeconds() int {
	return (price.SecondClassSingleFarePrice - basicFare) * costMultiplier * 60 * 1000
}

// NSPricesRepository is responsible for saving and loading the NSJourneyPrice for an journey
type NSPricesRepository interface {
	Store(price NSJourneyPrice) (err error)
	GetByHash(hash string) (price NSJourneyPrice, err error)
}

// NSStationsRepository is responsible for saving and loading NSStation struct
type NSStationsRepository interface {
	Store(stations []NSStation) (err error)
	GetByName(name string) (station NSStation, err error)
	GetByCode(code string) (station NSStation, err error)
}

// ErrorHandler is responsible for handling application errors
type ErrorHandler interface {
	HandleSoftError(err error)
	HandleHardError(err error)
}

// LFUCache implements a least frequently used cache
type LFUCache interface {
	Get(key interface{}) (value interface{}, err error)
	Set(key interface{}, value interface{}) (err error)
}

// NSStation contains info for an NSStation
type NSStation struct {
	DBTimestamp
	Name          string  `bson:"name"`
	Code          string  `bson:"code"`
	Country       string  `bson:"country"`
	EVACode       string  `bson:"eva_code"`
	Latitude      float64 `bson:"latitude"`
	Longitude     float64 `bson:"longitude"`
	StartIngDate  string  `bson:"starting_date"`
	UICCode       string  `bson:"UICCode"`
	IsDepreciated bool    `bson:"is_depreciated"`
	CurrentName   string  `bson:"current_name"`
}

// ToLower converts the NSStation struct values to lowercase
func (station NSStation) ToLower() NSStation {
	return NSStation{
		Name:          strings.ToLower(station.Name),
		Code:          strings.ToLower(station.Code),
		Country:       strings.ToLower(station.Country),
		EVACode:       station.EVACode,
		Latitude:      station.Latitude,
		Longitude:     station.Longitude,
		StartIngDate:  station.StartIngDate,
		UICCode:       station.UICCode,
		IsDepreciated: station.IsDepreciated,
		CurrentName:   strings.ToLower(station.CurrentName),
	}
}

// NSJourney are options for fetching the price of a journey
type NSJourney struct {
	Year            string `bson:"year"`
	FromStationCode string `bson:"from_station_code"`
	ToStationCode   string `bson:"to_station_code"`
	date            time.Time
}

// NewNSJourney creates a new NSJourney instance
func NewNSJourney(timestamp time.Time, fromStationCode, toStationCode string) NSJourney {
	return NSJourney{
		Year:            timestamp.Format(yearFormat),
		FromStationCode: fromStationCode,
		ToStationCode:   toStationCode,
		date:            timestamp,
	}
}

// ToMap converts  the JS journey struct to a `map[string]string` map
func (journey NSJourney) ToMap() map[string]string {
	date := journey.date.Format(dateFormat)
	if journey.date.Year() < time.Now().Year() {
		parsed, err := time.Parse(dateFormat, strconv.Itoa(time.Now().Year()-1)+"-12-30")
		if err == nil {
			log.Printf("cannot parse %f as date", strconv.Itoa(time.Now().Year()-1)+"-12-30")
			date = parsed.Format(dateFormat)
		}
	}

	log.Printf("date =  %s\n", date)
	return map[string]string{
		"date":        date,
		"toStation":   journey.FromStationCode,
		"fromStation": journey.ToStationCode,
	}
}

// NSPriceHash gets the hash for an ns journey used to determine the price of the journey
func (journey NSJourney) NSPriceHash() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(journey.FromStationCode+hashSeparator+journey.ToStationCode+hashSeparator+journey.Year)))
}

//////////////////////////
// Enrichment Service
//////////////////////////

// EnrichedRecord represents an enriched record.
type EnrichedRecord struct {
	DBTimestamp      `bson:"db_timestamp,omitempty"`
	RawRecordID      *TransactionID     `bson:"raw_record_id"`
	TransactionID    *TransactionID     `bson:"transaction_id"`
	ID               *TransactionID     `bson:"id"`
	StartTime        TimeInMilliSeconds `bson:"start_time"`
	EndTime          TimeInMilliSeconds `bson:"end_time"`
	StartTimeIsExact bool               `bson:"start_time_is_exact"`
	FromStationCode  string             `bson:"from_station_code"`
	ToStationCode    string             `bson:"to_station_code"`
	CompanyName      CompanyName        `bson:"company_name"`
	TransactionType  TransactionType    `bson:"transaction_type"`
	Duration         int64              `bson:"duration"`
}

// RawRecordsEnrichmentService is the interface for filtering raw records
type RawRecordsEnrichmentService interface {
	Enrich(records []RawRecord) RawRecordsEnrichmentResults
}

// RawRecordsEnrichmentResults is the results of the raw records filter
type RawRecordsEnrichmentResults struct {
	ValidRecords []EnrichedRecord
	Error        RawRecordsEnrichmentError
}

// HasError determines if there are error results in the filter results
func (results RawRecordsEnrichmentResults) HasError() bool {
	return len(results.Error.ErrorRecords) > 0
}

// GetRawRecordsOptions are settings that are passed to the RawRecordsRepository
type GetRawRecordsOptions struct {
	TransactionID TransactionID
	SortBy        string
	SortDirection string
}

// EnrichedRecordsRepository fetches enriched records.
type EnrichedRecordsRepository interface {
	Store(records []EnrichedRecord) (err error)
}
