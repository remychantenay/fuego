package fuego

import (
	"cloud.google.com/go/firestore"
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

// Document returns a new document.
//
// path: a path is a sequence of IDs separated by slashes.
// E.g. "users/user123/bookmarks".
func (f *Fuego) Document(path, documentID string) document.Doc {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	return &document.Document{
		FirestoreClient: f.firestore,
		Path:      path,
		ID:        documentID,
	}
}

// DocumentWithGeneratedID returns a new document without ID.
//
// path: a path is a sequence of IDs separated by slashes.
// E.g. "users/user123/bookmarks".
func (f *Fuego) DocumentWithGeneratedID(path string) document.Doc {
	return f.Document(path, "")
}
