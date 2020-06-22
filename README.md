# fuego (Work In Progress 🏗️)
[![Build Status](https://travis-ci.org/remychantenay/fuego.svg?branch=master)](https://travis-ci.org/remychantenay/fuego)
[![Go Report Card](https://goreportcard.com/badge/github.com/remychantenay/fuego)](https://goreportcard.com/report/github.com/remychantenay/fuego)
[![codebeat badge](https://codebeat.co/badges/60d273d3-08e6-4f48-9c35-86ab75fc1924)](https://codebeat.co/projects/github-com-remychantenay-fuego-master)
[![GoDoc](https://godoc.org/github.com/remychantenay/fuego?status.svg)](https://godoc.org/github.com/remychantenay/fuego)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Fuego is a Go client library for Google Firestore.

Fuego does not do anything crazy or magic – its purpose is to reduce the amount of boilerplate by abstracting the ceremony that comes with interacting with Firestore.

![Fuego](https://raw.githubusercontent.com/remychantenay/fuego/master/art/fuego.jpg)

## Features
### Done
* CRUD operations
* Retrieving specific fields only
* Simple Exists check
* Collections and Queries

### WIP
* Simple maps and arrays updates

### TODO
* Write Batches

## Usage
### Import
```bash
go get -u github.com/remychantenay/fuego
```

### Document
Bear in mind that the examples below are not including how to initialize Firebase and the Firestore client. You will find this information in Firebase's documentation.
#### Create
```go
import "github.com/remychantenay/fuego"

type User struct {
    FirstName string `firestore:"FirstName"`
    LastName  string `firestore:"LastName"`
    EmailAddress  string `firestore:"EmailAddress"`
}

func main() {

    fuego := New(firestoreClient) // firestoreClient needs to be created beforehand.
        
    user := User{
        FirstName: "John",
        LastName:  "Smith",
        EmailAddress:  "jsmith@email.com",
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

### Fields
#### Retrieve
At times, you may want to retrieve the value of only one field:
```go
value, err := fuego.Document("users", "jsmith").Field("FirstName").Retrieve(ctx)
if err != nil {
    panic(err.Error())
}

// Note the required type assertion
fmt.Println("FirstName: ", value.(string)) // prints John
```

#### Update
Same goes for updating a specific field:
```go
err := fuego.Document("users", "jsmith").Field("FirstName").Update(ctx, "Mike")
if err != nil {
    panic(err.Error())
}
```

#### Update (Map)
Fuego provides with different ways to update a field that contains a Map:
* Merge
* Override
* Append

```go
newMap := map[string]interface{}{
    "Android": "dPtQzw_6YU0WctLu0kHye-:APA91bEDAUcMhLB3XHK...",
}

err := fuego.Document("users", "jsmith").
    Field("Tokens").
    UpdateMap(ctx, newMap, document.Merge) // See document/option.go for more info
if err != nil {
    panic(err.Error())
}
```

### Collections
#### Retrieve
Retrieving all the documents in a collection is straight-forward:
```go
users, err := fuego.Collection("users").Retrieve(ctx, &User{})
if err != nil {
    panic(err.Error())
}

// Note the required type assertion
fmt.Println("FirstName: ", users[0].(*User).FirstName) // prints John
```

#### Query
When it comes to performing more complex queries, it just works like the firestore client. Fuego's collection struct embeds a `firestore.Query`. This allow to directly use its methods:
```go
collection := fuego.Collection("users")
query := collection.Where("FirstName", "==", "John").Limit(50)
users, err := collection.RetrieveWith(ctx, &User{}, query)
if err != nil {
    panic(err.Error())
}

// Note the required type assertion
fmt.Println("FirstName: ", users[0].(*User).FirstName) // prints John
```

More info [here](https://firebase.google.com/docs/firestore/query-data/queries) on how to use queries.

## Transitive Dependencies
* Firebase: `firebase.google.com/go`
* Firestore: `cloud.google.com/go/firestore`

More info [here](https://godoc.org/github.com/remychantenay/fuego?imports)

## Resources
[google-cloud-go on GitHub](https://github.com/googleapis/google-cloud-go/tree/master/firestore)

## License
Apache License Version 2.0