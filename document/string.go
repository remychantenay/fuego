package document

import (
	"cloud.google.com/go/firestore"
	"context"
)

// StringField provides the necessary to interact with a Firestore document field of type String.
type StringField interface {

	// Retrieve returns the value of a specific field containing an string.
	Retrieve(ctx context.Context) (string, error)

	// Update updates the value of a specific field containing a string.
	Update(ctx context.Context, with string) error
}

// String represents a document field of type String.
type String struct {
	Document Document
	Name     string
}

// Retrieve returns the content of a specific field for a given document.
func (s *String) Retrieve(ctx context.Context) (string, error) {
	snapshot, err := s.Document.GetDocumentRef().Get(ctx)
	if err != nil {
		return "", err
	}

	if !snapshot.Exists() {
		return "", ErrDocumentNotExist
	}

	result, err := snapshot.DataAt(s.Name)
	if err != nil {
		return "", err
	}

	return result.(string), nil
}

// Update updates the value of a specific field of type String.
func (s *String) Update(ctx context.Context, with string) error {

	_, err := s.Document.GetDocumentRef().Set(ctx, map[string]interface{}{
		s.Name: with,
	}, firestore.MergeAll)
	return err
}
