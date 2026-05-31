package todoservice

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
)

type FilePersistance struct {
	path string
}

func NewFilePersitanceWithPath(path string) *FilePersistance {
	return &FilePersistance{
		path: path,
	}
}

func (p *FilePersistance) Load() ([]Todo, error) {
	var todos []Todo

	file, err := os.Open(p.path)
	if err != nil {
		// it's OK if the file does not exist.
		if errors.Is(err, os.ErrNotExist) {
			return []Todo{}, nil
		}
		return todos, fmt.Errorf("error opening persisted file %s: %w", p.path, err)
	}

	if err := gob.NewDecoder(file).Decode(&todos); err != nil {
		return todos, fmt.Errorf("error decoding persisted file %s: %w", p.path, err)
	}

	return todos, nil
}

func (p *FilePersistance) Store(items []Todo) error {
	file, err := os.OpenFile(p.path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error opening persisted file %s: %w", p.path, err)
	}

	if err := gob.NewEncoder(file).Encode(items); err != nil {
		return fmt.Errorf("error encoding to persisted file %s: %w", p.path, err)
	}

	return nil
}
