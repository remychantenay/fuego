package document

import (
	"cloud.google.com/go/firestore"
	"context"
)

// Doc provides the necessary to interact with a Firestore document.
type Doc interface {

	// GetDocumentRef returns a Document Reference (DocumentRef).
	GetDocumentRef() *firestore.DocumentRef

	// Create creates a document.
	Create(ctx context.Context, from interface{}) error

	// Retrieve populate the destination passed as argument.
	//
	// note: the `to` parameter has to be a pointer.
	Retrieve(ctx context.Context, to interface{}) error

	// Exists returns true if the document exists, false otherwise.
	//
	// note: if an error occurs, false is also returned.
	Exists(ctx context.Context) bool

	// Delete removes a document from Firestore.
	Delete(ctx context.Context) error

	// Field returns a specific DocumentField.
	Field(name string) DocumentField
}

// Document provides features related to Firestore documents.
type Document struct {
	FirestoreClient *firestore.Client
	Path            string
	ID              string
}

// GetDocumentRef returns a document reference.
func (d *Document) GetDocumentRef() *firestore.DocumentRef {

	collectionRef := d.FirestoreClient.Collection(d.Path)

	if len(d.ID) == 0 {
		return collectionRef.NewDoc()
	}

	return collectionRef.Doc(d.ID)
}

// Create a document in Firestore.
func (d *Document) Create(ctx context.Context, from interface{}) error {
	_, err := d.GetDocumentRef().Set(ctx, from)
	return err
}

// Retrieve a document from Firestore.
//
// to: the destination must be a pointer.
func (d *Document) Retrieve(ctx context.Context, to interface{}) error {
	s, err := d.GetDocumentRef().Get(ctx)
	if err != nil {
		return err
	}

	if !s.Exists() {
		return ErrDocumentNotExist
	}

	if err := s.DataTo(to); err != nil {
		return err
	}

	return nil
}

// Exists returns true if a given document exists, false otherwise.
func (d *Document) Exists(ctx context.Context) bool {
	s, err := d.GetDocumentRef().Get(ctx)
	if err != nil {
		return false
	}

	if !s.Exists() {
		return false
	}

	return true
}

// Delete removes a document from Firestore.
func (d *Document) Delete(ctx context.Context) error {
	_, err := d.GetDocumentRef().Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Field returns a new Field.
func (d *Document) Field(name string) DocumentField {
	return &Field{
		Document: d,
		Name:     name,
	}
}
