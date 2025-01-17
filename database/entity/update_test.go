package entity

import (
	"errors"
	"testing"

	entitymock "github.com/pakkasys/fluidapi/database/entity/mock"
	"github.com/pakkasys/fluidapi/database/query"
	"github.com/pakkasys/fluidapi/database/util"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUpdateEntities_NormalOperation tests the normal operation where updates
// are successfully applied.
func TestUpdateEntities_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Test table name, updates, and selectors
	tableName := "user"
	updateFields := []query.UpdateField{
		{Field: "name", Value: "Alice"},
	}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("RowsAffected").Return(int64(1), nil)

	rowsAffected, err :=
		Update(mockDB, tableName, selectors, updateFields, mockSQLUtil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestUpdateEntities_NoUpdates tests the case where no updates are provided.
func TestUpdateEntities_NoUpdates(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Test table name and selectors
	tableName := "user"
	updateFields := []query.UpdateField{}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	rowsAffected, err :=
		Update(mockDB, tableName, selectors, updateFields, mockSQLUtil)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), rowsAffected)
	mockDB.AssertExpectations(t)
}

// TestUpdateEntities_Error tests the case where an error occurs during the
// update process.
func TestUpdateEntities_Error(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Test table name, updates, and selectors
	tableName := "user"
	updateFields := []query.UpdateField{
		{Field: "name", Value: "Alice"},
	}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	// Simulate an error during Prepare
	mockDB.On("Prepare", mock.Anything).
		Return(nil, errors.New("prepare error"))
	mockSQLUtil.On("CheckDBError", mock.Anything).
		Return(errors.New("prepare error"))

	rowsAffected, err :=
		Update(mockDB, tableName, selectors, updateFields, mockSQLUtil)

	assert.Equal(t, int64(0), rowsAffected)
	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
	mockSQLUtil.AssertExpectations(t)
}

// TestCheckUpdateResult_NormalOperation tests the normal operation where rows
// are affected.
func TestCheckUpdateResult_NormalOperation(t *testing.T) {
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Setup mock expectations
	mockResult.On("RowsAffected").Return(int64(1), nil)

	rowsAffected, err := checkUpdateResult(mockResult, nil, mockSQLUtil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)
	mockResult.AssertExpectations(t)
}

// TestCheckUpdateResult_OtherError tests the case where a non-MySQL error
// occurs.
func TestCheckUpdateResult_OtherError(t *testing.T) {
	mockSQLUtil := new(entitymock.MockSQLUtil)

	otherErr := errors.New("some other error")
	mockSQLUtil.On("CheckDBError", otherErr).Return(otherErr)

	rowsAffected, err := checkUpdateResult(nil, otherErr, mockSQLUtil)

	assert.Equal(t, int64(0), rowsAffected)
	assert.EqualError(t, err, "some other error")
	mockSQLUtil.AssertExpectations(t)
}

// TestCheckUpdateResult_RowsAffectedError tests the case where an error occurs
// when retrieving rows affected.
func TestCheckUpdateResult_RowsAffectedError(t *testing.T) {
	mockResult := new(utilmock.MockResult)
	mockSQLUtil := new(entitymock.MockSQLUtil)

	// Simulate an error when retrieving RowsAffected
	mockResult.On("RowsAffected").
		Return(int64(0), errors.New("rows affected error"))

	rowsAffected, err := checkUpdateResult(mockResult, nil, mockSQLUtil)

	assert.Equal(t, int64(0), rowsAffected)
	assert.EqualError(t, err, "rows affected error")
	mockResult.AssertExpectations(t)
}

// TestUpdate_NormalOperation tests the normal operation of the update function.
func TestUpdate_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)

	// Test table name, updates, and selectors
	tableName := "user"
	updateFields := []query.UpdateField{
		{Field: "name", Value: "Alice"},
	}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)

	result, err := update(mockDB, tableName, updateFields, selectors)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestUpdate_PrepareError tests the case where an error occurs during the
// preparation of the statement.
func TestUpdate_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test table name, updates, and selectors
	tableName := "user"
	updateFields := []query.UpdateField{
		{Field: "name", Value: "Alice"},
	}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	// Simulate an error during Prepare
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	result, err := update(mockDB, tableName, updateFields, selectors)

	assert.Nil(t, result)
	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestUpdate_ExecError tests the case where an error occurs during the
// execution of the statement.
func TestUpdate_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test table name, updates, and selectors
	tableName := "user"
	updates := []query.UpdateField{
		{Field: "name", Value: "Alice"},
	}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
	mockStmt.On("Close").Return(nil)

	result, err := update(mockDB, tableName, updates, selectors)

	assert.Nil(t, result)
	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}

// TestUpdate_EmptyUpdates tests the case where no updates are provided.
func TestUpdate_EmptyUpdates(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(utilmock.MockResult)

	// Test table name and selectors
	tableName := "user"
	updateFields := []query.UpdateField{}
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)

	result, err := update(mockDB, tableName, updateFields, selectors)

	assert.NotNil(t, result)
	assert.Nil(t, err)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}
