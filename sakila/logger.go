package sakila

// Logger defines the operations for a service logger.
type Logger interface {
	Error(err error)
	Info(args ...interface{})
}
