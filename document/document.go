// Package document provides generic methods to interact with Firestore documents.
package document

import (
	"cloud.google.com/go/firestore"
	"context"
)

// Document provides the necessary to interact with a Firestore document.
type Document interface {

	// Create a document.
	Create(ctx context.Context, from interface{}) error

	// Retrieve populate the destination passed as parameter.
	//
	// note: the `to` parameter has to be a pointer.
	Retrieve(ctx context.Context, to interface{}) error

	// Exists returns true if the document exists, false otherwise.
	//
	// note: if an error occurs, false is also returned.
	Exists(ctx context.Context) bool

	// Delete removes a document from Firestore.
	Delete(ctx context.Context) error

	// Array returns a specific Array field.
	Array(name string) *Array

	// String returns a specific String field.
	String(name string) *String

	// Number returns a specific Number field.
	Number(name string) *Number

	// Boolean returns a specific Boolean field.
	Boolean(name string) *Boolean

	// Map returns a specific Map field.
	Map(name string) *Map

	// Timestamp returns a specific Timestamp field.
	Timestamp(name string) *Timestamp

	// GetDocumentRef returns a Document Reference (DocumentRef).
	GetDocumentRef() *firestore.DocumentRef

	// InBatch returns true if a WriteBatch has been started, false otherwise.
	InBatch() bool

	// Batch returns the pointer to the WriteBatch (if any), nil oherwise.
	Batch() *firestore.WriteBatch
}

// FirestoreDocument provides features related to Firestore documents.
type FirestoreDocument struct {

	// ColRef (firestore.CollectionRef) is a reference to the collection.
	ColRef *firestore.CollectionRef

	// ID is the ID of the document
	ID string

	// writeBatch will be nil if not started with fuego.StartBatch() or cancelled with fuego.CancelBatch().
	writeBatch *firestore.WriteBatch

	firestore *firestore.Client
}

// New creates and returns a new FirestoreDocument.
func New(fs *firestore.Client, path, documentID string, wb *firestore.WriteBatch) *FirestoreDocument {
	r := fs.Collection(path)
	return &FirestoreDocument{
		ColRef:     r,
		ID:         documentID,
		firestore:  fs,
		writeBatch: wb,
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
	ref := d.GetDocumentRef()
	if d.InBatch() {
		d.Batch().Set(ref, from)
		return nil
	}
	_, err := ref.Set(ctx, from)
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
	if err != nil || !s.Exists() {
		return false
	}

	return true
}

// Delete removes a document from Firestore.
func (d *FirestoreDocument) Delete(ctx context.Context) error {
	ref := d.GetDocumentRef()
	if d.InBatch() {
		d.Batch().Delete(ref)
		return nil
	}
	_, err := ref.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Array returns a new Array.
func (d *FirestoreDocument) Array(name string) *Array {
	return &Array{
		Document:  d,
		Name:      name,
		firestore: d.firestore,
	}
}

// String returns a new String.
func (d *FirestoreDocument) String(name string) *String {
	return &String{
		Document: d,
		Name:     name,
	}
}

// Number returns a new Number.
func (d *FirestoreDocument) Number(name string) *Number {
	return &Number{
		Document:  d,
		Name:      name,
		firestore: d.firestore,
	}
}

// Boolean returns a new Boolean.
func (d *FirestoreDocument) Boolean(name string) *Boolean {
	return &Boolean{
		Document: d,
		Name:     name,
	}
}

// Map returns a new Map.
func (d *FirestoreDocument) Map(name string) *Map {
	return &Map{
		Document:  d,
		Name:      name,
		firestore: d.firestore,
	}
}

// Timestamp returns a new Timestamp.
func (d *FirestoreDocument) Timestamp(name string) *Timestamp {
	return &Timestamp{
		Document: d,
		Name:     name,
	}
}

// InBatch returns true if a WriteBatch has been started, false otherwise.
func (d *FirestoreDocument) InBatch() bool {
	return d.writeBatch != nil
}

// Batch returns the pointer to the WriteBatch (if any), nil oherwise.
func (d *FirestoreDocument) Batch() *firestore.WriteBatch {
	return d.writeBatch
}
