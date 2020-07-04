package fuego

import (
	"context"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/remychantenay/fuego/collection"
	"github.com/remychantenay/fuego/document"
)

// Fuego is a wrapper for the Firestore client.
// Contains the Firestore client.
type Fuego struct {

	// FirestoreClient is a ptr to a firestore client.
	FirestoreClient *firestore.Client

	// WriteBatch is a ptr to a writebatch.
	// will be nil if not started with StartBatch() or cancelled with CancelBatch().
	WriteBatch *firestore.WriteBatch
}

// New creates and returns a Fuego wrapper.
func New(fs *firestore.Client) *Fuego {
	return &Fuego{
		FirestoreClient: fs,
		WriteBatch:      nil,
	}
}

// StartBatch starts a write batch.
// Write operations that follow will be added to the batch and processed when CommitBatch is called.
func (f *Fuego) StartBatch() {
	f.WriteBatch = f.FirestoreClient.Batch()
}

// CommitBatch commits the write batch previously started with StartBatch().
func (f *Fuego) CommitBatch(ctx context.Context) ([]*firestore.WriteResult, error) {
	if f.WriteBatch == nil {
		return nil, ErrBatchWriteNotStarted
	}
	return f.WriteBatch.Commit(ctx)
}

// CancelBatch cancels an on-going write batch (if any).
func (f *Fuego) CancelBatch() {
	f.WriteBatch = nil
}

// Document returns a new FirestoreDocument.
func (f *Fuego) Document(path, documentID string) *document.FirestoreDocument {
	return document.New(f.FirestoreClient, cleanPath(path), documentID, f.WriteBatch)
}

// DocumentWithGeneratedID returns a new FirestoreDocument without ID.
func (f *Fuego) DocumentWithGeneratedID(path string) *document.FirestoreDocument {
	return f.Document(path, "")
}

// Collection returns a new FirestoreCollection.
func (f *Fuego) Collection(path string) *collection.FirestoreCollection {
	return collection.New(f.FirestoreClient, cleanPath(path))
}

// cleanPath cleans and returns a given path.
func cleanPath(path string) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	return path
}
