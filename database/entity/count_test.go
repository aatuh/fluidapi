package entity

import (
	"errors"
	"testing"

	databasemock "github.com/pakkasys/fluidapi/database/mock"
	"github.com/pakkasys/fluidapi/database/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCount tests the Count function
func TestCount(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(
			mockDB *databasemock.MockDB,
			mockStmt *databasemock.MockStmt,
			mockRow *databasemock.MockRow,
		)
		tableName     string
		dbOptions     *query.CountOptions
		expectedCount int
		expectedError string
	}{
		{
			name: "Normal Operation",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRow *databasemock.MockRow,
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("Close").Return(nil)
				mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
				mockRow.On("Scan", mock.Anything).Return(nil)
			},
			tableName:     "test_table",
			dbOptions:     &query.CountOptions{},
			expectedCount: 0,
			expectedError: "",
		},
		{
			name: "Prepare Error",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRow *databasemock.MockRow,
			) {
				mockDB.On("Prepare", mock.Anything).
					Return(nil, errors.New("prepare error"))
			},
			tableName:     "test_table",
			dbOptions:     &query.CountOptions{},
			expectedCount: 0,
			expectedError: "prepare error",
		},
		{
			name: "QueryRow Error",
			setupMocks: func(
				mockDB *databasemock.MockDB,
				mockStmt *databasemock.MockStmt,
				mockRow *databasemock.MockRow,
			) {
				mockDB.On("Prepare", mock.Anything).Return(mockStmt, nil)
				mockStmt.On("Close").Return(nil)
				mockStmt.On("QueryRow", mock.Anything).Return(mockRow)
				mockRow.On("Scan", mock.Anything).Return(errors.New("query row error"))
			},
			tableName:     "test_table",
			dbOptions:     &query.CountOptions{},
			expectedCount: 0,
			expectedError: "query row error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDB := new(databasemock.MockDB)
			mockStmt := new(databasemock.MockStmt)
			mockRow := new(databasemock.MockRow)
			if tt.setupMocks != nil {
				tt.setupMocks(mockDB, mockStmt, mockRow)
			}

			// Act
			count, err := Count(mockDB, tt.tableName, tt.dbOptions)

			// Assert
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedCount, count)

			// Verify mock expectations
			mockDB.AssertExpectations(t)
			mockStmt.AssertExpectations(t)
			mockRow.AssertExpectations(t)
		})
	}
}
