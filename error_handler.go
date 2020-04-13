package main

import "log"

// SentryErrorHandler is an implementation of the error handler which sends errors to sentry
type SentryErrorHandler struct {
}

// HandleSoftError is responsible for handling non fatal errors
func (handler SentryErrorHandler) HandleSoftError(err error) {
	log.Printf("%+v", err.Error())
}

// HandleHardError is responsible for handling fatal errors
func (handler SentryErrorHandler) HandleHardError(err error) {
	log.Fatalf("%+v", err.Error())
}
