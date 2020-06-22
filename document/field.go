package document

import (
	"cloud.google.com/go/firestore"
	"context"
)

// DocumentField provides the necessary to interact with a Firestore document field.
type DocumentField interface {

	// Retrieve returns the value of a specifc field.
	//
	// note : the returned data will require a type assertion.
	Retrieve(ctx context.Context) (interface{}, error)

	// UpdateField updates the value of a specific field.
	Update(ctx context.Context, with interface{}) error

	// MergeMapWith will merge the provided data with the existing data (if any) of a Map field.
	// Note: this is the default behaviour with Firestore.
	MergeMapWith(ctx context.Context, data map[string]interface{}) error

	// OverrideMapWith will override the existing data (if any) of a Map field.
	OverrideMapWith(ctx context.Context, data map[string]interface{}) error

	// AppendArray will append the provided data to the existing data (if any) of an Array field.
	AppendArray(ctx context.Context, data []interface{}) error

	// OverrideArray will override the existing data (if any) of an Array field.
	// Note: this is the default behaviour with Firestore.
	OverrideArray(ctx context.Context, data []interface{}) error
}

// Field represents a document field.
type Field struct {
	Document Document
	Name     string

	firestore *firestore.Client
}

// Retrieve returns the content of a specific field for a given document.
func (f *Field) Retrieve(ctx context.Context) (interface{}, error) {
	s, err := f.Document.GetDocumentRef().Get(ctx)
	if err != nil {
		return nil, err
	}

	if !s.Exists() {
		return nil, ErrDocumentNotExist
	}

	result, err := s.DataAt(f.Name)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Update updates the value of a specific field.
func (f *Field) Update(ctx context.Context, with interface{}) error {

	var m = map[string]interface{}{
		f.Name: with,
	}
	_, err := f.Document.GetDocumentRef().Set(ctx, m, firestore.MergeAll)
	return err
}

// MergeMapWith merges the value of a specific Map field.
func (f *Field) MergeMapWith(ctx context.Context, data map[string]interface{}) error {
	var m = map[string]interface{}{
		f.Name: data,
	}
	_, err := f.Document.GetDocumentRef().Set(ctx, m, firestore.MergeAll)
	return err
}

// OverrideMapWith simply update (override) the field with a given Map.
func (f *Field) OverrideMapWith(ctx context.Context, data map[string]interface{}) error {
	var m = map[string]interface{}{
		f.Name: data,
	}
	_, err := f.Document.GetDocumentRef().Set(ctx, m, firestore.MergeAll)
	return err
}

// OverrideArray will override the existing data (if any) of an Array field.
func (f *Field) OverrideArray(ctx context.Context, data []interface{}) error {
	_, err := f.Document.GetDocumentRef().Set(ctx, map[string]interface{}{
		f.Name: data,
	}, firestore.MergeAll)
	return err
}

// AppendArray will append the provided data to the existing data (if any) of an Array field.
//
// The update will be done inside a transaction as we need to read the value beforehand.
func (f *Field) AppendArray(ctx context.Context, data []interface{}) error {

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
