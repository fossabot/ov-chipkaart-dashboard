package main

import "github.com/pkg/errors"

//////////////////////////
//Filters Service
//////////////////////////

// RawRecordsEnrichmentError returns the errors that are recorded during filter operations
type RawRecordsEnrichmentError struct {
	ErrorRecords []ErrorRawRecord
}

// Error returns the raw records filter error as a string
func (error RawRecordsEnrichmentError) Error() string {
	err := errors.New("RawRecordsEnrichmentError")
	for _, errorRecord := range error.ErrorRecords {
		err = errors.Wrapf(errorRecord.Error, "Record = %+v", errorRecord.Record)
	}
	return err.Error()
}

// ErrorRawRecord represents the record together with the error
type ErrorRawRecord struct {
	Record RawRecord
	Error  error
}

///////////////////////////////
// Price Calculation Service //
///////////////////////////////

// EnrichedRecordsError returns the errors that are recorded during filter operations
type EnrichedRecordsError struct {
	ErrorRecords []ErrorEnrichedRecord
}

// Error returns the raw records filter error as a string
func (error EnrichedRecordsError) Error() string {
	err := errors.New("EnrichedRecordsError")
	for _, errorRecord := range error.ErrorRecords {
		err = errors.Wrapf(errorRecord.Error, "Record = %+v", errorRecord.Record)
	}
	return err.Error()
}

// ErrorEnrichedRecord represents the record together with the error
type ErrorEnrichedRecord struct {
	Record EnrichedRecord
	Error  error
}
