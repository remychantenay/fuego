/*
Package collection provides generic methods to interact with Firestore collections.

A Fuego collection provides useful functions allowing to easily manipulate Firestore collections.

Usage

Retrieving a collection is straightforward:

	users, err := fuego.Collection("users").Retrieve(ctx, &User{})
	if err != nil {
		panic(err)
	}

	// Note the required type assertion
	fmt.Println("FirstName: ", users[0].(*User).FirstName) // prints John

When it comes to performing more complex queries, it just works like the firestore client.
Fuego's collection struct embeds a firestore.Query (https://firebase.google.com/docs/firestore/query-data/queries).
This allow to directly use its methods:

	collection := fuego.Collection("users")
	query := collection.Where("FirstName", "==", "John").Limit(50)
	users, err := collection.RetrieveWith(ctx, &User{}, query)

Fuego also provide with the ability to set a value for a field for all documents within a given collection:

	err := fuego.Collection("users").SetForAll(ctx, "Premium", true) // Yay!

Or simply delete all documents in the collection:

	err := fuego.Collection("users").DeleteAll(ctx)

*/
package collection
