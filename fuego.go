package fuego

import (
	"cloud.google.com/go/firestore"
	"github.com/remychantenay/fuego/collection"
	"github.com/remychantenay/fuego/document"
	"strings"
)

// Fuego is a wrapper for the Firestore client
// It contains the Firestore client
type Fuego struct {
	firestore *firestore.Client
}

// New creates and returns a Fuego wrapper
func New(fc *firestore.Client) *Fuego {
	return &Fuego{
		firestore: fc,
	}
}

// Document returns a new FirestoreDocument.
func (f *Fuego) Document(path, documentID string) *document.FirestoreDocument {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	return document.New(f.firestore, path, documentID)
}

// DocumentWithGeneratedID returns a new FirestoreDocument without ID.
func (f *Fuego) DocumentWithGeneratedID(path string) *document.FirestoreDocument {
	return f.Document(path, "")
}

// Collection returns a new FirestoreCollection.
func (f *Fuego) Collection(path string) *collection.FirestoreCollection {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	return collection.New(f.firestore, path)
}
