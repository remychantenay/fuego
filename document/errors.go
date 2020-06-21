package document

import "errors"

var (
	// ErrDocumentNotExist indicates that the requested document doesn't exist.
	ErrDocumentNotExist = errors.New("document: doesn't exist")
)
