package logger

import (
	"io"

	"github.com/go-kit/kit/log"
)

// GoKitLogger is an instance of a logger using the go-kit library
type GoKitLogger struct {
	client log.Logger
}

// NewGoKitLogger creates an instance of the go kit logger
func NewGoKitLogger(writer io.Writer) Logger {
	return &GoKitLogger{
		client: log.NewLogfmtLogger(writer),
	}
}

// Log is responsible for logging stuff
func (logger *GoKitLogger) Log(keyVals ...interface{}) error {
	err := logger.client.Log(keyVals)
	if err != nil {
		println(err.Error())
	}

	return err
}
