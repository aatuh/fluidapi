package entity

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/pakkasys/fluidapi/database/query"
	"github.com/pakkasys/fluidapi/database/util"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRowScanner is a mock implementation of RowScanner.
type MockRowScanner[T any] struct {
	mock.Mock
}

func (m *MockRowScanner[T]) Scan(row util.Row, entity *T) error {
	return row.Scan(entity)
}

// MockRowScannerMultiple is a mock implementation of RowScannerMultiple.
type MockRowScannerMultiple[T any] struct {
	mock.Mock
}

func (m *MockRowScannerMultiple[T]) Scan(rows util.Rows, entity *T) error {
	return rows.Scan(entity)
}

type TestEntity struct {
	ID   int
	Name string
	Age  int
}

// TestGetEntity_NormalOperation tests the normal operation where a single
// entity is successfully retrieved.
func TestGetEntity_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test table name and options
	tableName := "user"
	dbOptions := &query.GetOptions{}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockRow.On("Scan", []any{&TestEntity{}}).Return(nil).Once()
	mockStmt.On("Close").Return(nil)
	mockRow.On("Err").Return(nil)

	entity, err := Get(tableName, mockScanner.Scan, mockDB, dbOptions)

	assert.NoError(t, err)
	assert.NotNil(t, entity)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntity_QueryError tests the case where an error occurs during query
// execution.
func TestGetEntity_QueryError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test table name and options
	tableName := "user"
	dbOptions := &query.GetOptions{}

	// Simulate an error during query execution
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("query error"))

	entity, err := Get[TestEntity](tableName, nil, mockDB, dbOptions)

	assert.Nil(t, entity)
	assert.EqualError(t, err, "query error")
	mockDB.AssertExpectations(t)
}

// TestGetEntity_NoRows tests the case where sql.ErrNoRows is returned.
func TestGetEntity_NoRows(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test table name and options
	tableName := "user"
	dbOptions := &query.GetOptions{}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockRow.On("Scan", []any{&TestEntity{}}).Return(nil).Once()
	mockStmt.On("Close").Return(nil)
	mockRow.On("Err").Return(sql.ErrNoRows).Once()

	entity, err := Get(tableName, mockScanner.Scan, mockDB, dbOptions)

	assert.NoError(t, err)
	assert.Nil(t, entity) // No entity should be returned
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntity_RowScannerError tests the case where an error occurs during row
// scanning.
func TestGetEntity_RowScannerError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test table name and options
	tableName := "user"
	dbOptions := &query.GetOptions{}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockRow.On("Scan", []any{&TestEntity{}}).
		Return(errors.New("row scanner error")).Once()
	mockStmt.On("Close").Return(nil)

	entity, err := Get(tableName, mockScanner.Scan, mockDB, dbOptions)

	assert.Nil(t, entity)
	assert.EqualError(t, err, "row scanner error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntities_NormalOperation tests normal operation where multiple
// entities are successfully retrieved.
func TestGetEntities_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and options
	tableName := "user"
	dbOptions := &query.GetOptions{}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Query", mock.Anything).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once() // Simulate a row read
	mockRows.On("Scan", []any{&TestEntity{}}).Return(nil).Once()
	mockRows.On("Next").Return(false).Once() // No more rows
	mockRows.On("Err").Return(nil)

	entities, err := GetMany(tableName, mockScanner.Scan, mockDB, dbOptions)

	assert.NoError(t, err)
	assert.NotNil(t, entities)
	assert.Len(t, entities, 1)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntities_QueryError tests the case where an error occurs during query
// execution.
func TestGetEntities_QueryError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test query and options
	tableName := "user"
	dbOptions := &query.GetOptions{}

	// Simulate an error during RowsQuery
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("query error"))

	entities, err := GetMany[TestEntity](tableName, nil, mockDB, dbOptions)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "query error")
	mockDB.AssertExpectations(t)
}

// TestGetEntities_NoRows tests the case where sql.ErrNoRows is returned.
func TestGetEntities_NoRows(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and options
	tableName := "user"
	dbOptions := &query.GetOptions{}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Query", mock.Anything).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(false).Once() // No rows found
	mockRows.On("Err").Return(sql.ErrNoRows).Once()

	entities, err := GetMany(tableName, mockScanner.Scan, mockDB, dbOptions)

	assert.NoError(t, err)
	assert.Len(t, entities, 0) // No entities should be returned
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestGetEntities_RowScannerError tests the case where an error occurs during
// row scanning.
func TestGetEntities_RowScannerError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and options
	tableName := "user"
	dbOptions := &query.GetOptions{}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Query", mock.Anything).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once() // Simulate a row read
	mockRows.On("Scan", []any{&TestEntity{}}).
		Return(errors.New("row scanner error")).Once()

	entities, err := GetMany(tableName, mockScanner.Scan, mockDB, dbOptions)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "row scanner error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestQueryMultiple_NormalOperation tests normal operation where multiple
// entities are successfully retrieved.
func TestQueryMultiple_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM user WHERE active = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("Query", params).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once() // Simulate a row read
	mockRows.On("Scan", []any{&TestEntity{}}).Return(nil).Once()
	mockRows.On("Next").Return(false).Once() // No more rows
	mockRows.On("Err").Return(nil)

	entities, err := queryMultiple(mockDB, query, params, mockScanner.Scan)

	assert.NoError(t, err)
	assert.NotNil(t, entities)
	assert.Len(t, entities, 1)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestQueryMultiple_RowsQueryError tests the case where an error occurs during
// the query execution.
func TestQueryMultiple_RowsQueryError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test query and parameters
	query := "SELECT * FROM user WHERE active = ?"
	params := []any{1}

	// Simulate an error during RowsQuery
	mockDB.On("Prepare", query).Return(nil, errors.New("query error"))

	entities, err := queryMultiple[TestEntity](mockDB, query, params, nil)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "query error")
	mockDB.AssertExpectations(t)
}

// TestQueryMultiple_RowScannerError tests the case where an error occurs during
// row scanning.
func TestQueryMultiple_RowScannerError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM user WHERE active = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("Query", params).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once() // Simulate a row read
	mockRows.On("Scan", []any{&TestEntity{}}).
		Return(errors.New("row scanner error")).Once()

	entities, err := queryMultiple(mockDB, query, params, mockScanner.Scan)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "row scanner error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestQueryMultiple_RowsErr tests the case where rows.Err() returns an error.
func TestQueryMultiple_RowsErr(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRows := new(utilmock.MockRows)
	mockScanner := new(MockRowScannerMultiple[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM user WHERE active = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("Query", params).Return(mockRows, nil)
	mockStmt.On("Close").Return(nil)
	mockRows.On("Close").Return(nil)
	mockRows.On("Next").Return(true).Once() // Simulate a row read
	mockRows.On("Scan", []any{&TestEntity{}}).Return(nil).Once()
	mockRows.On("Next").Return(false).Once() // No more rows
	mockRows.On("Err").Return(errors.New("rows error"))

	entities, err := queryMultiple(mockDB, query, params, mockScanner.Scan)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "rows error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRows.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestQuerySingle_NormalOperation tests normal operation where a single entity
// is successfully retrieved.
func TestQuerySingle_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM user WHERE id = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("QueryRow", params).Return(mockRow)
	mockRow.On("Scan", []any{&TestEntity{}}).Return(nil).Once()
	mockStmt.On("Close").Return(nil)
	mockRow.On("Err").Return(nil)

	entity, err := querySingle(mockDB, query, params, mockScanner.Scan)

	assert.NoError(t, err)
	assert.NotNil(t, entity)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestQuerySingle_PrepareError tests the case where an error occurs during
// query preparation.
func TestQuerySingle_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test query and parameters
	query := "SELECT * FROM user WHERE id = ?"
	params := []any{1}

	// Simulate an error during Prepare
	mockDB.On("Prepare", query).Return(nil, errors.New("prepare error"))

	entity, err := querySingle[TestEntity](mockDB, query, params, nil)

	assert.Nil(t, entity)
	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestQuerySingle_RowScannerError tests the case where the row scanner returns
// an error.
func TestQuerySingle_RowScannerError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM user WHERE id = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("QueryRow", params).Return(mockRow)
	mockRow.On("Scan", []any{&TestEntity{}}).
		Return(errors.New("row scanner error")).Once()

	mockStmt.On("Close").Return(nil)

	entity, err := querySingle(mockDB, query, params, mockScanner.Scan)

	assert.Nil(t, entity)
	assert.EqualError(t, err, "row scanner error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestQuerySingle_RowErr tests the case where the row.Err() method returns an
// error.
func TestQuerySingle_RowErr(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)
	mockScanner := new(MockRowScanner[TestEntity])

	// Test query and parameters
	query := "SELECT * FROM user WHERE id = ?"
	params := []any{1}

	// Setup the mock DB expectations
	mockDB.On("Prepare", query).Return(mockStmt, nil)
	mockStmt.On("QueryRow", params).Return(mockRow)
	mockRow.On("Scan", []any{&TestEntity{}}).Return(nil).Once()
	mockStmt.On("Close").Return(nil)

	// Simulate an error returned by row.Err()
	mockRow.On("Err").Return(errors.New("row error"))

	entity, err := querySingle(mockDB, query, params, mockScanner.Scan)

	assert.Nil(t, entity)
	assert.EqualError(t, err, "row error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
	mockScanner.AssertExpectations(t)
}

// TestRowsToEntities_NoRowScannerMultiple tests the case where there is no
// RowScannerMultiple provided.
func TestRowsToEntities_NoRowScannerMultiple(t *testing.T) {
	entities, err := rowsToEntities[any](nil, nil)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "must provide rowScannerMultiple")
}

// TestRowsToEntities_NormalOperation tests normal operation with multiple rows.
func TestRowsToEntities_NormalOperation(t *testing.T) {
	mockRows := new(utilmock.MockRows)
	mockRowScanner := new(MockRowScannerMultiple[TestEntity])

	// Setup the row scanning behavior
	testEntity := TestEntity{}
	mockRows.On("Next").Return(true).Once() // Simulate the first row read
	mockRows.On("Scan", []any{&testEntity}).Return(nil).Once()
	mockRows.On("Next").Return(true).Once() // Simulate the second row read
	mockRows.On("Scan", []any{&testEntity}).Return(nil).Once()
	mockRows.On("Next").Return(false).Once() // No more rows
	mockRows.On("Err").Return(nil)           // No error in rows

	entities, err := rowsToEntities(mockRows, mockRowScanner.Scan)

	assert.NoError(t, err)
	assert.Len(t, entities, 2)
	mockRows.AssertExpectations(t)
	mockRowScanner.AssertExpectations(t)
}

// TestRowsToEntities_NoRows tests the case where there are no rows.
func TestRowsToEntities_NoRows(t *testing.T) {
	mockRows := new(utilmock.MockRows)
	mockRowScanner := new(MockRowScannerMultiple[TestEntity])

	// Setup the row scanning behavior
	mockRows.On("Next").Return(false).Once() // No rows
	mockRows.On("Err").Return(nil)           // No error in rows

	entities, err := rowsToEntities(mockRows, mockRowScanner.Scan)

	assert.NoError(t, err)
	assert.Len(t, entities, 0)
	mockRows.AssertExpectations(t)
	mockRowScanner.AssertExpectations(t)
}

// TestRowsToEntities_RowScannerError tests the case where the row scanner
// returns an error.
func TestRowsToEntities_RowScannerError(t *testing.T) {
	mockRows := new(utilmock.MockRows)
	mockRowScanner := new(MockRowScannerMultiple[TestEntity])

	// Setup the row scanning behavior
	mockRows.On("Next").Return(true).Once() // Simulate the first row read
	mockRows.On("Scan", []any{&TestEntity{}}).
		Return(errors.New("row scanner error")).Once()

	entities, err := rowsToEntities(mockRows, mockRowScanner.Scan)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "row scanner error")
	mockRows.AssertExpectations(t)
	mockRowScanner.AssertExpectations(t)
}

// TestRowsToEntities_RowsError tests the case where rows return an error.
func TestRowsToEntities_RowsError(t *testing.T) {
	mockRows := new(utilmock.MockRows)
	mockRowScanner := new(MockRowScannerMultiple[TestEntity])

	// Setup the row scanning behavior
	mockRows.On("Next").Return(true).Once() // Simulate the first row read
	mockRows.On("Scan", []any{&TestEntity{}}).Return(nil).Once()
	mockRows.On("Next").Return(false).Once() // No more rows
	mockRows.On("Err").Return(errors.New("rows error")).Once()

	entities, err := rowsToEntities(mockRows, mockRowScanner.Scan)

	assert.Nil(t, entities)
	assert.EqualError(t, err, "rows error")
	mockRows.AssertExpectations(t)
	mockRowScanner.AssertExpectations(t)
}
