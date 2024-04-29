package utils

import (
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
