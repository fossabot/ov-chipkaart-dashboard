package main

import "time"

const (
	companyNS = "NS"
)

const timeoutAPIRequest = 100 * time.Millisecond

const dateFormat = "2006-01-02"

const yearFormat = "2006"

// basic fare for all transport using NS. This is the 2020 fare
const basicFare = 98

// Used to get the time of the journey. so journey time = cost multiplier * (journey price - base fare)
const costMultiplier = 5
