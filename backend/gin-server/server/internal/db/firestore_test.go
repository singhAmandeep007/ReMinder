package db_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/config"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/constants"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Example struct for testing
type TestItem struct {
	ID          string    `firestore:"-"`
	Name        string    `firestore:"name"`
	Value       int       `firestore:"value"`
	IsActive    bool      `firestore:"is_active"`
	Tags        []string  `firestore:"tags,omitempty"`
	CreatedAt   time.Time `firestore:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at,omitempty"`
	Description string    `firestore:"description,omitempty"`
}

type NestedTestItem struct {
	ID       string `firestore:"-"`
	Name     string `firestore:"name"`
	Metadata struct {
		CreatedBy string `firestore:"created_by"`
		Version   int    `firestore:"version"`
	} `firestore:"metadata"`
}

func setupTestEnvironment(t *testing.T) (*db.FirestoreDatabase, func()) {
	// Set up test environment
	os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")

	// Create config and logger
	// Create config and logger
	cfg := &config.Config{
		DBType:              constants.Firestore,
		UseFirebaseEmulator: true,
		FirebaseProjectID:   "test-project",
	}

	log := logger.New()

	// Create database instance
	database, err := db.NewFirestoreDatabase(cfg, log)
	require.NoError(t, err)

	ctx := context.Background()

	// Connect to database (emulator)
	err = database.Connect(ctx)
	require.NoError(t, err)

	// Return cleanup function
	cleanup := func() {
		database.Close(ctx)
	}

	return database.(*db.FirestoreDatabase), cleanup
}

// Helper function to clean up a collection
func cleanupCollection(t *testing.T, collection db.Collection, filter map[string]interface{}) {
	ctx := context.Background()
	var items []struct {
		ID string `firestore:"-"`
	}

	err := collection.GetAllByCondition(ctx, filter, &items)
	if err != nil {
		t.Logf("Error cleaning up collection: %v", err)
		return
	}

	for _, item := range items {
		collection.DeleteById(ctx, item.ID)
	}
}

func TestFirestoreCreate(t *testing.T) {
	// Set up test environment
	database, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get collection
	collName := "test_collection_create"
	collection := database.Collection(collName)

	// Test context
	ctx := context.Background()

	// Clean up before testing
	cleanupCollection(t, collection, map[string]interface{}{})

	t.Run("Create Simple Item", func(t *testing.T) {
		item := TestItem{
			Name:      "Test Item",
			Value:     42,
			IsActive:  true,
			CreatedAt: time.Now(),
		}

		id, err := collection.Create(ctx, item)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		// Clean up
		err = collection.DeleteById(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Create With Tags", func(t *testing.T) {
		item := TestItem{
			Name:      "Tagged Item",
			Value:     100,
			IsActive:  true,
			Tags:      []string{"test", "important", "new"},
			CreatedAt: time.Now(),
		}

		id, err := collection.Create(ctx, item)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		// Verify tags were saved
		var retrieved TestItem
		err = collection.GetById(ctx, id, &retrieved)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(retrieved.Tags))
		assert.Contains(t, retrieved.Tags, "important")

		// Clean up
		err = collection.DeleteById(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Create With Nested Data", func(t *testing.T) {
		item := NestedTestItem{
			Name: "Nested Item",
		}
		item.Metadata.CreatedBy = "test_user"
		item.Metadata.Version = 1

		id, err := collection.Create(ctx, item)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		// Verify nested data was saved
		var retrieved NestedTestItem
		err = collection.GetById(ctx, id, &retrieved)
		assert.NoError(t, err)
		assert.Equal(t, "test_user", retrieved.Metadata.CreatedBy)
		assert.Equal(t, 1, retrieved.Metadata.Version)

		// Clean up
		err = collection.DeleteById(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Create With Map Data", func(t *testing.T) {
		item := map[string]interface{}{
			"name":       "Map Item",
			"value":      77,
			"is_active":  true,
			"created_at": time.Now(),
		}

		id, err := collection.Create(ctx, item)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		// Verify map data was saved
		var retrieved map[string]interface{}
		err = collection.GetById(ctx, id, &retrieved)
		assert.NoError(t, err)
		assert.Equal(t, "Map Item", retrieved["name"])
		assert.Equal(t, int64(77), retrieved["value"])

		// Clean up
		err = collection.DeleteById(ctx, id)
		assert.NoError(t, err)
	})
}

func TestFirestoreRead(t *testing.T) {
	// Set up test environment
	database, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get collection
	collName := "test_collection_read"
	collection := database.Collection(collName)

	// Test context
	ctx := context.Background()

	// Clean up before testing
	cleanupCollection(t, collection, map[string]interface{}{})

	// Create test items
	items := []TestItem{
		{
			Name:        "Read Item 1",
			Value:       10,
			IsActive:    true,
			Tags:        []string{"read", "test"},
			CreatedAt:   time.Now(),
			Description: "First test item",
		},
		{
			Name:        "Read Item 2",
			Value:       20,
			IsActive:    true,
			Tags:        []string{"read", "important"},
			CreatedAt:   time.Now(),
			Description: "Second test item",
		},
		{
			Name:        "Read Item 3",
			Value:       30,
			IsActive:    false,
			Tags:        []string{"read", "archived"},
			CreatedAt:   time.Now(),
			Description: "Third test item",
		},
	}

	// Create all items and store IDs
	for i := range items {
		id, err := collection.Create(ctx, items[i])
		require.NoError(t, err)
		items[i].ID = id
	}

	// Run tests and clean up after all tests
	defer func() {
		for _, item := range items {
			collection.DeleteById(ctx, item.ID)
		}
	}()

	t.Run("Get By ID", func(t *testing.T) {
		var retrieved TestItem
		err := collection.GetById(ctx, items[0].ID, &retrieved)
		assert.NoError(t, err)
		assert.Equal(t, items[0].Name, retrieved.Name)
		assert.Equal(t, items[0].Value, retrieved.Value)
		assert.Equal(t, 2, len(retrieved.Tags))
	})

	t.Run("Get By ID - Not Found", func(t *testing.T) {
		var retrieved TestItem
		err := collection.GetById(ctx, "non-existent-id", &retrieved)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Get One By Condition", func(t *testing.T) {
		var retrieved TestItem
		err := collection.GetOne(ctx, map[string]interface{}{"name": "Read Item 2"}, &retrieved)
		assert.NoError(t, err)
		assert.Equal(t, items[1].Name, retrieved.Name)
		assert.Equal(t, 20, retrieved.Value)
	})

	t.Run("Get One By Condition - Multiple Fields", func(t *testing.T) {
		var retrieved TestItem
		err := collection.GetOne(ctx, map[string]interface{}{
			"value":     30,
			"is_active": false,
		}, &retrieved)
		assert.NoError(t, err)
		assert.Equal(t, items[2].Name, retrieved.Name)
		assert.Equal(t, 30, retrieved.Value)
	})

	t.Run("Get One By Condition - Not Found", func(t *testing.T) {
		var retrieved TestItem
		err := collection.GetOne(ctx, map[string]interface{}{"name": "Non-existent Item"}, &retrieved)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), db.ErrNotFound.Error())
	})

	t.Run("Get All By Condition", func(t *testing.T) {
		var retrieved []TestItem
		err := collection.GetAllByCondition(ctx, map[string]interface{}{"is_active": true}, &retrieved)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(retrieved))
	})

	t.Run("Get All By Condition - Empty Result", func(t *testing.T) {
		var retrieved []TestItem
		err := collection.GetAllByCondition(ctx, map[string]interface{}{"name": "Non-existent"}, &retrieved)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(retrieved))
	})

	t.Run("Count Documents", func(t *testing.T) {
		count, err := collection.Count(ctx, map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	t.Run("Count Documents With Filter", func(t *testing.T) {
		count, err := collection.Count(ctx, map[string]interface{}{"is_active": false})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})
}

func TestFirestoreUpdate(t *testing.T) {
	// Set up test environment
	database, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get collection
	collName := "test_collection_update"
	collection := database.Collection(collName)

	// Test context
	ctx := context.Background()

	// Clean up before testing
	cleanupCollection(t, collection, map[string]interface{}{})

	t.Run("Update With Map - Partial", func(t *testing.T) {
		// Create test item
		item := TestItem{
			Name:        "Update Test Item",
			Value:       50,
			IsActive:    true,
			Tags:        []string{"update", "test"},
			CreatedAt:   time.Now(),
			Description: "Item for update testing",
		}

		id, err := collection.Create(ctx, item)
		require.NoError(t, err)

		// Update only specific fields
		updateData := map[string]interface{}{
			"value":     100,
			"is_active": false,
		}

		err = collection.UpdateById(ctx, id, updateData)
		assert.NoError(t, err)

		// Verify update
		var updated TestItem
		err = collection.GetById(ctx, id, &updated)
		assert.NoError(t, err)
		assert.Equal(t, 100, updated.Value)
		assert.False(t, updated.IsActive)
		assert.Equal(t, item.Name, updated.Name)               // Unchanged
		assert.Equal(t, item.Description, updated.Description) // Unchanged
		assert.Equal(t, 2, len(updated.Tags))                  // Unchanged

		// Clean up
		err = collection.DeleteById(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Update With Struct - Full", func(t *testing.T) {
		// Create test item
		item := TestItem{
			Name:        "Full Update Test Item",
			Value:       25,
			IsActive:    true,
			Tags:        []string{"update", "test"},
			CreatedAt:   time.Now(),
			Description: "Item for full update testing",
		}

		id, err := collection.Create(ctx, item)
		require.NoError(t, err)

		// Update entire item
		item.Value = 75
		item.IsActive = false
		item.Tags = []string{"updated"}
		item.Description = "Updated description"
		item.UpdatedAt = time.Now()

		err = collection.UpdateById(ctx, id, item)
		assert.NoError(t, err)

		// Verify update
		var updated TestItem
		err = collection.GetById(ctx, id, &updated)
		assert.NoError(t, err)
		assert.Equal(t, 75, updated.Value)
		assert.False(t, updated.IsActive)
		assert.Equal(t, "Updated description", updated.Description)
		assert.Equal(t, 1, len(updated.Tags))
		assert.Equal(t, "updated", updated.Tags[0])

		// Clean up
		err = collection.DeleteById(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Update With Nested Data", func(t *testing.T) {
		// Create test item
		item := NestedTestItem{
			Name: "Nested Update Item",
		}
		item.Metadata.CreatedBy = "original_user"
		item.Metadata.Version = 1

		id, err := collection.Create(ctx, item)
		require.NoError(t, err)

		// Update nested data
		updateData := map[string]interface{}{
			"metadata": map[string]interface{}{
				"version": 2,
			},
		}

		err = collection.UpdateById(ctx, id, updateData)
		assert.NoError(t, err)

		// Verify update
		var updated NestedTestItem
		err = collection.GetById(ctx, id, &updated)
		assert.NoError(t, err)
		assert.Equal(t, 2, updated.Metadata.Version)
		assert.Equal(t, "original_user", updated.Metadata.CreatedBy) // Unchanged

		// Clean up
		err = collection.DeleteById(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Update Non-existent Document", func(t *testing.T) {
		updateData := map[string]interface{}{
			"value": 999,
		}

		err := collection.UpdateById(ctx, "non-existent-id", updateData)
		assert.NoError(t, err) // Firestore doesn't return an error when updating non-existent docs
	})
}

func TestFirestoreDelete(t *testing.T) {
	// Set up test environment
	database, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get collection
	collName := "test_collection_delete"
	collection := database.Collection(collName)

	// Test context
	ctx := context.Background()

	// Clean up before testing
	cleanupCollection(t, collection, map[string]interface{}{})

	t.Run("Delete Existing Document", func(t *testing.T) {
		// Create test item
		item := TestItem{
			Name:      "Delete Test Item",
			Value:     100,
			IsActive:  true,
			CreatedAt: time.Now(),
		}

		id, err := collection.Create(ctx, item)
		require.NoError(t, err)

		// Delete the item
		err = collection.DeleteById(ctx, id)
		assert.NoError(t, err)

		// Verify deletion
		var retrieved TestItem
		err = collection.GetById(ctx, id, &retrieved)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Delete Non-existent Document", func(t *testing.T) {
		err := collection.DeleteById(ctx, "non-existent-id")
		assert.NoError(t, err) // Firestore doesn't return an error when deleting non-existent docs
	})

	t.Run("Delete And Verify Count", func(t *testing.T) {
		// Create multiple items
		for i := 0; i < 5; i++ {
			item := TestItem{
				Name:      "Count Test Item",
				Value:     i,
				IsActive:  true,
				CreatedAt: time.Now(),
			}
			_, err := collection.Create(ctx, item)
			require.NoError(t, err)
		}

		// Verify initial count
		initialCount, err := collection.Count(ctx, map[string]interface{}{"name": "Count Test Item"})
		assert.NoError(t, err)
		assert.Equal(t, int64(5), initialCount)

		// Get all items
		var items []TestItem
		err = collection.GetAllByCondition(ctx, map[string]interface{}{"name": "Count Test Item"}, &items)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(items))

		// Delete first item
		err = collection.DeleteById(ctx, items[0].ID)
		assert.NoError(t, err)

		// Verify updated count
		updatedCount, err := collection.Count(ctx, map[string]interface{}{"name": "Count Test Item"})
		assert.NoError(t, err)
		assert.Equal(t, int64(4), updatedCount)

		// Clean up remaining items
		for i := 1; i < len(items); i++ {
			collection.DeleteById(ctx, items[i].ID)
		}
	})
}

func TestFirestoreEdgeCases(t *testing.T) {
	// Set up test environment
	database, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get collection
	collName := "test_collection_edge"
	collection := database.Collection(collName)

	// Test context
	ctx := context.Background()

	// Clean up before testing
	cleanupCollection(t, collection, map[string]interface{}{})

	t.Run("Empty Struct Fields", func(t *testing.T) {
		item := TestItem{
			Name:      "",
			CreatedAt: time.Now(),
		}

		id, err := collection.Create(ctx, item)
		assert.NoError(t, err)

		var retrieved TestItem
		err = collection.GetById(ctx, id, &retrieved)
		assert.NoError(t, err)
		assert.Equal(t, "", retrieved.Name)
		assert.Equal(t, 0, retrieved.Value)
		assert.False(t, retrieved.IsActive)

		// Clean up
		err = collection.DeleteById(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Special Characters in Document Data", func(t *testing.T) {
		item := TestItem{
			Name:        "Special Chars: @#$%^&*()[]{}!?",
			Description: "Line 1\nLine 2\tTabbed\r\nWindows",
			CreatedAt:   time.Now(),
		}

		id, err := collection.Create(ctx, item)
		assert.NoError(t, err)

		var retrieved TestItem
		err = collection.GetById(ctx, id, &retrieved)
		assert.NoError(t, err)
		assert.Equal(t, item.Name, retrieved.Name)
		assert.Equal(t, item.Description, retrieved.Description)

		// Clean up
		err = collection.DeleteById(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Empty Arrays and Maps", func(t *testing.T) {
		item := TestItem{
			Name:      "Empty Arrays Test",
			Tags:      []string{},
			CreatedAt: time.Now(),
		}

		id, err := collection.Create(ctx, item)
		assert.NoError(t, err)

		var retrieved TestItem
		err = collection.GetById(ctx, id, &retrieved)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(retrieved.Tags))

		// Clean up
		err = collection.DeleteById(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("Query with Multiple Conditions", func(t *testing.T) {
		// Create test items with various conditions
		for i := 0; i < 10; i++ {
			active := i%2 == 0
			value := i * 10
			item := TestItem{
				Name:      "Filter Test Item",
				Value:     value,
				IsActive:  active,
				CreatedAt: time.Now(),
			}
			_, err := collection.Create(ctx, item)
			require.NoError(t, err)
		}

		// Query with multiple conditions
		var results []TestItem
		filter := map[string]interface{}{
			"name":      "Filter Test Item",
			"is_active": true,
		}

		err := collection.GetAllByCondition(ctx, filter, &results)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(results))

		for _, item := range results {
			assert.True(t, item.IsActive)
			assert.Equal(t, "Filter Test Item", item.Name)
		}

		// Clean up
		cleanupCollection(t, collection, map[string]interface{}{"name": "Filter Test Item"})
	})

	t.Run("Concurrent Updates", func(t *testing.T) {
		// Create test item
		item := TestItem{
			Name:      "Concurrent Item",
			Value:     1,
			IsActive:  true,
			CreatedAt: time.Now(),
		}

		id, err := collection.Create(ctx, item)
		require.NoError(t, err)

		// Perform concurrent updates
		done := make(chan bool)
		errors := make(chan error, 5)

		for i := 0; i < 5; i++ {
			go func(val int) {
				updateData := map[string]interface{}{
					"value": val,
				}

				err := collection.UpdateById(ctx, id, updateData)
				if err != nil {
					errors <- err
				}
				done <- true
			}(i * 10)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 5; i++ {
			<-done
		}

		// Check if any errors occurred
		select {
		case err := <-errors:
			assert.Fail(t, "Error in concurrent update", err)
		default:
			// No errors
		}

		// Verify the item was updated
		var updated TestItem
		err = collection.GetById(ctx, id, &updated)
		assert.NoError(t, err)
		assert.Equal(t, "Concurrent Item", updated.Name)

		// Clean up
		err = collection.DeleteById(ctx, id)
		assert.NoError(t, err)
	})
}

func TestFirestoreTransactions(t *testing.T) {
	// Set up test environment
	database, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Get collection
	collName := "test_collection_transaction"
	collection := database.Collection(collName)

	// Test context
	ctx := context.Background()

	// Clean up before testing
	cleanupCollection(t, collection, map[string]interface{}{})

	// This test case demonstrates how you might test transactions
	// However, actual transaction implementation would be in your database methods
	t.Run("Atomic Counter Update", func(t *testing.T) {
		// Create a counter document
		counter := map[string]interface{}{
			"name":  "test_counter",
			"value": 0,
		}

		counterId, err := collection.Create(ctx, counter)
		require.NoError(t, err)

		// Update counter multiple times
		for i := 0; i < 5; i++ {
			// Get current value
			var currentCounter map[string]interface{}
			err = collection.GetById(ctx, counterId, &currentCounter)
			assert.NoError(t, err)

			currentValue := int64(0)
			if v, ok := currentCounter["value"].(int64); ok {
				currentValue = v
			}

			// Increment value
			updateData := map[string]interface{}{
				"value": currentValue + 1,
			}

			err = collection.UpdateById(ctx, counterId, updateData)
			assert.NoError(t, err)
		}

		// Verify final value
		var finalCounter map[string]interface{}
		err = collection.GetById(ctx, counterId, &finalCounter)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), finalCounter["value"])

		// Clean up
		err = collection.DeleteById(ctx, counterId)
		assert.NoError(t, err)
	})
}
