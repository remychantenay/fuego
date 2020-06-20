package fuego

import (
	"cloud.google.com/go/firestore"
)

// Fuego is a wrapper for the Firestore client
// It contains different "Services" (e.g. Document, Collection, etc.)
type Fuego struct {
	Document DocumentService
}

// New creates and returns a Fuego wrapper
func New(fc *firestore.Client) *Fuego {
	return &Fuego{
		Document: &Document{
			firestore: fc,
		},
	}
}
