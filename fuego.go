package fuego

import (
	"cloud.google.com/go/firestore"
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

// Document returns a new document.
//
// path: a path is a sequence of IDs separated by slashes.
// E.g. "users/user123/bookmarks".
func (f *Fuego) Document(path, documentID string) Doc {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	return &Document{
		firestore: f.firestore,
		Path:      path,
		ID:        documentID,
	}
}

// DocumentWithGeneratedID returns a new document without ID.
//
// path: a path is a sequence of IDs separated by slashes.
// E.g. "users/user123/bookmarks".
func (f *Fuego) DocumentWithGeneratedID(path string) Doc {
	return f.Document(path, "")
}

// Field returns a new Field.
func (d *Document) Field(name string) DocumentField {
	return &Field{
		Document: d,
		Name:     name,
	}
}
