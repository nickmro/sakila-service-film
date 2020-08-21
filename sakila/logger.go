package sakila

// Logger defines the operations for a service logger.
type Logger interface {
	Error(args ...interface{})
	Info(args ...interface{})
}
