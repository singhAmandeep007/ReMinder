# BigQuery Go Package

A extendable Go package for performing BigQuery operations with a simple and intuitive API.

## Features

- **Thread-safe**: All operations are protected by read-write mutexes
- **Extendable**: Easy to extend with custom functionality
- **Comprehensive**: Supports all basic BigQuery operations
- **Error handling**: Proper error wrapping and context
- **Configurable**: Flexible configuration options
- **Production-ready**: Includes timeouts, proper resource management

## Installation

```bash
go get github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/bq
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/bq"
)

func main() {
    ctx := context.Background()

    // Configure the client
    config := bigquery.Config{
        ProjectID:           "your-project-id",
        CredentialsPath:     "/path/to/service-account.json",
        Location:            "US",
        QueryTimeoutSeconds: 300,
    }

    // Create client
    client, err := bigquery.NewClient(ctx, config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Execute a query
    result, err := client.Query(ctx, "SELECT 1 as test")
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Query returned %d rows", len(result.Rows))
}
```

## Configuration Options

- `ProjectID`: Google Cloud Project ID (required)
- `CredentialsPath`: Path to service account JSON file
- `CredentialsJSON`: Service account JSON as byte slice
- `Location`: Default location for datasets and jobs (default: "US")
- `QueryTimeoutSeconds`: Default timeout for queries (default: 300)

### Dataset Operations

- `CreateDataset(ctx, datasetID)` - Create a new dataset
- `DeleteDataset(ctx, datasetID, deleteContents)` - Delete a dataset
- `ListTables(ctx, datasetID)` - List tables in a dataset
- `ListTableIDs(ctx, datasetID)` - List table IDs in a dataset

### Table Operations

- `CreateTable(ctx, datasetID, tableID, schema)` - Create a table
- `DeleteTable(ctx, datasetID, tableID)` - Delete a table
- `GetTableInfo(ctx, datasetID, tableID)` - Get table metadata
- `InsertRows(ctx, datasetID, tableID, rows)` - Insert data

### Query Operations

- `Query(ctx, sql)` - Execute a SELECT query
- `ExecuteDML(ctx, sql)` - Execute DML statements (INSERT, UPDATE, DELETE)

## Schema Helper Functions

The package includes helper functions for creating field schemas:

```go
schema := []*bigquery.FieldSchema{
    bigquery.StringField("name", true),      // required string field
    bigquery.IntegerField("age", false),     // optional integer field
    bigquery.FloatField("score", false),     // optional float field
    bigquery.BooleanField("active", true),   // required boolean field
    bigquery.TimestampField("created", true), // required timestamp field
    bigquery.DateField("birthday", false),   // optional date field
}
```

```go
schema := bq.CreateSchemaFromFields(
  &bq.FieldSchema{
    Name:     "name",
    Type:     bq.FieldTypeString,
    Required: true,
  },
  &bq.FieldSchema{
    Name:     "age",
    Type:     bq.FieldTypeInteger,
    Required: false,
  },
)
```

## Thread Safety

The client is fully thread-safe and can be used concurrently from multiple goroutines. All operations are protected by appropriate locking mechanisms.

## Error Handling

All methods return descriptive errors with context. Use Go's error wrapping features to handle errors appropriately:

```go
result, err := client.Query(ctx, sql)
if err != nil {
    log.Printf("Query failed: %v", err)
    return
}
```

## Best Practices

1. **Reuse the client**: Create one client instance and reuse it throughout your application
2. **Use contexts**: Always pass appropriate contexts for cancellation and timeouts
3. **Handle errors**: Check and handle all returned errors
4. **Close resources**: Always defer `client.Close()` after creating a client
5. **Use connection pooling**: The underlying BigQuery client handles connection pooling automatically

## Testing

### Run unit tests only (fast)
`go test -short ./...`

### Run all tests including integration tests
`go test ./...`

### Run benchmarks
`go test -bench=. ./...`

### Run with coverage
`go test -cover ./...`


## Examples

```go
  bqClient, err := bq.NewClient(ctx, config)
	if err != nil {
		// Handle error
	}
	defer bqClient.Close()

  err = bqClient.CreateDataset(ctx, "test_dataset", "Test dataset for application")
	if err != nil {
		// Handle error
	} else {
		// Handle success
	}

  err = bqClient.CreateTable(ctx, "test_dataset", "test_table", "Test table for application", bq.CreateSchemaFromFields(
    &bq.FieldSchema{Name: "name", Type: bq.StringFieldType},
    &bq.FieldSchema{Name: "age", Type: bq.IntegerFieldType},
  ))

  tableInfo, err := bqClient.GetTableInfo(ctx, "test_dataset", "test_table")
	if err != nil {
		// Handle error
	} else {
		// Handle success
	}

  err = bqClient.InsertRows(ctx, "test_dataset", "test_table", []map[string]interface{}{
		{"name": "John Doe", "age": 30},
		{"name": "Jane Smith", "age": 25}})

  sqlQuery := `
		SELECT name, age
		FROM ` + "`projectID.test_dataset.test_table`" + `
		WHERE age > 25
		ORDER BY age DESC
	`

	result, err := bqClient.Query(ctx, sqlQuery)
	if err != nil {
		// Handle error
	} else {
    // Handle success
  }

  updateSQL := `
	  ALTER TABLE ` + "`projectID.test_dataset.test_table`" + `
	  ADD COLUMN address STRING
	`

	affectedRows, err := bqClient.ExecuteDML(ctx, updateSQL)
	if err != nil {
		// Handle error
	} else {
		// Handle success
  }

  err = bqClient.DeleteTable(ctx, "test_dataset", "test_table")

  err = bqClient.DeleteDataset(ctx, "test_dataset", true)
```
