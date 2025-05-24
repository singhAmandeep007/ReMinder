package bq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type (
	// FieldSchema is an alias for bigquery.FieldSchema
	FieldSchema = bigquery.FieldSchema
	// Schema is an alias for bigquery.Schema
	Schema = bigquery.Schema
	// FieldType is an alias for bigquery.FieldType
	FieldType = bigquery.FieldType
)

const (
	StringFieldType     = bigquery.StringFieldType
	IntegerFieldType    = bigquery.IntegerFieldType
	FloatFieldType      = bigquery.FloatFieldType
	BooleanFieldType    = bigquery.BooleanFieldType
	TimestampFieldType  = bigquery.TimestampFieldType
	DateFieldType       = bigquery.DateFieldType
	TimeFieldType       = bigquery.TimeFieldType
	DateTimeFieldType   = bigquery.DateTimeFieldType
	RecordFieldType     = bigquery.RecordFieldType
	GeographyFieldType  = bigquery.GeographyFieldType
	NumericFieldType    = bigquery.NumericFieldType
	BigNumericFieldType = bigquery.BigNumericFieldType
	JSONFieldType       = bigquery.JSONFieldType
)

// Config holds the configuration for BigQuery client
type Config struct {
	ProjectID           string
	CredentialsPath     string
	CredentialsJSON     []byte
	Location            string // Default location for datasets/jobs
	QueryTimeoutSeconds int    // Default query timeout
}

// Client represents a thread-safe BigQuery client
type Client struct {
	bqClient *bigquery.Client
	config   Config
	mu       sync.RWMutex
}

// QueryResult represents the result of a BigQuery query
type QueryResult struct {
	Rows      []map[string]interface{}
	Schema    []*FieldSchema
	JobID     string
	TotalRows int64
}

// TableInfo contains metadata about a BigQuery table
type TableInfo struct {
	ProjectID   string
	DatasetID   string
	TableID     string
	Schema      []*FieldSchema
	Description string // Optional description of the table
	NumRows     int64
	NumBytes    int64
	Created     time.Time
	Modified    time.Time
}

// NewClient creates a new BigQuery client instance
func NewClient(ctx context.Context, config Config) (*Client, error) {
	if config.ProjectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}

	// Set default values
	if config.Location == "" {
		config.Location = "US"
	}
	if config.QueryTimeoutSeconds <= 0 {
		config.QueryTimeoutSeconds = 300 // 5 minutes default
	}

	var opts []option.ClientOption

	if config.CredentialsPath != "" {
		opts = append(opts, option.WithCredentialsFile(config.CredentialsPath))
	} else if len(config.CredentialsJSON) > 0 {
		opts = append(opts, option.WithCredentialsJSON(config.CredentialsJSON))
	}

	bqClient, err := bigquery.NewClient(ctx, config.ProjectID, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create BigQuery client: %w", err)
	}

	return &Client{
		bqClient: bqClient,
		config:   config,
	}, nil
}

// Close closes the BigQuery client connection
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.bqClient != nil {
		return c.bqClient.Close()
	}
	return nil
}

// CreateDataset creates a new dataset
func (c *Client) CreateDataset(ctx context.Context, datasetID string, description string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dataset := c.bqClient.Dataset(datasetID)
	// https://pkg.go.dev/cloud.google.com/go/bigquery#DatasetMetadata
	meta := &bigquery.DatasetMetadata{
		Location:    c.config.Location,
		Description: description,
	}

	if err := dataset.Create(ctx, meta); err != nil {
		return fmt.Errorf("failed to create dataset %s: %w", datasetID, err)
	}

	return nil
}

// DeleteDataset deletes a dataset and optionally its contents
func (c *Client) DeleteDataset(ctx context.Context, datasetID string, deleteContents bool) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dataset := c.bqClient.Dataset(datasetID)

	if err := dataset.DeleteWithContents(ctx); err != nil && deleteContents {
		return fmt.Errorf("failed to delete dataset %s with contents: %w", datasetID, err)
	} else if err := dataset.Delete(ctx); err != nil && !deleteContents {
		return fmt.Errorf("failed to delete dataset %s: %w", datasetID, err)
	}

	return nil
}

// CreateSchemaFromFields creates a BigQuery schema from a list of FieldSchema
func CreateSchemaFromFields(fields ...*FieldSchema) Schema {
	return Schema(fields)
}

// CreateTable creates a new table with the specified schema
func (c *Client) CreateTable(ctx context.Context, datasetID, tableID, description string, schema []*FieldSchema) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dataset := c.bqClient.Dataset(datasetID)
	table := dataset.Table(tableID)

	meta := &bigquery.TableMetadata{
		// https://pkg.go.dev/cloud.google.com/go/bigquery#Schema
		Schema:      schema,
		Location:    c.config.Location,
		Description: description,
	}

	if err := table.Create(ctx, meta); err != nil {
		return fmt.Errorf("failed to create table %s.%s: %w", datasetID, tableID, err)
	}

	return nil
}

// DeleteTable deletes a table
func (c *Client) DeleteTable(ctx context.Context, datasetID, tableID string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dataset := c.bqClient.Dataset(datasetID)
	table := dataset.Table(tableID)

	if err := table.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete table %s.%s: %w", datasetID, tableID, err)
	}

	return nil
}

// Query executes a SQL query and returns the results
func (c *Client) Query(ctx context.Context, sql string) (*QueryResult, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query := c.bqClient.Query(sql)
	query.Location = c.config.Location

	// Set query timeout
	queryCtx, cancel := context.WithTimeout(ctx, time.Duration(c.config.QueryTimeoutSeconds)*time.Second)
	defer cancel()

	job, err := query.Run(queryCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %w", err)
	}

	status, err := job.Wait(queryCtx)
	if err != nil {
		return nil, fmt.Errorf("query job failed: %w", err)
	}

	if status.Err() != nil {
		return nil, fmt.Errorf("query execution error: %w", status.Err())
	}

	it, err := job.Read(queryCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to read query results: %w", err)
	}

	var rows []map[string]interface{}
	for {
		var row map[string]bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate over results: %w", err)
		}

		// Convert bigquery.Value to interface{}
		convertedRow := make(map[string]interface{})
		for k, v := range row {
			convertedRow[k] = convertBigQueryValue(v)
		}
		rows = append(rows, convertedRow)
	}

	return &QueryResult{
		Rows:      rows,
		Schema:    it.Schema,
		JobID:     job.ID(),
		TotalRows: int64(it.TotalRows),
	}, nil
}

// GetTableInfo retrieves metadata about a table
func (c *Client) GetTableInfo(ctx context.Context, datasetID, tableID string) (*TableInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dataset := c.bqClient.Dataset(datasetID)
	table := dataset.Table(tableID)

	meta, err := table.Metadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get table metadata for %s.%s: %w", datasetID, tableID, err)
	}

	return &TableInfo{
		ProjectID:   c.config.ProjectID,
		DatasetID:   datasetID,
		TableID:     tableID,
		Schema:      meta.Schema,
		NumRows:     int64(meta.NumRows),
		NumBytes:    meta.NumBytes,
		Created:     meta.CreationTime,
		Modified:    meta.LastModifiedTime,
		Description: meta.Description,
	}, nil
}

// InsertRows inserts rows into a table
func (c *Client) InsertRows(ctx context.Context, datasetID, tableID string, rows []map[string]interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dataset := c.bqClient.Dataset(datasetID)
	table := dataset.Table(tableID)
	inserter := table.Inserter()

	// Convert rows to BigQuery ValueSaver format
	var valueSavers []bigquery.ValueSaver
	for _, row := range rows {
		vs := &mapValueSaver{values: row}
		valueSavers = append(valueSavers, vs)
	}

	if err := inserter.Put(ctx, valueSavers); err != nil {
		return fmt.Errorf("failed to insert rows into %s.%s: %w", datasetID, tableID, err)
	}

	return nil
}

// ListTablesIDs lists all tables in a dataset
func (c *Client) ListTablesIDs(ctx context.Context, datasetID string) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dataset := c.bqClient.Dataset(datasetID)
	it := dataset.Tables(ctx)

	var tableIDs []string

	for {
		table, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate tables: %v", err)
		}
		tableIDs = append(tableIDs, table.TableID)
	}

	return tableIDs, nil
}

// ListTables lists all tables in a dataset with their metadata
func (c *Client) ListTables(ctx context.Context, datasetID string) ([]*TableInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dataset := c.bqClient.Dataset(datasetID)
	it := dataset.Tables(ctx)

	var tables []*TableInfo

	for {
		table, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate tables: %v", err)
		}

		meta, err := table.Metadata(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get metadata for table %s: %w", table.TableID, err)
		}

		tables = append(tables, &TableInfo{
			ProjectID:   c.config.ProjectID,
			DatasetID:   datasetID,
			TableID:     table.TableID,
			Schema:      meta.Schema,
			NumRows:     int64(meta.NumRows),
			NumBytes:    meta.NumBytes,
			Created:     meta.CreationTime,
			Modified:    meta.LastModifiedTime,
			Description: meta.Description,
		})
	}

	return tables, nil
}

// ExecuteDML executes Data Manipulation Language (DML) statements
func (c *Client) ExecuteDML(ctx context.Context, sql string) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	query := c.bqClient.Query(sql)
	query.Location = c.config.Location

	queryCtx, cancel := context.WithTimeout(ctx, time.Duration(c.config.QueryTimeoutSeconds)*time.Second)
	defer cancel()

	job, err := query.Run(queryCtx)
	if err != nil {
		return 0, fmt.Errorf("failed to run DML query: %w", err)
	}

	status, err := job.Wait(queryCtx)
	if err != nil {
		return 0, fmt.Errorf("DML job failed: %w", err)
	}

	if status.Err() != nil {
		return 0, fmt.Errorf("DML execution error: %w", status.Err())
	}

	// Get the number of affected rows
	if status.Statistics != nil && status.Statistics.Details != nil {
		if queryStats, ok := status.Statistics.Details.(*bigquery.QueryStatistics); ok {
			return queryStats.NumDMLAffectedRows, nil
		}
	}

	return 0, nil
}

// Helper types and functions

// mapValueSaver implements bigquery.ValueSaver for map[string]interface{}
type mapValueSaver struct {
	values map[string]interface{}
}

func (mvs *mapValueSaver) Save() (map[string]bigquery.Value, string, error) {
	bqValues := make(map[string]bigquery.Value)
	for k, v := range mvs.values {
		bqValues[k] = v
	}
	return bqValues, "", nil
}

// convertBigQueryValue converts bigquery.Value to a standard Go type
func convertBigQueryValue(v bigquery.Value) interface{} {
	switch val := v.(type) {
	case nil:
		return nil
	case bool, int64, float64, string, time.Time:
		return val
	case []bigquery.Value:
		result := make([]interface{}, len(val))
		for i, item := range val {
			result[i] = convertBigQueryValue(item)
		}
		return result
	case map[string]bigquery.Value:
		result := make(map[string]interface{})
		for k, item := range val {
			result[k] = convertBigQueryValue(item)
		}
		return result
	default:
		// For any other types, try to JSON marshal/unmarshal
		if data, err := json.Marshal(val); err == nil {
			var result interface{}
			if err := json.Unmarshal(data, &result); err == nil {
				return result
			}
		}
		return fmt.Sprintf("%v", val)
	}
}

// Utility functions for schema creation

// StringField creates a STRING field schema
func StringField(name string, required bool) *FieldSchema {
	return &FieldSchema{
		Name:     name,
		Type:     bigquery.StringFieldType,
		Required: required,
	}
}

// IntegerField creates an INTEGER field schema
func IntegerField(name string, required bool) *FieldSchema {
	return &FieldSchema{
		Name:     name,
		Type:     bigquery.IntegerFieldType,
		Required: required,
	}
}

// FloatField creates a FLOAT field schema
func FloatField(name string, required bool) *FieldSchema {
	return &FieldSchema{
		Name:     name,
		Type:     bigquery.FloatFieldType,
		Required: required,
	}
}

// BooleanField creates a BOOLEAN field schema
func BooleanField(name string, required bool) *FieldSchema {
	return &FieldSchema{
		Name:     name,
		Type:     bigquery.BooleanFieldType,
		Required: required,
	}
}

// TimestampField creates a TIMESTAMP field schema
func TimestampField(name string, required bool) *FieldSchema {
	return &FieldSchema{
		Name:     name,
		Type:     bigquery.TimestampFieldType,
		Required: required,
	}
}

// DateField creates a DATE field schema
func DateField(name string, required bool) *FieldSchema {
	return &FieldSchema{
		Name:     name,
		Type:     bigquery.DateFieldType,
		Required: required,
	}
}
