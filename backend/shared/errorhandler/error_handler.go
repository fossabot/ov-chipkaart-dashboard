package errorhandler

import "context"

// ErrorHandler is responsible for handling errors
type ErrorHandler interface {
	CaptureError(ctx context.Context, err error)
}
