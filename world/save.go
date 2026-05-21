package world

import (
	"encoding/gob"
	"io"
)

func Encode(w *World, out io.Writer) error {
	return gob.NewEncoder(out).Encode(w)
}

func Decode(in io.Reader) (*World, error) {
	var w World
	if err := gob.NewDecoder(in).Decode(&w); err != nil {
		return nil, err
	}
	return &w, nil
}
