package document

import (
	"cloud.google.com/go/firestore"
	"context"
	"time"
)

// TimestampField provides the necessary to interact with a Firestore document field of type Timestamp.
type TimestampField interface {

	// Retrieve returns the value of a specific field containing a timestamp (time.Time).
	//
	// location needs to be a value from the IANA Time Zone database (e.g. "America/Los_Angeles")
	// A time.Time zero value will be returned if an error occurs.
	Retrieve(ctx context.Context, location string) (time.Time, error)

	// Update updates the value of a specific field containing a timestamp (time.Time).
	Update(ctx context.Context, with time.Time) error
}

// Timestamp represents a document field of type Timestamp.
type Timestamp struct {
	Document Document
	Name     string
}

// Retrieve returns the content of a specific field for a given document.
//
// location needs to be a value from the IANA Time Zone database (e.g. "America/Los_Angeles")
// A time.Time zero value will be returned if an error occurs.
func (t *Timestamp) Retrieve(ctx context.Context, location string) (time.Time, error) {
	snapshot, err := t.Document.GetDocumentRef().Get(ctx)
	if err != nil {
		return time.Time{}, err
	}

	if !snapshot.Exists() {
		return time.Time{}, ErrDocumentNotExist
	}

	result, err := snapshot.DataAt(t.Name)
	if err != nil {
		return time.Time{}, err
	}

	loc, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, err
	}
	return result.(time.Time).In(loc), nil
}

// Update updates the value of a specific field of type Timestamp.
func (t *Timestamp) Update(ctx context.Context, with time.Time) error {

	_, err := t.Document.GetDocumentRef().Set(ctx, map[string]interface{}{
		t.Name: with,
	}, firestore.MergeAll)
	return err
}
