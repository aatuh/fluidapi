package entity

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/pakkasys/fluidapi/database"
	databasemock "github.com/pakkasys/fluidapi/database/mock"
	"github.com/pakkasys/fluidapi/database/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRowScanner is a mock implementation of RowScanner.
type MockRowScanner[T any] struct {
	mock.Mock
}

func (m *MockRowScanner[T]) Scan(row database.Row, entity *T) error {
	return row.Scan(entity)
}

// MockRowScannerMultiple is a mock implementation of RowScannerMultiple.
type MockRowScannerMultiple[T any] struct {
	mock.Mock
}

func (m *MockRowScannerMultiple[T]) Scan(rows database.Rows, entity *T) error {
	return rows.Scan(entity)
}

type GetTestStruct struct {
	ID   int
	Name string
	Age  int
}

// TestGet tests the Get function.
func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(
			mockDB *databasemock.MockDB,
			mockStmt *databasemock.MockStmt,
			mockRow *databasemock.MockRow,
			mockScanner *MockRowScanner[GetTestStruct],
		)
		tableName     string
		options       *query.GetOptions
		expectedError string
		expectEntity  bool
	}{
		{
			name: "Normal Operation",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRow *databasemock.MockRow,
				mockScanner *MockRowScanner[GetTestStruct],
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
				mockRow.On("Scan", mock.Anything).Return(nil).Once()
				mockRow.On("Err").Return(nil)
				mockStmt.On("Close").Return(nil)
			},
			tableName:     "user",
			options:       &query.GetOptions{},
			expectedError: "",
			expectEntity:  true,
		},
		{
			name: "Query Error",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRow *databasemock.MockRow,
				mockScanner *MockRowScanner[GetTestStruct],
			) {
				mockDB.On("Prepare", mock.Anything).
					Return(nil, errors.New("query error"))
			},
			tableName:     "user",
			options:       &query.GetOptions{},
			expectedError: "query error",
			expectEntity:  false,
		},
		{
			name: "No Rows",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRow *databasemock.MockRow,
				mockScanner *MockRowScanner[GetTestStruct],
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
				mockRow.On("Scan", mock.Anything).Return(nil).Once()
				mockRow.On("Err").Return(sql.ErrNoRows).Once()
				mockStmt.On("Close").Return(nil)
			},
			tableName:     "user",
			options:       &query.GetOptions{},
			expectedError: "",
			expectEntity:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(databasemock.MockDB)
			mockStmt := new(databasemock.MockStmt)
			mockRow := new(databasemock.MockRow)
			mockScanner := new(MockRowScanner[GetTestStruct])

			if tt.setupMocks != nil {
				tt.setupMocks(mockDB, mockStmt, mockRow, mockScanner)
			}

			// Act
			entity, err := Get(
				tt.tableName,
				mockScanner.Scan,
				mockDB,
				tt.options,
			)

			// Assert
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, entity)
			} else {
				assert.NoError(t, err)
				if tt.expectEntity {
					assert.NotNil(t, entity)
				} else {
					assert.Nil(t, entity)
				}
			}

			// Verify mock expectations
			mockDB.AssertExpectations(t)
			mockStmt.AssertExpectations(t)
			mockRow.AssertExpectations(t)
			mockScanner.AssertExpectations(t)
		})
	}
}

// TestGetMany tests the GetMany function
func TestGetMany(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(
			mockDB *databasemock.MockDB,
			mockStmt *databasemock.MockStmt,
			mockRows *databasemock.MockRows,
			mockScanner *MockRowScannerMultiple[GetTestStruct],
		)
		tableName      string
		options        *query.GetOptions
		expectedError  string
		expectEntities int
	}{
		{
			name: "Normal Operation",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRows *databasemock.MockRows,
				mockScanner *MockRowScannerMultiple[GetTestStruct],
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("Query", mock.Anything).Return(mockRows, nil)
				mockStmt.On("Close").Return(nil)
				mockRows.On("Close").Return(nil)
				mockRows.On("Next").Return(true).Once()
				mockRows.On("Scan", mock.Anything).Return(nil).Once()
				mockRows.On("Next").Return(false).Once()
				mockRows.On("Err").Return(nil)
			},
			tableName:      "user",
			options:        &query.GetOptions{},
			expectedError:  "",
			expectEntities: 1,
		},
		{
			name: "Query Error",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRows *databasemock.MockRows,
				mockScanner *MockRowScannerMultiple[GetTestStruct],
			) {
				mockDB.On("Prepare", mock.Anything).
					Return(nil, errors.New("query error"))
			},
			tableName:      "user",
			options:        &query.GetOptions{},
			expectedError:  "query error",
			expectEntities: 0,
		},
		{
			name: "No Rows",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRows *databasemock.MockRows,
				mockScanner *MockRowScannerMultiple[GetTestStruct],
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("Query", mock.Anything).Return(mockRows, nil)
				mockStmt.On("Close").Return(nil)
				mockRows.On("Close").Return(nil)
				mockRows.On("Next").Return(false).Once()
				mockRows.On("Err").Return(sql.ErrNoRows).Once()
			},
			tableName:      "user",
			options:        &query.GetOptions{},
			expectedError:  "",
			expectEntities: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(databasemock.MockDB)
			mockStmt := new(databasemock.MockStmt)
			mockRows := new(databasemock.MockRows)
			mockScanner := new(MockRowScannerMultiple[GetTestStruct])

			if tt.setupMocks != nil {
				tt.setupMocks(mockDB, mockStmt, mockRows, mockScanner)
			}

			// Act
			entities, err := GetMany(
				tt.tableName,
				mockScanner.Scan,
				mockDB,
				tt.options,
			)

			// Assert
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, entities)
			} else {
				assert.NoError(t, err)
				assert.Len(t, entities, tt.expectEntities)
			}

			// Verify mock expectations
			mockDB.AssertExpectations(t)
			mockStmt.AssertExpectations(t)
			mockRows.AssertExpectations(t)
			mockScanner.AssertExpectations(t)
		})
	}
}

// TestQueryMultiple tests the QueryMultiple function
func TestQueryMultiple(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(
			mockDB *databasemock.MockDB,
			mockStmt *databasemock.MockStmt,
			mockRows *databasemock.MockRows,
			mockScanner *MockRowScannerMultiple[GetTestStruct],
		)
		query          string
		params         []any
		expectedError  string
		expectEntities int
	}{
		{
			name: "Normal Operation",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRows *databasemock.MockRows,
				mockScanner *MockRowScannerMultiple[GetTestStruct],
			) {
				mockDB.On("Prepare", "SELECT * FROM user WHERE active = ?").Return(mockStmt, nil)
				mockStmt.On("Query", []any{1}).Return(mockRows, nil)
				mockStmt.On("Close").Return(nil)
				mockRows.On("Close").Return(nil)
				mockRows.On("Next").Return(true).Once()
				mockRows.On("Scan", mock.Anything).Return(nil).Once()
				mockRows.On("Next").Return(false).Once()
				mockRows.On("Err").Return(nil)
			},
			query:          "SELECT * FROM user WHERE active = ?",
			params:         []any{1},
			expectedError:  "",
			expectEntities: 1,
		},
		{
			name: "Query Error",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRows *databasemock.MockRows,
				mockScanner *MockRowScannerMultiple[GetTestStruct],
			) {
				mockDB.On("Prepare", "SELECT * FROM user WHERE active = ?").
					Return(nil, errors.New("query error"))
			},
			query:          "SELECT * FROM user WHERE active = ?",
			params:         []any{1},
			expectedError:  "query error",
			expectEntities: 0,
		},
		{
			name: "Rows Error",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRows *databasemock.MockRows,
				mockScanner *MockRowScannerMultiple[GetTestStruct],
			) {
				mockDB.On("Prepare", "SELECT * FROM user WHERE active = ?").Return(mockStmt, nil)
				mockStmt.On("Query", []any{1}).Return(mockRows, nil)
				mockStmt.On("Close").Return(nil)
				mockRows.On("Close").Return(nil)
				mockRows.On("Next").Return(true).Once()
				mockRows.On("Scan", mock.Anything).Return(nil).Once()
				mockRows.On("Next").Return(false).Once()
				mockRows.On("Err").Return(errors.New("rows error"))
			},
			query:          "SELECT * FROM user WHERE active = ?",
			params:         []any{1},
			expectedError:  "rows error",
			expectEntities: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(databasemock.MockDB)
			mockStmt := new(databasemock.MockStmt)
			mockRows := new(databasemock.MockRows)
			mockScanner := new(MockRowScannerMultiple[GetTestStruct])

			if tt.setupMocks != nil {
				tt.setupMocks(mockDB, mockStmt, mockRows, mockScanner)
			}

			// Act
			entities, err := queryMultiple(
				mockDB,
				tt.query,
				tt.params,
				mockScanner.Scan,
			)

			// Assert
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, entities)
			} else {
				assert.NoError(t, err)
				assert.Len(t, entities, tt.expectEntities)
			}

			// Verify mock expectations
			mockDB.AssertExpectations(t)
			mockStmt.AssertExpectations(t)
			mockRows.AssertExpectations(t)
			mockScanner.AssertExpectations(t)
		})
	}
}

// TestQuerySingle tests the QuerySingle function
func TestQuerySingle(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(
			mockDB *databasemock.MockDB,
			mockStmt *databasemock.MockStmt,
			mockRow *databasemock.MockRow,
			mockScanner *MockRowScanner[GetTestStruct],
		)
		query         string
		params        []any
		expectedError string
		expectEntity  bool
	}{
		{
			name: "Normal Operation",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRow *databasemock.MockRow,
				mockScanner *MockRowScanner[GetTestStruct],
			) {
				mockDB.On("Prepare", "SELECT * FROM user WHERE id = ?").Return(mockStmt, nil)
				mockStmt.On("QueryRow", []any{1}).Return(mockRow)
				mockRow.On("Scan", mock.Anything).Return(nil).Once()
				mockStmt.On("Close").Return(nil)
				mockRow.On("Err").Return(nil)
			},
			query:         "SELECT * FROM user WHERE id = ?",
			params:        []any{1},
			expectedError: "",
			expectEntity:  true,
		},
		{
			name: "Prepare Error",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRow *databasemock.MockRow,
				mockScanner *MockRowScanner[GetTestStruct],
			) {
				mockDB.On("Prepare", "SELECT * FROM user WHERE id = ?").
					Return(nil, errors.New("prepare error"))
			},
			query:         "SELECT * FROM user WHERE id = ?",
			params:        []any{1},
			expectedError: "prepare error",
			expectEntity:  false,
		},
		{
			name: "Row Scanner Error",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRow *databasemock.MockRow,
				mockScanner *MockRowScanner[GetTestStruct],
			) {
				mockDB.On("Prepare", "SELECT * FROM user WHERE id = ?").Return(mockStmt, nil)
				mockStmt.On("QueryRow", []any{1}).Return(mockRow)
				mockRow.On("Scan", mock.Anything).Return(errors.New("row scanner error")).Once()
				mockStmt.On("Close").Return(nil)
			},
			query:         "SELECT * FROM user WHERE id = ?",
			params:        []any{1},
			expectedError: "row scanner error",
			expectEntity:  false,
		},
		{
			name: "Row Error",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRow *databasemock.MockRow,
				mockScanner *MockRowScanner[GetTestStruct],
			) {
				mockDB.On("Prepare", "SELECT * FROM user WHERE id = ?").Return(mockStmt, nil)
				mockStmt.On("QueryRow", []any{1}).Return(mockRow)
				mockRow.On("Scan", mock.Anything).Return(nil).Once()
				mockStmt.On("Close").Return(nil)
				mockRow.On("Err").Return(errors.New("row error"))
			},
			query:         "SELECT * FROM user WHERE id = ?",
			params:        []any{1},
			expectedError: "row error",
			expectEntity:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(databasemock.MockDB)
			mockStmt := new(databasemock.MockStmt)
			mockRow := new(databasemock.MockRow)
			mockScanner := new(MockRowScanner[GetTestStruct])

			if tt.setupMocks != nil {
				tt.setupMocks(mockDB, mockStmt, mockRow, mockScanner)
			}

			// Act
			entity, err := querySingle(
				mockDB,
				tt.query,
				tt.params,
				mockScanner.Scan,
			)

			// Assert
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, entity)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, entity)
			}

			// Verify mock expectations
			mockDB.AssertExpectations(t)
			mockStmt.AssertExpectations(t)
			mockRow.AssertExpectations(t)
			mockScanner.AssertExpectations(t)
		})
	}
}

// TestRowsToEntities tests the RowsToEntities function
func TestRowsToEntities(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(
			mockRows *databasemock.MockRows,
			mockScanner *MockRowScannerMultiple[GetTestStruct],
		)
		rowScanner     RowScannerMultiple[GetTestStruct]
		expectedError  string
		expectedLength int
	}{
		{
			name: "No RowScannerMultiple",
			setupMocks: func(
				mockRows *databasemock.MockRows,
				mockScanner *MockRowScannerMultiple[GetTestStruct],
			) {
				// No mocks needed for this test case
			},
			rowScanner:     nil,
			expectedError:  "must provide rowScannerMultiple",
			expectedLength: 0,
		},
		{
			name: "Normal Operation",
			setupMocks: func(
				mockRows *databasemock.MockRows,
				mockScanner *MockRowScannerMultiple[GetTestStruct],
			) {
				testEntity := GetTestStruct{}
				mockRows.On("Next").Return(true).Once() // Simulate the first row read
				mockRows.On("Scan", []any{&testEntity}).Return(nil).Once()
				mockRows.On("Next").Return(true).Once() // Simulate the second row read
				mockRows.On("Scan", []any{&testEntity}).Return(nil).Once()
				mockRows.On("Next").Return(false).Once() // No more rows
				mockRows.On("Err").Return(nil)
			},
			rowScanner:     new(MockRowScannerMultiple[GetTestStruct]).Scan,
			expectedError:  "",
			expectedLength: 2,
		},
		{
			name: "No Rows",
			setupMocks: func(
				mockRows *databasemock.MockRows,
				mockScanner *MockRowScannerMultiple[GetTestStruct],
			) {
				mockRows.On("Next").Return(false).Once() // No rows
				mockRows.On("Err").Return(nil)
			},
			rowScanner:     new(MockRowScannerMultiple[GetTestStruct]).Scan,
			expectedError:  "",
			expectedLength: 0,
		},
		{
			name: "Row Scanner Error",
			setupMocks: func(
				mockRows *databasemock.MockRows,
				mockScanner *MockRowScannerMultiple[GetTestStruct],
			) {
				mockRows.On("Next").Return(true).Once() // Simulate the first row read
				mockRows.On("Scan", []any{&GetTestStruct{}}).
					Return(errors.New("row scanner error")).Once()
			},
			rowScanner:     new(MockRowScannerMultiple[GetTestStruct]).Scan,
			expectedError:  "row scanner error",
			expectedLength: 0,
		},
		{
			name: "Rows Error",
			setupMocks: func(
				mockRows *databasemock.MockRows,
				mockScanner *MockRowScannerMultiple[GetTestStruct],
			) {
				mockRows.On("Next").Return(true).Once() // Simulate the first row read
				mockRows.On("Scan", []any{&GetTestStruct{}}).Return(nil).Once()
				mockRows.On("Next").Return(false).Once() // No more rows
				mockRows.On("Err").Return(errors.New("rows error")).Once()
			},
			rowScanner:     new(MockRowScannerMultiple[GetTestStruct]).Scan,
			expectedError:  "rows error",
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRows := new(databasemock.MockRows)
			mockScanner := new(MockRowScannerMultiple[GetTestStruct])

			if tt.setupMocks != nil {
				tt.setupMocks(mockRows, mockScanner)
			}

			// Act
			entities, err := rowsToEntities(mockRows, tt.rowScanner)

			// Assert
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, entities)
			} else {
				assert.NoError(t, err)
				assert.Len(t, entities, tt.expectedLength)
			}

			// Verify mock expectations
			mockRows.AssertExpectations(t)
			mockScanner.AssertExpectations(t)
		})
	}
}
