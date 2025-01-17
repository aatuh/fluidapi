package entity

import (
	"errors"
	"testing"

	"github.com/pakkasys/fluidapi/database/query"
	"github.com/pakkasys/fluidapi/database/util"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestDeleteEntities_NormalOperation tests the normal operation of
// DeleteEntities.
func TestDeleteEntities_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(MockSQLResult)

	// Test selectors and delete options
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}
	opts := query.DeleteOptions{
		Limit:  5,
		Orders: nil,
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)
	mockResult.On("RowsAffected").Return(int64(2), nil)

	rowsAffected, err := Delete(mockDB, "user", selectors, &opts)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), rowsAffected)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestDeleteEntities_DeleteError tests the case where an error occurs during
// the delete operation.
func TestDeleteEntities_DeleteError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test selectors and delete options
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}
	opts := query.DeleteOptions{
		Limit:  5,
		Orders: nil,
	}

	// Simulate an error during the delete operation
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("delete error"))

	_, err := Delete(mockDB, "user", selectors, &opts)

	assert.EqualError(t, err, "delete error")
	mockDB.AssertExpectations(t)
}

// TestDeleteEntities_RowsAffectedError tests the case where an error occurs
// when getting the number of rows affected.
func TestDeleteEntities_RowsAffectedError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(MockSQLResult)

	// Test selectors and delete options
	selectors := []util.Selector{
		{Table: "user", Field: "id", Predicate: "=", Value: 1},
	}
	opts := query.DeleteOptions{
		Limit:  5,
		Orders: nil,
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)

	// Simulate an error when calling RowsAffected
	mockResult.On("RowsAffected").Return(int64(0), errors.New("rows affected error"))

	_, err := Delete(mockDB, "user", selectors, &opts)

	assert.EqualError(t, err, "rows affected error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestDelete_NormalOperation tests the normal operation of the delete function.
func TestDelete_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(MockSQLResult)

	// Test selectors and options
	selectors := []util.Selector{
		{Field: "id", Predicate: "=", Value: 1},
	}
	opts := query.DeleteOptions{
		Limit: 10,
		Orders: []util.Order{
			{Table: "user", Field: "name", Direction: "ASC"},
		},
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)

	result, err := delete(mockDB, "user", selectors, &opts)

	assert.NoError(t, err)
	assert.Equal(t, mockResult, result)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestDelete_NoSelectors tests the case where no selectors are provided.
func TestDelete_NoSelectors(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockResult := new(MockSQLResult)

	// Empty selectors and options
	selectors := []util.Selector{}
	opts := query.DeleteOptions{
		Limit:  0,
		Orders: nil,
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
	mockStmt.On("Close").Return(nil)

	result, err := delete(mockDB, "user", selectors, &opts)

	assert.NoError(t, err)
	assert.Equal(t, mockResult, result)
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockResult.AssertExpectations(t)
}

// TestDelete_PrepareError tests the case where an error occurs during SQL
// preparation.
func TestDelete_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Test selectors and options
	selectors := []util.Selector{
		{Field: "id", Predicate: "=", Value: 1},
	}
	opts := query.DeleteOptions{
		Limit:  0,
		Orders: nil,
	}

	// Simulate an error during Prepare
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	_, err := delete(mockDB, "user", selectors, &opts)

	assert.EqualError(t, err, "prepare error")
	mockDB.AssertExpectations(t)
}

// TestDelete_ExecError tests the case where an error occurs during SQL
// execution.
func TestDelete_ExecError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)

	// Test selectors and options
	selectors := []util.Selector{
		{Field: "id", Predicate: "=", Value: 1},
	}
	opts := query.DeleteOptions{
		Limit:  0,
		Orders: nil,
	}

	// Setup the mock DB expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	// Simulate an error during Exec
	mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
	mockStmt.On("Close").Return(nil)

	_, err := delete(mockDB, "user", selectors, &opts)

	assert.EqualError(t, err, "exec error")
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
}
