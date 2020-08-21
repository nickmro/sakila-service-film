package http

import (
	"net/http"
	"sakila/sakila-film-service/sakila"
)

var errorStatuses = map[error]int{
	sakila.ErrorNotFound: http.StatusNotFound,
	sakila.ErrorInternal: http.StatusInternalServerError,
}

func statusForError(err error) int {
	if status := errorStatuses[err]; status != 0 {
		return status
	}

	return http.StatusInternalServerError
}
