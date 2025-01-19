package entity

import (
	"errors"
	"testing"

	databasemock "github.com/pakkasys/fluidapi/database/mock"
	"github.com/stretchr/testify/assert"
)

// TestExec tests the Exec method.
func TestExec(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockResult *MockSQLResult)
		query         string
		parameters    []any
		expectedError string
		expectResult  bool
	}{
		{
			name: "Normal Operation",
			setupMocks: func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockResult *MockSQLResult) {
				mockDB.On("Prepare", "UPDATE users SET name = ? WHERE id = ?").Return(mockStmt, nil)
				mockStmt.On("Exec", []any{"Alice", 1}).Return(mockResult, nil)
				mockStmt.On("Close").Return(nil)
			},
			query:         "UPDATE users SET name = ? WHERE id = ?",
			parameters:    []any{"Alice", 1},
			expectedError: "",
			expectResult:  true,
		},
		{
			name: "Prepare Error",
			setupMocks: func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockResult *MockSQLResult) {
				mockDB.On("Prepare", "UPDATE users SET name = ? WHERE id = ?").Return(nil, errors.New("prepare error"))
			},
			query:         "UPDATE users SET name = ? WHERE id = ?",
			parameters:    []any{"Alice", 1},
			expectedError: "prepare error",
			expectResult:  false,
		},
		{
			name: "Exec Error",
			setupMocks: func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockResult *MockSQLResult) {
				mockDB.On("Prepare", "UPDATE users SET name = ? WHERE id = ?").Return(mockStmt, nil)
				mockStmt.On("Exec", []any{"Alice", 1}).Return(nil, errors.New("exec error"))
				mockStmt.On("Close").Return(nil)
			},
			query:         "UPDATE users SET name = ? WHERE id = ?",
			parameters:    []any{"Alice", 1},
			expectedError: "exec error",
			expectResult:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(databasemock.MockDB)
			mockStmt := new(databasemock.MockStmt)
			mockResult := new(MockSQLResult)

			if tt.setupMocks != nil {
				tt.setupMocks(mockDB, mockStmt, mockResult)
			}

			// Act
			result, err := Exec(mockDB, tt.query, tt.parameters)

			// Assert
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectResult {
					assert.NotNil(t, result)
				} else {
					assert.Nil(t, result)
				}
			}

			// Verify mock expectations
			mockDB.AssertExpectations(t)
			mockStmt.AssertExpectations(t)
			mockResult.AssertExpectations(t)
		})
	}
}

// TestQuery tests the Query function
func TestQuery(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(
			mockDB *databasemock.MockDB,
			mockStmt *databasemock.MockStmt,
			mockRows *databasemock.MockRows,
		)
		query         string
		parameters    []any
		expectedError string
		expectRows    bool
		expectStmt    bool
	}{
		{
			name: "Normal Operation",
			setupMocks: func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockRows *databasemock.MockRows) {
				mockDB.On("Prepare", "SELECT * FROM users WHERE id = ?").Return(mockStmt, nil)
				mockStmt.On("Query", []any{1}).Return(mockRows, nil)
			},
			query:         "SELECT * FROM users WHERE id = ?",
			parameters:    []any{1},
			expectedError: "",
			expectRows:    true,
			expectStmt:    true,
		},
		{
			name: "Prepare Error",
			setupMocks: func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockRows *databasemock.MockRows) {
				mockDB.On("Prepare", "SELECT * FROM users WHERE id = ?").Return(nil, errors.New("prepare error"))
			},
			query:         "SELECT * FROM users WHERE id = ?",
			parameters:    []any{1},
			expectedError: "prepare error",
			expectRows:    false,
			expectStmt:    false,
		},
		{
			name: "Query Error",
			setupMocks: func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockRows *databasemock.MockRows) {
				mockDB.On("Prepare", "SELECT * FROM users WHERE id = ?").Return(mockStmt, nil)
				mockStmt.On("Query", []any{1}).Return(nil, errors.New("query error"))
				mockStmt.On("Close").Return(nil)
			},
			query:         "SELECT * FROM users WHERE id = ?",
			parameters:    []any{1},
			expectedError: "query error",
			expectRows:    false,
			expectStmt:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(databasemock.MockDB)
			mockStmt := new(databasemock.MockStmt)
			mockRows := new(databasemock.MockRows)

			if tt.setupMocks != nil {
				tt.setupMocks(mockDB, mockStmt, mockRows)
			}

			// Act
			rows, stmt, err := Query(mockDB, tt.query, tt.parameters)

			// Assert
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, rows)
				assert.Nil(t, stmt)
			} else {
				assert.NoError(t, err)
				if tt.expectRows {
					assert.NotNil(t, rows)
				} else {
					assert.Nil(t, rows)
				}
				if tt.expectStmt {
					assert.NotNil(t, stmt)
				} else {
					assert.Nil(t, stmt)
				}
			}

			// Verify mock expectations
			mockDB.AssertExpectations(t)
			mockStmt.AssertExpectations(t)
			mockRows.AssertExpectations(t)
		})
	}
}
