package types

import (
	"golang.org/x/mod/sumdb/dirhash"
)

// Hash1 calculate H1 hash for local provider zip archive
func Hash1(providerArchiveLoc string) (string, error) {
	s, err := dirhash.HashZip(providerArchiveLoc, dirhash.Hash1)
	if err != nil {
		return "", err
	}
	return s, nil
}
