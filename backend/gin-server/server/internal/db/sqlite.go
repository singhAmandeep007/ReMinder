package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/config"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// SQLiteDatabase implements the Database interface for SQLite
type SQLiteDatabase struct {
	conn   *sql.DB
	config *config.Config
	logger *logger.Logger
}

// SQLiteCollection implements the Collection interface for SQLite tables
type SQLiteCollection struct {
	db         *SQLiteDatabase
	tableName  string
	primaryKey string
}

// SQLiteTransaction implements the Transaction interface for SQLite
type SQLiteTransaction struct {
	tx *sql.Tx
}

// NewSQLiteDatabase creates a new SQLite database instance
func NewSQLiteDatabase(config *config.Config, logger *logger.Logger) (Database, error) {
	if config.SQLiteFile == "" {
		return nil, errors.New("database file path is required for SQLite")
	}

	return &SQLiteDatabase{
		config: config,
		logger: logger,
	}, nil
}

// Collection returns a collection/table handler for the given name
func (s *SQLiteDatabase) Collection(name string) Collection {
	return &SQLiteCollection{
		db:         s,
		tableName:  name,
		primaryKey: "id",
	}
}

// Connect establishes a connection to the SQLite database
func (s *SQLiteDatabase) Connect(ctx context.Context) error {
	s.logger.Infof("Connecting to SQLite database: %s", s.config.SQLiteFile)

	conn, err := sql.Open("sqlite3", s.config.SQLiteFile)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Set connection pool settings
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(time.Hour)

	// Enable foreign key constraints - add this line
	// SQLite Foreign Key Constraints are Disabled by Default: You must enable them with PRAGMA foreign_keys = ON
	if _, err := conn.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		conn.Close()
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Test the connection
	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		return fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	s.conn = conn
	s.logger.Infof("Successfully connected to SQLite database: %s", s.config.SQLiteFile)
	return nil
}

// Close closes the database connection
func (s *SQLiteDatabase) Close(ctx context.Context) error {
	if s.conn != nil {
		s.logger.Infof("Closing SQLite database connection")
		return s.conn.Close()
	}
	return nil
}

// Ping checks if the database is accessible
func (s *SQLiteDatabase) Ping(ctx context.Context) error {
	if s.conn == nil {
		return errors.New("database connection not established")
	}
	return s.conn.PingContext(ctx)
}

// Transaction defines the interface for transaction operations
type Transaction interface {
	// Commit commits the transaction
	Commit() error

	// Rollback rolls back the transaction
	Rollback() error

	// GetTx returns the underlying transaction object
	GetTx() interface{}
}

// BeginTransaction starts a new transaction
func (s *SQLiteDatabase) BeginTransaction(ctx context.Context) (Transaction, error) {
	if s.conn == nil {
		return nil, errors.New("database connection not established")
	}

	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return &SQLiteTransaction{tx: tx}, nil
}

// GetConn returns the database connection
func (s *SQLiteDatabase) GetConn(ctx context.Context) (*sql.DB, error) {
	if s.conn == nil {
		return nil, errors.New("database connection not established")
	}
	return s.conn, nil
}

// Commit commits the transaction
func (t *SQLiteTransaction) Commit() error {
	if t.tx == nil {
		return errors.New("transaction not started")
	}
	return t.tx.Commit()
}

// Rollback rolls back the transaction
func (t *SQLiteTransaction) Rollback() error {
	if t.tx == nil {
		return errors.New("transaction not started")
	}
	return t.tx.Rollback()
}

// GetTx returns the underlying transaction
func (t *SQLiteTransaction) GetTx() interface{} {
	return t.tx
}

// Migrate runs database migrations
func (s *SQLiteDatabase) Migrate(ctx context.Context) error {
	s.logger.Infof("Running SQLite migrations")

	// Create basic tables if they don't exist
	// Note the order here: users first, then groups, then reminders (respects dependencies)
	// Create parent tables before child tables (users → reminder_groups → reminders)
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS reminder_groups (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			user_id TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS reminders (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			is_pinned BOOLEAN DEFAULT FALSE,
			user_id TEXT NOT NULL,
			reminder_group_id TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      --- if a row in the users table (the parent table) is deleted, all corresponding rows in the current table (the table with the foreign key) that have a matching user_id will also be automatically deleted.
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
      FOREIGN KEY (reminder_group_id) REFERENCES reminder_groups(id) ON DELETE SET NULL

		)`,
	}

	for _, query := range queries {
		_, err := s.conn.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	s.logger.Infof("SQLite migrations completed successfully")
	return nil
}

// Seed populates the database with initial data
func (s *SQLiteDatabase) Seed(ctx context.Context) error {
	s.logger.Infof("Seeding SQLite database")

	// Check if users table is empty
	var count int
	err := s.conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check users table: %w", err)
	}

	if count == 0 {
		// Insert admin user
		_, err := s.conn.ExecContext(ctx,
			"INSERT INTO users (id, username, email, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			"admin-uuid", "admin", "admin@example.com",
			"$2a$10$zgbBOT.6IbXjZEFCJdCgeubIm4LQfy9jAEhTjkxPLAfCzer9SZape", // password: admin123
			"admin",
			time.Now().UTC(), time.Now().UTC())
		if err != nil {
			return fmt.Errorf("failed to seed admin user: %w", err)
		}

		s.logger.Infof("Admin user seeded successfully")
	}

	return nil
}

// Create inserts a new document/record into the collection/table
func (c *SQLiteCollection) Create(ctx context.Context, data interface{}) (string, error) {
	// Get database connection
	conn, err := c.db.GetConn(ctx)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Extract field names, values, and ID
	columns, placeholders, values, id, err := extractFieldsForInsert(data, c.primaryKey)
	if err != nil {
		return "", err
	}

	// Generate SQL query
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		c.tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	// Execute query
	_, err = conn.ExecContext(ctx, query, values...)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return "", fmt.Errorf("%w: %v", ErrDuplicate, err)
		}
		return "", fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return id, nil
}

// GetById retrieves a document/record by ID
func (c *SQLiteCollection) GetById(ctx context.Context, id string, result interface{}) error {
	// Get database connection
	conn, err := c.db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Validate result is a pointer to struct
	if err := validateResultType(result); err != nil {
		return err
	}

	// Prepare the query
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", c.tableName, c.primaryKey)

	// Execute the query
	rows, err := conn.QueryContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer rows.Close()

	// Ensure at least one row is found
	if !rows.Next() {
		return fmt.Errorf("%w: id %s", ErrNotFound, id)
	}

	// Map results to the struct
	if err := mapRowToStruct(rows, result); err != nil {
		return err
	}

	return nil
}

// GetAllByCondition fetches all records from the collection based on filter criteria
func (c *SQLiteCollection) GetAllByCondition(ctx context.Context, filter map[string]interface{}, results interface{}) error {
	// Get database connection
	conn, err := c.db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Validate results is a pointer to slice of structs
	resultsValue := reflect.ValueOf(results)
	if resultsValue.Kind() != reflect.Ptr || resultsValue.IsNil() {
		return fmt.Errorf("%w: results must be a non-nil pointer", ErrInvalidInput)
	}

	sliceValue := resultsValue.Elem()
	if sliceValue.Kind() != reflect.Slice {
		return fmt.Errorf("%w: results must be a pointer to a slice", ErrInvalidInput)
	}

	// Get the element type of the slice
	elemType := sliceValue.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		return fmt.Errorf("%w: slice elements must be structs", ErrInvalidInput)
	}

	// Build query
	query := fmt.Sprintf("SELECT * FROM %s", c.tableName)
	whereClause, values := buildWhereClause(filter)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	// Execute query
	rows, err := conn.QueryContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer rows.Close()

	// Clear the slice before populating
	sliceValue.Set(reflect.MakeSlice(sliceValue.Type(), 0, 0))

	// Process results
	for rows.Next() {
		// Create a new instance of the struct
		newElemPtr := reflect.New(elemType)

		// Map row to struct
		if err := mapRowToStruct(rows, newElemPtr.Interface()); err != nil {
			return err
		}

		// Append to results slice
		sliceValue.Set(reflect.Append(sliceValue, newElemPtr.Elem()))
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return nil
}

// GetOne fetches a single record from the collection based on filter criteria
func (c *SQLiteCollection) GetOne(ctx context.Context, filter map[string]interface{}, result interface{}) error {
	// Get database connection
	conn, err := c.db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Validate result is a pointer to a struct
	if err := validateResultType(result); err != nil {
		return err
	}

	// Build query
	query := fmt.Sprintf("SELECT * FROM %s", c.tableName)
	whereClause, values := buildWhereClause(filter)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}
	query += " LIMIT 1"

	// Execute query
	rows, err := conn.QueryContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}
	defer rows.Close()

	// Check if we have a result
	if !rows.Next() {
		return ErrNotFound
	}

	// Map row to struct
	if err := mapRowToStruct(rows, result); err != nil {
		return err
	}

	return rows.Err()
}

// Update updates a document/record by ID
func (c *SQLiteCollection) UpdateById(ctx context.Context, id string, data interface{}) error {
	// Get database connection
	conn, err := c.db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Extract fields to update
	updateFields, values, err := extractFieldsForUpdate(data)
	if err != nil {
		return err
	}

	// Add ID to values for WHERE clause
	values = append(values, id)

	// Build query
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = ?",
		c.tableName,
		strings.Join(updateFields, ", "),
		c.primaryKey,
	)

	// Execute query
	result, err := conn.ExecContext(ctx, query, values...)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("%w: %v", ErrDuplicate, err)
		}
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: id %s", ErrNotFound, id)
	}

	return nil
}

// Delete removes a document/record by ID
func (c *SQLiteCollection) DeleteById(ctx context.Context, id string) error {
	// Get database connection
	conn, err := c.db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Build query
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", c.tableName, c.primaryKey)

	// Execute query
	result, err := conn.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: id %s", ErrNotFound, id)
	}

	return nil
}

// Count returns the number of documents/records that match the filter
func (c *SQLiteCollection) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	// Get database connection
	conn, err := c.db.GetConn(ctx)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Build query
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", c.tableName)
	whereClause, values := buildWhereClause(filter)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	// Execute query
	var count int64
	err = conn.QueryRowContext(ctx, query, values...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return count, nil
}

// Helper function to check if a value is zero/empty
func isZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// Cache for field mappings to avoid repeated camelCase to snake_case conversions
var fieldNameCache sync.Map

// Helper function to convert camelCase to snake_case with caching
func camelToSnake(s string) string {
	// Check cache first
	if cached, ok := fieldNameCache.Load(s); ok {
		return cached.(string)
	}

	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	snakeCase := result.String()
	// Store in cache
	fieldNameCache.Store(s, snakeCase)
	return snakeCase
}

// Extract field names and values for INSERT operation
func extractFieldsForInsert(data interface{}, primaryKey string) ([]string, []string, []interface{}, string, error) {
	// Reflect on the data to extract field names and values
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, nil, nil, "", fmt.Errorf("%w: data must be a struct or pointer to struct", ErrInvalidInput)
	}

	t := v.Type()
	fieldCount := v.NumField()

	// Prepare column names and placeholders for values
	columns := make([]string, 0, fieldCount)
	placeholders := make([]string, 0, fieldCount)
	values := make([]interface{}, 0, fieldCount)
	var id string

	// Extract field names and values
	for i := 0; i < fieldCount; i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip zero/empty values
		if isZeroOfUnderlyingType(fieldValue.Interface()) {
			continue
		}

		// Use field tag if available, otherwise use field name converted to snake_case
		columnName := field.Tag.Get("db")
		if columnName == "" {
			columnName = camelToSnake(field.Name)
		}

		// Get the ID field for returning
		if columnName == primaryKey {
			if str, ok := fieldValue.Interface().(string); ok {
				id = str
			}
		}

		columns = append(columns, columnName)
		placeholders = append(placeholders, "?")
		values = append(values, fieldValue.Interface())
	}

	return columns, placeholders, values, id, nil
}

// Extract field names and values for UPDATE operation
func extractFieldsForUpdate(data interface{}) ([]string, []interface{}, error) {
	// Reflect on the data to extract field names and values
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("%w: data must be a struct or pointer to struct", ErrInvalidInput)
	}

	t := v.Type()
	fieldCount := v.NumField()

	// Prepare update fields and values
	updateFields := make([]string, 0, fieldCount)
	values := make([]interface{}, 0, fieldCount)

	// Extract field names and values
	for i := 0; i < fieldCount; i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip zero/empty values
		if isZeroOfUnderlyingType(fieldValue.Interface()) {
			continue
		}

		// Skip primary key (id) field - we use it in the WHERE clause
		columnName := field.Tag.Get("db")
		if columnName == "" {
			columnName = camelToSnake(field.Name)
		}

		if columnName == "id" {
			continue
		}

		updateFields = append(updateFields, fmt.Sprintf("%s = ?", columnName))
		values = append(values, fieldValue.Interface())
	}

	// Add updated_at if exists
	if _, found := t.FieldByName("Updated_at"); found {
		updateFields = append(updateFields, "updated_at = ?")
		values = append(values, time.Now().UTC())
	} else if _, found := t.FieldByName("UpdatedAt"); found {
		updateFields = append(updateFields, "updated_at = ?")
		values = append(values, time.Now().UTC())
	}

	return updateFields, values, nil
}

// Build WHERE clause and values for a filter
func buildWhereClause(filter map[string]interface{}) (string, []interface{}) {
	if len(filter) == 0 {
		return "", nil
	}

	var whereClauses []string
	var values []interface{}

	for k, v := range filter {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", k))
		values = append(values, v)
	}

	return strings.Join(whereClauses, " AND "), values
}

// Validate that result is a pointer to a struct
func validateResultType(result interface{}) error {
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Ptr || resultValue.IsNil() {
		return fmt.Errorf("%w: result must be a non-nil pointer", ErrInvalidInput)
	}

	resultElem := resultValue.Elem()
	if resultElem.Kind() != reflect.Struct {
		return fmt.Errorf("%w: result must be a pointer to a struct", ErrInvalidInput)
	}

	return nil
}

// Map a database row to a struct
func mapRowToStruct(rows *sql.Rows, result interface{}) error {
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Create a slice of interface{} to hold the values
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// Scan the result into the valuePtrs
	if err := rows.Scan(valuePtrs...); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Create a map of column name to value
	columnMap := make(map[string]interface{})
	for i, column := range columns {
		columnMap[column] = values[i]
	}

	// Map values to the result struct
	resultValue := reflect.ValueOf(result).Elem()
	resultType := resultValue.Type()

	for i := 0; i < resultValue.NumField(); i++ {
		field := resultType.Field(i)
		fieldValue := resultValue.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		// Use field tag if available, otherwise use field name converted to snake_case
		columnName := field.Tag.Get("db")
		if columnName == "" {
			columnName = camelToSnake(field.Name)
		}

		if val, ok := columnMap[columnName]; ok {
			if val == nil {
				continue
			}

			setFieldValue(fieldValue, val)
		}
	}

	return nil
}

// Set a field value based on its type
func setFieldValue(fieldValue reflect.Value, val interface{}) {
	switch fieldValue.Kind() {
	case reflect.String:
		if str, ok := val.(string); ok {
			fieldValue.SetString(str)
		} else if bytes, ok := val.([]byte); ok {
			fieldValue.SetString(string(bytes))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if num, ok := val.(int64); ok {
			fieldValue.SetInt(num)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if num, ok := val.(int64); ok {
			fieldValue.SetUint(uint64(num))
		}
	case reflect.Float32, reflect.Float64:
		if num, ok := val.(float64); ok {
			fieldValue.SetFloat(num)
		}
	case reflect.Bool:
		if b, ok := val.(bool); ok {
			fieldValue.SetBool(b)
		} else if b, ok := val.(int64); ok {
			fieldValue.SetBool(b != 0)
		}
	case reflect.Struct:
		// Handle time.Time specially
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			if timeStr, ok := val.(string); ok {
				if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
					fieldValue.Set(reflect.ValueOf(t))
				}
			} else if timeBytes, ok := val.([]byte); ok {
				if t, err := time.Parse(time.RFC3339, string(timeBytes)); err == nil {
					fieldValue.Set(reflect.ValueOf(t))
				}
			}
		}
	}
}
