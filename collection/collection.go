package collection

import (
	"cloud.google.com/go/firestore"
	"context"
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
}

// FirestoreCollection provides features related to Firestore collections.
type FirestoreCollection struct {

	// Ref (firestore.CollectionRef) is a reference to the collection.
	Ref *firestore.CollectionRef

	// Query (firestore.Query) is embedded so its methods can conveniently be used directly.
	firestore.Query
}

// New creates and returns a new FirestoreCollection.
func New(fs *firestore.Client, path string) *FirestoreCollection {
	r := fs.Collection(path)
	return &FirestoreCollection{
		Ref:   r,
		Query: r.Query,
	}
}

// Retrieve retrieve all the documents from a collection.
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
func (c *FirestoreCollection) RetrieveWith(ctx context.Context, sample interface{}, query firestore.Query) ([]interface{}, error) {
	c.Query = query // replacing the embedded query
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
