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

### Set Up
```go
import "github.com/remychantenay/fuego"

func main() {
    fuegoClient := fuego.New(firestoreClient) // firestoreClient needs to be created beforehand.
    /// Your code here...
}
```

### Document
Bear in mind that the examples below are not including how to initialize Firebase and nor create the Firestore client. Check Firebase's documentation for more info.

#### Create
```go
type User struct {
    FirstName       string              `firestore:"FirstName"`
    LastName        string              `firestore:"LastName"`
    EmailAddress    string              `firestore:"EmailAddress"`
    Address         string              `firestore:"Address"`
    Tokens          map[string]string   `firestore:"Tokens"`
}

func main() {
    user := User{
        FirstName:      "John",
        LastName:       "Smith",
        EmailAddress:   "jsmith@email.com",
        Address:        []string{"123 Street", "2nd Building"},
        Tokens: map[string]string{
            "Android": "cHzDGYfl5G8vnNCd9xQsbZ:APA...",
        },
    }

    err := fuegoClient.Document("users", "jsmith").Create(ctx, user)
    if err != nil {
        panic(err)
    }
}
```
#### Retrieve
```go
user := User{}
err := fuegoClient.Document("users", "jsmith").Retrieve(ctx, &user)
if err != nil {
    panic(err)
}

fmt.Println("LastName: ", user.LastName) // prints: Smith
```

#### Exists
You also may want to only check if a given document exists without providing a struct:
```go
// Note: false will be returned if an error occurs as well
value := fuegoClient.Document("users", "jsmith").Exists(ctx)
fmt.Println("Exists: ", value)
```

Please read the [doc](https://godoc.org/github.com/remychantenay/fuego/document) to see all the documents related operations.

### Collections
Below the list of operations available for collections (read the doc for more details):
| Operation | Description | Additional Information |
| ------ | ------ | ------ |
| Retrieve | Retrieve all documents from a collection. | |
| RetrieveWith | Retrieve documents from a collection (if they meet the criteras). | Uses firestore.Query. |
| SetForAll | Set a field value for all documents in the collection | Uses Write Batches. |
| DeleteAll | Removes all documents from a collection. | Uses Write Batches. |

Please read the [doc](https://godoc.org/github.com/remychantenay/fuego/collection) to see all the collections related operations.

### Write Batches
Write Batches allow to group writes together to avoid multiple round trips. They are **NOT** transactions.
More info [here](https://firebase.google.com/docs/firestore/manage-data/transactions).

**IMPORTANT**: not all operations are compatible with Write Batches.
```go
fuegoClient.StartBatch()

// All **supported** operations executed between here and commit will be batched.
fuegoClient.Document("users", "jsmith").Number("Age").Update(33)
fuegoClient.Document("users", "enorton").String("FirstName").Update("Eddy")
fuegoClient.Document("users", "jdoe").Boolean("Premium").Update(true)

wr, err := fuegoClient.CommitBatch(ctx)
if err != nil {
    panic(err)
}
```

## Dependencies
* Firebase: `firebase.google.com/go`
* Firestore: `cloud.google.com/go/firestore`

More info [here](https://godoc.org/github.com/remychantenay/fuego?imports)

## License
Apache License Version 2.0