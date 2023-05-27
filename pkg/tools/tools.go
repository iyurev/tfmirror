package tools

import (
	"golang.org/x/mod/sumdb/dirhash"
	"os"
)

// Hash1 calculate H1 hash for local provider zip archive
func Hash1(providerArchiveLoc string) (string, error) {
	s, err := dirhash.HashZip(providerArchiveLoc, dirhash.Hash1)
	if err != nil {
		return "", err
	}
	return s, nil
}

func IsExists(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
