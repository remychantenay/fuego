package fuego

import (
	"cloud.google.com/go/firestore"
	"context"
	"strings"
)

// Document provides features related to Firestore documents
type Document struct {
	firestore *firestore.Client
}

// DocumentService interfaces with Firestore's document features
type DocumentService interface {

	// GetDocumentRef returns a Document Reference (DocumentRef).
	GetDocumentRef(path, documentID string) *firestore.DocumentRef

	// Create creates a document.
	Create(ctx context.Context, path, documentID string, doc interface{}) error

	// Retrieve populate the destination passed as argument.
	//
	// note: the `to` parameter has to be a pointer.
	Retrieve(ctx context.Context, path, documentID string, to interface{}) error

	// Field returns the value of a specifc field.
	//
	// note : the returned data will require a type assertion.
	Field(ctx context.Context, path, documentID, fieldName string) (interface{}, error)

	// Exists returns true if the document exists, false otherwise.
	//
	// note: if an error occurs, false is also returned.
	Exists(ctx context.Context, path, documentID string) bool
}

// GetDocumentRef returns a document reference
//
// path: a path is a sequence of IDs separated by slashes.
// E.g. "users/user123/bookmarks"
func (d *Document) GetDocumentRef(path, documentID string) *firestore.DocumentRef {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	return d.firestore.Collection(path).Doc(documentID)
}

// Create a document in Firestore
func (d *Document) Create(ctx context.Context, path, documentID string, doc interface{}) error {
	_, err := d.GetDocumentRef(path, documentID).Set(ctx, doc)
	return err
}

// Retrieve a document from Firestore
//
// to: the destination must be a pointer.
func (d *Document) Retrieve(ctx context.Context, path, documentID string, to interface{}) error {
	s, err := d.GetDocumentRef(path, documentID).Get(ctx)
	if err != nil {
		return err
	}

	if !s.Exists() {
		return ErrDocumentNotExist
	}

	if err := s.DataTo(to); err != nil {
		return err
	}

	return nil
}

// Field returns the content of a specific field for a given document
func (d *Document) Field(ctx context.Context, path, documentID, fieldName string) (interface{}, error) {
	s, err := d.GetDocumentRef(path, documentID).Get(ctx)
	if err != nil {
		return nil, err
	}

	if !s.Exists() {
		return nil, ErrDocumentNotExist
	}

	result, err := s.DataAt(fieldName)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Exists returns true if a given document exists, false otherwise
func (d *Document) Exists(ctx context.Context, path, documentID string) bool {
	s, err := d.GetDocumentRef(path, documentID).Get(ctx)
	if err != nil {
		return false
	}

	if !s.Exists() {
		return false
	}

	return true
}
