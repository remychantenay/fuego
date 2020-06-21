# fuego (Work In Progress 🏗️)
[![Build Status](https://travis-ci.org/remychantenay/fuego.svg?branch=master)](https://travis-ci.org/remychantenay/fuego)
[![Go Report Card](https://goreportcard.com/badge/github.com/remychantenay/fuego)](https://goreportcard.com/report/github.com/remychantenay/fuego)
[![codebeat badge](https://codebeat.co/badges/60d273d3-08e6-4f48-9c35-86ab75fc1924)](https://codebeat.co/projects/github-com-remychantenay-fuego-master)
[![GoDoc](https://godoc.org/github.com/remychantenay/fuego?status.svg)](https://godoc.org/github.com/remychantenay/fuego)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Fuego is a Go client library for Google Firestore.

Fuego does not do anything crazy or magic – its purpose is to reduce the amount of boilerplate by abstracting the ceremony that comes with interacting with Firestore.

## Features
* CRUD operations
* Write Batches
* Simple maps and arrays updates
* Retrieving specific fields only

## Usage
### Import
```bash
go get -u github.com/remychantenay/fuego
```

### Examples
Bear in mind that the examples below are not including how to initialize Firebase and the Firestore client. You will find this information in Firebase's documentation.
#### Create Document
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
#### Retrieve Document
```go
    user := User{}
    err := fuego.Document("users", "jsmith").Retrieve(ctx, &user)
    if err != nil {
        panic(err.Error())
     }

    fmt.Println("LastName: ", user.FirstName) // prints Smith
```

#### Retrieve Field
At times, you may want to retrieve the value of only one field:
```go
    value, err := fuego.Document("users", "jsmith").Field("FirstName").Retrieve(ctx)
    if err != nil {
        panic(err.Error())
    }

    // Note the required type assertion
    fmt.Println("FirstName: ", value.(string)) // prints John
```

#### Update Field
Same goes for updating a specific field:
```go
    err := fuego.Document("users", "jsmith").Field("FirstName").Update(ctx, "Mike")
    if err != nil {
        panic(err.Error())
    }
```

#### Exists
You also may want to only check if a given document exists without providing a struct:
```go
    // Note: false will be returned if an error occurs as well
    value := fuego.Document("users", "jsmith").Exists(ctx)

    fmt.Println("Exists: ", value)
```

## Transitive Dependencies
* Firebase: `firebase.google.com/go`
* Firestore: `cloud.google.com/go/firestore`

More info [here](https://godoc.org/github.com/remychantenay/fuego?imports)

## Resources
[google-cloud-go on GitHub](https://github.com/googleapis/google-cloud-go/tree/master/firestore)

## License
Apache License Version 2.0