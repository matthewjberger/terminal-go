package world

import (
	"encoding/gob"
	"fmt"
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
	if w.Version != SaveVersion {
		return nil, fmt.Errorf("save version %d does not match runtime version %d", w.Version, SaveVersion)
	}
	return &w, nil
}
