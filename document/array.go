package document

import (
	"cloud.google.com/go/firestore"
	"context"
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
	Document Document
	Name     string

	firestore *firestore.Client
}

// Retrieve returns the content of a specific field for a given document.
func (a *Array) Retrieve(ctx context.Context) ([]interface{}, error) {
	s, err := a.Document.GetDocumentRef().Get(ctx)
	if err != nil {
		return nil, err
	}

	if !s.Exists() {
		return nil, ErrDocumentNotExist
	}

	result, err := s.DataAt(a.Name)
	if err != nil {
		return nil, err
	}

	return result.([]interface{}), nil
}

// Override will override the existing data (if any) of an Array field.
func (a *Array) Override(ctx context.Context, data []interface{}) error {
	_, err := a.Document.GetDocumentRef().Set(ctx, map[string]interface{}{
		a.Name: data,
	}, firestore.MergeAll)
	return err
}

// Append will append the provided data to the existing data (if any) of an Array field.
//
// The update will be done inside a transaction as we need to read the value beforehand.
func (a *Array) Append(ctx context.Context, data []interface{}) error {

	return a.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		document, err := tx.Get(a.Document.GetDocumentRef())
		if err != nil {
			return err
		}

		value, err := document.DataAt(a.Name)
		if err != nil {
			return err
		}

		value = append(value.([]interface{}), data...)

		err = tx.Set(a.Document.GetDocumentRef(), map[string]interface{}{
			a.Name: value,
		}, firestore.MergeAll)
		if err != nil {
			return err
		}

		return nil // Success, no errors
	})
}
