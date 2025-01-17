package entity

import (
	"errors"
	"testing"

	"github.com/pakkasys/fluidapi/database/query"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCountEntities_NormalOperation tests the CountEntities function.
func TestCountEntities_NormalOperation(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)

	// Setup mock expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Close").Return(nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(nil)

	// Example table name and dbOptions
	tableName := "test_table"
	dbOptions := &query.CountOptions{}

	count, err := Count(mockDB, tableName, dbOptions)

	assert.NoError(t, err)
	assert.Equal(t, 0, count) // Adjust as per the test case

	// Verify that all expectations were met
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

// TestCountEntities_PrepareError tests the case where an error occurs during
// prepare call.
func TestCountEntities_PrepareError(t *testing.T) {
	mockDB := new(utilmock.MockDB)

	// Setup mock expectations for Prepare error
	mockDB.On("Prepare", mock.Anything).Return(nil, errors.New("prepare error"))

	// Example table name and dbOptions
	tableName := "test_table"
	dbOptions := &query.CountOptions{}

	count, err := Count(mockDB, tableName, dbOptions)

	assert.Equal(t, 0, count)
	assert.EqualError(t, err, "prepare error")

	// Verify that all expectations were met
	mockDB.AssertExpectations(t)
}

// TestCountEntities_QueryRowError tests the case where an error occurs during
// query row call.
func TestCountEntities_QueryRowError(t *testing.T) {
	mockDB := new(utilmock.MockDB)
	mockStmt := new(utilmock.MockStmt)
	mockRow := new(utilmock.MockRow)

	// Setup mock expectations
	mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
	mockStmt.On("Close").Return(nil)
	mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
	mockRow.On("Scan", mock.Anything).Return(errors.New("query row error"))

	// Example table name and dbOptions
	tableName := "test_table"
	dbOptions := &query.CountOptions{}

	count, err := Count(mockDB, tableName, dbOptions)

	assert.Equal(t, 0, count)
	assert.EqualError(t, err, "query row error")

	// Verify that all expectations were met
	mockDB.AssertExpectations(t)
	mockStmt.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}
