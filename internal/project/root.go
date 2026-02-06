package project

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	configFileName = "gob.yaml"
)

type Root struct {
	Path string
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Find(start string) (*Root, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return nil, err
	}
	for {
		if exists(filepath.Join(dir, configFileName)) {
			return &Root{Path: dir}, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return nil, errors.New("gob.yaml not found: not a GoBrain project")
}
