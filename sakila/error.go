package sakila

// Error is a service error.
type Error string

const (
	// ErrorInternal is an internal server error.
	ErrorInternal = Error("internal")
	// ErrorNotFound is a resource not found error.
	ErrorNotFound = Error("not_found")
)

func (e Error) Error() string {
	return string(e)
}
