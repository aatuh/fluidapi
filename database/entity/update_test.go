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

// TestUpdate tests the Update function
func TestUpdate(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(
			mockDB *databasemock.MockDB,
			mockStmt *databasemock.MockStmt,
			mockResult *databasemock.MockResult,
			mockErrorChecker *entitymock.MockErrorChecker,
		)
		updateFields  []query.UpdateField
		selectors     []clause.Selector
		expectedRows  int64
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
				mockResult.On("RowsAffected").Return(int64(2), nil)
			},
			updateFields: []query.UpdateField{
				{Field: "name", Value: "Alice"},
			},
			selectors: []clause.Selector{
				{Field: "id", Value: 1},
			},
			expectedRows:  2,
			expectedError: "",
		},
		{
			name:          "No Updates",
			setupMocks:    nil,
			updateFields:  []query.UpdateField{},
			selectors:     []clause.Selector{{Field: "id", Value: 1}},
			expectedRows:  0,
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
				mockErrorChecker.On("Check", mock.Anything).Return(errors.New("exec error"))
			},
			updateFields: []query.UpdateField{
				{Field: "name", Value: "Alice"},
			},
			selectors: []clause.Selector{
				{Field: "id", Value: 1},
			},
			expectedRows:  0,
			expectedError: "exec error",
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
			rows, err := Update(
				mockDB,
				"user_table",
				tt.selectors,
				tt.updateFields,
				mockErrorChecker,
			)

			// Assert
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedRows, rows)

			// Verify mock expectations
			mockDB.AssertExpectations(t)
			mockStmt.AssertExpectations(t)
			mockResult.AssertExpectations(t)
			mockErrorChecker.AssertExpectations(t)
		})
	}
}

// TestCheckUpdateResult tests the CheckUpdateResult function
func TestCheckUpdateResult(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(
			mockResult *databasemock.MockResult,
			mockErrorChecker *entitymock.MockErrorChecker,
		)
		inputErr      error
		expectedRows  int64
		expectedError string
	}{
		{
			name: "Normal Operation",
			setupMocks: func(
				mockResult *databasemock.MockResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				mockResult.On("RowsAffected").Return(int64(3), nil)
			},
			inputErr:      nil,
			expectedRows:  3,
			expectedError: "",
		},
		{
			name: "Exec Error",
			setupMocks: func(
				mockResult *databasemock.MockResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				mockErrorChecker.On("Check", errors.New("exec error")).
					Return(errors.New("exec error"))
			},
			inputErr:      errors.New("exec error"),
			expectedRows:  0,
			expectedError: "exec error",
		},
		{
			name: "RowsAffected Error",
			setupMocks: func(
				mockResult *databasemock.MockResult,
				mockErrorChecker *entitymock.MockErrorChecker,
			) {
				mockResult.On("RowsAffected").
					Return(int64(0), errors.New("rows affected error"))
			},
			inputErr:      nil,
			expectedRows:  0,
			expectedError: "rows affected error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockResult := new(databasemock.MockResult)
			mockErrorChecker := new(entitymock.MockErrorChecker)

			if tt.setupMocks != nil {
				tt.setupMocks(mockResult, mockErrorChecker)
			}

			// Act
			rows, err := checkUpdateResult(
				mockResult,
				tt.inputErr,
				mockErrorChecker,
			)

			// Assert
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedRows, rows)

			// Verify mock expectations
			mockResult.AssertExpectations(t)
			mockErrorChecker.AssertExpectations(t)
		})
	}
}
