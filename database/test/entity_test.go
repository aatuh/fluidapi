package test

// import (
// 	"database/sql"
// 	"errors"
// 	"testing"

// 	databasemock "github.com/pakkasys/fluidapi/database/mock"
// 	"github.com/pakkasys/fluidapi/database/types"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // MockSQLResult is a mock implementation of the sql.Result interface.
// type MockSQLResult struct {
// 	mock.Mock
// }

// func (m *MockSQLResult) LastInsertId() (int64, error) {
// 	args := m.Called()
// 	return args.Get(0).(int64), args.Error(1)
// }

// func (m *MockSQLResult) RowsAffected() (int64, error) {
// 	args := m.Called()
// 	return args.Get(0).(int64), args.Error(1)
// }

// type CreateTestStruct struct {
// 	ID   int
// 	Name string
// 	Age  int
// }

// // TestInsert tests the Insert function.
// func TestInsert(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setupMocks func(
// 			mockDB *databasemock.MockDB,
// 			mockStmt *databasemock.MockStmt,
// 			mockResult *databasemock.MockResult,
// 			mockErrorChecker *databasemock.MockErrorChecker,
// 		)
// 		entity        *CreateTestStruct
// 		expectedID    int64
// 		expectedError string
// 	}{
// 		{
// 			name: "Normal Operation",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
// 				mockStmt.On("Close").Return(nil)
// 				mockResult.On("LastInsertId").Return(int64(1), nil)
// 			},
// 			entity:        &CreateTestStruct{ID: 1, Name: "Alice"},
// 			expectedID:    1,
// 			expectedError: "",
// 		},
// 		{
// 			name: "Exec Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
// 				mockStmt.On("Close").Return(nil)
// 				mockErrorChecker.On("Check", mock.Anything).Return(errors.New("exec error"))
// 			},
// 			entity:        &CreateTestStruct{ID: 1, Name: "Alice"},
// 			expectedID:    0,
// 			expectedError: "exec error",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockDB := new(databasemock.MockDB)
// 			mockStmt := new(databasemock.MockStmt)
// 			mockResult := new(databasemock.MockResult)
// 			mockErrorChecker := new(databasemock.MockErrorChecker)

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockDB, mockStmt, mockResult, mockErrorChecker)
// 			}

// 			inserter := func(entity *CreateTestStruct) ([]string, []any) {
// 				return []string{"id", "name"}, []any{entity.ID, entity.Name}
// 			}

// 			id, err := Insert(
// 				tt.entity,
// 				mockDB,
// 				"user",
// 				inserter,
// 				mockErrorChecker,
// 			)

// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 			assert.Equal(t, tt.expectedID, id)

// 			mockDB.AssertExpectations(t)
// 			mockStmt.AssertExpectations(t)
// 			mockResult.AssertExpectations(t)
// 			mockErrorChecker.AssertExpectations(t)
// 		})
// 	}
// }

// // TestInsertMany tests the InsertMany function.
// func TestInsertMany(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		entities   []*CreateTestStruct
// 		setupMocks func(
// 			mockDB *databasemock.MockDB,
// 			mockStmt *databasemock.MockStmt,
// 			mockResult *databasemock.MockResult,
// 			mockErrorChecker *databasemock.MockErrorChecker,
// 		)
// 		expectedID    int64
// 		expectedError string
// 	}{
// 		{
// 			name: "Normal Operation",
// 			entities: []*CreateTestStruct{
// 				{ID: 1, Name: "Alice"},
// 				{ID: 2, Name: "Bob"},
// 			},
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
// 				mockStmt.On("Close").Return(nil)
// 				mockResult.On("LastInsertId").Return(int64(1), nil)
// 			},
// 			expectedID:    1,
// 			expectedError: "",
// 		},
// 		{
// 			name:          "Empty Entities",
// 			entities:      []*CreateTestStruct{},
// 			setupMocks:    nil,
// 			expectedID:    0,
// 			expectedError: "",
// 		},
// 		{
// 			name: "Exec Error",
// 			entities: []*CreateTestStruct{
// 				{ID: 1, Name: "Alice"},
// 				{ID: 2, Name: "Bob"},
// 			},
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
// 				mockStmt.On("Close").Return(nil)
// 				mockErrorChecker.On("Check", mock.Anything).Return(errors.New("exec error"))
// 			},
// 			expectedID:    0,
// 			expectedError: "exec error",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Arrange
// 			mockDB := new(databasemock.MockDB)
// 			mockStmt := new(databasemock.MockStmt)
// 			mockResult := new(databasemock.MockResult)
// 			mockErrorChecker := new(databasemock.MockErrorChecker)

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockDB, mockStmt, mockResult, mockErrorChecker)
// 			}

// 			// Define a sample Inserter function
// 			inserter := func(entity *CreateTestStruct) ([]string, []any) {
// 				if entity.ID == 1 {
// 					return []string{"id", "name"}, []any{1, "Alice"}
// 				}
// 				return []string{"id", "name"}, []any{2, "Bob"}
// 			}

// 			// Act
// 			id, err := InsertMany(
// 				tt.entities,
// 				mockDB,
// 				"user",
// 				inserter,
// 				mockErrorChecker,
// 			)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 			assert.Equal(t, tt.expectedID, id)

// 			// Verify mock expectations
// 			mockDB.AssertExpectations(t)
// 			mockStmt.AssertExpectations(t)
// 			mockResult.AssertExpectations(t)
// 			mockErrorChecker.AssertExpectations(t)
// 		})
// 	}
// }

// // TestCheckInsertResult tests the CheckInsertResult function.
// func TestCheckInsertResult(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setupMocks func(
// 			mockResult *MockSQLResult,
// 			mockErrorChecker *databasemock.MockErrorChecker,
// 		)
// 		inputErr      error
// 		expectedID    int64
// 		expectedError string
// 	}{
// 		{
// 			name: "No Error with valid LastInsertId",
// 			setupMocks: func(
// 				mockResult *MockSQLResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				// Mock the LastInsertId call to return 123 and no error
// 				mockResult.On("LastInsertId").Return(int64(123), nil)
// 			},
// 			inputErr:      nil,
// 			expectedID:    123,
// 			expectedError: "",
// 		},
// 		{
// 			name: "General Error from inputErr",
// 			setupMocks: func(
// 				mockResult *MockSQLResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				// When a general error is present, we rely on mockErrorChecker
// 				mockErrorChecker.On("Check", errors.New("some other error")).
// 					Return(errors.New("some other error"))
// 			},
// 			inputErr:      errors.New("some other error"),
// 			expectedID:    0,
// 			expectedError: "some other error",
// 		},
// 		{
// 			name: "LastInsertId Error",
// 			setupMocks: func(
// 				mockResult *MockSQLResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				// Mock the LastInsertId call to return an error
// 				mockResult.On("LastInsertId").
// 					Return(int64(0), errors.New("last insert ID error"))
// 			},
// 			inputErr:      nil,
// 			expectedID:    0,
// 			expectedError: "last insert ID error",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockResult := new(MockSQLResult)
// 			mockErrorChecker := new(databasemock.MockErrorChecker)

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockResult, mockErrorChecker)
// 			}

// 			id, err := checkGetResult(
// 				mockResult,
// 				tt.inputErr,
// 				mockErrorChecker,
// 			)

// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 			assert.Equal(t, tt.expectedID, id)

// 			// Verify that all the expectations on mocks were met
// 			mockResult.AssertExpectations(t)
// 			mockErrorChecker.AssertExpectations(t)
// 		})
// 	}
// }

// // TestCount tests the Count function
// func TestCount(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setupMocks func(
// 			mockDB *databasemock.MockDB,
// 			mockStmt *databasemock.MockStmt,
// 			mockRow *databasemock.MockRow,
// 		)
// 		tableName     string
// 		dbOptions     *CountOptions
// 		expectedCount int
// 		expectedError string
// 	}{
// 		{
// 			name: "Normal Operation",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRow *databasemock.MockRow,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Close").Return(nil)
// 				mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
// 				mockRow.On("Scan", mock.Anything).Return(nil)
// 			},
// 			tableName:     "test_table",
// 			dbOptions:     &CountOptions{},
// 			expectedCount: 0,
// 			expectedError: "",
// 		},
// 		{
// 			name: "Prepare Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRow *databasemock.MockRow,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).
// 					Return(nil, errors.New("prepare error"))
// 			},
// 			tableName:     "test_table",
// 			dbOptions:     &CountOptions{},
// 			expectedCount: 0,
// 			expectedError: "prepare error",
// 		},
// 		{
// 			name: "QueryRow Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRow *databasemock.MockRow,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Close").Return(nil)
// 				mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
// 				mockRow.On("Scan", mock.Anything).Return(errors.New("query row error"))
// 			},
// 			tableName:     "test_table",
// 			dbOptions:     &CountOptions{},
// 			expectedCount: 0,
// 			expectedError: "query row error",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Arrange
// 			mockDB := new(databasemock.MockDB)
// 			mockStmt := new(databasemock.MockStmt)
// 			mockRow := new(databasemock.MockRow)
// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockDB, mockStmt, mockRow)
// 			}

// 			// Act
// 			count, err := Count(mockDB, tt.tableName, tt.dbOptions)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 			assert.Equal(t, tt.expectedCount, count)

// 			// Verify mock expectations
// 			mockDB.AssertExpectations(t)
// 			mockStmt.AssertExpectations(t)
// 			mockRow.AssertExpectations(t)
// 		})
// 	}
// }

// type UpsertTestStruct struct {
// 	ID   int64
// 	Name string
// }

// // TestUpsert tests the Upsert function.
// func TestUpsert(t *testing.T) {
// 	// Use TestUpsertMany to cover functionality since Upsert is a wrapper
// 	// that delegates to UpsertMany.
// 	t.Run("Delegates to UpsertMany", func(t *testing.T) {
// 		mockDB := new(databasemock.MockDB)
// 		mockStmt := new(databasemock.MockStmt)
// 		mockResult := new(databasemock.MockResult)
// 		mockErrorChecker := new(databasemock.MockErrorChecker)

// 		mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 		mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
// 		mockStmt.On("Close").Return(nil)
// 		mockResult.On("LastInsertId").Return(int64(1), nil)

// 		inserter := func(entity *UpsertTestStruct) ([]string, []any) {
// 			return []string{"id", "name"}, []any{entity.ID, entity.Name}
// 		}

// 		id, err := Upsert(
// 			mockDB,
// 			"user_table",
// 			&UpsertTestStruct{ID: 1, Name: "Alice"},
// 			inserter,
// 			[]types.Projection{{Table: "name", Alias: "name_alias"}},
// 			mockErrorChecker,
// 		)

// 		assert.Nil(t, err)
// 		assert.Equal(t, int64(1), id)

// 		mockDB.AssertExpectations(t)
// 		mockStmt.AssertExpectations(t)
// 		mockResult.AssertExpectations(t)
// 		mockErrorChecker.AssertExpectations(t)
// 	})
// }

// // TestUpsertMany tests the UpsertMany function
// func TestUpsertMany(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setupMocks func(
// 			mockDB *databasemock.MockDB,
// 			mockStmt *databasemock.MockStmt,
// 			mockResult *databasemock.MockResult,
// 			mockErrorChecker *databasemock.MockErrorChecker,
// 		)
// 		entities          []*UpsertTestStruct
// 		updateProjections []types.Projection
// 		expectedID        int64
// 		expectedError     string
// 	}{
// 		{
// 			name: "Normal Operation",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
// 				mockStmt.On("Close").Return(nil)
// 				mockResult.On("LastInsertId").Return(int64(1), nil)
// 			},
// 			entities: []*UpsertTestStruct{
// 				{ID: 1, Name: "Alice"},
// 			},
// 			updateProjections: []types.Projection{
// 				{Table: "name", Alias: "name_alias"},
// 			},
// 			expectedID:    1,
// 			expectedError: "",
// 		},
// 		{
// 			name:       "No Entities",
// 			setupMocks: nil,
// 			entities:   []*UpsertTestStruct{},
// 			updateProjections: []types.Projection{
// 				{Table: "name", Alias: "name_alias"},
// 			},
// 			expectedID:    0,
// 			expectedError: "must provide entities to upsert",
// 		},
// 		{
// 			name:       "No Update Projections",
// 			setupMocks: nil,
// 			entities: []*UpsertTestStruct{
// 				{ID: 1, Name: "Alice"},
// 			},
// 			updateProjections: []types.Projection{},
// 			expectedID:        0,
// 			expectedError:     "must provide update projections",
// 		},
// 		{
// 			name:       "No Alias in Update Projections",
// 			setupMocks: nil,
// 			entities: []*UpsertTestStruct{
// 				{ID: 1, Name: "Alice"},
// 			},
// 			updateProjections: []types.Projection{
// 				{Table: "name", Alias: ""},
// 			},
// 			expectedID:    0,
// 			expectedError: "must provide update projections alias",
// 		},
// 		{
// 			name: "Exec Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
// 				mockStmt.On("Close").Return(nil)
// 				mockErrorChecker.On("Check", mock.Anything).Return(errors.New("exec error"))
// 			},
// 			entities: []*UpsertTestStruct{
// 				{ID: 1, Name: "Alice"},
// 			},
// 			updateProjections: []types.Projection{
// 				{Table: "name", Alias: "name_alias"},
// 			},
// 			expectedID:    0,
// 			expectedError: "exec error",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockDB := new(databasemock.MockDB)
// 			mockStmt := new(databasemock.MockStmt)
// 			mockResult := new(databasemock.MockResult)
// 			mockErrorChecker := new(databasemock.MockErrorChecker)

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockDB, mockStmt, mockResult, mockErrorChecker)
// 			}

// 			// Define inserter function
// 			inserter := func(entity *UpsertTestStruct) ([]string, []any) {
// 				return []string{"id", "name"}, []any{entity.ID, entity.Name}
// 			}

// 			// Act
// 			id, err := UpsertMany(
// 				mockDB,
// 				"user_table",
// 				tt.entities,
// 				inserter,
// 				tt.updateProjections,
// 				mockErrorChecker,
// 			)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 			assert.Equal(t, tt.expectedID, id)

// 			// Verify mock expectations
// 			mockDB.AssertExpectations(t)
// 			mockStmt.AssertExpectations(t)
// 			mockResult.AssertExpectations(t)
// 			mockErrorChecker.AssertExpectations(t)
// 		})
// 	}
// }

// // MockRowScanner is a mock implementation of RowScanner.
// type MockRowScanner[T any] struct {
// 	mock.Mock
// }

// func (m *MockRowScanner[T]) Scan(row Row, entity *T) error {
// 	return row.Scan(entity)
// }

// // MockRowScannerMultiple is a mock implementation of RowScannerMultiple.
// type MockRowScannerMultiple[T any] struct {
// 	mock.Mock
// }

// func (m *MockRowScannerMultiple[T]) Scan(rows Rows, entity *T) error {
// 	return rows.Scan(entity)
// }

// type GetTestStruct struct {
// 	ID   int
// 	Name string
// 	Age  int
// }

// // TestGet tests the Get function.
// func TestGet(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setupMocks func(
// 			mockDB *databasemock.MockDB,
// 			mockStmt *databasemock.MockStmt,
// 			mockRow *databasemock.MockRow,
// 			mockScanner *MockRowScanner[GetTestStruct],
// 		)
// 		tableName     string
// 		options       *GetOptions
// 		expectedError string
// 		expectEntity  bool
// 	}{
// 		{
// 			name: "Normal Operation",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRow *databasemock.MockRow,
// 				mockScanner *MockRowScanner[GetTestStruct],
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
// 				mockRow.On("Scan", mock.Anything).Return(nil).Once()
// 				mockRow.On("Err").Return(nil)
// 				mockStmt.On("Close").Return(nil)
// 			},
// 			tableName:     "user",
// 			options:       &GetOptions{},
// 			expectedError: "",
// 			expectEntity:  true,
// 		},
// 		{
// 			name: "Query Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRow *databasemock.MockRow,
// 				mockScanner *MockRowScanner[GetTestStruct],
// 			) {
// 				mockDB.On("Prepare", mock.Anything).
// 					Return(nil, errors.New("query error"))
// 			},
// 			tableName:     "user",
// 			options:       &GetOptions{},
// 			expectedError: "query error",
// 			expectEntity:  false,
// 		},
// 		{
// 			name: "No Rows",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRow *databasemock.MockRow,
// 				mockScanner *MockRowScanner[GetTestStruct],
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
// 				mockRow.On("Scan", mock.Anything).Return(nil).Once()
// 				mockRow.On("Err").Return(sql.ErrNoRows).Once()
// 				mockStmt.On("Close").Return(nil)
// 			},
// 			tableName:     "user",
// 			options:       &GetOptions{},
// 			expectedError: "",
// 			expectEntity:  false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockDB := new(databasemock.MockDB)
// 			mockStmt := new(databasemock.MockStmt)
// 			mockRow := new(databasemock.MockRow)
// 			mockScanner := new(MockRowScanner[GetTestStruct])

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockDB, mockStmt, mockRow, mockScanner)
// 			}

// 			// Act
// 			entity, err := Get(
// 				tt.tableName,
// 				mockScanner.Scan,
// 				mockDB,
// 				tt.options,
// 			)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 				assert.Nil(t, entity)
// 			} else {
// 				assert.Nil(t, err)
// 				if tt.expectEntity {
// 					assert.NotNil(t, entity)
// 				} else {
// 					assert.Nil(t, entity)
// 				}
// 			}

// 			// Verify mock expectations
// 			mockDB.AssertExpectations(t)
// 			mockStmt.AssertExpectations(t)
// 			mockRow.AssertExpectations(t)
// 			mockScanner.AssertExpectations(t)
// 		})
// 	}
// }

// // TestGetMany tests the GetMany function
// func TestGetMany(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setupMocks func(
// 			mockDB *databasemock.MockDB,
// 			mockStmt *databasemock.MockStmt,
// 			mockRows *databasemock.MockRows,
// 			mockScanner *MockRowScannerMultiple[GetTestStruct],
// 		)
// 		tableName      string
// 		options        *GetOptions
// 		expectedError  string
// 		expectEntities int
// 	}{
// 		{
// 			name: "Normal Operation",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRows *databasemock.MockRows,
// 				mockScanner *MockRowScannerMultiple[GetTestStruct],
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Query", mock.Anything).Return(mockRows, nil)
// 				mockStmt.On("Close").Return(nil)
// 				mockRows.On("Close").Return(nil)
// 				mockRows.On("Next").Return(true).Once()
// 				mockRows.On("Scan", mock.Anything).Return(nil).Once()
// 				mockRows.On("Next").Return(false).Once()
// 				mockRows.On("Err").Return(nil)
// 			},
// 			tableName:      "user",
// 			options:        &GetOptions{},
// 			expectedError:  "",
// 			expectEntities: 1,
// 		},
// 		{
// 			name: "Query Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRows *databasemock.MockRows,
// 				mockScanner *MockRowScannerMultiple[GetTestStruct],
// 			) {
// 				mockDB.On("Prepare", mock.Anything).
// 					Return(nil, errors.New("query error"))
// 			},
// 			tableName:      "user",
// 			options:        &GetOptions{},
// 			expectedError:  "query error",
// 			expectEntities: 0,
// 		},
// 		{
// 			name: "No Rows",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRows *databasemock.MockRows,
// 				mockScanner *MockRowScannerMultiple[GetTestStruct],
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Query", mock.Anything).Return(mockRows, nil)
// 				mockStmt.On("Close").Return(nil)
// 				mockRows.On("Close").Return(nil)
// 				mockRows.On("Next").Return(false).Once()
// 				mockRows.On("Err").Return(sql.ErrNoRows).Once()
// 			},
// 			tableName:      "user",
// 			options:        &GetOptions{},
// 			expectedError:  "",
// 			expectEntities: 0,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockDB := new(databasemock.MockDB)
// 			mockStmt := new(databasemock.MockStmt)
// 			mockRows := new(databasemock.MockRows)
// 			mockScanner := new(MockRowScannerMultiple[GetTestStruct])

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockDB, mockStmt, mockRows, mockScanner)
// 			}

// 			// Act
// 			entities, err := GetMany(
// 				tt.tableName,
// 				mockScanner.Scan,
// 				mockDB,
// 				tt.options,
// 			)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 				assert.Nil(t, entities)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Len(t, entities, tt.expectEntities)
// 			}

// 			// Verify mock expectations
// 			mockDB.AssertExpectations(t)
// 			mockStmt.AssertExpectations(t)
// 			mockRows.AssertExpectations(t)
// 			mockScanner.AssertExpectations(t)
// 		})
// 	}
// }

// // TestQueryMultiple tests the QueryMultiple function
// func TestQueryMultiple(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setupMocks func(
// 			mockDB *databasemock.MockDB,
// 			mockStmt *databasemock.MockStmt,
// 			mockRows *databasemock.MockRows,
// 			mockScanner *MockRowScannerMultiple[GetTestStruct],
// 		)
// 		query          string
// 		params         []any
// 		expectedError  string
// 		expectEntities int
// 	}{
// 		{
// 			name: "Normal Operation",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRows *databasemock.MockRows,
// 				mockScanner *MockRowScannerMultiple[GetTestStruct],
// 			) {
// 				mockDB.On("Prepare", "SELECT * FROM user WHERE active = ?").Return(mockStmt, nil)
// 				mockStmt.On("Query", []any{1}).Return(mockRows, nil)
// 				mockStmt.On("Close").Return(nil)
// 				mockRows.On("Close").Return(nil)
// 				mockRows.On("Next").Return(true).Once()
// 				mockRows.On("Scan", mock.Anything).Return(nil).Once()
// 				mockRows.On("Next").Return(false).Once()
// 				mockRows.On("Err").Return(nil)
// 			},
// 			query:          "SELECT * FROM user WHERE active = ?",
// 			params:         []any{1},
// 			expectedError:  "",
// 			expectEntities: 1,
// 		},
// 		{
// 			name: "Query Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRows *databasemock.MockRows,
// 				mockScanner *MockRowScannerMultiple[GetTestStruct],
// 			) {
// 				mockDB.On("Prepare", "SELECT * FROM user WHERE active = ?").
// 					Return(nil, errors.New("query error"))
// 			},
// 			query:          "SELECT * FROM user WHERE active = ?",
// 			params:         []any{1},
// 			expectedError:  "query error",
// 			expectEntities: 0,
// 		},
// 		{
// 			name: "Rows Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRows *databasemock.MockRows,
// 				mockScanner *MockRowScannerMultiple[GetTestStruct],
// 			) {
// 				mockDB.On("Prepare", "SELECT * FROM user WHERE active = ?").Return(mockStmt, nil)
// 				mockStmt.On("Query", []any{1}).Return(mockRows, nil)
// 				mockStmt.On("Close").Return(nil)
// 				mockRows.On("Close").Return(nil)
// 				mockRows.On("Next").Return(true).Once()
// 				mockRows.On("Scan", mock.Anything).Return(nil).Once()
// 				mockRows.On("Next").Return(false).Once()
// 				mockRows.On("Err").Return(errors.New("rows error"))
// 			},
// 			query:          "SELECT * FROM user WHERE active = ?",
// 			params:         []any{1},
// 			expectedError:  "rows error",
// 			expectEntities: 0,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockDB := new(databasemock.MockDB)
// 			mockStmt := new(databasemock.MockStmt)
// 			mockRows := new(databasemock.MockRows)
// 			mockScanner := new(MockRowScannerMultiple[GetTestStruct])

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockDB, mockStmt, mockRows, mockScanner)
// 			}

// 			// Act
// 			entities, err := queryMultiple(
// 				mockDB,
// 				tt.query,
// 				tt.params,
// 				mockScanner.Scan,
// 			)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 				assert.Nil(t, entities)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Len(t, entities, tt.expectEntities)
// 			}

// 			// Verify mock expectations
// 			mockDB.AssertExpectations(t)
// 			mockStmt.AssertExpectations(t)
// 			mockRows.AssertExpectations(t)
// 			mockScanner.AssertExpectations(t)
// 		})
// 	}
// }

// // TestQuerySingle tests the QuerySingle function
// func TestQuerySingle(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setupMocks func(
// 			mockDB *databasemock.MockDB,
// 			mockStmt *databasemock.MockStmt,
// 			mockRow *databasemock.MockRow,
// 			mockScanner *MockRowScanner[GetTestStruct],
// 		)
// 		query         string
// 		params        []any
// 		expectedError string
// 		expectEntity  bool
// 	}{
// 		{
// 			name: "Normal Operation",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRow *databasemock.MockRow,
// 				mockScanner *MockRowScanner[GetTestStruct],
// 			) {
// 				mockDB.On("Prepare", "SELECT * FROM user WHERE id = ?").Return(mockStmt, nil)
// 				mockStmt.On("QueryRow", []any{1}).Return(mockRow)
// 				mockRow.On("Scan", mock.Anything).Return(nil).Once()
// 				mockStmt.On("Close").Return(nil)
// 				mockRow.On("Err").Return(nil)
// 			},
// 			query:         "SELECT * FROM user WHERE id = ?",
// 			params:        []any{1},
// 			expectedError: "",
// 			expectEntity:  true,
// 		},
// 		{
// 			name: "Prepare Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRow *databasemock.MockRow,
// 				mockScanner *MockRowScanner[GetTestStruct],
// 			) {
// 				mockDB.On("Prepare", "SELECT * FROM user WHERE id = ?").
// 					Return(nil, errors.New("prepare error"))
// 			},
// 			query:         "SELECT * FROM user WHERE id = ?",
// 			params:        []any{1},
// 			expectedError: "prepare error",
// 			expectEntity:  false,
// 		},
// 		{
// 			name: "Row Scanner Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRow *databasemock.MockRow,
// 				mockScanner *MockRowScanner[GetTestStruct],
// 			) {
// 				mockDB.On("Prepare", "SELECT * FROM user WHERE id = ?").Return(mockStmt, nil)
// 				mockStmt.On("QueryRow", []any{1}).Return(mockRow)
// 				mockRow.On("Scan", mock.Anything).Return(errors.New("row scanner error")).Once()
// 				mockStmt.On("Close").Return(nil)
// 			},
// 			query:         "SELECT * FROM user WHERE id = ?",
// 			params:        []any{1},
// 			expectedError: "row scanner error",
// 			expectEntity:  false,
// 		},
// 		{
// 			name: "Row Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockRow *databasemock.MockRow,
// 				mockScanner *MockRowScanner[GetTestStruct],
// 			) {
// 				mockDB.On("Prepare", "SELECT * FROM user WHERE id = ?").Return(mockStmt, nil)
// 				mockStmt.On("QueryRow", []any{1}).Return(mockRow)
// 				mockRow.On("Scan", mock.Anything).Return(nil).Once()
// 				mockStmt.On("Close").Return(nil)
// 				mockRow.On("Err").Return(errors.New("row error"))
// 			},
// 			query:         "SELECT * FROM user WHERE id = ?",
// 			params:        []any{1},
// 			expectedError: "row error",
// 			expectEntity:  false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockDB := new(databasemock.MockDB)
// 			mockStmt := new(databasemock.MockStmt)
// 			mockRow := new(databasemock.MockRow)
// 			mockScanner := new(MockRowScanner[GetTestStruct])

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockDB, mockStmt, mockRow, mockScanner)
// 			}

// 			// Act
// 			entity, err := querySingle(
// 				mockDB,
// 				tt.query,
// 				tt.params,
// 				mockScanner.Scan,
// 			)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 				assert.Nil(t, entity)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.NotNil(t, entity)
// 			}

// 			// Verify mock expectations
// 			mockDB.AssertExpectations(t)
// 			mockStmt.AssertExpectations(t)
// 			mockRow.AssertExpectations(t)
// 			mockScanner.AssertExpectations(t)
// 		})
// 	}
// }

// // TestRowsToEntities tests the RowsToEntities function
// func TestRowsToEntities(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setupMocks func(
// 			mockRows *databasemock.MockRows,
// 			mockScanner *MockRowScannerMultiple[GetTestStruct],
// 		)
// 		rowScanner     RowScannerMultiple[GetTestStruct]
// 		expectedError  string
// 		expectedLength int
// 	}{
// 		{
// 			name: "No RowScannerMultiple",
// 			setupMocks: func(
// 				mockRows *databasemock.MockRows,
// 				mockScanner *MockRowScannerMultiple[GetTestStruct],
// 			) {
// 				// No mocks needed for this test case
// 			},
// 			rowScanner:     nil,
// 			expectedError:  "must provide rowScannerMultiple",
// 			expectedLength: 0,
// 		},
// 		{
// 			name: "Normal Operation",
// 			setupMocks: func(
// 				mockRows *databasemock.MockRows,
// 				mockScanner *MockRowScannerMultiple[GetTestStruct],
// 			) {
// 				testEntity := GetTestStruct{}
// 				mockRows.On("Next").Return(true).Once() // Simulate the first row read
// 				mockRows.On("Scan", []any{&testEntity}).Return(nil).Once()
// 				mockRows.On("Next").Return(true).Once() // Simulate the second row read
// 				mockRows.On("Scan", []any{&testEntity}).Return(nil).Once()
// 				mockRows.On("Next").Return(false).Once() // No more rows
// 				mockRows.On("Err").Return(nil)
// 			},
// 			rowScanner:     new(MockRowScannerMultiple[GetTestStruct]).Scan,
// 			expectedError:  "",
// 			expectedLength: 2,
// 		},
// 		{
// 			name: "No Rows",
// 			setupMocks: func(
// 				mockRows *databasemock.MockRows,
// 				mockScanner *MockRowScannerMultiple[GetTestStruct],
// 			) {
// 				mockRows.On("Next").Return(false).Once() // No rows
// 				mockRows.On("Err").Return(nil)
// 			},
// 			rowScanner:     new(MockRowScannerMultiple[GetTestStruct]).Scan,
// 			expectedError:  "",
// 			expectedLength: 0,
// 		},
// 		{
// 			name: "Row Scanner Error",
// 			setupMocks: func(
// 				mockRows *databasemock.MockRows,
// 				mockScanner *MockRowScannerMultiple[GetTestStruct],
// 			) {
// 				mockRows.On("Next").Return(true).Once() // Simulate the first row read
// 				mockRows.On("Scan", []any{&GetTestStruct{}}).
// 					Return(errors.New("row scanner error")).Once()
// 			},
// 			rowScanner:     new(MockRowScannerMultiple[GetTestStruct]).Scan,
// 			expectedError:  "row scanner error",
// 			expectedLength: 0,
// 		},
// 		{
// 			name: "Rows Error",
// 			setupMocks: func(
// 				mockRows *databasemock.MockRows,
// 				mockScanner *MockRowScannerMultiple[GetTestStruct],
// 			) {
// 				mockRows.On("Next").Return(true).Once() // Simulate the first row read
// 				mockRows.On("Scan", []any{&GetTestStruct{}}).Return(nil).Once()
// 				mockRows.On("Next").Return(false).Once() // No more rows
// 				mockRows.On("Err").Return(errors.New("rows error")).Once()
// 			},
// 			rowScanner:     new(MockRowScannerMultiple[GetTestStruct]).Scan,
// 			expectedError:  "rows error",
// 			expectedLength: 0,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockRows := new(databasemock.MockRows)
// 			mockScanner := new(MockRowScannerMultiple[GetTestStruct])

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockRows, mockScanner)
// 			}

// 			// Act
// 			entities, err := rowsToEntities(mockRows, tt.rowScanner)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 				assert.Nil(t, entities)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Len(t, entities, tt.expectedLength)
// 			}

// 			// Verify mock expectations
// 			mockRows.AssertExpectations(t)
// 			mockScanner.AssertExpectations(t)
// 		})
// 	}
// }

// // TestUpdate tests the Update function
// func TestUpdate(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setupMocks func(
// 			mockDB *databasemock.MockDB,
// 			mockStmt *databasemock.MockStmt,
// 			mockResult *databasemock.MockResult,
// 			mockErrorChecker *databasemock.MockErrorChecker,
// 		)
// 		updates  []types.Update
// 		selectors     []types.Selector
// 		expectedRows  int64
// 		expectedError string
// 	}{
// 		{
// 			name: "Normal Operation",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
// 				mockStmt.On("Close").Return(nil)
// 				mockResult.On("RowsAffected").Return(int64(2), nil)
// 			},
// 			updates: []types.Update{
// 				{Field: "name", Value: "Alice"},
// 			},
// 			selectors: []types.Selector{
// 				{Field: "id", Value: 1},
// 			},
// 			expectedRows:  2,
// 			expectedError: "",
// 		},
// 		{
// 			name:          "No Updates",
// 			setupMocks:    nil,
// 			updates:  []types.Update{},
// 			selectors:     []types.Selector{{Field: "id", Value: 1}},
// 			expectedRows:  0,
// 			expectedError: "",
// 		},
// 		{
// 			name: "Exec Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
// 				mockStmt.On("Close").Return(nil)
// 				mockErrorChecker.On("Check", mock.Anything).Return(errors.New("exec error"))
// 			},
// 			updates: []types.Update{
// 				{Field: "name", Value: "Alice"},
// 			},
// 			selectors: []types.Selector{
// 				{Field: "id", Value: 1},
// 			},
// 			expectedRows:  0,
// 			expectedError: "exec error",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockDB := new(databasemock.MockDB)
// 			mockStmt := new(databasemock.MockStmt)
// 			mockResult := new(databasemock.MockResult)
// 			mockErrorChecker := new(databasemock.MockErrorChecker)

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockDB, mockStmt, mockResult, mockErrorChecker)
// 			}

// 			// Act
// 			rows, err := Update(
// 				mockDB,
// 				"user_table",
// 				tt.selectors,
// 				tt.updates,
// 				mockErrorChecker,
// 			)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 			assert.Equal(t, tt.expectedRows, rows)

// 			// Verify mock expectations
// 			mockDB.AssertExpectations(t)
// 			mockStmt.AssertExpectations(t)
// 			mockResult.AssertExpectations(t)
// 			mockErrorChecker.AssertExpectations(t)
// 		})
// 	}
// }

// // TestCheckUpdateResult tests the CheckUpdateResult function
// func TestCheckUpdateResult(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setupMocks func(
// 			mockResult *databasemock.MockResult,
// 			mockErrorChecker *databasemock.MockErrorChecker,
// 		)
// 		inputErr      error
// 		expectedRows  int64
// 		expectedError string
// 	}{
// 		{
// 			name: "Normal Operation",
// 			setupMocks: func(
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockResult.On("RowsAffected").Return(int64(3), nil)
// 			},
// 			inputErr:      nil,
// 			expectedRows:  3,
// 			expectedError: "",
// 		},
// 		{
// 			name: "Exec Error",
// 			setupMocks: func(
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockErrorChecker.On("Check", errors.New("exec error")).
// 					Return(errors.New("exec error"))
// 			},
// 			inputErr:      errors.New("exec error"),
// 			expectedRows:  0,
// 			expectedError: "exec error",
// 		},
// 		{
// 			name: "RowsAffected Error",
// 			setupMocks: func(
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockResult.On("RowsAffected").
// 					Return(int64(0), errors.New("rows affected error"))
// 			},
// 			inputErr:      nil,
// 			expectedRows:  0,
// 			expectedError: "rows affected error",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockResult := new(databasemock.MockResult)
// 			mockErrorChecker := new(databasemock.MockErrorChecker)

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockResult, mockErrorChecker)
// 			}

// 			// Act
// 			rows, err := checkUpdateResult(
// 				mockResult,
// 				tt.inputErr,
// 				mockErrorChecker,
// 			)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 			assert.Equal(t, tt.expectedRows, rows)

// 			// Verify mock expectations
// 			mockResult.AssertExpectations(t)
// 			mockErrorChecker.AssertExpectations(t)
// 		})
// 	}
// }

// // TestDelete tests the Delete function
// func TestDelete(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setupMocks func(
// 			mockDB *databasemock.MockDB,
// 			mockStmt *databasemock.MockStmt,
// 			mockResult *databasemock.MockResult,
// 			mockErrorChecker *databasemock.MockErrorChecker,
// 		)
// 		selectors     []types.Selector
// 		opts          *DeleteOptions
// 		expectedCount int64
// 		expectedError string
// 	}{
// 		{
// 			name: "Normal Operation",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
// 				mockStmt.On("Close").Return(nil)
// 				mockResult.On("RowsAffected").Return(int64(3), nil)
// 			},
// 			selectors:     []types.Selector{{Field: "id", Value: 1}},
// 			opts:          nil,
// 			expectedCount: 3,
// 			expectedError: "",
// 		},
// 		{
// 			name: "Exec Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
// 				mockStmt.On("Close").Return(nil)
// 			},
// 			selectors:     []types.Selector{{Field: "id", Value: 1}},
// 			opts:          nil,
// 			expectedCount: 0,
// 			expectedError: "exec error",
// 		},
// 		{
// 			name: "RowsAffected Error",
// 			setupMocks: func(
// 				mockDB *databasemock.MockDB,
// 				mockStmt *databasemock.MockStmt,
// 				mockResult *databasemock.MockResult,
// 				mockErrorChecker *databasemock.MockErrorChecker,
// 			) {
// 				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
// 				mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
// 				mockStmt.On("Close").Return(nil)
// 				mockResult.On("RowsAffected").Return(int64(0), errors.New("rows affected error"))
// 			},
// 			selectors:     []types.Selector{{Field: "id", Value: 1}},
// 			opts:          nil,
// 			expectedCount: 0,
// 			expectedError: "rows affected error",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockDB := new(databasemock.MockDB)
// 			mockStmt := new(databasemock.MockStmt)
// 			mockResult := new(databasemock.MockResult)
// 			mockErrorChecker := new(databasemock.MockErrorChecker)

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockDB, mockStmt, mockResult, mockErrorChecker)
// 			}

// 			// Act
// 			count, err := Delete(
// 				mockDB, // implements Preparer
// 				"user_table",
// 				tt.selectors,
// 				tt.opts,
// 			)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 			assert.Equal(t, tt.expectedCount, count)

// 			// Verify mock expectations
// 			mockDB.AssertExpectations(t)
// 			mockResult.AssertExpectations(t)
// 		})
// 	}
// }

// // TestExec tests the Exec method.
// func TestExec(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		setupMocks    func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockResult *MockSQLResult)
// 		query         string
// 		parameters    []any
// 		expectedError string
// 		expectResult  bool
// 	}{
// 		{
// 			name: "Normal Operation",
// 			setupMocks: func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockResult *MockSQLResult) {
// 				mockDB.On("Prepare", "UPDATE users SET name = ? WHERE id = ?").Return(mockStmt, nil)
// 				mockStmt.On("Exec", []any{"Alice", 1}).Return(mockResult, nil)
// 				mockStmt.On("Close").Return(nil)
// 			},
// 			query:         "UPDATE users SET name = ? WHERE id = ?",
// 			parameters:    []any{"Alice", 1},
// 			expectedError: "",
// 			expectResult:  true,
// 		},
// 		{
// 			name: "Prepare Error",
// 			setupMocks: func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockResult *MockSQLResult) {
// 				mockDB.On("Prepare", "UPDATE users SET name = ? WHERE id = ?").Return(nil, errors.New("prepare error"))
// 			},
// 			query:         "UPDATE users SET name = ? WHERE id = ?",
// 			parameters:    []any{"Alice", 1},
// 			expectedError: "prepare error",
// 			expectResult:  false,
// 		},
// 		{
// 			name: "Exec Error",
// 			setupMocks: func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockResult *MockSQLResult) {
// 				mockDB.On("Prepare", "UPDATE users SET name = ? WHERE id = ?").Return(mockStmt, nil)
// 				mockStmt.On("Exec", []any{"Alice", 1}).Return(nil, errors.New("exec error"))
// 				mockStmt.On("Close").Return(nil)
// 			},
// 			query:         "UPDATE users SET name = ? WHERE id = ?",
// 			parameters:    []any{"Alice", 1},
// 			expectedError: "exec error",
// 			expectResult:  false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockDB := new(databasemock.MockDB)
// 			mockStmt := new(databasemock.MockStmt)
// 			mockResult := new(MockSQLResult)

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockDB, mockStmt, mockResult)
// 			}

// 			// Act
// 			result, err := Exec(mockDB, tt.query, tt.parameters)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 				assert.Nil(t, result)
// 			} else {
// 				assert.Nil(t, err)
// 				if tt.expectResult {
// 					assert.NotNil(t, result)
// 				} else {
// 					assert.Nil(t, result)
// 				}
// 			}

// 			// Verify mock expectations
// 			mockDB.AssertExpectations(t)
// 			mockStmt.AssertExpectations(t)
// 			mockResult.AssertExpectations(t)
// 		})
// 	}
// }

// // TestQuery tests the Query function
// func TestQuery(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		setupMocks func(
// 			mockDB *databasemock.MockDB,
// 			mockStmt *databasemock.MockStmt,
// 			mockRows *databasemock.MockRows,
// 		)
// 		query         string
// 		parameters    []any
// 		expectedError string
// 		expectRows    bool
// 		expectStmt    bool
// 	}{
// 		{
// 			name: "Normal Operation",
// 			setupMocks: func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockRows *databasemock.MockRows) {
// 				mockDB.On("Prepare", "SELECT * FROM users WHERE id = ?").Return(mockStmt, nil)
// 				mockStmt.On("Query", []any{1}).Return(mockRows, nil)
// 			},
// 			query:         "SELECT * FROM users WHERE id = ?",
// 			parameters:    []any{1},
// 			expectedError: "",
// 			expectRows:    true,
// 			expectStmt:    true,
// 		},
// 		{
// 			name: "Prepare Error",
// 			setupMocks: func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockRows *databasemock.MockRows) {
// 				mockDB.On("Prepare", "SELECT * FROM users WHERE id = ?").Return(nil, errors.New("prepare error"))
// 			},
// 			query:         "SELECT * FROM users WHERE id = ?",
// 			parameters:    []any{1},
// 			expectedError: "prepare error",
// 			expectRows:    false,
// 			expectStmt:    false,
// 		},
// 		{
// 			name: "Query Error",
// 			setupMocks: func(mockDB *databasemock.MockDB, mockStmt *databasemock.MockStmt, mockRows *databasemock.MockRows) {
// 				mockDB.On("Prepare", "SELECT * FROM users WHERE id = ?").Return(mockStmt, nil)
// 				mockStmt.On("Query", []any{1}).Return(nil, errors.New("query error"))
// 				mockStmt.On("Close").Return(nil)
// 			},
// 			query:         "SELECT * FROM users WHERE id = ?",
// 			parameters:    []any{1},
// 			expectedError: "query error",
// 			expectRows:    false,
// 			expectStmt:    false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockDB := new(databasemock.MockDB)
// 			mockStmt := new(databasemock.MockStmt)
// 			mockRows := new(databasemock.MockRows)

// 			if tt.setupMocks != nil {
// 				tt.setupMocks(mockDB, mockStmt, mockRows)
// 			}

// 			// Act
// 			rows, stmt, err := Query(mockDB, tt.query, tt.parameters)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.EqualError(t, err, tt.expectedError)
// 				assert.Nil(t, rows)
// 				assert.Nil(t, stmt)
// 			} else {
// 				assert.Nil(t, err)
// 				if tt.expectRows {
// 					assert.NotNil(t, rows)
// 				} else {
// 					assert.Nil(t, rows)
// 				}
// 				if tt.expectStmt {
// 					assert.NotNil(t, stmt)
// 				} else {
// 					assert.Nil(t, stmt)
// 				}
// 			}

// 			// Verify mock expectations
// 			mockDB.AssertExpectations(t)
// 			mockStmt.AssertExpectations(t)
// 			mockRows.AssertExpectations(t)
// 		})
// 	}
// }
