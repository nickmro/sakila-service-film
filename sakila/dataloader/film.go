package dataloader

import (
	"strconv"

	"github.com/graph-gophers/dataloader"
)

// FilmKey returns a film key.
func FilmKey(id int) dataloader.Key {
	key := strconv.Itoa(id)

	return dataloader.StringKey(key)
}
