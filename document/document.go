package document

import (
	"cloud.google.com/go/firestore"
	"context"
)

// Document provides the necessary to interact with a Firestore document.
type Document interface {

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

// FirestoreDocument provides features related to Firestore documents.
type FirestoreDocument struct {

	// ColRef (firestore.CollectionRef) is a reference to the collection.
	ColRef *firestore.CollectionRef

	// ID is the ID of the document
	ID string

	firestore *firestore.Client
}

// New creates and returns a new FirestoreDocument.
func New(fs *firestore.Client, path, documentID string) *FirestoreDocument {
	r := fs.Collection(path)
	return &FirestoreDocument{
		ColRef:    r,
		ID:        documentID,
		firestore: fs,
	}
}

// GetDocumentRef returns a document reference.
func (d *FirestoreDocument) GetDocumentRef() *firestore.DocumentRef {

	if len(d.ID) == 0 {
		return d.ColRef.NewDoc()
	}

	return d.ColRef.Doc(d.ID)
}

// Create a document in Firestore.
func (d *FirestoreDocument) Create(ctx context.Context, from interface{}) error {
	_, err := d.GetDocumentRef().Set(ctx, from)
	return err
}

// Retrieve a document from Firestore.
//
// to: the destination must be a pointer.
func (d *FirestoreDocument) Retrieve(ctx context.Context, to interface{}) error {
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
func (d *FirestoreDocument) Exists(ctx context.Context) bool {
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
func (d *FirestoreDocument) Delete(ctx context.Context) error {
	_, err := d.GetDocumentRef().Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Field returns a new Field.
func (d *FirestoreDocument) Field(name string) DocumentField {
	return &Field{
		Document:  d,
		Name:      name,
		firestore: d.firestore,
	}
}
