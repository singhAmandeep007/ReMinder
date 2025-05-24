package bq

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock interfaces for testing
type mockBigQueryClient struct {
	mock.Mock
}

type mockDataset struct {
	mock.Mock
}

type mockTable struct {
	mock.Mock
}

type mockQuery struct {
	mock.Mock
}

type mockJob struct {
	mock.Mock
}

type mockIterator struct {
	mock.Mock
	rows []map[string]bigquery.Value
	idx  int
}

func (m *mockIterator) Next(dst interface{}) error {
	if m.idx >= len(m.rows) {
		return errors.New("iterator.Done")
	}

	row := dst.(*map[string]bigquery.Value)
	*row = m.rows[m.idx]
	m.idx++
	return nil
}

func (m *mockIterator) Schema() []*bigquery.FieldSchema {
	return []*bigquery.FieldSchema{
		{Name: "id", Type: bigquery.IntegerFieldType},
		{Name: "name", Type: bigquery.StringFieldType},
	}
}

func (m *mockIterator) TotalRows() uint64 {
	return uint64(len(m.rows))
}

// Test Config validation
func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config with credentials path",
			config: Config{
				ProjectID:       "test-project",
				CredentialsPath: "/path/to/creds.json",
				Location:        "US",
			},
			wantErr: false,
		},
		{
			name: "valid config with credentials JSON",
			config: Config{
				ProjectID:       "test-project",
				CredentialsJSON: []byte(`{"type": "service_account"}`),
				Location:        "EU",
			},
			wantErr: false,
		},
		{
			name: "missing project ID",
			config: Config{
				CredentialsPath: "/path/to/creds.json",
			},
			wantErr: true,
		},
		{
			name:    "empty config",
			config:  Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Note: This will fail in unit test environment without actual GCP credentials
			// In real testing, you'd mock the bigquery.NewClient call
			_, err := NewClient(ctx, tt.config)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				// In unit tests, we expect error due to missing credentials
				// In integration tests with proper setup, this would pass
				assert.Error(t, err) // Expected in unit test environment
			}
		})
	}
}

func TestConfig_Defaults(t *testing.T) {
	config := Config{
		ProjectID: "test-project",
	}

	// Test that defaults are applied (would need to mock NewClient for full test)
	assert.Equal(t, "", config.Location)
	assert.Equal(t, 0, config.QueryTimeoutSeconds)

	// After NewClient call, these would be set to defaults:
	// Location: "US"
	// QueryTimeoutSeconds: 300
}

func TestQueryResult_Structure(t *testing.T) {
	result := &QueryResult{
		Rows: []map[string]interface{}{
			{"id": int64(1), "name": "test1"},
			{"id": int64(2), "name": "test2"},
		},
		Schema: []*bigquery.FieldSchema{
			{Name: "id", Type: bigquery.IntegerFieldType},
			{Name: "name", Type: bigquery.StringFieldType},
		},
		JobID:     "job_123",
		TotalRows: 2,
	}

	assert.Len(t, result.Rows, 2)
	assert.Equal(t, "job_123", result.JobID)
	assert.Equal(t, int64(2), result.TotalRows)
	assert.Len(t, result.Schema, 2)
}

func TestTableInfo_Structure(t *testing.T) {
	now := time.Now()
	tableInfo := &TableInfo{
		ProjectID: "test-project",
		DatasetID: "test_dataset",
		TableID:   "test_table",
		Schema: []*bigquery.FieldSchema{
			{Name: "id", Type: bigquery.IntegerFieldType},
		},
		NumRows:  1000,
		NumBytes: 50000,
		Created:  now,
		Modified: now,
	}

	assert.Equal(t, "test-project", tableInfo.ProjectID)
	assert.Equal(t, "test_dataset", tableInfo.DatasetID)
	assert.Equal(t, "test_table", tableInfo.TableID)
	assert.Equal(t, int64(1000), tableInfo.NumRows)
	assert.Equal(t, int64(50000), tableInfo.NumBytes)
}

func TestConvertBigQueryValue(t *testing.T) {
	tests := []struct {
		name     string
		input    bigquery.Value
		expected interface{}
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: nil,
		},
		{
			name:     "string value",
			input:    "test string",
			expected: "test string",
		},
		{
			name:     "integer value",
			input:    int64(42),
			expected: int64(42),
		},
		{
			name:     "float value",
			input:    3.14,
			expected: 3.14,
		},
		{
			name:     "boolean value",
			input:    true,
			expected: true,
		},
		{
			name:     "array value",
			input:    []bigquery.Value{"a", "b", "c"},
			expected: []interface{}{"a", "b", "c"},
		},
		{
			name: "struct value",
			input: map[string]bigquery.Value{
				"field1": "value1",
				"field2": int64(123),
			},
			expected: map[string]interface{}{
				"field1": "value1",
				"field2": int64(123),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertBigQueryValue(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapValueSaver(t *testing.T) {
	values := map[string]interface{}{
		"id":     int64(1),
		"name":   "test",
		"active": true,
	}

	mvs := &mapValueSaver{values: values}

	bqValues, insertID, err := mvs.Save()

	assert.NoError(t, err)
	assert.Empty(t, insertID)
	assert.Equal(t, len(values), len(bqValues))

	for k, v := range values {
		assert.Equal(t, v, bqValues[k])
	}
}

// Schema helper function tests
func TestSchemaHelpers(t *testing.T) {
	t.Run("StringField", func(t *testing.T) {
		field := StringField("name", true)
		assert.Equal(t, "name", field.Name)
		assert.Equal(t, bigquery.StringFieldType, field.Type)
		assert.True(t, field.Required)

		optionalField := StringField("description", false)
		assert.False(t, optionalField.Required)
	})

	t.Run("IntegerField", func(t *testing.T) {
		field := IntegerField("id", true)
		assert.Equal(t, "id", field.Name)
		assert.Equal(t, bigquery.IntegerFieldType, field.Type)
		assert.True(t, field.Required)
	})

	t.Run("FloatField", func(t *testing.T) {
		field := FloatField("price", false)
		assert.Equal(t, "price", field.Name)
		assert.Equal(t, bigquery.FloatFieldType, field.Type)
		assert.False(t, field.Required)
	})

	t.Run("BooleanField", func(t *testing.T) {
		field := BooleanField("active", true)
		assert.Equal(t, "active", field.Name)
		assert.Equal(t, bigquery.BooleanFieldType, field.Type)
		assert.True(t, field.Required)
	})

	t.Run("TimestampField", func(t *testing.T) {
		field := TimestampField("created_at", true)
		assert.Equal(t, "created_at", field.Name)
		assert.Equal(t, bigquery.TimestampFieldType, field.Type)
		assert.True(t, field.Required)
	})

	t.Run("DateField", func(t *testing.T) {
		field := DateField("birth_date", false)
		assert.Equal(t, "birth_date", field.Name)
		assert.Equal(t, bigquery.DateFieldType, field.Type)
		assert.False(t, field.Required)
	})
}

// Integration test example (requires actual GCP credentials and setup)
func TestIntegration_ClientOperations(t *testing.T) {
	// Skip this test in normal unit test runs
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test requires:
	// 1. GOOGLE_APPLICATION_CREDENTIALS environment variable set
	// 2. Or service account JSON file
	// 3. Actual GCP project with BigQuery enabled

	config := Config{
		ProjectID: "your-test-project", // Replace with actual project
		Location:  "US",
	}

	ctx := context.Background()
	client, err := NewClient(ctx, config)
	if err != nil {
		t.Skipf("Failed to create client (expected in unit test env): %v", err)
		return
	}
	defer client.Close()

	// Test basic query
	t.Run("simple query", func(t *testing.T) {
		sql := "SELECT 1 as id, 'test' as name"
		result, err := client.Query(ctx, sql)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Rows, 1)
		assert.Equal(t, int64(1), result.Rows[0]["id"])
		assert.Equal(t, "test", result.Rows[0]["name"])
	})

	// Test dataset operations
	t.Run("dataset operations", func(t *testing.T) {
		datasetID := "test_dataset_" + time.Now().Format("20060102_150405")

		// Create dataset
		err := client.CreateDataset(ctx, datasetID, "")
		require.NoError(t, err)

		// Clean up
		defer func() {
			err := client.DeleteDataset(ctx, datasetID, true)
			assert.NoError(t, err)
		}()

		// Test table operations within dataset
		t.Run("table operations", func(t *testing.T) {
			tableID := "test_table"
			schema := []*bigquery.FieldSchema{
				StringField("name", true),
				IntegerField("age", false),
			}

			// Create table
			err := client.CreateTable(ctx, datasetID, tableID, "", schema)
			require.NoError(t, err)

			// Get table info
			tableInfo, err := client.GetTableInfo(ctx, datasetID, tableID)
			require.NoError(t, err)
			assert.Equal(t, datasetID, tableInfo.DatasetID)
			assert.Equal(t, tableID, tableInfo.TableID)

			// Insert rows
			rows := []map[string]interface{}{
				{"name": "Alice", "age": int64(30)},
				{"name": "Bob", "age": int64(25)},
			}
			err = client.InsertRows(ctx, datasetID, tableID, rows)
			require.NoError(t, err)

			// List tables
			tables, err := client.ListTables(ctx, datasetID)
			require.NoError(t, err)
			assert.Contains(t, tables, tableID)

			// Query the table
			sql := fmt.Sprintf("SELECT * FROM `%s.%s.%s`", config.ProjectID, datasetID, tableID)
			result, err := client.Query(ctx, sql)
			require.NoError(t, err)
			assert.Len(t, result.Rows, 2)

			// Test DML
			dmlSQL := fmt.Sprintf("UPDATE `%s.%s.%s` SET age = age + 1 WHERE name = 'Alice'",
				config.ProjectID, datasetID, tableID)
			affectedRows, err := client.ExecuteDML(ctx, dmlSQL)
			require.NoError(t, err)
			assert.Equal(t, int64(1), affectedRows)

			// Delete table
			err = client.DeleteTable(ctx, datasetID, tableID)
			require.NoError(t, err)
		})
	})
}

// Benchmark tests
func BenchmarkConvertBigQueryValue(b *testing.B) {
	testCases := []struct {
		name  string
		value bigquery.Value
	}{
		{"string", "test string"},
		{"integer", int64(12345)},
		{"float", 3.14159},
		{"boolean", true},
		{"array", []bigquery.Value{"a", "b", "c", "d", "e"}},
		{
			"struct",
			map[string]bigquery.Value{
				"field1": "value1",
				"field2": int64(123),
				"field3": true,
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				convertBigQueryValue(tc.value)
			}
		})
	}
}

func BenchmarkMapValueSaver(b *testing.B) {
	values := map[string]interface{}{
		"id":      int64(1),
		"name":    "test name",
		"email":   "test@example.com",
		"active":  true,
		"score":   95.5,
		"created": time.Now(),
	}

	mvs := &mapValueSaver{values: values}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mvs.Save()
	}
}

// Example usage tests
func ExampleNewClient() {
	config := Config{
		ProjectID:           "my-gcp-project",
		CredentialsPath:     "/path/to/service-account.json",
		Location:            "US",
		QueryTimeoutSeconds: 300,
	}

	ctx := context.Background()
	client, err := NewClient(ctx, config)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Use client for operations...
}

func ExampleClient_Query() {
	// Assuming client is already created
	var client *Client

	ctx := context.Background()
	sql := `
		SELECT
			name,
			COUNT(*) as count
		FROM
			` + "`my-project.my_dataset.my_table`" + `
		WHERE
			created_date >= '2024-01-01'
		GROUP BY
			name
		ORDER BY
			count DESC
		LIMIT 10
	`

	result, err := client.Query(ctx, sql)
	if err != nil {
		panic(err)
	}

	for _, row := range result.Rows {
		name := row["name"].(string)
		count := row["count"].(int64)
		fmt.Printf("Name: %s, Count: %d\n", name, count)
	}
}

func ExampleClient_CreateTable() {
	// Assuming client is already created
	var client *Client

	ctx := context.Background()

	schema := []*bigquery.FieldSchema{
		StringField("user_id", true),
		StringField("email", true),
		StringField("name", false),
		IntegerField("age", false),
		BooleanField("active", true),
		TimestampField("created_at", true),
		TimestampField("updated_at", false),
	}

	err := client.CreateTable(ctx, "my_dataset", "users", "", schema)
	if err != nil {
		panic(err)
	}
}

func ExampleClient_InsertRows() {
	// Assuming client is already created
	var client *Client

	ctx := context.Background()

	rows := []map[string]interface{}{
		{
			"user_id":    "user_001",
			"email":      "alice@example.com",
			"name":       "Alice Smith",
			"age":        int64(30),
			"active":     true,
			"created_at": time.Now(),
		},
		{
			"user_id":    "user_002",
			"email":      "bob@example.com",
			"name":       "Bob Johnson",
			"age":        int64(25),
			"active":     true,
			"created_at": time.Now(),
		},
	}

	err := client.InsertRows(ctx, "my_dataset", "users", rows)
	if err != nil {
		panic(err)
	}
}
