package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/config"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
)

type FirestoreDatabase struct {
	config     *config.Config
	logger     *logger.Logger
	client     *firestore.Client
	app        *firebase.App
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewFirestoreDatabase(config *config.Config, logger *logger.Logger) (Database, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &FirestoreDatabase{
		config:     config,
		logger:     logger,
		ctx:        ctx,
		cancelFunc: cancel,
	}, nil
}

func (f *FirestoreDatabase) Connect(ctx context.Context) error {
	f.logger.Infof("Connecting to Firestore database")

	var app *firebase.App
	var err error
	var opts []option.ClientOption

	// Configure Firestore client based on environment
	fmt.Println("Using Firestore emulator:", f.config.UseFirebaseEmulator)
	if f.config.UseFirebaseEmulator {
		// Using emulator for local development
		emulatorHost := f.config.FirebaseEmulatorHost
		if emulatorHost == "" {
			emulatorHost = "localhost:8081"
		}

		f.logger.Infof("Using Firestore emulator at %s", emulatorHost)
		opts = append(opts, option.WithoutAuthentication())

		// Set FIRESTORE_EMULATOR_HOST environment variable programmatically
		err = f.setEmulatorEnv(emulatorHost)
		if err != nil {
			return fmt.Errorf("failed to set emulator environment: %v", err)
		}
	} else {
		// Using production Firestore
		if f.config.FirebaseGoogleAppCredentials != "" {
			f.logger.Infof("Using credentials file: %s", f.config.FirebaseGoogleAppCredentials)
			opts = append(opts, option.WithCredentialsFile(f.config.FirebaseGoogleAppCredentials))
		}
	}

	// Configure and initialize Firebase app
	conf := &firebase.Config{
		ProjectID: f.config.FirebaseProjectID,
	}

	app, err = firebase.NewApp(f.ctx, conf, opts...)
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase app: %v", err)
	}
	f.app = app

	// Get Firestore client
	client, err := app.Firestore(f.ctx)
	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %v", err)
	}
	f.client = client

	// Test connection
	err = f.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("failed to connect to Firestore: %v", err)
	}

	f.logger.Infof("Successfully connected to Firestore database")
	return nil
}

func (f *FirestoreDatabase) setEmulatorEnv(host string) error {
	// In production code, you might want to use os.Setenv instead of relying on environment variables
	// This is for the FIRESTORE_EMULATOR_HOST
	// For simplicity in this implementation, we're assuming the emulator host is properly set in the environment
	os.Setenv("FIRESTORE_EMULATOR_HOST", host)
	f.logger.Infof("Emulator host set to: %s", host)
	return nil
}

func (f *FirestoreDatabase) Close(ctx context.Context) error {
	f.logger.Infof("Closing Firestore database connection")
	if f.client != nil {
		err := f.client.Close()
		if err != nil {
			return fmt.Errorf("failed to close Firestore client: %v", err)
		}
	}
	f.cancelFunc() // Cancel the context
	return nil
}

func (f *FirestoreDatabase) Ping(ctx context.Context) error {
	f.logger.Infof("Pinging Firestore database")
	if f.client == nil {
		return errors.New("firestore client is not initialized")
	}

	// Try to list collections as a ping
	iter := f.client.Collections(ctx)
	_, err := iter.Next()
	if err != nil && err != iterator.Done {
		return fmt.Errorf("failed to ping Firestore: %v", err)
	}

	return nil
}

func (f *FirestoreDatabase) Migrate(ctx context.Context) error {
	f.logger.Infof("Running Firestore migrations")
	// Firestore is schemaless, so no migrations are needed
	// This could be used to create initial collections, indexes, etc.
	return nil
}

func (f *FirestoreDatabase) Seed(ctx context.Context) error {
	f.logger.Infof("Seeding Firestore database")
	// Implement seeding logic if needed
	// For example, adding initial data to collections
	return nil
}

func (f *FirestoreDatabase) Collection(name string) Collection {
	f.logger.Infof("Getting Firestore collection: %s", name)
	return &FirestoreCollection{
		db:             f,
		collectionName: name,
	}
}

type FirestoreCollection struct {
	db             *FirestoreDatabase
	collectionName string
}

func (f *FirestoreCollection) Create(ctx context.Context, data interface{}) (string, error) {
	f.db.logger.Infof("Creating document in Firestore collection: %s", f.collectionName)

	if f.db.client == nil {
		return "", errors.New("firestore client is not initialized")
	}
	// Extract the ID from the struct
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return "", fmt.Errorf("data is not a struct")
	}

	// Look for an ID field
	var docID string
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return "", fmt.Errorf("struct does not have an ID field")
	}

	// Make sure ID field is a string and not empty
	if idField.Kind() != reflect.String {
		return "", fmt.Errorf("ID field must be a string")
	}

	docID = idField.String()
	if docID == "" {
		return "", fmt.Errorf("ID field cannot be empty")
	}

	// Convert struct to map for Firestore
	dataMap, err := structToMap(data)
	if err != nil {
		return "", fmt.Errorf("failed to convert struct to map: %v", err)
	}

	// Use the ID from the struct to create the document
	_, err = f.db.client.Collection(f.collectionName).Doc(docID).Set(ctx, dataMap)
	if err != nil {
		return "", fmt.Errorf("failed to create document: %v", err)
	}

	f.db.logger.Infof(fmt.Sprintf("Created document in collection %s with ID: %s", f.collectionName, docID))
	return docID, nil
}

func (f *FirestoreCollection) GetById(ctx context.Context, id string, result interface{}) error {
	f.db.logger.Infof("Getting document by ID from Firestore collection: %s", f.collectionName)

	if f.db.client == nil {
		return errors.New("firestore client is not initialized")
	}

	// Get document by ID
	docSnap, err := f.db.client.Collection(f.collectionName).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("document %s not found", id)
		}
		return fmt.Errorf("failed to get document: %v", err)
	}

	// Map the document data to the result interface
	err = docSnap.DataTo(result)
	if err != nil {
		return fmt.Errorf("failed to map document data: %v", err)
	}

	return nil
}

func (f *FirestoreCollection) GetOne(ctx context.Context, filter map[string]interface{}, result interface{}) error {
	f.db.logger.Infof("Getting one document from Firestore collection: %s with filter: %v", f.collectionName, filter)

	if f.db.client == nil {
		return errors.New("firestore client is not initialized")
	}

	// Create a query from the filter
	query := f.db.client.Collection(f.collectionName).Query
	for field, value := range filter {
		query = query.Where(field, "==", value)
	}

	// Limit to one result
	iter := query.Limit(1).Documents(ctx)
	defer iter.Stop()

	// Get the first document
	doc, err := iter.Next()
	if err == iterator.Done {
		return ErrNotFound
	}
	if err != nil {
		return ErrInternal
	}

	// Map the document data to the result interface
	err = doc.DataTo(result)
	if err != nil {
		return fmt.Errorf("failed to map document data: %v", err)
	}

	return nil
}

func (f *FirestoreCollection) GetAllByCondition(ctx context.Context, filter map[string]interface{}, results interface{}) error {
	f.db.logger.Infof("Getting all documents from Firestore collection: %s with filter: %v", f.collectionName, filter)

	if f.db.client == nil {
		return errors.New("firestore client is not initialized")
	}

	// Create a query from the filter
	query := f.db.client.Collection(f.collectionName).Query
	for field, value := range filter {
		query = query.Where(field, "==", value)
	}

	// Execute the query
	iter := query.Documents(ctx)
	defer iter.Stop()

	// Collect all documents
	var docs []*firestore.DocumentSnapshot
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate documents: %v", err)
		}
		docs = append(docs, doc)
	}

	// Handle the results
	resultsVal := reflect.ValueOf(results)
	if resultsVal.Kind() != reflect.Ptr || resultsVal.Elem().Kind() != reflect.Slice {
		return errors.New("results parameter must be a pointer to a slice")
	}

	sliceVal := resultsVal.Elem()
	elemType := sliceVal.Type().Elem()

	for _, doc := range docs {
		// Create a new item of the appropriate type
		item := reflect.New(elemType).Interface()

		// Map the document data to the item
		if err := doc.DataTo(item); err != nil {
			return fmt.Errorf("failed to map document data: %v", err)
		}

		// Add item ID if it has an ID field
		if elemType.Kind() == reflect.Struct {
			if field := reflect.ValueOf(item).Elem().FieldByName("ID"); field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
				field.SetString(doc.Ref.ID)
			}
		}

		// Append the item to the results slice
		sliceVal = reflect.Append(sliceVal, reflect.ValueOf(item).Elem())
	}

	// Update the results slice
	resultsVal.Elem().Set(sliceVal)

	return nil
}

func (f *FirestoreCollection) UpdateById(ctx context.Context, id string, data interface{}) error {
	f.db.logger.Infof("Updating document by ID in Firestore collection: %s", f.collectionName)

	if f.db.client == nil {
		return errors.New("firestore client is not initialized")
	}

	// Check data type and handle accordingly
	switch data.(type) {
	case map[string]interface{}:
		// If it's a map, we can use MergeAll directly
		_, err := f.db.client.Collection(f.collectionName).Doc(id).Set(ctx, data, firestore.MergeAll)
		if err != nil {
			return fmt.Errorf("failed to update document: %v", err)
		}
	default:
		// For structs or other types, we need to handle differently
		// Option 1: Convert to a map first
		dataMap, err := structToMap(data)
		if err != nil {
			// Option 2: If conversion fails, use Replace instead of Merge
			_, err := f.db.client.Collection(f.collectionName).Doc(id).Set(ctx, data)
			if err != nil {
				return fmt.Errorf("failed to update document: %v", err)
			}
			return nil
		}
		// Use the converted map with MergeAll
		_, err = f.db.client.Collection(f.collectionName).Doc(id).Set(ctx, dataMap, firestore.MergeAll)
		if err != nil {
			return fmt.Errorf("failed to update document: %v", err)
		}
	}

	return nil
}

func (f *FirestoreCollection) DeleteById(ctx context.Context, id string) error {
	f.db.logger.Infof("Deleting document by ID from Firestore collection: %s", f.collectionName)

	if f.db.client == nil {
		return errors.New("firestore client is not initialized")
	}

	// Delete the document by ID
	_, err := f.db.client.Collection(f.collectionName).Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete document: %v", err)
	}

	return nil
}

func (f *FirestoreCollection) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	f.db.logger.Infof("Counting documents in Firestore collection: %s with filter: %v", f.collectionName, filter)

	if f.db.client == nil {
		return 0, errors.New("firestore client is not initialized")
	}

	// Create a query from the filter
	query := f.db.client.Collection(f.collectionName).Query
	for field, value := range filter {
		query = query.Where(field, "==", value)
	}

	// Execute the query
	iter := query.Documents(ctx)
	defer iter.Stop()

	// Count documents
	var count int64
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("failed to iterate documents: %v", err)
		}
		count++
	}

	return count, nil
}

func structToMap(data interface{}) (map[string]interface{}, error) {
	// Use reflection to get field values
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("data is not a struct")
	}

	result := make(map[string]interface{})
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Get the firestore tag
		tag := field.Tag.Get("firestore")
		if tag == "-" {
			// Skip fields tagged with firestore:"-"
			continue
		}

		// Parse the tag to handle options like "omitempty"
		parts := strings.Split(tag, ",")
		name := parts[0]

		// If tag is empty, use field name
		if name == "" {
			name = field.Name
		}

		// Check if field should be omitted when empty
		omitEmpty := false
		for _, opt := range parts[1:] {
			if opt == "omitempty" {
				omitEmpty = true
				break
			}
		}

		// Skip empty fields if omitempty is specified
		if omitEmpty {
			// Check for zero values based on field type
			if isZeroValue(fieldValue) {
				continue
			}
		}

		// Add field to map
		result[name] = fieldValue.Interface()
	}

	return result, nil
}

// Helper function to check if a value is the zero value for its type
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		// For structs like time.Time, we can use IsZero() method if available
		if v.Type() == reflect.TypeOf(time.Time{}) {
			return v.Interface().(time.Time).IsZero()
		}
		// For other structs, check if all fields are zero values
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
	return false
}
