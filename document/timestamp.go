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
	Document Document
	Name     string
}

// Retrieve returns the content of a specific field for a given document.
//
// location needs to be a value from the IANA Time Zone database
// A time.Time zero value will be returned if an error occurs.
//  val, err := fuego.Document("users", "jsmith").Timestamp("LastSeenAt").Retrieve(ctx, "America/Los_Angeles")
func (t *Timestamp) Retrieve(ctx context.Context, location string) (time.Time, error) {
	value, err := internal.RetrieveFieldValue(ctx, t.Document.GetDocumentRef(), t.Name)
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
func (t *Timestamp) Update(ctx context.Context, with time.Time) error {

	_, err := t.Document.GetDocumentRef().Set(ctx, map[string]interface{}{
		t.Name: with,
	}, firestore.MergeAll)
	return err
}
