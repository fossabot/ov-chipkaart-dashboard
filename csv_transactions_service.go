package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const csvFileColumns = 12
const timestampFormat = "02-01-2006 15:04:05"
const csvDateFormat = "02-01-2006"

// TransactionFetcherCSVService is the container for the Transaction Fetcher CSV Service
type TransactionFetcherCSVService struct {
	csvFileReader CSVFileReader
}

// CSVFileReader reads a CSV file into an array of strings
type CSVFileReader interface {
	ReadAll(file string) (records [][]string, err error)
}

// NewTransactionFetcherCSVService initializes the TransactionFetcherCSVService
func NewTransactionFetcherCSVService(reader CSVFileReader) *TransactionFetcherCSVService {
	return &TransactionFetcherCSVService{
		csvFileReader: reader,
	}
}

// CSVTransactionFetchOptions is config for fetching a records from a CSV file
type CSVTransactionFetchOptions struct {
	fileID     string
	cardNumber string
	startDate  time.Time
	endDate    time.Time
}

// FetchTransactionRecords returns an array of records from a CSV file.
func (service TransactionFetcherCSVService) FetchTransactionRecords(config CSVTransactionFetchOptions) (results []RawRecord, err error) {
	records, err := service.csvFileReader.ReadAll(config.fileID)
	if err != nil {
		return results, errors.Wrapf(err, "cannot read csv file")
	}

	for _, record := range records {
		err = service.validateLine(record)
		if err != nil {
			return results, errors.Wrap(err, "record line is invalid")
		}

		if service.getCardNumber(record) != config.cardNumber {
			continue
		}

		recordNotWithinLimit, err := service.recordIsNotWithinTimeLimit(record, config.startDate, config.endDate)
		if err != nil {
			return results, errors.Wrapf(err, "could not check if error is withing time limit")
		}
		if recordNotWithinLimit {
			continue
		}

		fare, err := service.getFare(record)
		if err != nil {
			return results, errors.Wrapf(err, "cannot get fare as string")
		}

		timestamp, err := service.getTransactionDateTime(record)
		if err != nil {
			return results, errors.Wrapf(err, "cannot parse date into string")
		}

		source := rawRecordSourceCSV
		transactionID := NewTransactionID()
		results = append(results, RawRecord{
			CheckInInfo:         service.getCheckInInfo(record),
			CheckInText:         service.getCheckInText(record),
			Fare:                fare,
			ProductInfo:         service.getProductInfo(record),
			TransactionDateTime: timestamp,
			TransactionInfo:     service.getTransactionInfo(record),
			TransactionName:     service.getTransactionName(record),
			Source:              &source,
			ID:                  &transactionID,
		})
	}

	return results, err
}

func (service TransactionFetcherCSVService) recordIsNotWithinTimeLimit(record []string, start time.Time, end time.Time) (result bool, err error) {
	date, err := time.Parse(csvDateFormat, record[0])
	if err != nil {
		return result, errors.Wrapf(err, "cannot parse date %s using format  %s", record[0], csvDateFormat)
	}

	return start.Unix() > date.Unix() || date.Unix() > end.Unix(), nil
}

func (service TransactionFetcherCSVService) getCardNumber(record []string) string {
	return strings.Replace(record[11], " ", "", -1)
}

func (service TransactionFetcherCSVService) getTransactionName(record []string) TransactionName {
	return TransactionName(record[6])
}

func (service TransactionFetcherCSVService) getTransactionInfo(record []string) string {
	return record[4]
}

func (service TransactionFetcherCSVService) getProductInfo(record []string) string {
	return record[8]
}

// This returns the datetime in milliseconds to make it compatible with the API dateTime
func (service TransactionFetcherCSVService) getTransactionDateTime(record []string) (timestamp int64, err error) {
	dateString := record[0]
	if service.isCheckInTransaction(record) {
		dateString += " " + record[1]
	} else {
		dateString += " " + record[3]
	}
	dateString += ":00"

	date, err := time.Parse(timestampFormat, dateString)
	if err != nil {
		return timestamp, errors.Wrapf(err, "cannot parse date %s using format  %s", dateString, timestampFormat)
	}

	return date.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)), err
}

func (service TransactionFetcherCSVService) getFare(record []string) (fare *float64, err error) {
	fareRecord := record[5]
	if fareRecord == "" {
		return fare, err
	}

	result, err := strconv.ParseFloat(strings.Replace(fareRecord, ",", ".", -1), 64)
	if err != nil {
		return fare, errors.Wrapf(err, "cannot convert fare %s to float", fareRecord)
	}

	return &result, err
}

func (service TransactionFetcherCSVService) getCheckInInfo(record []string) string {
	return record[2]
}

func (service TransactionFetcherCSVService) isCheckInTransaction(record []string) bool {
	return service.getTransactionName(record) != "Check-uit"

}
func (service TransactionFetcherCSVService) getCheckInText(record []string) string {
	if service.isCheckInTransaction(record) {
		return ""
	}

	return "Check-in"
}

func (service TransactionFetcherCSVService) validateLine(record []string) (err error) {
	if len(record) != csvFileColumns {
		return errors.New(fmt.Sprintf("the csv row contains %d columns instead of %d: record = %#+v", len(record), csvFileColumns, record))
	}

	return nil
}
