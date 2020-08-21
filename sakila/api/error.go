package api

// Error is an API service error.
type Error struct {
	Err     error  `json:"error"`
	Message string `json:"message,omitempty"`
}

// NewError returns a new API error.
func NewError(err error, msg string) Error {
	return Error{Err: err, Message: msg}
}

func (e Error) Error() string {
	return e.Err.Error()
}
