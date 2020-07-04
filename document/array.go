package document

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/remychantenay/fuego/document/internal"
)

// ArrayField provides the necessary to interact with a Firestore document field of type Array.
type ArrayField interface {

	// Retrieve returns the value of a specific field containing an array.
	//
	// note : the returned data will require a type assertion.
	Retrieve(ctx context.Context) ([]interface{}, error)

	// Append will append the provided data to the existing data (if any) of an Array field.
	Append(ctx context.Context, data []interface{}) error

	// Override will override the existing data (if any) of an Array field.
	// Note: this is the default behaviour with Firestore.
	Override(ctx context.Context, data []interface{}) error
}

// Array represents a document field of type Array.
type Array struct {

	// Document is the underlying document (incl. ID and ref).
	Document Document

	// Name is the name of the field.
	Name string

	firestore *firestore.Client
}

// Retrieve returns the content of a specific field for a given document.
//  values, err := fuego.Document("users", "jsmith").Array("Address").Retrieve(ctx)
func (f *Array) Retrieve(ctx context.Context) ([]interface{}, error) {
	value, err := internal.RetrieveFieldValue(ctx, f.Document.GetDocumentRef(), f.Name)
	if err != nil {
		return nil, err
	}

	return value.([]interface{}), nil
}

// Override will override the existing data (if any) of an Array field.
//  values, err := fuego.Document("users", "jsmith").Array("Address").Override(ctx, []interface{}{"New Street", "New Building"})
func (f *Array) Override(ctx context.Context, data []interface{}) error {

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

// Append will append the provided data to the existing data (if any) of an Array field.
//
// The update will be executed inside a transaction.
//  values, err := fuego.Document("users", "jsmith").Array("Address").Append(ctx, []interface{}{"More info"})
func (f *Array) Append(ctx context.Context, data []interface{}) error {

	return f.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		document, err := tx.Get(f.Document.GetDocumentRef())
		if err != nil {
			return err
		}

		value, err := document.DataAt(f.Name)
		if err != nil {
			return err
		}

		value = append(value.([]interface{}), data...)

		err = tx.Set(f.Document.GetDocumentRef(), map[string]interface{}{
			f.Name: value,
		}, firestore.MergeAll)
		if err != nil {
			return err
		}

		return nil // Success, no errors
	})
}
