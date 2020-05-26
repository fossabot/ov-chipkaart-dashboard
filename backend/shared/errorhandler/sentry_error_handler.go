package errorhandler

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
)

// SentryErrorHandler implements and error handler with Sentry
type SentryErrorHandler struct {
	hub *sentry.Hub
}

// NewSentryErrorHandler creates a new sentry error handler
func NewSentryErrorHandler(options sentry.ClientOptions) (ErrorHandler, error) {
	err := sentry.Init(options)
	if err != nil {
		return nil, err
	}

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(2 * time.Second)

	return &SentryErrorHandler{hub: sentry.CurrentHub()}, nil
}

// CaptureError captures an error
func (sentry *SentryErrorHandler) CaptureError(_ context.Context, err error) {
	sentry.hub.CaptureException(err)
}
