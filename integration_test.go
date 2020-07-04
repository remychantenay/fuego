package fuego

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	firebase "firebase.google.com/go"
)

var fuego *Fuego

type TestedStruct struct {
	FirstName    string            `firestore:"FirstName"`
	LastName     string            `firestore:"LastName"`
	EmailAddress string            `firestore:"EmailAddress"`
	Tokens       map[string]string `firestore:"Tokens"`
	Address      []string          `firestore:"Address"`
	Age          int64             `firestore:"Age"`
	LastSeenAt   time.Time         `firestore:"LastSeenAt"`
	Premium      bool              `firestore:"Premium"`
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
		Address:    []string{"123 Street", "2nd Building"},
		Age:        29,
		LastSeenAt: time.Now(),
		Premium:    false,
	}

	err := fuego.Document("users", "jsmith").Create(ctx, user)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestIntegration_Document_Exists(t *testing.T) {
	ctx := context.Background()

	value := fuego.Document("users", "jsmith").Exists(ctx)
	if !value {
		t.Fatalf("Wanted %v but got %v", true, value)
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

func TestIntegration_String_Retrieve(t *testing.T) {
	ctx := context.Background()

	value, err := fuego.Document("users", "jsmith").
		String("FirstName").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println("FirstName: ", value)
}

func TestIntegration_String_Update(t *testing.T) {
	ctx := context.Background()

	expectedNewValue := "Mike"

	err := fuego.Document("users", "jsmith").
		String("FirstName").
		Update(ctx, expectedNewValue)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, err := fuego.Document("users", "jsmith").
		String("FirstName").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if value != expectedNewValue {
		t.Fatalf("Got %s but expected %s", value, expectedNewValue)
	}

	fmt.Println("New FirstName: ", value)
}

func TestIntegration_Boolean_Retrieve(t *testing.T) {
	ctx := context.Background()

	value, err := fuego.Document("users", "jsmith").
		Boolean("Premium").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println("Premium: ", value)
}

func TestIntegration_Boolean_Update(t *testing.T) {
	ctx := context.Background()

	expectedNewValue := true

	err := fuego.Document("users", "jsmith").
		Boolean("Premium").
		Update(ctx, expectedNewValue)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, err := fuego.Document("users", "jsmith").
		Boolean("Premium").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if value != expectedNewValue {
		t.Fatalf("Got %t but expected %t", value, expectedNewValue)
	}

	fmt.Println("New Premium: ", value)
}

func TestIntegration_Number_Retrieve(t *testing.T) {
	ctx := context.Background()

	value, err := fuego.Document("users", "jsmith").
		Number("Age").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println("Age: ", value)
}

func TestIntegration_Number_Update(t *testing.T) {
	ctx := context.Background()

	expectedNewValue := int64(31)

	err := fuego.Document("users", "jsmith").
		Number("Age").
		Update(ctx, expectedNewValue)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, err := fuego.Document("users", "jsmith").
		Number("Age").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if value != expectedNewValue {
		t.Fatalf("Got %d but expected %d", value, expectedNewValue)
	}

	fmt.Println("New Age: ", value)
}

func TestIntegration_Number_Increment(t *testing.T) {
	ctx := context.Background()

	expectedNewValue := int64(32)

	err := fuego.Document("users", "jsmith").
		Number("Age").
		Increment(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, err := fuego.Document("users", "jsmith").
		Number("Age").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if value != expectedNewValue {
		t.Fatalf("Got %d but expected %d", value, expectedNewValue)
	}

	fmt.Println("Incremented Age: ", value)
}

func TestIntegration_Number_Decrement(t *testing.T) {
	ctx := context.Background()

	expectedNewValue := int64(31)

	err := fuego.Document("users", "jsmith").
		Number("Age").
		Decrement(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, err := fuego.Document("users", "jsmith").
		Number("Age").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if value != expectedNewValue {
		t.Fatalf("Got %d but expected %d", value, expectedNewValue)
	}

	fmt.Println("Decremented Age: ", value)
}

func TestIntegration_Timestamp_Retrieve(t *testing.T) {
	ctx := context.Background()

	value, err := fuego.Document("users", "jsmith").
		Timestamp("LastSeenAt").
		Retrieve(ctx, "Europe/Berlin")
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Println("LastSeenAt: ", value)
}

func TestIntegration_Timestamp_Update(t *testing.T) {
	ctx := context.Background()

	expectedNewTime := time.Now()

	err := fuego.Document("users", "jsmith").
		Timestamp("LastSeenAt").
		Update(ctx, expectedNewTime)
	if err != nil {
		t.Fatalf(err.Error())
	}

	value, err := fuego.Document("users", "jsmith").
		Timestamp("LastSeenAt").
		Retrieve(ctx, "Europe/Berlin")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if !value.Equal(expectedNewTime) {
		t.Fatalf("Got %v but expected %v", value, expectedNewTime)
	}

	fmt.Println("New LastSeenAt: ", value)
}

func TestIntegration_Map_Merge(t *testing.T) {
	ctx := context.Background()

	expectedAndroidToken := "AND_456"
	expectedNewMap := map[string]interface{}{
		"Android": expectedAndroidToken,
	}

	err := fuego.Document("users", "jsmith").
		Map("Tokens").
		Merge(ctx, expectedNewMap)
	if err != nil {
		t.Fatalf(err.Error())
	}

	resultMap, err := fuego.Document("users", "jsmith").
		Map("Tokens").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Making sure the map has been merged
	if len(resultMap) != 2 {
		t.Fatalf("The map should have been merged, got %d elements", len(resultMap))
	}

	androidToken := resultMap["Android"]
	if androidToken != expectedAndroidToken {
		t.Fatalf("Got %s but expected %s", androidToken, expectedAndroidToken)
	}
}

/* TODO: Finish the work on maps
func TestIntegration_Field_OverrideMapWith(t *testing.T) {
	ctx := context.Background()

	expectedAndroidToken := "AND_789"
	expectedNewMap := map[string]interface{}{
		"Android": expectedAndroidToken,
	}

	err := fuego.Document("users", "jsmith").
		Field("Tokens").
		OverrideMapWith(ctx, expectedNewMap)
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
}*/

func TestIntegration_Array_Append(t *testing.T) {
	ctx := context.Background()

	expectedArraySizeAfterWrite := 3

	err := fuego.Document("users", "jsmith").
		Array("Address").
		Append(ctx, []interface{}{"4th Floor"})
	if err != nil {
		t.Fatalf(err.Error())
	}

	values, err := fuego.Document("users", "jsmith").
		Array("Address").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Making sure the data has been added
	if len(values) != expectedArraySizeAfterWrite {
		t.Fatalf("The array is expected to have %d elements, got %d.", expectedArraySizeAfterWrite, len(values))
	}
}

func TestIntegration_Array_Override(t *testing.T) {
	ctx := context.Background()

	expectedArraySizeAfterWrite := 1

	err := fuego.Document("users", "jsmith").
		Array("Address").
		Override(ctx, []interface{}{"4th Floor"})
	if err != nil {
		t.Fatalf(err.Error())
	}

	values, err := fuego.Document("users", "jsmith").
		Array("Address").
		Retrieve(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Making sure the data has been overridden
	if len(values) != expectedArraySizeAfterWrite {
		t.Fatalf("The array is expected to have %d elements, got %d.", expectedArraySizeAfterWrite, len(values))
	}
}

func TestIntegration_Document_Delete(t *testing.T) {
	ctx := context.Background()

	err := fuego.Document("users", "jsmith").Delete(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestIntegration_Document_Create_Many(t *testing.T) {
	ctx := context.Background()

	for i := 0; i < 800; i++ {
		user := TestedStruct{
			FirstName:    "John",
			LastName:     "Smith",
			EmailAddress: "jsmith@email.com",
			Tokens: map[string]string{
				"Android": "AND_123",
				"IOS":     "IOS_123",
			},
			Address:    []string{"123 Street", "2nd Building"},
			Age:        30,
			LastSeenAt: time.Now(),
			Premium:    false,
		}

		err := fuego.DocumentWithGeneratedID("users").Create(ctx, user)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
}

func TestIntegration_Collection_SetForAll(t *testing.T) {
	ctx := context.Background()

	err := fuego.Collection("users").SetForAll(ctx, "LastName", "Doe")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestIntegration_Collection_DeleteAll(t *testing.T) {
	ctx := context.Background()

	err := fuego.Collection("users").DeleteAll(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
