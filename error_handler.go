package main

import (
	"log"

	"github.com/pkg/errors"
)

// SentryErrorHandler is an implementation of the error handler which sends errors to sentry
type SentryErrorHandler struct {
}

// NewSentryErrorHandler creates a new instance fo an error handler which sends errors to sentry
func NewSentryErrorHandler() SentryErrorHandler {
	return SentryErrorHandler{}
}

// HandleSoftError is responsible for handling non fatal errors
func (handler SentryErrorHandler) HandleSoftError(err error) {
	log.Printf(errors.Wrapf(err, "%+v", err).Error())
}

// HandleHardError is responsible for handling fatal errors
func (handler SentryErrorHandler) HandleHardError(err error) {
	log.Panicf(errors.Wrapf(err, "%+v", err).Error())
}
