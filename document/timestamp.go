package document

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/remychantenay/fuego/document/internal"
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

	// Document is the underlying document (incl. ID and ref).
	Document Document

	// Name is the name of the field.
	Name string
}

// Retrieve returns the content of a specific field for a given document.
//
// location needs to be a value from the IANA Time Zone database
// A time.Time zero value will be returned if an error occurs.
//  val, err := fuego.Document("users", "jsmith").Timestamp("LastSeenAt").Retrieve(ctx, "America/Los_Angeles")
func (f *Timestamp) Retrieve(ctx context.Context, location string) (time.Time, error) {
	value, err := internal.RetrieveFieldValue(ctx, f.Document.GetDocumentRef(), f.Name)
	if err != nil {
		return time.Time{}, err
	}

	loc, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, err
	}
	return value.(time.Time).In(loc), nil
}

// Update updates the value of a specific field of type Timestamp.
//  err := fuego.Document("users", "jsmith").Timestamp("LastSeenAt").Update(ctx, time.Now())
func (f *Timestamp) Update(ctx context.Context, with time.Time) error {

	ref := f.Document.GetDocumentRef()
	m := map[string]interface{}{
		f.Name: with,
	}

	if f.Document.InBatch() {
		f.Document.Batch().Set(ref, m, firestore.MergeAll)
		return nil
	}

	_, err := ref.Set(ctx, m, firestore.MergeAll)
	return err
}
