package main

import (
	"strings"
	"time"
)

// CompanyName is the company to which a transaction belongs
type CompanyName string

// String returns the company name as a string
func (companyName CompanyName) String() string {
	return string(companyName)
}

const (
	companyNameNS  = CompanyName("NS")
	companyNameRET = CompanyName("RET")
)

// TransactionType represents the type of transaction
type TransactionType string

// String returns the transaction type as a string
func (transactionType TransactionType) String() string {
	return string(transactionType)
}

const (
	transactionTypeTravel     = TransactionType("Travel")
	transactionTypeSupplement = TransactionType("Supplement")
)

// TimeInMilliSeconds represents time in milliseconds
type TimeInMilliSeconds int

// ToTime converts time in milliseconds to a time object
func (t TimeInMilliSeconds) ToTime() time.Time {
	return time.Unix(0, int64(t)*1000000)
}

// ToInt64 converts time in milliseconds to an int64 value
func (t TimeInMilliSeconds) ToInt64() int64 {
	return int64(t)
}

// TransactionName  represents the various transaction names
type TransactionName string

// String returns the transaction name as a string
func (name TransactionName) String() string {
	return string(name)
}

// IsTheSameAs is used to compare 2 transaction names
func (name TransactionName) IsTheSameAs(comp TransactionName) bool {
	return strings.ToLower(name.String()) == strings.ToLower(comp.String())
}

const (
	transactionNameCheckIn                  = TransactionName("Check-in")
	transactionNameCheckOut                 = TransactionName("Check-uit")
	transactionNameIntercityDirectSurcharge = TransactionName("Toeslag Intercity Direct")
)

// RawRecordSource is the source for a raw record entry
type RawRecordSource string

// String converts a raw record source to a string
func (source RawRecordSource) String() string {
	return string(source)
}

const (
	rawRecordSourceAPI = RawRecordSource("API")
	rawRecordSourceCSV = RawRecordSource("CSV")
)

const timeoutAPIRequest = 100 * time.Millisecond

const dateFormat = "2006-01-02"

const yearFormat = "2006"

const hashSeparator = "-"

// basic fare for all transport using NS. This is the 2020 fare
const basicFare = 98

// Used to get the time of the journey. so journey time = cost multiplier * (journey price - base fare)
const costMultiplier = 5
