# fuego (Work In Progress 🏗️)
[![Build Status](https://travis-ci.org/remychantenay/fuego.svg?branch=master)](https://travis-ci.org/remychantenay/fuego)
[![Go Report Card](https://goreportcard.com/badge/github.com/remychantenay/fuego)](https://goreportcard.com/report/github.com/remychantenay/fuego)
[![codebeat badge](https://codebeat.co/badges/60d273d3-08e6-4f48-9c35-86ab75fc1924)](https://codebeat.co/projects/github-com-remychantenay-fuego-master)
[![GoDoc](https://godoc.org/github.com/remychantenay/fuego?status.svg)](https://godoc.org/github.com/remychantenay/fuego)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Fuego is a Go client library for Google Firestore. It can be used in App Engine, Cloud Run, Google Cloud Functions and more.

It does not do anything crazy or magic – the purpose is to reduce the amount of boilerplate by hiding the ceremony that comes with interacting with Firestore.

## Context
While working on a project running on multiple services and serverless functions, I quickly realised that the amount of repetitive and boilerplate code related to Firestore increased exponentially.

I decided to extract this code in a thin layer on top of the Firestore admin client. That's how fuego came to life.

![Fuego](https://raw.githubusercontent.com/remychantenay/fuego/master/art/fuego.jpg)

## Features
### Documents
* CRUD operations
* Retrieving/Updating specific fields only
* Simple Exists check
* Write Batches

### Collections
* Batch processing (update field for all, delete all, etc.)
* Better flexibility with Arrays and Maps (append/merge or override data)

## Usage
### Import
```bash
go get github.com/remychantenay/fuego
```

### Document
Bear in mind that the examples below are not including how to initialize Firebase and nor create the Firestore client. Check Firebase's documentation for more info.
#### Create
```go
import "github.com/remychantenay/fuego"

type User struct {
    FirstName       string              `firestore:"FirstName"`
    LastName        string              `firestore:"LastName"`
    EmailAddress    string              `firestore:"EmailAddress"`
    Address         string              `firestore:"Address"`
    Tokens          map[string]string   `firestore:"Tokens"`
}

func main() {

    fuego := New(firestoreClient) // firestoreClient needs to be created beforehand.
        
    user := User{
        FirstName:      "John",
        LastName:       "Smith",
        EmailAddress:   "jsmith@email.com",
        Address:        []string{"123 Street", "2nd Building"},
        Tokens: map[string]string{
			"Android": "cHzDGYfl5G8vnNCd9xQsbZ:APA...",
		},
    }

    err := fuego.Document("users", "jsmith").Create(ctx, user)
    if err != nil {
        panic(err.Error())
    }

    // Or if you wish to let Firestore generate an ID for you:
    err := fuego.DocumentWithGeneratedID("users").Create(ctx, user)
    if err != nil {
        panic(err.Error())
    }
}
```
#### Retrieve
```go
user := User{}
err := fuego.Document("users", "jsmith").Retrieve(ctx, &user)
if err != nil {
    panic(err.Error())
}

fmt.Println("LastName: ", user.FirstName) // prints Smith
```

#### Exists
You also may want to only check if a given document exists without providing a struct:
```go
// Note: false will be returned if an error occurs as well
value := fuego.Document("users", "jsmith").Exists(ctx)

fmt.Println("Exists: ", value)
```

### Field
#### Retrieve
At times, you may want to retrieve the value of only one field:
```go
value, err := fuego.Document("users", "jsmith").String("FirstName").Retrieve(ctx)
if err != nil {
    panic(err)
}

fmt.Println("FirstName: ", value // prints: John
```

#### Update
Same goes for updating a specific field:
```go
err := fuego.Document("users", "jsmith").String("FirstName").Update(ctx, "Mike")
if err != nil {
    panic(err)
}
```

The example above shows how to retrieve and update a field of type `String` but other types are supported.

### Arrays
##### Retrieving
```go
values, err := fuego.Document("users", "jsmith").Array("Address").Retrieve(ctx)
if err != nil {
    panic(err)
}

// Note the required type assertion
fmt.Println("First Element: ", values[0].(string)) // prints: 123 Street
```

As you can see above, as an array can contain different types of data, type assertion is **required**.

Fuego provides different ways to update an array in a document:

##### Appending
```go
err := fuego.Document("users", "jsmith").
    Array("Address").
    Append(ctx, []interface{}{"4th Floor"})
if err != nil {
    panic(err)
}
```

##### Overriding
```go
err := fuego.Document("users", "jsmith").
    Array("Address").
    Override(ctx, []interface{}{"4th Floor"})
if err != nil {
    panic(err)
}
```

### Maps

##### Retrieving
```go
values, err := fuego.Document("users", "jsmith").Map("Tokens").Retrieve(ctx)
if err != nil {
    panic(err)
}

// Note the required type assertion
fmt.Println("Android Token: ", values["Android"].(string)) // prints: cHzDGYfl5G8vnNCd9xQsbZ:APA...
```

Just like with arrays, type assertion is **required** with maps.

Fuego provides different ways to update a map in a document:

##### Merging
```go
err := fuego.Document("users", "jsmith").
    Map("Token").
    Merge(ctx, map[string]interface{}{
		"Android": "aVxDGYfl5G8vnNCd9xQsbZ:EPE...",
	})
if err != nil {
    panic(err)
}
```

##### Overriding (TODO)
```go
TODO
```

### Collections
#### Retrieve
Retrieving all the documents in a collection is straight-forward:
```go
users, err := fuego.Collection("users").Retrieve(ctx, &User{})
if err != nil {
    panic(err)
}

// Note the required type assertion
fmt.Println("FirstName: ", users[0].(*User).FirstName) // prints John
```

#### Query
When it comes to performing more complex queries, it just works like the firestore client. Fuego's collection struct embeds a [firestore.Query](https://firebase.google.com/docs/firestore/query-data/queries). This allow to directly use its methods:
```go
collection := fuego.Collection("users")
query := collection.Where("FirstName", "==", "John").Limit(50)
users, err := collection.RetrieveWith(ctx, &User{}, query)
if err != nil {
    panic(err)
}

// Note the required type assertion
fmt.Println("FirstName: ", users[0].(*User).FirstName) // prints John
```

### Write Batches
Write Batches allow to group writes together to avoid multiple round trips. They are **NOT** transactions.
More info [here](https://firebase.google.com/docs/firestore/manage-data/transactions).

**IMPORTANT**: not all operations are compatible with Write Batches.
```go
fuego.StartBatch()

// All supported operations executed between here and commit will be batched.
fuego.Document("users", "jsmith").Number("Age").Update(33)
fuego.Document("users", "jdoe").String("FirstName").Update("Jane")
fuego.Document("users", "jdoe").Boolean("Premium").Update(true)

wr, err := fuego.CommitBatch(ctx)
if err != nil {
    panic(err)
}
```

## More Reading
The [doc](https://godoc.org/github.com/remychantenay/fuego) contains examples for more use cases.

## Dependencies
* Firebase: `firebase.google.com/go`
* Firestore: `cloud.google.com/go/firestore`

More info [here](https://godoc.org/github.com/remychantenay/fuego?imports)

## License
Apache License Version 2.0