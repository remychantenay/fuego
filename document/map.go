package document

import (
	"cloud.google.com/go/firestore"
	"context"
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
	Document Document
	Name     string

	firestore *firestore.Client
}

// Retrieve returns the content of a specific field for a given document.
func (m *Map) Retrieve(ctx context.Context) (map[string]interface{}, error) {
	s, err := m.Document.GetDocumentRef().Get(ctx)
	if err != nil {
		return nil, err
	}

	if !s.Exists() {
		return nil, ErrDocumentNotExist
	}

	result, err := s.DataAt(m.Name)
	if err != nil {
		return nil, err
	}

	return result.(map[string]interface{}), nil
}

// Merge merges the value of a specific Map field.
func (m *Map) Merge(ctx context.Context, data map[string]interface{}) error {
	_, err := m.Document.GetDocumentRef().Set(ctx, map[string]interface{}{
		m.Name: data,
	}, firestore.MergeAll)
	return err
}

// Override simply update (override) the field with a given Map.
func (m *Map) Override(ctx context.Context, data map[string]interface{}) error {
	_, err := m.Document.GetDocumentRef().Set(ctx, map[string]interface{}{
		m.Name: data,
	}, firestore.MergeAll)
	return err
}
