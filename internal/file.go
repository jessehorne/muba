package internal

import (
	"io"
	"os"
)

func readFromFile(path string) []byte {
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		data, readErr := io.ReadAll(f)
		if readErr == nil {
			return data
		}
	}
	return nil
}
