package entity

import (
	"errors"
	"testing"

	"github.com/pakkasys/fluidapi/database/clause"
	entitymock "github.com/pakkasys/fluidapi/database/entity/mock"
	databasemock "github.com/pakkasys/fluidapi/database/mock"
	"github.com/pakkasys/fluidapi/database/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestDelete tests the Delete function
func TestDelete(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(
			mockDB *databasemock.MockDB,
			mockStmt *databasemock.MockStmt,
			mockResult *databasemock.MockResult,
			mockErrorChecker *entitymock.MockErrorChecker,
		)
		selectors     []clause.Selector
		opts          *query.DeleteOptions
		expectedCount int64
		expectedError string
	}{
		{
			name: "Normal Operation",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockResult *databasemock.MockResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
				mockStmt.On("Close").Return(nil)
				mockResult.On("RowsAffected").Return(int64(3), nil)
			},
			selectors:     []clause.Selector{{Field: "id", Value: 1}},
			opts:          nil,
			expectedCount: 3,
			expectedError: "",
		},
		{
			name: "Exec Error",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockResult *databasemock.MockResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("Exec", mock.Anything).Return(nil, errors.New("exec error"))
				mockStmt.On("Close").Return(nil)
			},
			selectors:     []clause.Selector{{Field: "id", Value: 1}},
			opts:          nil,
			expectedCount: 0,
			expectedError: "exec error",
		},
		{
			name: "RowsAffected Error",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockResult *databasemock.MockResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("Exec", mock.Anything).Return(mockResult, nil)
				mockStmt.On("Close").Return(nil)
				mockResult.On("RowsAffected").Return(int64(0), errors.New("rows affected error"))
			},
			selectors:     []clause.Selector{{Field: "id", Value: 1}},
			opts:          nil,
			expectedCount: 0,
			expectedError: "rows affected error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(databasemock.MockDB)
			mockStmt := new(databasemock.MockStmt)
			mockResult := new(databasemock.MockResult)
			mockErrorChecker := new(entitymock.MockErrorChecker)

			if tt.setupMocks != nil {
				tt.setupMocks(mockDB, mockStmt, mockResult, mockErrorChecker)
			}

			// Act
			count, err := Delete(
				mockDB, // implements Preparer
				"user_table",
				tt.selectors,
				tt.opts,
			)

			// Assert
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedCount, count)

			// Verify mock expectations
			mockDB.AssertExpectations(t)
			mockResult.AssertExpectations(t)
		})
	}
}
