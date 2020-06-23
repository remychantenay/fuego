package internal

import (
	"cloud.google.com/go/firestore"
	"context"
)

// RetrieveFieldValue returns the value of a field.
func RetrieveFieldValue(ctx context.Context, ref *firestore.DocumentRef, fieldName string) (interface{}, error) {
	s, err := ref.Get(ctx)
	if err != nil {
		return nil, err
	}

	return s.DataAt(fieldName)
}
