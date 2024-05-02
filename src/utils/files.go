package utils

import (
	"errors"
	"os"
	"strings"
)

func TryCreateDirectory(path string) error {
	err := os.MkdirAll(path, 0750)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		return err
	}
	return err
}

func TryCreateDirectoryPanic(path string) {
	err := TryCreateDirectory(path)
	if err != nil {
		panic(err)
	}
}

func CheckIfFileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		// Schrödinger: file may or may not exist. See err for details.
		panic(err)
	}
}
