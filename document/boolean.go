package document

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/remychantenay/fuego/document/internal"
)

// BooleanField provides the necessary to interact with a Firestore document field of type Boolean.
type BooleanField interface {

	// Retrieve returns the value of a specific field containing a Boolean (bool).
	Retrieve(ctx context.Context) (bool, error)

	// Update updates the value of a specific field containing a Boolean (bool).
	Update(ctx context.Context, with bool) error
}

// Boolean represents a document field of type Boolean.
type Boolean struct {
	Document Document
	Name     string
}

// Retrieve returns the content of a specific field for a given document.
func (b *Boolean) Retrieve(ctx context.Context) (bool, error) {
	value, err := internal.RetrieveFieldValue(ctx, b.Document.GetDocumentRef(), b.Name)
	if err != nil {
		return false, err
	}

	return value.(bool), nil
}

// Update updates the value of a specific field of type Boolean.
func (b *Boolean) Update(ctx context.Context, with bool) error {

	_, err := b.Document.GetDocumentRef().Set(ctx, map[string]bool{
		b.Name: with,
	}, firestore.MergeAll)
	return err
}
