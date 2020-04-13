package main

import (
	"crypto/md5"
	"net/http"
	"strings"
	"time"
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
	CheckInInfo            string   `json:"checkInInfo" bson:"check_in_info"`
	CheckInText            string   `json:"checkInText" bson:"check_in_text"`
	Fare                   *float64 `json:"fare" bson:"fare"`
	FareCalculation        string   `json:"fareCalculation" bson:"fare_calculation"`
	FareText               string   `json:"fareText" bson:"fare_text"`
	ModalType              string   `json:"modalType" bson:"modal_type"`
	ProductInfo            string   `json:"productInfo" bson:"product_info"`
	ProductText            string   `json:"productText" bson:"product_text"`
	Pto                    string   `json:"pto" bson:"product_text"`
	TransactionDateTime    int64    `json:"transactionDateTime" bson:"transaction_timestamp"`
	TransactionInfo        string   `json:"transactionInfo" bson:"transaction_info"`
	TransactionName        string   `json:"transactionName" bson:"transaction_name"`
	EPurseMut              *float64 `json:"ePurseMut" bson:"e_purse_mut"`
	EPurseMutInfo          string   `json:"ePurseMutInfo" bson:"e_purse_mut_info"`
	TransactionExplanation string   `json:"transactionExplanation" bson:"transaction_explanation"`
	TransactionPriority    string   `json:"transactionPriority" bson:"transaction_priority"`
}

// HTTPClient is the class used to perform http requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RawRecordsRepository is used to persist raw transaction records
type RawRecordsRepository interface {
	Store(records []RawRecord, id TransactionID) (err error)
}

// NSJourneyPrice represents the price for an NS journey
type NSJourneyPrice struct {
	DBTimestamp
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

// NSPricesRepository is responsible for saving and loading the NSJourneyPrice for an journey
type NSPricesRepository interface {
	Store(price NSJourneyPrice) (err error)
	GetByHash(hash string) (price NSJourneyPrice, err error)
}

// NSStationsRepository is responsible for saving and loading NSStation structs
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
	Put(key interface{}, value interface{})
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
	Year            string `bson:"Year"`
	FromStationCode string `bson:"from_station_code"`
	ToStationCode   string `bson:"to_station_code"`
	date            time.Time
}

// ToMap converts  the JS journey struct to a `map[string]string` map
func (journey NSJourney) ToMap() map[string]string {
	return map[string]string{
		"date":        journey.date.Format(dateFormat),
		"toStation":   journey.FromStationCode,
		"fromStation": journey.ToStationCode,
	}
}

// NSPriceHash gets the hash for an ns journey used to determine the price of the journey
func (journey NSJourney) NSPriceHash() string {
	hashBytes := md5.Sum([]byte(journey.FromStationCode + hashSeparator + journey.ToStationCode + hashSeparator + journey.Year))
	return string(hashBytes[:])
}
