package fuego

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

	// UpdateMap updates the value of a specific Map field.
	//
	// See option.go for more details in regards to Write options.
	// Default behaviour is equivalent to Merge (default in Firestore)
	UpdateMap(ctx context.Context, with map[string]interface{}, opts ...WriteOption) error

	mergeMap(ctx context.Context, with map[string]interface{}) error
	overrideMap(ctx context.Context, with map[string]interface{}) error
}

// Field represents a document field.
type Field struct {
	Document *Document
	Name     string
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

// UpdateMap updates the value of a specific Map field.
func (f *Field) UpdateMap(ctx context.Context, with map[string]interface{}, opts ...WriteOption) error {
	writeOption := Merge // Default
	if len(opts) > 0 {
		writeOption = opts[0]
	}

	switch writeOption {
	case Override:
		return f.overrideMap(ctx, with)

	case Append:
		// TODO
	}

	return f.mergeMap(ctx, with)
}

// mergeMap simply update (merge) the field with a given Map.
func (f *Field) mergeMap(ctx context.Context, with map[string]interface{}) error {
	return f.Update(ctx, with)
}

// overrideMap simply update (override) the field with a given Map.
func (f *Field) overrideMap(ctx context.Context, with map[string]interface{}) error {
	var m = map[string]interface{}{
		f.Name: with,
	}
	_, err := f.Document.GetDocumentRef().Set(ctx, m)
	return err
}
