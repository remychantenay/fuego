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

	// Document is the underlying document (incl. ID and ref).
	Document Document

	// Name is the name of the field.
	Name string
}

// Retrieve returns the content of a specific field for a given document.
//  str, err := fuego.Document("users", "jsmith").String("FirstName").Retrieve(ctx)
func (f *String) Retrieve(ctx context.Context) (string, error) {
	value, err := internal.RetrieveFieldValue(ctx, f.Document.GetDocumentRef(), f.Name)
	if err != nil {
		return "", err
	}

	return value.(string), nil
}

// Update updates the value of a specific field of type String.
//  err := fuego.Document("users", "jsmith").String("FirstName").Update(ctx, "Jane")
func (f *String) Update(ctx context.Context, with string) error {

	ref := f.Document.GetDocumentRef()
	m := map[string]string{
		f.Name: with,
	}

	if f.Document.InBatch() {
		f.Document.Batch().Set(ref, m, firestore.MergeAll)
		return nil
	}
	_, err := ref.Set(ctx, m, firestore.MergeAll)
	return err
}
