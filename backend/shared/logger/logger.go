package logger

// Logger is the interface used for implementing loggers
type Logger interface {
	Log(...interface{}) error
}
