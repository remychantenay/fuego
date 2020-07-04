package document

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/remychantenay/fuego/document/internal"
)

// NumberField provides the necessary to interact with a Firestore document field of type Number.
type NumberField interface {

	// Retrieve returns the value of a specific field containing a number (int64).
	Retrieve(ctx context.Context) (int64, error)

	// Update the value of a specific field containing a number (int64).
	Update(ctx context.Context, with int64) error

	// Increment the value of a specific field containing a number (int64).
	// If the field doesn't exist, it will be set to 1.
	Increment(ctx context.Context) error

	// Decrement the value of a specific field containing a number (int64).
	// If the field doesn't exist, it will be set to 0.
	Decrement(ctx context.Context) error
}

// Number represents a document field of type Number.
type Number struct {

	// Document is the underlying document (incl. ID and ref).
	Document Document

	// Name is the name of the field.
	Name string

	firestore *firestore.Client
}

// Retrieve returns the content of a specific field for a given document.
//  nb, err := fuego.Document("users", "jsmith").Number("Age").Retrieve(ctx)
func (f *Number) Retrieve(ctx context.Context) (int64, error) {
	value, err := internal.RetrieveFieldValue(ctx, f.Document.GetDocumentRef(), f.Name)
	if err != nil {
		return 0, err
	}

	return value.(int64), nil
}

// Update the value of a specific field of type Number.
//  err := fuego.Document("users", "jsmith").Number("Age").Update(ctx, 42).
func (f *Number) Update(ctx context.Context, with int64) error {

	ref := f.Document.GetDocumentRef()
	m := map[string]int64{
		f.Name: with,
	}

	if f.Document.InBatch() {
		f.Document.Batch().Set(ref, m, firestore.MergeAll)
		return nil
	}

	_, err := ref.Set(ctx, m, firestore.MergeAll)
	return err
}

// Increment the value of a specific field of type Number.
//
// The update will be executed inside a transaction.
// If the field doesn't exist, it will be set to 1.
//  err := fuego.Document("users", "jsmith").Number("Age").Increment(ctx)
func (f *Number) Increment(ctx context.Context) error {

	return f.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		ref := f.Document.GetDocumentRef()
		document, err := tx.Get(ref)
		if err != nil {
			return err
		}

		newValue := int64(1)
		value, err := document.DataAt(f.Name)
		if err == nil {
			newValue = value.(int64) + 1
		}

		err = tx.Set(ref, map[string]int64{
			f.Name: newValue,
		}, firestore.MergeAll)
		if err != nil {
			return err
		}

		return nil // Success, no errors
	})
}

// Decrement the value of a specific field of type Number.
//
// The update will be executed inside a transaction.
// If the field doesn't exist, it will be set to 0.
//  err := fuego.Document("users", "jsmith").Number("Age").Decrement(ctx)
func (f *Number) Decrement(ctx context.Context) error {

	return f.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		ref := f.Document.GetDocumentRef()
		document, err := tx.Get(ref)
		if err != nil {
			return err
		}

		newValue := int64(0)
		value, err := document.DataAt(f.Name)
		if err == nil {
			newValue = value.(int64) - 1
		}

		err = tx.Set(ref, map[string]int64{
			f.Name: newValue,
		}, firestore.MergeAll)
		if err != nil {
			return err
		}

		return nil // Success, no errors
	})
}
