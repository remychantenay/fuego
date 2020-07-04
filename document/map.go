package document

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/remychantenay/fuego/document/internal"
)

// MapField provides the necessary to interact with a Firestore document field of type Map.
type MapField interface {

	// Retrieve returns the value of a specific field containing a map.
	//
	// note : the returned values will require type assertion.
	Retrieve(ctx context.Context) (map[string]interface{}, error)

	// Merge will merge the provided data with the existing data (if any) of a Map field.
	// Note: this is the default behaviour with Firestore.
	Merge(ctx context.Context, data map[string]interface{}) error

	// Override will override the existing data (if any)  withthe provided data.
	Override(ctx context.Context, data map[string]interface{}) error
}

// Map represents a document field of type Map.
type Map struct {

	// Document is the underlying document (incl. ID and ref).
	Document Document

	// Name is the name of the field.
	Name string

	firestore *firestore.Client
}

// Retrieve returns the content of a specific field for a given document.
func (f *Map) Retrieve(ctx context.Context) (map[string]interface{}, error) {
	value, err := internal.RetrieveFieldValue(ctx, f.Document.GetDocumentRef(), f.Name)
	if err != nil {
		return nil, err
	}

	return value.(map[string]interface{}), nil
}

// Merge merges the value of a specific Map field.
func (f *Map) Merge(ctx context.Context, data map[string]interface{}) error {

	ref := f.Document.GetDocumentRef()
	m := map[string]interface{}{
		f.Name: data,
	}

	if f.Document.InBatch() {
		f.Document.Batch().Set(ref, m, firestore.MergeAll)
		return nil
	}

	_, err := ref.Set(ctx, m, firestore.MergeAll)
	return err
}

// Override simply update (override) the field with a given Map.
func (f *Map) Override(ctx context.Context, data map[string]interface{}) error {

	ref := f.Document.GetDocumentRef()
	m := map[string]interface{}{
		f.Name: data,
	}

	if f.Document.InBatch() {
		f.Document.Batch().Set(ref, m, firestore.MergeAll)
		return nil
	}

	_, err := ref.Set(ctx, m, firestore.MergeAll)
	return err
}
