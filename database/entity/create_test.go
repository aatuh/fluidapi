package entity

import (
	"errors"
	"testing"

	entitymock "github.com/pakkasys/fluidapi/database/entity/mock"
	utilmock "github.com/pakkasys/fluidapi/database/util/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSQLResult is a mock implementation of the sql.Result interface.
type MockSQLResult struct {
	mock.Mock
}

func (m *MockSQLResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSQLResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

type CreateTestStruct struct {
	ID   int
	Name string
	Age  int
}

// TestInsert tests the Insert function.
func TestInsert(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(
			mockDB *utilmock.MockDB,
			mockStmt *utilmock.MockStmt,
			mockResult *utilmock.MockResult,
			mockErrorChecker *entitymock.MockErrorChecker,
		)
		entity        *CreateTestStruct
		expectedID    int64
		expectedError string
	}{
		{
			name: "Normal Operation",
			setupMocks: func(
				mockDB *utilmock.MockDB,
				mockStmt *utilmock.MockStmt,
				mockResult *utilmock.MockResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
				mockStmt.On("Close").Return(nil)
				mockResult.On("LastInsertId").Return(int64(1), nil)
			},
			entity:        &CreateTestStruct{ID: 1, Name: "Alice"},
			expectedID:    1,
			expectedError: "",
		},
		{
			name: "Exec Error",
			setupMocks: func(
				mockDB *utilmock.MockDB,
				mockStmt *utilmock.MockStmt,
				mockResult *utilmock.MockResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
				mockStmt.On("Close").Return(nil)
				mockErrorChecker.On("Check", mock.Anything).Return(errors.New("exec error"))
			},
			entity:        &CreateTestStruct{ID: 1, Name: "Alice"},
			expectedID:    0,
			expectedError: "exec error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(utilmock.MockDB)
			mockStmt := new(utilmock.MockStmt)
			mockResult := new(utilmock.MockResult)
			mockErrorChecker := new(entitymock.MockErrorChecker)

			if tt.setupMocks != nil {
				tt.setupMocks(mockDB, mockStmt, mockResult, mockErrorChecker)
			}

			inserter := func(entity *CreateTestStruct) ([]string, []any) {
				return []string{"id", "name"}, []any{entity.ID, entity.Name}
			}

			id, err := Insert(
				tt.entity,
				mockDB,
				"user",
				inserter,
				mockErrorChecker,
			)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedID, id)

			mockDB.AssertExpectations(t)
			mockStmt.AssertExpectations(t)
			mockResult.AssertExpectations(t)
			mockErrorChecker.AssertExpectations(t)
		})
	}
}

// TestInsertMany tests the InsertMany function.
func TestInsertMany(t *testing.T) {
	tests := []struct {
		name       string
		entities   []*CreateTestStruct
		setupMocks func(
			mockDB *utilmock.MockDB,
			mockStmt *utilmock.MockStmt,
			mockResult *utilmock.MockResult,
			mockErrorChecker *entitymock.MockErrorChecker,
		)
		expectedID    int64
		expectedError string
	}{
		{
			name: "Normal Operation",
			entities: []*CreateTestStruct{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			setupMocks: func(
				mockDB *utilmock.MockDB,
				mockStmt *utilmock.MockStmt,
				mockResult *utilmock.MockResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
				mockStmt.On("Close").Return(nil)
				mockResult.On("LastInsertId").Return(int64(1), nil)
			},
			expectedID:    1,
			expectedError: "",
		},
		{
			name:          "Empty Entities",
			entities:      []*CreateTestStruct{},
			setupMocks:    nil,
			expectedID:    0,
			expectedError: "",
		},
		{
			name: "Exec Error",
			entities: []*CreateTestStruct{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			setupMocks: func(
				mockDB *utilmock.MockDB,
				mockStmt *utilmock.MockStmt,
				mockResult *utilmock.MockResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
				mockStmt.On("Close").Return(nil)
				mockErrorChecker.On("Check", mock.Anything).Return(errors.New("exec error"))
			},
			expectedID:    0,
			expectedError: "exec error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDB := new(utilmock.MockDB)
			mockStmt := new(utilmock.MockStmt)
			mockResult := new(utilmock.MockResult)
			mockErrorChecker := new(entitymock.MockErrorChecker)

			if tt.setupMocks != nil {
				tt.setupMocks(mockDB, mockStmt, mockResult, mockErrorChecker)
			}

			// Define a sample Inserter function
			inserter := func(entity *CreateTestStruct) ([]string, []any) {
				if entity.ID == 1 {
					return []string{"id", "name"}, []any{1, "Alice"}
				}
				return []string{"id", "name"}, []any{2, "Bob"}
			}

			// Act
			id, err := InsertMany(
				tt.entities,
				mockDB,
				"user",
				inserter,
				mockErrorChecker,
			)

			// Assert
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedID, id)

			// Verify mock expectations
			mockDB.AssertExpectations(t)
			mockStmt.AssertExpectations(t)
			mockResult.AssertExpectations(t)
			mockErrorChecker.AssertExpectations(t)
		})
	}
}

// TestCheckInsertResult tests the CheckInsertResult function.
func TestCheckInsertResult(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(
			mockResult *MockSQLResult,
			mockErrorChecker *entitymock.MockErrorChecker,
		)
		inputErr      error
		expectedID    int64
		expectedError string
	}{
		{
			name: "No Error with valid LastInsertId",
			setupMocks: func(
				mockResult *MockSQLResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				// Mock the LastInsertId call to return 123 and no error
				mockResult.On("LastInsertId").Return(int64(123), nil)
			},
			inputErr:      nil,
			expectedID:    123,
			expectedError: "",
		},
		{
			name: "General Error from inputErr",
			setupMocks: func(
				mockResult *MockSQLResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				// When a general error is present, we rely on mockErrorChecker
				mockErrorChecker.On("Check", errors.New("some other error")).
					Return(errors.New("some other error"))
			},
			inputErr:      errors.New("some other error"),
			expectedID:    0,
			expectedError: "some other error",
		},
		{
			name: "LastInsertId Error",
			setupMocks: func(
				mockResult *MockSQLResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				// Mock the LastInsertId call to return an error
				mockResult.On("LastInsertId").
					Return(int64(0), errors.New("last insert ID error"))
			},
			inputErr:      nil,
			expectedID:    0,
			expectedError: "last insert ID error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockResult := new(MockSQLResult)
			mockErrorChecker := new(entitymock.MockErrorChecker)

			if tt.setupMocks != nil {
				tt.setupMocks(mockResult, mockErrorChecker)
			}

			id, err := checkInsertResult(
				mockResult,
				tt.inputErr,
				mockErrorChecker,
			)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedID, id)

			// Verify that all the expectations on mocks were met
			mockResult.AssertExpectations(t)
			mockErrorChecker.AssertExpectations(t)
		})
	}
}
