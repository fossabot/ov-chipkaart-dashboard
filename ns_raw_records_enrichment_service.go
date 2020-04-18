package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"go.uber.org/ratelimit"
)

// NSRawRecordsEnrichmentService enriches NS records
type NSRawRecordsEnrichmentService struct {
	stationsCodeService NSStationsCodeService
	priceFetcher        NSPriceFetcherService
}

// NewNSRawRecordsEnrichmentService creates a new instance of the NSRawRecordsEnrichmentService
func NewNSRawRecordsEnrichmentService(stationsCodeService NSStationsCodeService, priceFetcher NSPriceFetcherService) NSRawRecordsEnrichmentService {
	return NSRawRecordsEnrichmentService{stationsCodeService, priceFetcher}
}

// Enrich goes over all the raw records and enriches NS specific records.
func (service NSRawRecordsEnrichmentService) Enrich(records []RawRecord) (results RawRecordsEnrichmentResults) {
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

	//counter := 0
	//var newRecords []RawRecord
	//for _, record := range records {
	//	if record.IsNS(){
	//		newRecords = append(newRecords, record)
	//		counter++
	//	}
	//}
	//spew.Dump(counter)
	var prev RawRecord

	rateLimiter := ratelimit.New(3)
	for _, record := range records {
		rateLimiter.Take()
		rawRecordID := record.ID
		transactionID := record.TransactionID
		enrichedRecordID := NewTransactionID()

		// Check if record is check-in record Record
		if record.IsCheckIn() {
			prev = record
			continue
		}

		// We're not setting prev = current here because you normally have to check in before taking the inter-city
		// supplement we don't want to lose that record.
		if record.IsNSSupplement() {
			enrichedRecords = append(enrichedRecords, EnrichedRecord{
				RawRecordID:      rawRecordID,
				TransactionID:    transactionID,
				ID:               &enrichedRecordID,
				StartTime:        record.TransactionDateTime,
				StartTimeIsExact: true,
				CompanyName:      companyNameNS,
				TransactionType:  transactionTypeSupplement,
			})
			continue
		}

		// Record belongs to the RET company and thus it's not an NS record.
		if record.IsRET() {
			prev = record
			continue
		}

		// Record is a checkout record meaning we can calculate the price
		if record.IsCheckOut() {
			if record.IsNS() {
				enrichedRecord, errorRecord := service.getEnrichedNsRecord(prev, record, rawRecordID, transactionID, &enrichedRecordID)
				if errorRecord.Error != nil {
					spew.Dump(errorRecord.Error)
					enrichmentErrors.ErrorRecords = append(enrichmentErrors.ErrorRecords, errorRecord)
				} else {
					enrichedRecords = append(enrichedRecords, enrichedRecord)
				}
			} else {
				// If we're here we know that the record is a check-out record but we don't know the company to which it belongs.
				// We'll check if the journey can be an NS journey and if that's the case, we'll enrich it.
				// If the journey is not a valid NSJourney we don't bother about it.
				enrichedRecord, errorRecord := service.getEnrichedNsRecord(prev, record, rawRecordID, transactionID, &enrichedRecordID)
				if errorRecord.Error == nil {
					enrichedRecords = append(enrichedRecords, enrichedRecord)
				}
			}
		}

		// Record is not check-in/check-out or supplement.
		// It may be something else like adding money into your ov-chipkaart
		prev = record
	}

	log.Println(prev)
	return RawRecordsEnrichmentResults{
		ValidRecords: enrichedRecords,
		Error:        enrichmentErrors,
	}
}

func (service NSRawRecordsEnrichmentService) getEnrichedNsRecord(prev, record RawRecord, rawRecordID, transactionID, newTransactionID *TransactionID) (enrichedRecord EnrichedRecord, errorRecord ErrorRecord) {
	var startTime int64
	var startTimeIsExact = false
	if (prev != RawRecord{} && prev.IsCheckIn() && record.IsNS() && prev.TransactionInfo == record.CheckInInfo) {
		startTime = prev.TransactionDateTime.ToInt64()
		startTimeIsExact = false
	} else {
		startTimeIsExact = false
	}

	spew.Dump(record)

	log.Printf("Fetching station name for %s ", record.CheckInInfo)
	fromStation, err := service.stationsCodeService.GetCodeForStationName(record.CheckInInfo)
	if err != nil {
		return enrichedRecord, ErrorRecord{
			Record: record,
			Error:  errors.Wrapf(err, "cannot get code for station: %s", record.CheckInInfo),
		}
	}

	log.Printf("Fetching station name for %s ", record.TransactionInfo)
	toStation, err := service.stationsCodeService.GetCodeForStationName(record.TransactionInfo)
	if err != nil {
		return enrichedRecord, ErrorRecord{
			Record: record,
			Error:  errors.Wrapf(err, "cannot get code for station: %s", record.TransactionInfo),
		}
	}

	journey := NewNSJourney(record.TransactionDateTime.ToTime(), fromStation.Code, toStation.Code)
	spew.Dump(journey)
	if !startTimeIsExact {
		price, err := service.priceFetcher.FetchPrice(journey)
		if err != nil {
			return enrichedRecord, ErrorRecord{
				Record: record,
				Error:  errors.Wrap(err, "cannot fetch price for journey"),
			}
		}
		startTime = record.TransactionDateTime.ToInt64() - int64(price.EstimatedDurationInMilliSeconds())
	}

	return EnrichedRecord{
		FromStationCode:  journey.FromStationCode,
		ToStationCode:    journey.ToStationCode,
		Duration:         record.TransactionDateTime.ToInt64() - startTime,
		RawRecordID:      rawRecordID,
		TransactionID:    transactionID,
		ID:               newTransactionID,
		StartTime:        TimeInMilliSeconds(startTime),
		StartTimeIsExact: startTimeIsExact,
		CompanyName:      companyNameNS,
		TransactionType:  transactionTypeTravel,
	}, errorRecord
}
