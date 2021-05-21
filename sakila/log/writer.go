package log

import (
	"fmt"
	"strings"

	"github.com/nickmro/sakila-service-film/sakila"

	"go.uber.org/zap"
)

// Writer writes to a log.
type Writer struct {
	*zap.Logger
}

// Environment represents the logger environment type.
type Environment string

const (
	// EnvironmentTest is the test logger type.
	EnvironmentTest = Environment("TEST")
	// EnvironmentDevelopment is the development logger type.
	EnvironmentDevelopment = Environment("DEVELOPMENT")
	// EnvironmentProduction is the production logger type.
	EnvironmentProduction = Environment("PRODUCTION")
)

// NewWriter returns a new logger based on the environment.
func NewWriter(e Environment) (w *Writer, err error) {
	var logger *zap.Logger

	options := []zap.Option{
		zap.AddCallerSkip(1),
	}

	switch e {
	case EnvironmentProduction:
		logger, err = zap.NewProduction(options...)
	case EnvironmentTest:
		logger = zap.NewNop()
	case EnvironmentDevelopment:
		fallthrough
	default:
		logger, err = zap.NewDevelopment(options...)
	}

	if err != nil {
		return nil, err
	}

	return &Writer{Logger: logger}, nil
}

// Error writes an error.
func (w *Writer) Error(err error) {
	if _, ok := err.(sakila.Error); !ok {
		w.Logger.Error(err.Error())
	}
}

// Info writes info.
func (w *Writer) Info(a ...interface{}) {
	w.Logger.Info(joinArgs(a...))
}

// Fatal writes an errpr and panics.
func (w *Writer) Fatal(e error) {
	w.Logger.Fatal(e.Error())
}

// Flush flushes the remaining logs in the writer.
func (w *Writer) Flush() {
	if err := w.Sync(); err != nil {
		w.Logger.Fatal(err.Error())
	}
}

func joinArgs(a ...interface{}) string {
	s := make([]string, len(a))
	for i := range a {
		s[i] = fmt.Sprintf("%v", a[i])
	}

	return strings.Join(s, " ")
}
