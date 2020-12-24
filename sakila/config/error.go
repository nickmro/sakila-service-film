package config

// Error is a config error.
type Error string

// ErrorMissing returns a missing config error.
const ErrorMissing = Error("missing")

// Error returns the error as a string.
func (e Error) Error() string {
	return string(e)
}
