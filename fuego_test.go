package fuego

import (
	"fmt"
	"os"
	//"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	//"google.golang.org/api/option"
	"context"
	"testing"
)

var fuego *Fuego

type TestedStruct struct {
	FirstName string `firestore:"FirstName"`
	LastName  string `firestore:"LastName"`
}

func TestMain(m *testing.M) {
	err := setup()
	if err == nil {
		exit := m.Run()
		tearDown()
		os.Exit(exit)
	} else {
		panic(err)
	}
}

func setup() error {
	ctx := context.Background()

	projectID, present := os.LookupEnv("FIRESTORE_PROJECT_ID")
	if !present {
		projectID = "fuego-test"
	}

	// 1. Setting up the Firebase project
	firebaseApp, err := getFirebaseAppSession(ctx, projectID)
	if err != nil {
		return err
	}

	// 2. Setting up the Firestore Client
	fsClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		return err
	}

	fuego = New(fsClient)
	return nil
}

func tearDown() {
	//form3Client.Close()
}

// getFirebaseAppSession - get a firebase app session (initialize if there is none)
func getFirebaseAppSession(ctx context.Context, projectID string) (*firebase.App, error) {
	var err error

	conf := &firebase.Config{
		ProjectID: projectID,
	}
	session, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, err
	}

	return session, err
}

func TestIntegration_CreateDocument(t *testing.T) {
	ctx := context.Background()

	user := TestedStruct{
		FirstName: "John",
		LastName:  "Smith",
	}

	err := fuego.Document.Create(ctx, "users", "jsmith", user)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestIntegration_RetrieveDocument(t *testing.T) {
	ctx := context.Background()

	user := TestedStruct{}
	err := fuego.Document.Retrieve(ctx, "users", "jsmith", &user)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestIntegration_RetrieveField(t *testing.T) {
	ctx := context.Background()

	value, err := fuego.Document.Field(ctx, "users", "jsmith", "FirstName")
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println("FirstName: ", value.(string))
}
