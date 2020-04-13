package main

import (
	"time"

	"github.com/pkg/errors"
)

// NSRawRecordsEnrichmentService enriches NS records
type NSRawRecordsEnrichmentService struct {
	stationsCodeService NSStationsCodeService
	priceFetcher        NSPriceFetcherService
}

// NewNSRecordsFilter creates a new instance of the NSRawRecordsEnrichmentService
func NewNSRecordsFilter(stationsCodeService NSStationsCodeService, priceFetcher NSPriceFetcherService) NSRawRecordsEnrichmentService {
	return NSRawRecordsEnrichmentService{stationsCodeService, priceFetcher}
}

// Enrich goes over all the raw records and enriches NS specific records.
func (filter NSRawRecordsEnrichmentService) Enrich(records []RawRecord) (results RawRecordsEnrichmentResults) {
	var (
		enrichedRecords  []EnrichedRecord
		enrichmentErrors RawRecordsEnrichmentError
	)

	if len(records) == 0 {
		return RawRecordsEnrichmentResults{
			ValidRecords: enrichedRecords,
			Error:        enrichmentErrors,
		}
	}

	var prev RawRecord
	for _, record := range records {
		rawRecordID := *record.ID
		transactionID := *record.TransactionID

		// Check if record is NS Record
		if record.IsCheckIn() {
			prev = record
			continue
		}

		if record.IsNSSupplement() {
			// We're not setting prev = current here because you normally have to check in before taking the inter-city
			// supplement we don't want to lose that record.
			enrichedRecords = append(enrichedRecords, EnrichedRecord{
				RawRecordID:      rawRecordID,
				TransactionID:    transactionID,
				ID:               NewTransactionID(),
				StartTime:        record.TransactionDateTime,
				StartTimeIsExact: true,
				CompanyName:      companyNameNS,
				TransactionType:  transactionTypeSupplement,
			})
			continue
		}

		if record.IsCheckOut() {
			if record.Pto == "NS" || record.ModalType == "Trein" {
				enrichedRecord, errorRecord := filter.getEnrichedNsRecord(record, prev, rawRecordID, transactionID)
				if errorRecord.Error != nil {
					enrichmentErrors.ErrorRecords = append(enrichmentErrors.ErrorRecords, errorRecord)
					prev = record
					continue
				}

				enrichedRecords = append(enrichedRecords, enrichedRecord)
				prev = record
				continue
			}

			// set prev to current
			prev = record

			// Check if from station and to station exist and if there's a valid NS journey between those stations

			// check if fromStation is a valid station
			fromStation, err := filter.stationsCodeService.GetCodeForStationName(record.CheckInInfo)
			if err != nil {
				continue
			}

			// check if toStation is a valid station
			toStation, err := filter.stationsCodeService.GetCodeForStationName(record.TransactionInfo)
			if err != nil {
				continue
			}

			// Checking if a journey exists between both stations
			journey := NewNSJourney(time.Unix(record.TransactionDateTime, 0), fromStation.Code, toStation.Code)
			price, err := filter.priceFetcher.FetchPrice(journey)
			if err != nil {
				continue
			}

			// Calculating start time based on price
			startTime := record.TransactionDateTime - int64(price.EstimateDurationInMilliSeconds())

			// Creating the enriched record.
			enrichedRecords = append(enrichedRecords, EnrichedRecord{
				RawRecordID:      rawRecordID,
				TransactionID:    transactionID,
				ID:               NewTransactionID(),
				StartTime:        startTime,
				StartTimeIsExact: false,
				CompanyName:      companyNameNS,
				TransactionType:  transactionTypeTravel,
			})

			prev = record
			continue
		}

		prev = record
	}

	return RawRecordsEnrichmentResults{
		ValidRecords: enrichedRecords,
		Error:        enrichmentErrors,
	}
}

func (filter NSRawRecordsEnrichmentService) getEnrichedNsRecord(prev, record RawRecord, rawRecordID, transactionID TransactionID) (enrichedRecord EnrichedRecord, errorRecord ErrorRecord) {
	var startTime int64
	var startTimeIsExact = false
	if (prev != RawRecord{} && prev.IsCheckIn() && prev.TransactionInfo == record.CheckInInfo) {
		startTime = prev.TransactionDateTime
		startTimeIsExact = false
	} else {
		fromStation, err := filter.stationsCodeService.GetCodeForStationName(record.CheckInInfo)
		if err != nil {
			return enrichedRecord, ErrorRecord{
				Record: record,
				Error:  errors.Wrap(err, "cannot get code for station"),
			}
		}

		toStation, err := filter.stationsCodeService.GetCodeForStationName(record.TransactionInfo)
		if err != nil {
			return enrichedRecord, ErrorRecord{
				Record: record,
				Error:  errors.Wrap(err, "cannot get code for station"),
			}
		}

		journey := NewNSJourney(time.Unix(record.TransactionDateTime, 0), fromStation.Code, toStation.Code)
		price, err := filter.priceFetcher.FetchPrice(journey)
		if err != nil {
			return enrichedRecord, ErrorRecord{
				Record: record,
				Error:  errors.Wrap(err, "cannot fetch price for journey"),
			}
		}

		startTimeIsExact = false
		startTime = record.TransactionDateTime - int64(price.EstimateDurationInMilliSeconds())
	}

	return EnrichedRecord{
		RawRecordID:      rawRecordID,
		TransactionID:    transactionID,
		ID:               NewTransactionID(),
		StartTime:        startTime,
		StartTimeIsExact: startTimeIsExact,
		CompanyName:      companyNameNS,
		TransactionType:  transactionTypeTravel,
	}, errorRecord
}
