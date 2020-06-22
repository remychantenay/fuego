package fuego

import (
	"fmt"
	"os"
	//"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"github.com/remychantenay/fuego/document"
	"testing"
)

var fuego *Fuego

type TestedStruct struct {
	FirstName    string            `firestore:"FirstName"`
	LastName     string            `firestore:"LastName"`
	EmailAddress string            `firestore:"EmailAddress"`
	Tokens       map[string]string `firestore:"Tokens"`
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

func TestIntegration_Document_Create(t *testing.T) {
	ctx := context.Background()

	user := TestedStruct{
		FirstName:    "John",
		LastName:     "Smith",
		EmailAddress: "jsmith@email.com",
		Tokens: map[string]string{
			"Android": "AND_123",
			"IOS":     "IOS_123",
		},
	}

	err := fuego.Document("users", "jsmith").Create(ctx, user)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestIntegration_Document_Retrieve(t *testing.T) {
	ctx := context.Background()

	user := TestedStruct{}
	err := fuego.Document("users", "jsmith").Retrieve(ctx, &user)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println("FirstName: ", user.FirstName)
}

func TestIntegration_Collection_Retrieve(t *testing.T) {
	ctx := context.Background()

	res, err := fuego.Collection("users").Retrieve(ctx, &TestedStruct{})
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println("FirstName: ", res[0].(*TestedStruct).FirstName)
}

func TestIntegration_Collection_RetrieveWith(t *testing.T) {
	ctx := context.Background()

	collection := fuego.Collection("users")
	query := collection.Where("FirstName", "==", "John").Limit(50)

	// 1. Success
	res, err := collection.RetrieveWith(ctx, &TestedStruct{}, query)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println("FirstName: ", res[0].(*TestedStruct).FirstName)

	// 2. Empty
	query = collection.Where("FirstName", "==", "Jane").Limit(50)
	res, err = collection.RetrieveWith(ctx, &TestedStruct{}, query)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(res) != 0 {
		t.Fatalf("Expected an empty slice.")
	}
}

func TestIntegration_Field_Retrieve(t *testing.T) {
	ctx := context.Background()

	value, err := fuego.Document("users", "jsmith").
		Field("FirstName").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println("FirstName: ", value.(string))
}

func TestIntegration_Field_Update(t *testing.T) {
	ctx := context.Background()

	expectedNewValue := "Mike"

	err := fuego.Document("users", "jsmith").
		Field("FirstName").
		Update(ctx, expectedNewValue)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, err := fuego.Document("users", "jsmith").
		Field("FirstName").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if value.(string) != expectedNewValue {
		t.Fatalf("Got %s but expected %s", value, expectedNewValue)
	}

	fmt.Println("FirstName: ", value.(string))
}

func TestIntegration_Field_UpdateMap_Merge(t *testing.T) {
	ctx := context.Background()

	expectedAndroidToken := "AND_456"
	expectedNewMap := map[string]interface{}{
		"Android": expectedAndroidToken,
	}

	err := fuego.Document("users", "jsmith").
		Field("Tokens").
		UpdateMap(ctx, expectedNewMap, document.Merge)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, err := fuego.Document("users", "jsmith").
		Field("Tokens").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	m := value.(map[string]interface{})

	// Making sure the map has been merged
	if len(m) != 2 {
		t.Fatalf("The map should have been merged, got %d elements", len(m))
	}

	androidToken := m["Android"]
	if androidToken != expectedAndroidToken {
		t.Fatalf("Got %s but expected %s", androidToken, expectedAndroidToken)
	}
}

func TestIntegration_Field_UpdateMapOverride(t *testing.T) {
	ctx := context.Background()

	expectedAndroidToken := "AND_789"
	expectedNewMap := map[string]interface{}{
		"Android": expectedAndroidToken,
	}

	err := fuego.Document("users", "jsmith").
		Field("Tokens").
		UpdateMap(ctx, expectedNewMap, document.Override)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, err := fuego.Document("users", "jsmith").
		Field("Tokens").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	m := value.(map[string]interface{})

	// Making sure the map has been overridden
	if len(m) != 1 {
		t.Fatalf("The map should have been overridden, got %d elements", len(m))
	}

	androidToken := m["Android"]
	if androidToken != expectedAndroidToken {
		t.Fatalf("Got %s but expected %s", androidToken, expectedAndroidToken)
	}
}

func TestIntegration_Document_Exists(t *testing.T) {
	ctx := context.Background()

	value := fuego.Document("users", "jsmith").Exists(ctx)
	if !value {
		t.Fatalf("Wanted %v but got %v", true, value)
	}
}

func TestIntegration_Document_Delete(t *testing.T) {
	ctx := context.Background()

	err := fuego.Document("users", "jsmith").Delete(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
