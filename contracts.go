package main

import (
	"net/http"
	"time"
)

// TransactionFetchOptions are the options needed when fetching a list of transactions
type TransactionFetchOptions struct {
	Username   string
	Password   string
	CardNumber string
	StartDate  time.Time
	EndDate    time.Time
}

// HTTPClient is the class used to perform http requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RawRecordsRepository is used to persist raw transaction records
type RawRecordsRepository interface {
	Store(records []Record, id TransactionID) (err error)
}
