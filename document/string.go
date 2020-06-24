package document

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/remychantenay/fuego/document/internal"
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
//  str, err := fuego.Document("users", "jsmith").String("FirstName").Retrieve(ctx)
func (s *String) Retrieve(ctx context.Context) (string, error) {
	value, err := internal.RetrieveFieldValue(ctx, s.Document.GetDocumentRef(), s.Name)
	if err != nil {
		return "", err
	}

	return value.(string), nil
}

// Update updates the value of a specific field of type String.
//  err := fuego.Document("users", "jsmith").String("FirstName").Update(ctx, "Jane")
func (s *String) Update(ctx context.Context, with string) error {

	_, err := s.Document.GetDocumentRef().Set(ctx, map[string]string{
		s.Name: with,
	}, firestore.MergeAll)
	return err
}
