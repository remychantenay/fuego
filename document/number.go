package document

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/remychantenay/fuego/document/internal"
)

// NumberField provides the necessary to interact with a Firestore document field of type Number.
type NumberField interface {

	// Retrieve returns the value of a specific field containing a number (int64).
	Retrieve(ctx context.Context) (int64, error)

	// Update updates the value of a specific field containing a number (int64).
	Update(ctx context.Context, with int64) error
}

// Number represents a document field of type Number.
type Number struct {
	Document Document
	Name     string
}

// Retrieve returns the content of a specific field for a given document.
//  nb, err := fuego.Document("users", "jsmith").Number("Age").Retrieve(ctx)
func (n *Number) Retrieve(ctx context.Context) (int64, error) {
	value, err := internal.RetrieveFieldValue(ctx, n.Document.GetDocumentRef(), n.Name)
	if err != nil {
		return 0, err
	}

	return value.(int64), nil
}

// Update updates the value of a specific field of type Number.
//  err := fuego.Document("users", "jsmith").Number("Age").Update(ctx, 42)
func (n *Number) Update(ctx context.Context, with int64) error {

	_, err := n.Document.GetDocumentRef().Set(ctx, map[string]int64{
		n.Name: with,
	}, firestore.MergeAll)
	return err
}
