package document

import "errors"

var (
	// ErrDocumentNotExist indicates that the requested document doesn't exist.
	ErrDocumentNotExist = errors.New("document: doesn't exist")

	// ErrFieldRetrieve indicates that the requested field value could not be retrieved.
	ErrFieldRetrieve = errors.New("field: couldn't retrieve the field value")
)
