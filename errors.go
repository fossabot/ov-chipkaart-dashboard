package main

import "github.com/pkg/errors"

//////////////////////////
//Filters Service
//////////////////////////

// RawRecordsEnrichmentError returns the errors that are recorded during filter operations
type RawRecordsEnrichmentError struct {
	ErrorRecords []ErrorRecord
}

// Error returns the raw records filter error as a string
func (error RawRecordsEnrichmentError) Error() string {
	err := errors.New("RawRecordsEnrichmentError")
	for _, errorRecord := range error.ErrorRecords {
		err = errors.Wrapf(errorRecord.Error, "Record = %+v", errorRecord.Record)
	}
	return err.Error()
}

// ErrorRecord represents the record together with the error
type ErrorRecord struct {
	Record RawRecord
	Error  error
}
