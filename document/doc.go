/*
Package document provides generic methods to interact with Firestore documents.

A Fuego document provides useful functions allowing to easily manipulate Firestore documents.

Usage

Creating a document is straightforward:

	type User struct {
		FirstName       string              `firestore:"FirstName"`
		LastName        string              `firestore:"LastName"`
		EmailAddress    string              `firestore:"EmailAddress"`
		Address         string              `firestore:"Address"`
		Tokens          map[string]string   `firestore:"Tokens"`
	}

	err := fuego.Document("users", "jsmith").Create(ctx, user)
    if err != nil {
        panic(err)
    }

In some cases, you may want to create a document and let Firebase generated a unique ID for you.
Fuego supports this:

	err := fuego.DocumentWithGeneratedID("users").Create(ctx, user)

Retrieving a document as even simpler:

	user := User{}
	err := fuego.Document("users", "jsmith").Retrieve(ctx, &user)

You also may want to only check if a given document exists without providing a struct:

	// Note: false will be returned if an error occurs as well
	value := fuego.Document("users", "jsmith").Exists(ctx)

Fields - All types

Fuego also allow to easily manipulate specific fields of a document (e.g. retrieve, update, increment, ...)
It supports all the fields offered by Firestore at the exception of geopoints and references.

Usage

Retrieving a field is almost as simple as retrieving a document:

	value, err := fuego.Document("users", "jsmith").String("FirstName").Retrieve(ctx)

Updating a field:

	err := fuego.Document("users", "jsmith").String("FirstName").Update(ctx, "Mike")

Fields - Numbers

Numbers are stored in Firestore as int64. Fuego provides operations that are frequently performed with number fields.

For instance, you may want to increment or decrement a number field.
With Fuego, this can be done in one single operatipn. It uses Transactions underneath.

	err := fuego.Document("users", "jsmith").Number(Age).Increment(ctx)
	// or
	err := fuego.Document("users", "ben_button").Number(Age).Decrement(ctx)

Fields - Arrays

As you may know, dealing with arrays in Firestore documents can be somewhat of a burden.
When updating an array field, the entire field is overridden.
Using Fuego, you can decide to either override or append to the field:

	err := fuego.Document("users", "jsmith").
	Array("Address").
	Override(ctx, []interface{}{"4th Floor"})

	// or...

	err := fuego.Document("users", "jsmith").
	Array("Address").
	Append(ctx, []interface{}{"4th Floor"})

Of course, you can also retrieve an array as follow:

	values, err := fuego.Document("users", "jsmith").Array("Address").Retrieve(ctx)
	if err != nil {
		panic(err)
	}

	// Note the required type assertion
	fmt.Println("First Element: ", values[0].(string))

Fields - Timestamp

Timestamp fields work the same way as the other fields except that a timezone (IANA Time Zone)
needs to be provided at retrieval:

	val, err := fuego.Document("users", "jsmith").Timestamp("LastSeenAt").Retrieve(ctx, "America/Los_Angeles")

*/
package document
