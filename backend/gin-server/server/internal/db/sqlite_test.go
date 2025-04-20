package db

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/config"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/constants"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDBFile = "./test.db"
)

// TestUser represents a user for testing
type TestUser struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Active    bool      `db:"active"`
}

// TestReminder represents a reminder for testing
type TestReminder struct {
	ID              string    `db:"id"`
	Title           string    `db:"title"`
	Description     string    `db:"description"`
	IsPinned        bool      `db:"is_pinned"`
	UserID          string    `db:"user_id"`
	ReminderGroupID string    `db:"reminder_group_id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// TestReminderGroup represents a reminder group for testing
type TestReminderGroup struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// setupDatabase initializes a test database
func setupDatabase(t *testing.T) (Database, func()) {
	// Remove any existing test database
	os.Remove(testDBFile)

	// Configure test logger
	testLogger := logger.New()

	// Configure the SQLite database
	testConfig := &config.Config{
		DBType:     constants.SQLite,
		SQLiteFile: testDBFile,
	}

	// Create a new database instance
	db, err := NewSQLiteDatabase(testConfig, testLogger)
	require.NoError(t, err, "Failed to create SQLite database")

	// Connect to the database
	ctx := context.Background()
	err = db.Connect(ctx)
	require.NoError(t, err, "Failed to connect to SQLite database")

	// Run migrations
	err = db.Migrate(ctx)
	require.NoError(t, err, "Failed to run migrations")

	// Return cleanup function
	cleanup := func() {
		db.Close(ctx)
		os.Remove(testDBFile)
	}

	return db, cleanup
}

// TestSQLiteDatabaseConnection tests database connection functionality
func TestSQLiteDatabaseConnection(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	ctx := context.Background()

	// Test ping functionality
	err := db.Ping(ctx)
	assert.NoError(t, err, "Failed to ping database")

	// Test close and reconnect
	err = db.Close(ctx)
	assert.NoError(t, err, "Failed to close database")

	err = db.Connect(ctx)
	assert.NoError(t, err, "Failed to reconnect to database")
}

// TestSQLiteDatabaseMigrations tests the database migration process
func TestSQLiteDatabaseMigrations(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	ctx := context.Background()

	// Cast to SQLiteDatabase to access internal connection
	sqliteDB, ok := db.(*SQLiteDatabase)
	require.True(t, ok, "Failed to cast to SQLiteDatabase")

	// Verify tables exist by querying the SQLite master table
	tables := []string{"users", "reminders", "reminder_groups"}
	for _, table := range tables {
		var name string
		err := sqliteDB.conn.QueryRowContext(ctx,
			"SELECT name FROM sqlite_master WHERE type='table' AND name=?",
			table).Scan(&name)
		assert.NoError(t, err, "Table %s does not exist", table)
		assert.Equal(t, table, name, "Table name mismatch")
	}
}

// TestSQLiteDatabaseSeed tests the database seeding process
func TestSQLiteDatabaseSeed(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	ctx := context.Background()

	// Run seed
	err := db.Seed(ctx)
	assert.NoError(t, err, "Failed to seed database")

	// Verify admin user exists
	sqliteDB, ok := db.(*SQLiteDatabase)
	require.True(t, ok, "Failed to cast to SQLiteDatabase")

	var count int
	err = sqliteDB.conn.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	assert.NoError(t, err, "Failed to query users table")
	assert.Equal(t, 1, count, "Admin user not found")
}

// TestSQLiteCollectionCreate tests the Create method of SQLiteCollection
func TestSQLiteCollectionCreate(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	ctx := context.Background()

	// Get users collection
	usersCollection := db.Collection("users")

	// Create a test user
	now := time.Now().UTC()
	userID := uuid.New().String()
	user := TestUser{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Insert the user
	id, err := usersCollection.Create(ctx, user)
	assert.NoError(t, err, "Failed to create user")
	assert.Equal(t, userID, id, "Returned ID does not match")

	// Test duplicate error
	_, err = usersCollection.Create(ctx, user)
	assert.Error(t, err, "Expected error for duplicate user")
	assert.ErrorIs(t, err, ErrDuplicate, "Expected duplicate error")
}

// TestSQLiteCollectionGetById tests the GetById method of SQLiteCollection
func TestSQLiteCollectionGetById(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	ctx := context.Background()

	// Get users collection
	usersCollection := db.Collection("users")

	// Create a test user
	now := time.Now().UTC()
	userID := uuid.New().String()
	user := TestUser{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Insert the user
	_, err := usersCollection.Create(ctx, user)
	require.NoError(t, err, "Failed to create user")

	// Get the user by ID
	var foundUser TestUser
	err = usersCollection.GetById(ctx, userID, &foundUser)
	assert.NoError(t, err, "Failed to find user by ID")
	assert.Equal(t, user.ID, foundUser.ID, "User ID mismatch")
	assert.Equal(t, user.Username, foundUser.Username, "Username mismatch")
	assert.Equal(t, user.Email, foundUser.Email, "Email mismatch")

	// Test not found error
	err = usersCollection.GetById(ctx, "non-existent-id", &foundUser)
	assert.Error(t, err, "Expected error for non-existent user")
	assert.ErrorIs(t, err, ErrNotFound, "Expected not found error")
}

// TestSQLiteCollectionCount tests the Count method of SQLiteCollection
func TestSQLiteCollectionCount(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	ctx := context.Background()

	// Get users collection
	usersCollection := db.Collection("users")

	// Count before inserting (should be 0)
	count, err := usersCollection.Count(ctx, map[string]interface{}{})
	assert.NoError(t, err, "Failed to count users")
	assert.Equal(t, int64(0), count, "Expected 0 users")

	// Create multiple test users
	for i := 0; i < 3; i++ {
		user := TestUser{
			ID:        uuid.New().String(),
			Username:  uuid.New().String(),
			Email:     uuid.New().String() + "@example.com",
			Password:  "password123",
			Role:      "user",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		_, err := usersCollection.Create(ctx, user)
		require.NoError(t, err, "Failed to create user")
	}

	// Count after inserting (should be 3)
	count, err = usersCollection.Count(ctx, map[string]interface{}{})
	assert.NoError(t, err, "Failed to count users")
	assert.Equal(t, int64(3), count, "Expected 3 users")

	// Count with filter
	user := TestUser{
		ID:        uuid.New().String(),
		Username:  "special_user",
		Email:     "special@example.com",
		Password:  "password123",
		Role:      "admin",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_, err = usersCollection.Create(ctx, user)
	require.NoError(t, err, "Failed to create special user")

	// Count admin users (should be 1)
	count, err = usersCollection.Count(ctx, map[string]interface{}{"role": "admin"})
	assert.NoError(t, err, "Failed to count admin users")
	assert.Equal(t, int64(1), count, "Expected 1 admin user")
}

// TestRelationships tests the relationships between tables
func TestRelationships(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	ctx := context.Background()

	// Get collections
	usersCollection := db.Collection("users")
	groupsCollection := db.Collection("reminder_groups")
	remindersCollection := db.Collection("reminders")

	// Create a test user
	userID := uuid.New().String()
	user := TestUser{
		ID:        userID,
		Username:  "relationshiptest",
		Email:     "relationship@example.com",
		Password:  "password123",
		Role:      "user",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_, err := usersCollection.Create(ctx, user)
	require.NoError(t, err, "Failed to create user")

	// Create a test reminder group
	groupID := uuid.New().String()
	group := TestReminderGroup{
		ID:        groupID,
		Name:      "Test Group",
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_, err = groupsCollection.Create(ctx, group)
	require.NoError(t, err, "Failed to create reminder group")

	// Create a test reminder
	reminderID := uuid.New().String()
	reminder := TestReminder{
		ID:              reminderID,
		Title:           "Test Reminder",
		Description:     "This is a test reminder",
		IsPinned:        true,
		UserID:          userID,
		ReminderGroupID: groupID,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	_, err = remindersCollection.Create(ctx, reminder)
	require.NoError(t, err, "Failed to create reminder")

	// Verify counts
	count, err := remindersCollection.Count(ctx, map[string]interface{}{"user_id": userID})
	assert.NoError(t, err, "Failed to count reminders")
	assert.Equal(t, int64(1), count, "Expected 1 reminder")

	// Test SET NULL behavior when deleting a group
	sqliteDB, ok := db.(*SQLiteDatabase)
	require.True(t, ok, "Failed to cast to SQLiteDatabase")

	// Delete the reminder group
	_, err = sqliteDB.conn.ExecContext(ctx, "DELETE FROM reminder_groups WHERE id = ?", groupID)
	assert.NoError(t, err, "Failed to delete reminder group")

	// Verify reminder still exists but has NULL group_id
	var reminderResult TestReminder
	err = remindersCollection.GetById(ctx, reminderID, &reminderResult)
	assert.NoError(t, err, "Failed to find reminder after group delete")
	assert.Empty(t, reminderResult.ReminderGroupID, "Expected reminder_group_id to be NULL after group delete")

	// Delete the user and verify cascade
	_, err = sqliteDB.conn.ExecContext(ctx, "DELETE FROM users WHERE id = ?", userID)
	assert.NoError(t, err, "Failed to delete user")

	// Verify group and reminder are deleted (cascade)
	count, err = groupsCollection.Count(ctx, map[string]interface{}{"user_id": userID})
	assert.NoError(t, err, "Failed to count groups after cascade")
	assert.Equal(t, int64(0), count, "Expected 0 groups after cascade delete")

	count, err = remindersCollection.Count(ctx, map[string]interface{}{"user_id": userID})
	assert.NoError(t, err, "Failed to count reminders after cascade")
	assert.Equal(t, int64(0), count, "Expected 0 reminders after cascade delete")
}

// TestDatabaseFactory tests the Factory function
func TestDatabaseFactory(t *testing.T) {
	testLogger := logger.New()

	// Test SQLite factory
	sqliteConfig := &config.Config{
		DBType:     constants.SQLite,
		SQLiteFile: testDBFile,
	}
	db, err := NewDBManager(sqliteConfig, testLogger)
	assert.NoError(t, err, "Failed to create SQLite database")
	assert.IsType(t, &SQLiteDatabase{}, db.DB, "Expected SQLiteDatabase type")

	// Test unsupported database type
	invalidConfig := &config.Config{
		DBType: "unsupported",
	}
	_, err = NewDBManager(invalidConfig, testLogger)
	assert.Error(t, err, "Expected error for unsupported database type")
}

// TestErrorHandling tests various error scenarios
func TestErrorHandling(t *testing.T) {
	testLogger := logger.New()

	// Test empty SQLite file path
	emptyConfig := &config.Config{
		DBType: constants.SQLite,
	}
	_, err := NewSQLiteDatabase(emptyConfig, testLogger)
	assert.Error(t, err, "Expected error for empty SQLite file path")

	// Test using database without connecting
	validConfig := &config.Config{
		DBType:     constants.SQLite,
		SQLiteFile: testDBFile,
	}
	db, err := NewSQLiteDatabase(validConfig, testLogger)
	assert.NoError(t, err, "Failed to create SQLite database")

	ctx := context.Background()

	// Test ping without connection
	err = db.Ping(ctx)
	assert.Error(t, err, "Expected error when pinging without connection")

	// Test GetConn without connection
	sqliteDB, ok := db.(*SQLiteDatabase)
	require.True(t, ok, "Failed to cast to SQLiteDatabase")
	_, err = sqliteDB.GetConn(ctx)
	assert.Error(t, err, "Expected error when getting connection without connecting")

	// Test connecting to invalid database
	invalidConfig := &config.Config{
		DBType:     constants.SQLite,
		SQLiteFile: "/invalid/path/test.db",
	}
	invalidDB, err := NewSQLiteDatabase(invalidConfig, testLogger)
	assert.NoError(t, err, "Failed to create SQLite database with invalid path")
	err = invalidDB.Connect(ctx)
	assert.Error(t, err, "Expected error when connecting to invalid path")

	// Test GetById with invalid result type
	db, cleanup := setupDatabase(t)
	defer cleanup()

	usersCollection := db.Collection("users")

	var invalidResult int // Not a struct pointer
	err = usersCollection.GetById(ctx, "test-id", invalidResult)
	assert.Error(t, err, "Expected error for invalid result type")

	err = usersCollection.GetById(ctx, "test-id", &invalidResult)
	assert.Error(t, err, "Expected error for invalid result type pointer")
}

// TestInvalidInputs tests handling of invalid inputs
func TestInvalidInputs(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	ctx := context.Background()
	usersCollection := db.Collection("users")

	// Test Create with non-struct data
	_, err := usersCollection.Create(ctx, "not a struct")
	assert.Error(t, err, "Expected error when creating with non-struct data")

	// Test Create with nil
	_, err = usersCollection.Create(ctx, nil)
	assert.Error(t, err, "Expected error when creating with nil data")

	// Test GetById with empty ID
	var result TestUser
	err = usersCollection.GetById(ctx, "", &result)
	assert.Error(t, err, "Expected error when finding with empty ID")

	// Test Count with invalid filter key
	_, err = usersCollection.Count(ctx, map[string]interface{}{"invalid;column": "value"})
	assert.Error(t, err, "Expected error when counting with invalid filter key")
}

// TestSQLiteCollectionIntegration performs an end-to-end test of the SQLiteCollection
func TestSQLiteCollectionIntegration(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	ctx := context.Background()

	// Get collections
	usersCollection := db.Collection("users")
	groupsCollection := db.Collection("reminder_groups")
	remindersCollection := db.Collection("reminders")

	// 1. Create a user
	userID := uuid.New().String()
	user := TestUser{
		ID:        userID,
		Username:  "integrationtest",
		Email:     "integration@example.com",
		Password:  "password123",
		Role:      "user",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	createdUserID, err := usersCollection.Create(ctx, user)
	require.NoError(t, err, "Failed to create user")
	assert.Equal(t, userID, createdUserID, "User ID mismatch")

	// 2. Verify user exists
	var foundUser TestUser
	err = usersCollection.GetById(ctx, userID, &foundUser)
	assert.NoError(t, err, "Failed to find user")
	assert.Equal(t, user.Username, foundUser.Username, "Username mismatch")

	// 3. Create multiple reminder groups
	for i := 0; i < 3; i++ {
		group := TestReminderGroup{
			ID:        uuid.New().String(),
			Name:      "Group " + uuid.New().String(),
			UserID:    userID,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		_, err := groupsCollection.Create(ctx, group)
		require.NoError(t, err, "Failed to create group")
	}

	// 4. Verify group count
	groupCount, err := groupsCollection.Count(ctx, map[string]interface{}{"user_id": userID})
	assert.NoError(t, err, "Failed to count groups")
	assert.Equal(t, int64(3), groupCount, "Expected 3 groups")

	// 5. Create reminders with and without groups
	sqliteDB, ok := db.(*SQLiteDatabase)
	require.True(t, ok, "Failed to cast to SQLiteDatabase")

	// Get a group ID
	var groupID string
	err = sqliteDB.conn.QueryRowContext(ctx,
		"SELECT id FROM reminder_groups WHERE user_id = ? LIMIT 1", userID).Scan(&groupID)
	require.NoError(t, err, "Failed to get group ID")

	// Create reminder with group
	reminderWithGroup := TestReminder{
		ID:              uuid.New().String(),
		Title:           "Reminder With Group",
		Description:     "This reminder has a group",
		IsPinned:        true,
		UserID:          userID,
		ReminderGroupID: groupID,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
	_, err = remindersCollection.Create(ctx, reminderWithGroup)
	require.NoError(t, err, "Failed to create reminder with group")

	// Create reminder without group
	reminderWithoutGroup := TestReminder{
		ID:          uuid.New().String(),
		Title:       "Reminder Without Group",
		Description: "This reminder has no group",
		IsPinned:    false,
		UserID:      userID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	_, err = remindersCollection.Create(ctx, reminderWithoutGroup)
	require.NoError(t, err, "Failed to create reminder without group")

	// 6. Verify reminder counts
	reminderCount, err := remindersCollection.Count(ctx, map[string]interface{}{"user_id": userID})
	assert.NoError(t, err, "Failed to count reminders")
	assert.Equal(t, int64(2), reminderCount, "Expected 2 reminders")

	groupReminderCount, err := remindersCollection.Count(ctx,
		map[string]interface{}{"reminder_group_id": groupID})
	assert.NoError(t, err, "Failed to count reminders in group")
	assert.Equal(t, int64(1), groupReminderCount, "Expected 1 reminder in group")

	// 7. Test the foreign key constraint - deleting a group should set reminder_group_id to NULL
	_, err = sqliteDB.conn.ExecContext(ctx, "DELETE FROM reminder_groups WHERE id = ?", groupID)
	assert.NoError(t, err, "Failed to delete group")

	// Verify the reminder's group ID is set to NULL
	var reminderGroupID sql.NullString
	err = sqliteDB.conn.QueryRowContext(ctx,
		"SELECT reminder_group_id FROM reminders WHERE id = ?",
		reminderWithGroup.ID).Scan(&reminderGroupID)
	assert.NoError(t, err, "Failed to query reminder")
	assert.False(t, reminderGroupID.Valid, "Expected reminder_group_id to be NULL")
}

///

// TestSQLiteCollectionGet tests the Get method of SQLiteCollection
func TestSQLiteCollectionGet(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	ctx := context.Background()
	usersCollection := db.Collection("users")

	// Create multiple test users with different roles
	users := []TestUser{
		{
			ID:        uuid.New().String(),
			Username:  "admin1",
			Email:     "admin1@example.com",
			Password:  "password123",
			Role:      "admin",
			Active:    true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:        uuid.New().String(),
			Username:  "user1",
			Email:     "user1@example.com",
			Password:  "password123",
			Role:      "user",
			Active:    false,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:        uuid.New().String(),
			Username:  "user2",
			Email:     "user2@example.com",
			Password:  "password123",
			Role:      "user",
			Active:    true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:        uuid.New().String(),
			Username:  "user3",
			Email:     "user3@example.com",
			Password:  "password123",
			Role:      "user",
			Active:    true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}

	sqliteDB, ok := db.(*SQLiteDatabase)
	require.True(t, ok, "Failed to cast to SQLiteDatabase")
	_, err := sqliteDB.conn.ExecContext(ctx, "ALTER TABLE users ADD active BOOLEAN")
	assert.NoError(t, err, "Failed to add a new column active")

	// Insert test users
	for _, user := range users {
		_, err := usersCollection.Create(ctx, user)
		require.NoError(t, err, "Failed to create test user")
	}

	// Test finding all users
	var allUsers []TestUser
	err = usersCollection.GetAllByCondition(ctx, map[string]interface{}{}, &allUsers)
	assert.NoError(t, err, "Failed to find all users")
	assert.Equal(t, 4, len(allUsers), "Expected 3 users")

	// Verify all users were retrieved
	usernames := make(map[string]bool)
	for _, user := range allUsers {
		usernames[user.Username] = true
	}

	assert.True(t, usernames["admin1"], "Missing admin1 user")
	assert.True(t, usernames["user1"], "Missing user1 user")
	assert.True(t, usernames["user2"], "Missing user2 user")
	assert.True(t, usernames["user3"], "Missing user3 user")

	// Test finding users by role
	var adminUsers []TestUser
	err = usersCollection.GetAllByCondition(ctx, map[string]interface{}{"role": "admin"}, &adminUsers)

	assert.NoError(t, err, "Failed to find admin users")
	assert.Equal(t, 1, len(adminUsers), "Expected 1 admin user")
	assert.Equal(t, "admin", adminUsers[0].Role, "Expected admin role")
	assert.Equal(t, "admin1", adminUsers[0].Username, "Expected admin1 username")

	var regularUsers []TestUser
	err = usersCollection.GetAllByCondition(ctx, map[string]interface{}{"role": "user"}, &regularUsers)
	assert.NoError(t, err, "Failed to find regular users")
	assert.Equal(t, 3, len(regularUsers), "Expected 2 regular users")

	// Test finding with multiple conditions
	var activeUsers []TestUser
	err = usersCollection.GetAllByCondition(ctx, map[string]interface{}{
		"role":   "user",
		"active": true,
	}, &activeUsers)

	assert.NoError(t, err, "Failed to find active users")
	assert.Equal(t, 2, len(activeUsers), "Expected 1 active regular user")
	assert.Equal(t, "user2", activeUsers[0].Username, "Expected user2")
	assert.Equal(t, "user3", activeUsers[1].Username, "Expected user3")

	// Test finding with invalid filter
	err = usersCollection.GetAllByCondition(ctx, map[string]interface{}{"invalid_column": "value"}, &allUsers)
	assert.Error(t, err, "Expected error with invalid filter")

	// Test finding with nil result
	err = usersCollection.GetAllByCondition(ctx, map[string]interface{}{}, nil)
	assert.Error(t, err, "Expected error with nil result")
}

// TestSQLiteCollectionUpdateById tests the UpdateById method of SQLiteCollection
func TestSQLiteCollectionUpdateById(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	ctx := context.Background()
	usersCollection := db.Collection("users")

	// Create a test user
	userID := uuid.New().String()
	user := TestUser{
		ID:        userID,
		Username:  "updatetest",
		Email:     "update@example.com",
		Password:  "password123",
		Role:      "user",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Insert the user
	_, err := usersCollection.Create(ctx, user)
	require.NoError(t, err, "Failed to create test user")

	// Update the user
	updatedUser := TestUser{
		ID:        userID,
		Username:  "updated_user",
		Email:     "updated@example.com",
		Password:  "newpassword123",
		Role:      "admin",
		UpdatedAt: time.Now().UTC(),
	}

	err = usersCollection.UpdateById(ctx, userID, updatedUser)
	assert.NoError(t, err, "Failed to update user")

	// Verify the update
	var foundUser TestUser
	err = usersCollection.GetById(ctx, userID, &foundUser)
	assert.NoError(t, err, "Failed to find updated user")
	assert.Equal(t, updatedUser.Username, foundUser.Username, "Username not updated")
	assert.Equal(t, updatedUser.Email, foundUser.Email, "Email not updated")
	assert.Equal(t, updatedUser.Role, foundUser.Role, "Role not updated")

	// Test updating non-existent user
	err = usersCollection.UpdateById(ctx, "non-existent-id", updatedUser)
	assert.Error(t, err, "Expected error when updating non-existent user")
	assert.ErrorIs(t, err, ErrNotFound, "Expected not found error")

	// Test updating with invalid data
	err = usersCollection.UpdateById(ctx, userID, "invalid-data")
	assert.Error(t, err, "Expected error when updating with invalid data")
}

// TestSQLiteCollectionDeleteById tests the DeleteById method of SQLiteCollection
func TestSQLiteCollectionDeleteById(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	ctx := context.Background()
	usersCollection := db.Collection("users")
	remindersCollection := db.Collection("reminders")

	// Create a test user
	userID := uuid.New().String()
	user := TestUser{
		ID:        userID,
		Username:  "deletetest",
		Email:     "delete@example.com",
		Password:  "password123",
		Role:      "user",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Insert the user
	_, err := usersCollection.Create(ctx, user)
	require.NoError(t, err, "Failed to create test user")

	// Create a reminder for the user
	reminder := TestReminder{
		ID:          uuid.New().String(),
		Title:       "Test Reminder",
		Description: "This is a test reminder",
		UserID:      userID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	_, err = remindersCollection.Create(ctx, reminder)
	require.NoError(t, err, "Failed to create test reminder")

	// Delete the user
	err = usersCollection.DeleteById(ctx, userID)
	assert.NoError(t, err, "Failed to delete user")

	// Verify user is deleted
	var foundUser TestUser
	err = usersCollection.GetById(ctx, userID, &foundUser)
	assert.Error(t, err, "Expected error when finding deleted user")
	assert.ErrorIs(t, err, ErrNotFound, "Expected not found error")

	// Verify cascade deletion of reminders
	var foundReminder TestReminder
	err = remindersCollection.GetById(ctx, reminder.ID, &foundReminder)
	assert.Error(t, err, "Expected error when Geting reminder after user deletion")
	assert.ErrorIs(t, err, ErrNotFound, "Expected not found error")

	// Test deleting non-existent user
	err = usersCollection.DeleteById(ctx, "non-existent-id")
	assert.Error(t, err, "Expected error when deleting non-existent user")
	assert.ErrorIs(t, err, ErrNotFound, "Expected not found error")

	// Test deleting with empty ID
	err = usersCollection.DeleteById(ctx, "")
	assert.Error(t, err, "Expected error when deleting with empty ID")
}

// TestSQLiteCollectionGetOne tests the FindOne method of SQLiteCollection
func TestSQLiteCollectionGetOne(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	ctx := context.Background()

	// Get users collection
	usersCollection := db.Collection("users")

	// Create test users
	now := time.Now().UTC()
	user1 := TestUser{
		ID:        uuid.New().String(),
		Username:  "user1",
		Email:     "user1@example.com",
		Password:  "password123",
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	user2 := TestUser{
		ID:        uuid.New().String(),
		Username:  "user2",
		Email:     "user2@example.com",
		Password:  "password123",
		Role:      "admin",
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := usersCollection.Create(ctx, user1)
	require.NoError(t, err, "Failed to create user1")
	_, err = usersCollection.Create(ctx, user2)
	require.NoError(t, err, "Failed to create user2")

	// Test GetOne with a valid filter
	var foundUser TestUser
	err = usersCollection.GetOne(ctx, map[string]interface{}{"username": "user1"}, &foundUser)
	assert.NoError(t, err, "Failed to find user1")
	assert.Equal(t, user1.ID, foundUser.ID, "User ID mismatch")
	assert.Equal(t, user1.Username, foundUser.Username, "Username mismatch")
	assert.Equal(t, user1.Email, foundUser.Email, "Email mismatch")

	// Test GetOne with a valid filter
	var foundSecondUser TestUser
	err = usersCollection.GetOne(ctx, map[string]interface{}{"email": "user2@example.com"}, &foundSecondUser)
	assert.NoError(t, err, "Failed to find user2")
	assert.Equal(t, user2.ID, foundSecondUser.ID, "User ID mismatch")
	assert.Equal(t, user2.Username, foundSecondUser.Username, "Username mismatch")
	assert.Equal(t, user2.Email, foundSecondUser.Email, "Email mismatch")

	// Test GetOne with a filter that matches no records
	err = usersCollection.GetOne(ctx, map[string]interface{}{"username": "nonexistent"}, &foundUser)
	assert.Error(t, err, "Expected error for nonexistent user")
	assert.ErrorIs(t, err, ErrNotFound, "Expected not found error")

	// Test GetOne with an invalid result type
	var invalidResult int
	err = usersCollection.GetOne(ctx, map[string]interface{}{"username": "user1"}, invalidResult)
	assert.Error(t, err, "Expected error for invalid result type")

	err = usersCollection.GetOne(ctx, map[string]interface{}{"username": "user1"}, &invalidResult)
	assert.Error(t, err, "Expected error for invalid result type pointer")
}
