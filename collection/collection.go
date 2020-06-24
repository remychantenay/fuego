// Package collection provides generic methods to interact with Firestore collections.
package collection

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/remychantenay/fuego/collection/internal"
	"google.golang.org/api/iterator"
)

// Collection provides the necessary to interact with a Firestore collection.
type Collection interface {

	// Retrieve retrieve all the documents from a collection.
	//
	// note: the `sample` is a pointer to an initialized struct of the expected type.
	Retrieve(ctx context.Context, sample interface{}) ([]interface{}, error)

	// RetrieveWith retrieve documents from a collection using the provided Query.
	RetrieveWith(ctx context.Context, sample interface{}, query firestore.Query) ([]interface{}, error)

	// SetForAll sets a field with a given value for ALL documents in the collection.
	//
	// Note: uses Batched writes
	SetForAll(ctx context.Context, fieldName string, fieldValue interface{}) error

	// DeleteAll removes all items from the collection.
	//
	// Use precautiously ;)
	DeleteAll(ctx context.Context) error
}

// FirestoreCollection provides features related to Firestore collections.
type FirestoreCollection struct {

	// Ref (firestore.CollectionRef) is a reference to the collection.
	Ref *firestore.CollectionRef

	// Query (firestore.Query) is embedded so its methods can conveniently be used directly.
	firestore.Query

	fsClient *firestore.Client
}

// New creates and returns a new FirestoreCollection.
func New(fs *firestore.Client, path string) *FirestoreCollection {
	r := fs.Collection(path)
	return &FirestoreCollection{
		Ref:      r,
		Query:    r.Query,
		fsClient: fs,
	}
}

// Retrieve retrieve all the documents from a collection.
//  values, err := fuego.Collection("users").Retrieve(ctx, &User{})
func (c *FirestoreCollection) Retrieve(ctx context.Context, sample interface{}) ([]interface{}, error) {
	result := make([]interface{}, 0)
	it := c.Documents(ctx)
	for {
		doc, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		if err := doc.DataTo(sample); err != nil {
			return nil, err
		}

		result = append(result, sample)
	}
	return result, nil
}

// RetrieveWith retrieve documents from a collection using the provided Query.
//  values, err := fuego.Collection("users").RetrieveWith(ctx, &User{}, query)
func (c *FirestoreCollection) RetrieveWith(ctx context.Context, sample interface{}, query firestore.Query) ([]interface{}, error) {
	c.Query = query // replacing the embedded query
	return c.Retrieve(ctx, sample)
}

// SetForAll will set a field with a given value for ALL documents in the collection.
//  err := fuego.Collection("users").SetForAll(ctx, "NewField", "NewValue")
func (c *FirestoreCollection) SetForAll(ctx context.Context, fieldName string, fieldValue interface{}) error {
	it := c.Ref.DocumentRefs(ctx)
	documentRefs, err := it.GetAll()
	if err != nil {
		return nil
	}

	writeMap := map[string]interface{}{
		fieldName: fieldValue,
	}

	// 1. Preparing the batches
	batchCount := internal.CalculateRequiredBatches(len(documentRefs))
	batches := make([]*firestore.WriteBatch, batchCount, batchCount)
	for i := 0; i < len(batches); i++ {
		batches[i] = c.fsClient.Batch()
	}

	// 2. Writing in the batches
	currentBatch := 0
	currentBatchOpCount := 0
	for i := 0; i < len(documentRefs); i++ {
		ref := documentRefs[i]

		// If we reached the limit for the current batch, we move on to the next
		if currentBatchOpCount == internal.MaxOperationsPerBatchedWrite {
			currentBatch++
			currentBatchOpCount = 0
		}

		currentBatchOpCount++
		batches[currentBatch].Set(ref, writeMap, firestore.MergeAll)
	}

	// 3. Committing all batches
	for i := 0; i < len(batches); i++ {
		_, err := batches[i].Commit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAll removes all items from the collection.
//  err := fuego.Collection("users").DeleteAll(ctx)
func (c *FirestoreCollection) DeleteAll(ctx context.Context) error {
	it := c.Ref.DocumentRefs(ctx)
	documentRefs, err := it.GetAll()
	if err != nil {
		return nil
	}

	// 1. Preparing the batches
	batchCount := internal.CalculateRequiredBatches(len(documentRefs))
	batches := make([]*firestore.WriteBatch, batchCount, batchCount)
	for i := 0; i < len(batches); i++ {
		batches[i] = c.fsClient.Batch()
	}

	// 2. Writing in the batches
	currentBatch := 0
	currentBatchOpCount := 0
	for i := 0; i < len(documentRefs); i++ {
		ref := documentRefs[i]

		// If we reached the limit for the current batch, we move on to the next
		if currentBatchOpCount == internal.MaxOperationsPerBatchedWrite {
			currentBatch++
			currentBatchOpCount = 0
		}

		currentBatchOpCount++
		batches[currentBatch].Delete(ref)
	}

	// 3. Committing all batches
	for i := 0; i < len(batches); i++ {
		_, err := batches[i].Commit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
