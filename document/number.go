package document

import (
	"cloud.google.com/go/firestore"
	"context"
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
func (n *Number) Retrieve(ctx context.Context) (int64, error) {
	snapshot, err := n.Document.GetDocumentRef().Get(ctx)
	if err != nil {
		return 0, err
	}

	if !snapshot.Exists() {
		return 0, ErrDocumentNotExist
	}

	result, err := snapshot.DataAt(n.Name)
	if err != nil {
		return 0, err
	}

	return result.(int64), nil
}

// Update updates the value of a specific field of type Number.
func (n *Number) Update(ctx context.Context, with int64) error {

	_, err := n.Document.GetDocumentRef().Set(ctx, map[string]interface{}{
		n.Name: with,
	}, firestore.MergeAll)
	return err
}
